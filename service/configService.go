package service

import (
	"context"
	"kuiper/model"
	"kuiper/store"
	"log"

	"go.opentelemetry.io/otel/trace"
)

type ConfigService interface {
	CreateConfig(ctx context.Context, cfg model.Config) (string, error)
	GetConfig(ctx context.Context, id, ver string) (model.Config, error)
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
	return cs.store.CreateConfig(nCtx, cfg)
}

func (cs configService) GetConfig(ctx context.Context, id, ver string) (model.Config, error) {
	nCtx, span := cs.trace.Start(ctx, "configService.GetConfig")
	defer span.End()
	return cs.store.GetConfig(nCtx, id, ver)
}
