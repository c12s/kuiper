package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/c12s/kuiper/internal/domain"
	"github.com/c12s/kuiper/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"google.golang.org/protobuf/proto"
)

type ConfigGroupService struct {
	administrator *oortapi.AdministrationAsyncClient
	authorizer    *AuthZService
	store         domain.ConfigGroupStore
	placements    *PlacementService
}

func NewConfigGroupService(evaluator oortapi.OortEvaluatorClient, administrator *oortapi.AdministrationAsyncClient, authorizer *AuthZService, store domain.ConfigGroupStore, placements *PlacementService) *ConfigGroupService {
	return &ConfigGroupService{
		administrator: administrator,
		authorizer:    authorizer,
		store:         store,
		placements:    placements,
	}
}

func (s *ConfigGroupService) Put(ctx context.Context, config *domain.ConfigGroup) (*domain.ConfigGroup, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigPut, OortResOrg, string(config.Org())) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigPut))
	}
	config.SetCreatedAt(time.Now())
	err := s.store.Put(ctx, config)
	if err != nil {
		return nil, err
	}
	err2 := s.administrator.SendRequest(&oortapi.CreateInheritanceRelReq{
		From: &oortapi.Resource{
			Id:   string(config.Org()),
			Kind: OortResOrg,
		},
		To: &oortapi.Resource{
			Id:   OortConfigId(config.Type(), string(config.Org()), config.Name(), config.Version()),
			Kind: OortResConfig,
		},
	}, func(resp *oortapi.AdministrationAsyncResp) {
		log.Println(resp.Error)
	})
	if err2 != nil {
		log.Println(err2)
	}
	return config, nil
}

func (s *ConfigGroupService) Get(ctx context.Context, org domain.Org, name, version string) (*domain.ConfigGroup, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeGroup, string(org), name, version)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	return s.store.Get(ctx, org, name, version)
}

func (s *ConfigGroupService) List(ctx context.Context, org domain.Org) ([]*domain.ConfigGroup, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResOrg, string(org)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	return s.store.List(ctx, org)
}

func (s *ConfigGroupService) Delete(ctx context.Context, org domain.Org, name, version string) (*domain.ConfigGroup, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigPut, OortResConfig, OortConfigId(domain.ConfTypeGroup, string(org), name, version)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigPut))
	}
	return s.store.Delete(ctx, org, name, version)
}

func (s *ConfigGroupService) Diff(ctx context.Context, referenceOrg domain.Org, referenceName, referenceVersion string, diffOrg domain.Org, diffName, diffVersion string) (map[string][]domain.Diff, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeGroup, string(referenceOrg), referenceName, referenceVersion)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeGroup, string(diffOrg), diffName, diffVersion)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	reference, err := s.store.Get(ctx, referenceOrg, referenceName, referenceVersion)
	if err != nil {
		return nil, err
	}
	diff, err := s.store.Get(ctx, diffOrg, diffName, diffVersion)
	if err != nil {
		return nil, err
	}
	return reference.Diff(diff), nil
}

func (s *ConfigGroupService) Place(ctx context.Context, org domain.Org, name, version, namespace string, nodeQuery []*magnetarapi.Selector) ([]domain.PlacementTask, *domain.Error) {
	config, err := s.store.Get(ctx, org, name, version)
	if err != nil {
		return nil, err
	}
	return s.placements.Place(ctx, config, namespace, nodeQuery, func(taskId string) ([]byte, *domain.Error) {
		cmd := &api.ApplyConfigGroupCommand{
			TaskId:    taskId,
			Namespace: namespace,
			Group: &api.ConfigGroup{
				Organization: string(config.Org()),
				Name:         config.Name(),
				Version:      config.Version(),
				CreatedAt:    config.CreatedAtUTC().String(),
				ParamSets:    mapParamSets(config.ParamSets()),
			},
		}
		cmdMarshalled, err := proto.Marshal(cmd)
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
		}
		return cmdMarshalled, nil
	})
}

func (s *ConfigGroupService) ListPlacementTasks(ctx context.Context, org domain.Org, name, version string) ([]domain.PlacementTask, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeGroup, string(org), name, version)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	return s.placements.List(ctx, org, name, version, domain.ConfTypeGroup)
}

func mapParamSets(paramSets []domain.NamedParamSet) []*api.NamedParamSet {
	protoParamSets := make([]*api.NamedParamSet, 0)
	for _, paramSet := range paramSets {
		params := mapParamSet(paramSet.ParamSet())
		protoParamSets = append(protoParamSets, &api.NamedParamSet{Name: paramSet.Name(), ParamSet: params})
	}
	return protoParamSets
}
