package store

import (
	"context"
	"encoding/json"
	"fmt"
	"kuiper/model"
	"log"
	"time"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/trace"
)

//ConfigStore is used for persistance of configurations.
//Configs are uniquely identified by ID and version. Each ID can have multiple versions.
type ConfigStore interface {
	//CreateConfig() persists a config and returns it's ID as string.
	CreateConfig(ctx context.Context, cfg model.Config) (string, error)
	//GetConfig() finds a config by it's ID and version.
	GetConfig(ctx context.Context, id, ver string) (model.Config, error)
}

func NewConfigStore(cli clientv3.Client, logger log.Logger, trace trace.Tracer) ConfigStore {
	return configStore{cli: cli, logger: logger, trace: trace}
}

type configStore struct {
	logger log.Logger
	cli    clientv3.Client
	trace  trace.Tracer
}

func makeKey(id, ver string) string {
	key := fmt.Sprintf("config/%s/%s/", id, ver)
	return key
}

func (cStore configStore) CreateConfig(ctx context.Context, cfg model.Config) (string, error) {
	_, span := cStore.trace.Start(ctx, "configStoreEtcd.CreateConfig")
	defer span.End()
	id := uuid.NewString()

	key := makeKey(id, cfg.Version)

	jsonB, err := json.Marshal(cfg)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	_, err = cStore.cli.KV.Put(kvCtx, key, string(jsonB))
	cancel()
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	return id, nil
}

func (cStore configStore) GetConfig(ctx context.Context, id, ver string) (model.Config, error) {
	_, span := cStore.trace.Start(ctx, "configStoreEtcd.GetConfig")
	defer span.End()

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	res, err := cStore.cli.KV.Get(kvCtx, makeKey(id, ver))
	cancel()
	if err != nil {
		span.RecordError(err)
		return model.Config{}, err
	}

	kvs := res.Kvs
	kv := kvs[0]

	data := kv.Value
	var cfg model.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		span.RecordError(err)
		return model.Config{}, nil
	}

	return cfg, nil
}
