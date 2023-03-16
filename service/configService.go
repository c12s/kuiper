package service

import (
	"context"
	"kuiper/model"
	"kuiper/store"
)

type ConfigService interface {
	CreateConfig(ctx context.Context, config model.Config) (string, error)
	GetConfig(id string) (model.Config, error)
}

func NewConfigService(cs store.ConfigStore) ConfigService {
	return configService{store: cs}
}

type configService struct {
	store store.ConfigStore
}

func (cs configService) CreateConfig(ctx context.Context, config model.Config) (string, error) {
	return cs.store.CreateConfig(config, context.TODO())
}

func (cs configService) GetConfig(id string) (model.Config, error) {
	return cs.store.GetConfig(id, context.TODO())
}
