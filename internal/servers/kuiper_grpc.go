package servers

import (
	"context"
	"errors"
	"fmt"
	"github.com/c12s/kuiper/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slices"
	"log"
)

type KuiperGrpcServer struct {
	api.UnimplementedKuiperServer
	configs map[string]*api.ConfigGroup
	// kojim cvorovima su dodeljene koje konfiguracije
	nodeConfigs map[string][]string
	// kojim ns-ovima su dodeljene koje konfiguracije
	nsConfigs     map[string][]string
	conn          *nats.Conn
	magnetar      magnetarapi.MagnetarClient
	evaluator     oortapi.OortEvaluatorClient
	administrator *oortapi.AdministrationAsyncClient
}

func NewKuiperServer(conn *nats.Conn, magnetar magnetarapi.MagnetarClient, evaluator oortapi.OortEvaluatorClient, administrator *oortapi.AdministrationAsyncClient) api.KuiperServer {
	return &KuiperGrpcServer{
		configs:       make(map[string]*api.ConfigGroup),
		nodeConfigs:   make(map[string][]string),
		nsConfigs:     make(map[string][]string),
		conn:          conn,
		magnetar:      magnetar,
		evaluator:     evaluator,
		administrator: administrator,
	}
}

func (k KuiperGrpcServer) PutConfigGroup(ctx context.Context, req *api.PutConfigGroupReq) (*api.PutConfigGroupResp, error) {
	// proveri da li korisnik sme da dodaje/manja konfiguraciju unutar zadate organizacije
	resp, err := k.evaluator.Authorize(ctx, &oortapi.AuthorizationReq{
		Subject: &oortapi.Resource{
			Id:   req.SubId,
			Kind: req.SubKind,
		},
		Object: &oortapi.Resource{
			Id:   req.Group.OrgId,
			Kind: "org",
		},
		PermissionName: "config.put",
	})
	if err != nil {
		return nil, err
	}
	if !resp.Authorized {
		return nil, errors.New("unauthorized")
	}
	key := groupId(req.Group)
	// ne treba da postoji trenutna verzija
	_, ok := k.configs[key]
	if ok {
		return nil, errors.New("config group version already exists")
	}
	// treba da postoji prethodna verzija
	if req.Group.Version > 1 {
		prevGroup := proto.Clone(req.Group).(*api.ConfigGroup)
		prevGroup.Version = req.Group.Version - 1
		prevKey := groupId(prevGroup)
		_, ok := k.configs[prevKey]
		if !ok {
			return nil, errors.New("previous config version not found")
		}
	}
	k.configs[key] = req.Group
	// javi oort-u da je nova config dodata u org
	err = k.administrator.SendRequest(&oortapi.CreateInheritanceRelReq{
		From: &oortapi.Resource{
			Id:   req.Group.OrgId,
			Kind: "org",
		},
		To: &oortapi.Resource{
			Id:   key,
			Kind: "config",
		},
	}, func(resp *oortapi.AdministrationAsyncResp) {
		log.Println(resp.Error)
	})
	if err != nil {
		log.Println(err)
	}
	return &api.PutConfigGroupResp{
		Group: req.Group,
	}, nil
}

func (k KuiperGrpcServer) ApplyConfigGroup(ctx context.Context, req *api.ApplyConfigGroupReq) (*api.ApplyConfigGroupResp, error) {
	group := &api.ConfigGroup{
		Name:    req.GroupName,
		OrgId:   req.OrgId,
		Version: req.Version,
	}
	groupId := groupId(group)
	// authorize - da li sub sme da pristupi konfiguraciji
	resp, err := k.evaluator.Authorize(ctx, &oortapi.AuthorizationReq{
		Subject: &oortapi.Resource{
			Id:   req.SubId,
			Kind: req.SubKind,
		},
		Object: &oortapi.Resource{
			Id:   groupId,
			Kind: "config",
		},
		PermissionName: "config.get",
	})
	if err != nil {
		return nil, err
	}
	if !resp.Authorized {
		return nil, errors.New("unauthorized - config.get")
	}
	// check if config exists
	group, ok := k.configs[groupId]
	if !ok {
		return nil, errors.New("config not found")
	}
	// authorize - da li sme da menja namespace
	resp, err = k.evaluator.Authorize(ctx, &oortapi.AuthorizationReq{
		Subject: &oortapi.Resource{
			Id:   req.SubId,
			Kind: req.SubKind,
		},
		Object: &oortapi.Resource{
			Id:   req.Namespace,
			Kind: "namespace",
		},
		PermissionName: "namespace.putconfig",
	})
	if err != nil {
		return nil, err
	}
	if !resp.Authorized {
		return nil, errors.New("unauthorized - namespace.put")
	}
	// todo: check if namespace exists
	// dodaj config u ns ako vec nije postojao
	if !slices.Contains(k.nsConfigs[req.Namespace], groupId) {
		k.nsConfigs[req.Namespace] = append(k.nodeConfigs[req.Namespace], groupId)
	}
	// javi oort-u da je nova config dodata u ns
	err = k.administrator.SendRequest(&oortapi.CreateInheritanceRelReq{
		From: &oortapi.Resource{
			Id:   req.Namespace,
			Kind: "namespace",
		},
		To: &oortapi.Resource{
			Id:   groupId,
			Kind: "config",
		},
	}, func(resp *oortapi.AdministrationAsyncResp) {
		log.Println(resp.Error)
	})
	if err != nil {
		log.Println(err)
	}
	// query nodes
	queryReq := &magnetarapi.QueryOrgOwnedNodesReq{
		Org: req.OrgId,
	}
	// mora rucno da se kopira jedan po jedan selektor
	// todo: izmeni ovo ako je ikako moguce
	query := make([]*magnetarapi.Selector, 0)
	for _, selector := range req.Query {
		s := copySelector(*selector)
		query = append(query, &s)
	}
	queryReq.Query = query
	queryResp, err := k.magnetar.QueryOrgOwnedNodes(ctx, queryReq)
	if err != nil {
		return nil, err
	}
	// send config to each node
	cmd := api.ApplyConfigCommand{
		Id:      groupId,
		Configs: group.Configs,
	}
	cmdMarshalled, err := cmd.Marshal()
	if err != nil {
		return nil, err
	}
	// todo: scheduler notifies the nodes, we should only generate subjects and marshalled messages here
	for _, node := range queryResp.Nodes {
		// preskoci cvor ako on vec ima tu konfiguraciju
		if slices.Contains(k.nodeConfigs[node.Id], groupId) {
			continue
		}
		err = k.conn.Publish(api.Subject(node.Id), cmdMarshalled)
		if err != nil {
			log.Println(err)
		} else {
			// cvor sada ima tu konfiguraciju
			k.nodeConfigs[node.Id] = append(k.nodeConfigs[node.Id], groupId)
		}
	}
	return &api.ApplyConfigGroupResp{}, nil
}

func groupId(group *api.ConfigGroup) string {
	return fmt.Sprintf("%s/%s/v%d", group.OrgId, group.Name, group.Version)
}

func copySelector(selector magnetarapi.Selector) magnetarapi.Selector {
	return magnetarapi.Selector{
		LabelKey: selector.LabelKey,
		ShouldBe: selector.ShouldBe,
		Value:    selector.Value,
	}
}
