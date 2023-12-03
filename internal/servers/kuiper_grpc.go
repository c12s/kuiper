package servers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/kuiper/pkg/client/agent_queue"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slices"
)

type KuiperGrpcServer struct {
	api.UnimplementedKuiperServer
	configs map[string]*api.ConfigGroup
	// kojim cvorovima su dodeljene koje konfiguracije
	nodeConfigs map[string][]string
	// kojim ns-ovima su dodeljene koje konfiguracije
	nsConfigs        map[string][]string
	conn             *nats.Conn
	magnetar         magnetarapi.MagnetarClient
	agentQueueClient agent_queue.AgentQueueClient
	evaluator        oortapi.OortEvaluatorClient
	administrator    *oortapi.AdministrationAsyncClient
}

func NewKuiperServer(conn *nats.Conn, magnetar magnetarapi.MagnetarClient, evaluator oortapi.OortEvaluatorClient, administrator *oortapi.AdministrationAsyncClient, agentQueueClient agent_queue.AgentQueueClient) api.KuiperServer {
	return &KuiperGrpcServer{
		configs:          make(map[string]*api.ConfigGroup),
		nodeConfigs:      make(map[string][]string),
		nsConfigs:        make(map[string][]string),
		conn:             conn,
		magnetar:         magnetar,
		agentQueueClient: agentQueueClient,
		evaluator:        evaluator,
		administrator:    administrator,
	}
}

func (k KuiperGrpcServer) PutConfigGroup(ctx context.Context, req *api.PutConfigGroupReq) (*api.PutConfigGroupResp, error) {
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
	_, ok := k.configs[key]
	if ok {
		return nil, errors.New("config group version already exists")
	}
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
	log.Printf("NEW CONFIGS %v\n", k.configs)
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
	log.Printf("GROUP IF %s\n", groupId)
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
	if !slices.Contains(k.nsConfigs[req.Namespace], groupId) {
		k.nsConfigs[req.Namespace] = append(k.nodeConfigs[req.Namespace], groupId)
	}

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

	queryResp, err := k.magnetar.QueryNodes(ctx, &magnetarapi.QueryNodesReq{
		Queries: req.Queries,
	})
	log.Printf("Query Resp %+v\n", queryResp)
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
		log.Println(err)
		return nil, err
	}

	for _, node := range queryResp.Nodes {
		log.Printf("Inside with node %+v\n", node)
		if slices.Contains(k.nodeConfigs[node.Id], groupId) {
			log.Printf("Skipping node %s as configuration as already present.", node.Id)
			continue
		}

		err = deseminateConfig(ctx, node.Id, cmdMarshalled, k.agentQueueClient)
		if err != nil {
			log.Println(err)
		} else {
			k.nodeConfigs[node.Id] = append(k.nodeConfigs[node.Id], groupId)
		}
	}
	return &api.ApplyConfigGroupResp{}, nil
}

func deseminateConfig(ctx context.Context, nodeId string, cmd []byte, agentQueueClient agent_queue.AgentQueueClient) error {
	log.Printf("Deseminating to node %s", nodeId)
	_, err := agentQueueClient.DeseminateConfig(ctx, &agent_queue.DeseminateConfigRequest{
		NodeId: nodeId,
		Config: cmd,
	})

	return err
}

func groupId(group *api.ConfigGroup) string {
	return fmt.Sprintf("%s/%s/v%d", group.OrgId, group.Name, group.Version)
}
