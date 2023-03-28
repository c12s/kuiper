package service

import (
	"context"
	"errors"
	"kuiper/model"
	"kuiper/store"
	"log"

	"go.opentelemetry.io/otel/trace"
)

var NoVersionError = errors.New("Must supply version name when creating a new config")
var NoServiceNameError = errors.New("Must supply service name when creating a new config")
var DbError = errors.New("Error happened while connecting to database")

type ConfigService interface {
	//Checks if cfg is a valid config and tries to persist it.
	CreateConfig(ctx context.Context, cfg model.Config) (string, error)
	//Finds a config by id and version
	GetConfig(ctx context.Context, id, ver string) (map[string]string, error)
	//Creates a new version of already existing config
	CreateNewVersion(ctx context.Context, cfg model.Config, id string) error
	//Deletes config by id and version, returns error when config wasn't foun
	DeleteConfig(ctx context.Context, id, ver string) (cfg map[string]string, err error)
	//Deletes all configs with the given ID
	DeleteConfigsWithPrefix(ctx context.Context, id string) (deleted map[string]map[string]string, err error)
}

func NewConfigService(cs store.ConfigStore, logger log.Logger, trace trace.Tracer) ConfigService {
	return configService{store: cs, logger: logger, trace: trace}
}

type configService struct {
	store  store.ConfigStore
	logger log.Logger
	trace  trace.Tracer
}

func (cs configService) CreateConfig(ctx context.Context, cfg model.Config) (string, error) {
	nCtx, span := cs.trace.Start(ctx, "configService.CreateConfig")
	defer span.End()

	if len(cfg.Version) == 0 {
		span.RecordError(NoVersionError)
		return "", NoVersionError
	}
	if len(cfg.Service) == 0 {
		span.RecordError(NoServiceNameError)
		return "", NoServiceNameError
	}

	return cs.store.SaveConfig(nCtx, cfg)
}

func (cs configService) GetConfig(ctx context.Context, id, ver string) (map[string]string, error) {
	nCtx, span := cs.trace.Start(ctx, "configService.GetConfig")
	defer span.End()
	return cs.store.GetConfig(nCtx, id, ver)
}

func (cs configService) CreateNewVersion(ctx context.Context, cfg model.Config, id string) error {
	nCtx, span := cs.trace.Start(ctx, "configService.CreateNewVersion")
	defer span.End()

	if len(cfg.Version) == 0 {
		span.RecordError(NoVersionError)
		return NoVersionError
	}

	return cs.store.SaveVersion(nCtx, cfg, id)
}

func (cs configService) DeleteConfig(ctx context.Context, id, ver string) (cfg map[string]string, err error) {
	nCtx, span := cs.trace.Start(ctx, "configService.DeleteConfig")
	defer span.End()
	return cs.store.DeleteConfig(nCtx, id, ver)
}

func (cs configService) DeleteConfigsWithPrefix(ctx context.Context, id string) (deleted map[string]map[string]string, err error) {
	nCtx, span := cs.trace.Start(ctx, "configService.DeleteConfig")
	defer span.End()
	return cs.store.DeleteConfigsWithPrefix(nCtx, id)
}
