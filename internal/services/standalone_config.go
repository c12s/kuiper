package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/c12s/kuiper/internal/domain"
	"github.com/c12s/kuiper/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	quasarapi "github.com/c12s/quasar/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

type StandaloneConfigService struct {
	administrator *oortapi.AdministrationAsyncClient
	authorizer    *AuthZService
	store         domain.StandaloneConfigStore
	placements    *PlacementService
	quasar        quasarapi.ConfigSchemaServiceClient
}

func NewStandaloneConfigService(administrator *oortapi.AdministrationAsyncClient, authorizer *AuthZService, store domain.StandaloneConfigStore, placements *PlacementService, quasar quasarapi.ConfigSchemaServiceClient) *StandaloneConfigService {
	return &StandaloneConfigService{
		administrator: administrator,
		authorizer:    authorizer,
		store:         store,
		placements:    placements,
		quasar:        quasar,
	}
}

func (s *StandaloneConfigService) Put(ctx context.Context, config *domain.StandaloneConfig, schema *quasarapi.ConfigSchemaDetails) (*domain.StandaloneConfig, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigPut, OortResOrg, string(config.Org())) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigPut))
	}
	if schema != nil {
		configMap := make(map[string]map[string]string)
		configMap[config.Name()] = make(map[string]string)
		for key, value := range config.ParamSet() {
			configMap[config.Name()][key] = value
		}
		yamlBytes, err := yaml.Marshal(configMap)
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("no metadata in ctx when sending req to magnetar")
		} else {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		resp, err := s.quasar.ValidateConfiguration(ctx, &quasarapi.ValidateConfigurationRequest{
			SchemaDetails: schema,
			Configuration: string(yamlBytes),
		})
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeInternal, err.Error())
		}
		if !resp.IsValid {
			return nil, domain.NewError(domain.ErrTypeSchemaInvalid, resp.Message)
		}
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

func (s *StandaloneConfigService) Get(ctx context.Context, org domain.Org, name, version string) (*domain.StandaloneConfig, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeStandalone, string(org), name, version)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	return s.store.Get(ctx, org, name, version)
}

func (s *StandaloneConfigService) List(ctx context.Context, org domain.Org) ([]*domain.StandaloneConfig, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResOrg, string(org)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	return s.store.List(ctx, org)
}

func (s *StandaloneConfigService) Delete(ctx context.Context, org domain.Org, name, version string) (*domain.StandaloneConfig, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigPut, OortResConfig, OortConfigId(domain.ConfTypeStandalone, string(org), name, version)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigPut))
	}
	return s.store.Delete(ctx, org, name, version)
}

func (s *StandaloneConfigService) Diff(ctx context.Context, referenceOrg domain.Org, referenceName, referenceVersion string, diffOrg domain.Org, diffName, diffVersion string) ([]domain.Diff, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeStandalone, string(referenceOrg), referenceName, referenceVersion)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeStandalone, string(diffOrg), diffName, diffVersion)) {
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
	return diff.Diff(reference), nil
}

func (s *StandaloneConfigService) Place(ctx context.Context, org domain.Org, name, version, namespace string, strategy *api.PlaceReq_Strategy) ([]domain.PlacementTask, *domain.Error) {
	config, err := s.store.Get(ctx, org, name, version)
	if err != nil {
		return nil, err
	}
	return s.placements.Place(ctx, config, namespace, strategy, func(taskId string) ([]byte, *domain.Error) {
		config := &api.StandaloneConfig{
			Organization: string(config.Org()),
			Name:         config.Name(),
			Version:      config.Version(),
			CreatedAt:    config.CreatedAtUTC().String(),
			ParamSet:     mapParamSet(config.ParamSet()),
		}
		configMarshalled, err := proto.Marshal(config)
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
		}
		cmd := &api.ApplyConfigCommand{
			TaskId:    taskId,
			Namespace: namespace,
			Config:    configMarshalled,
			Type:      "standalone",
		}
		cmdMarshalled, err := proto.Marshal(cmd)
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
		}
		return cmdMarshalled, nil
	}, "/standalone")
}

func (s *StandaloneConfigService) ListPlacementTasks(ctx context.Context, org domain.Org, name, version string) ([]domain.PlacementTask, *domain.Error) {
	if !s.authorizer.Authorize(ctx, PermConfigGet, OortResConfig, OortConfigId(domain.ConfTypeStandalone, string(org), name, version)) {
		return nil, domain.NewError(domain.ErrTypeUnauthorized, fmt.Sprintf("Permission denied: %s", PermConfigGet))
	}
	return s.placements.List(ctx, org, name, version, domain.ConfTypeStandalone)
}

func mapParamSet(params map[string]string) []*api.Param {
	paramSet := make([]*api.Param, 0)
	for key, value := range params {
		paramSet = append(paramSet, &api.Param{Key: key, Value: value})
	}
	return paramSet
}
