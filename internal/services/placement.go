package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/c12s/kuiper/internal/domain"
	"github.com/c12s/kuiper/pkg/client/agent_queue"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type PlacementService struct {
	magnetar       magnetarapi.MagnetarClient
	aq             agent_queue.AgentQueueClient
	administrator  *oortapi.AdministrationAsyncClient
	authorizer     *AuthZService
	store          domain.PlacementStore
	webhookBaseUrl string
}

func NewPlacementStore(magnetar magnetarapi.MagnetarClient, aq agent_queue.AgentQueueClient, administrator *oortapi.AdministrationAsyncClient, authorizer *AuthZService, store domain.PlacementStore, webhookBaseUrl string) *PlacementService {
	return &PlacementService{
		magnetar:       magnetar,
		aq:             aq,
		administrator:  administrator,
		authorizer:     authorizer,
		store:          store,
		webhookBaseUrl: webhookBaseUrl,
	}
}

func (s *PlacementService) Place(ctx context.Context, config domain.Config, namespace string, nodeQuery []*magnetarapi.Selector, cmd func(taskId string) ([]byte, *domain.Error), webhookPath string) ([]domain.PlacementTask, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(config.Type(), string(config.Org()), config.Name(), config.Version())) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	if !s.authorizer.Authorize(ctx, PermNsPut, OortResNamespace, namespace) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermNsPut))
	}
	// todo: check if namespace exists
	oortErr := s.administrator.SendRequest(&oortapi.CreateInheritanceRelReq{
		From: &oortapi.Resource{
			Id:   namespace,
			Kind: OortResNamespace,
		},
		To: &oortapi.Resource{
			Id:   OortConfigId(config.Type(), string(config.Org()), config.Name(), config.Version()),
			Kind: OortResConfig,
		},
	}, func(resp *oortapi.AdministrationAsyncResp) {
		log.Println(resp.Error)
	})
	if oortErr != nil {
		log.Println(oortErr)
	}

	queryReq := &magnetarapi.QueryOrgOwnedNodesReq{
		Org: string(config.Org()),
	}
	query := make([]*magnetarapi.Selector, 0)
	for _, selector := range nodeQuery {
		s := copySelector(selector)
		query = append(query, &s)
	}
	queryReq.Query = query
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("no metadata in ctx when sending req to magnetar")
	} else {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	queryResp, err := s.magnetar.QueryOrgOwnedNodes(ctx, queryReq)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeInternal, err.Error())
	}

	tasks := make([]domain.PlacementTask, 0)
	for _, node := range queryResp.Nodes {
		taskId := uuid.New().String()
		acceptedTs := time.Now().Unix()
		task := domain.NewPlacementTask(taskId, domain.Node(node.Id), domain.Namespace(namespace), domain.PlacementTaskStatusAccepted, acceptedTs, acceptedTs)
		placeErr := s.store.Place(ctx, config, task)
		if placeErr != nil {
			log.Println(placeErr)
			continue
		}
		tasks = append(tasks, *task)
		cmdMarshalled, err := cmd(taskId)
		if err != nil {
			log.Println(err)
			continue
		}
		deseminateErr := deseminateConfig(ctx, node.Id, cmdMarshalled, s.aq, s.webhookBaseUrl+webhookPath)
		if deseminateErr != nil {
			log.Println(deseminateErr)
		}
	}
	return tasks, nil
}

func (s *PlacementService) List(ctx context.Context, org domain.Org, name, version, configType string) ([]domain.PlacementTask, *domain.Error) {
	return s.store.ListByConfig(ctx, org, name, version, configType)
}

func (s *PlacementService) UpdateStatus(ctx context.Context, org domain.Org, name, version, configType, taskId string, status domain.PlacementTaskStatus) *domain.Error {
	return s.store.UpdateStatus(ctx, org, name, version, configType, taskId, status)
}

func deseminateConfig(ctx context.Context, nodeId string, cmd []byte, agentQueueClient agent_queue.AgentQueueClient, whUrl string) error {
	log.Printf("diseminating to node %s", nodeId)
	_, err := agentQueueClient.DeseminateConfig(ctx, &agent_queue.DeseminateConfigRequest{
		NodeId:  nodeId,
		Config:  cmd,
		Webhook: whUrl,
	})
	return err
}

func copySelector(selector *magnetarapi.Selector) magnetarapi.Selector {
	return magnetarapi.Selector{
		LabelKey: selector.LabelKey,
		ShouldBe: selector.ShouldBe,
		Value:    selector.Value,
	}
}
