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

type ConfigService interface {
	//Checks if cfg is a valid config and tries to persist it.
	CreateConfig(ctx context.Context, cfg model.Config) (string, error)
	//Finds a config by id and version
	GetConfig(ctx context.Context, id, ver string) (map[string]string, error)
	//Creates a new version of already existing config
	CreateNewVersion(ctx context.Context, cfg model.Config, id string) error
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
