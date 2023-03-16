package store

import (
	"context"
	"kuiper/model"
)

type ConfigStore interface {
	CreateConfig(cfg model.Config, ctx context.Context) (string, error)
	GetConfig(id string, ctx context.Context) (model.Config, error)
}

func NewConfigStore() ConfigStore {
	return configStore{}
}

type configStore struct {
}

func (cStore configStore) CreateConfig(cfg model.Config, ctx context.Context) (string, error) {
	return "", nil
}

func (cStore configStore) GetConfig(id string, ctx context.Context) (model.Config, error) {
	return model.Config{}, nil
}
