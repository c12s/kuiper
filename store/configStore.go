package store

import (
	"context"
	"encoding/json"
	"errors"
	"kuiper/model"
	"log"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/trace"
)

var ErrorNotFound = errors.New("Config not found")
var KeyAlreadyExistsError = errors.New("Version already exists for the given ID")

// ConfigStore is used for persistance of configurations.
// Configs are uniquely identified by ID and version. Each ID can have multiple versions.
type ConfigStore interface {
	//SaveConfig() persists a config and returns it's ID as string.
	SaveConfig(ctx context.Context, cfg model.Config) (string, error)
	//GetConfig() finds a config by it's ID and version.
	GetConfig(ctx context.Context, id, ver string) (map[string]string, error)
	//Creates a new version for an already existing id
	SaveVersion(ctx context.Context, cfg model.Config, id string) error
	//Deletes a config and returns the config that was deleted
	DeleteConfig(ctx context.Context, id, ver string) (map[string]string, error)
	//Deletes all the configs with the given ID and returns them
	DeleteConfigsWithPrefix(ctx context.Context, id string) (map[string]model.Entries, error)
	//Gets all of the service's configs
	GetConfigsByService(ctx context.Context, id string) (map[string]model.Entries, error)
	//Gets only the latest config of a service
	GetLatestConfigByService(ctx context.Context, id string) (map[string]model.Entries, error)
}

func NewConfigStore(cli clientv3.Client, logger log.Logger, trace trace.Tracer) ConfigStore {
	return configStoreEtcd{cli: cli, logger: logger, trace: trace}
}

type configStoreEtcd struct {
	logger log.Logger
	cli    clientv3.Client
	trace  trace.Tracer
}

func (cStore configStoreEtcd) SaveConfig(ctx context.Context, cfg model.Config) (string, error) {
	_, span := cStore.trace.Start(ctx, "configStoreEtcd.CreateConfig")
	defer span.End()

	var id string
	id = cfg.Service

	key := makeKey(id, cfg.Version)

	jsonB, err := json.Marshal(cfg.Entries)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	op := clientv3.OpPut(key, string(jsonB))
	res, err := cStore.cli.Txn(kvCtx).If(clientv3.Compare(clientv3.Version(key), "=", 0)).Then(op).Commit()
	cancel()
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	if !res.Succeeded {
		err = KeyAlreadyExistsError
		span.RecordError(KeyAlreadyExistsError)
		return "", err
	}

	return id, nil
}

func (cStore configStoreEtcd) GetConfig(ctx context.Context, id, ver string) (map[string]string, error) {
	_, span := cStore.trace.Start(ctx, "configStoreEtcd.GetConfig")
	defer span.End()

	entries := make(map[string]string)

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	res, err := cStore.cli.KV.Get(kvCtx, makeKey(id, ver))
	cancel()
	if err != nil {
		span.RecordError(err)
		return entries, err
	}

	kvs := res.Kvs
	if len(kvs) > 0 {
		kv := kvs[0]
		data := kv.Value
		err = json.Unmarshal(data, &entries)
		if err != nil {
			span.RecordError(err)
			return entries, err
		}
		return entries, nil
	}

	return entries, ErrorNotFound
}

func (cStore configStoreEtcd) GetConfigsByService(ctx context.Context, id string) (map[string]model.Entries, error) {
	ctx, span := cStore.trace.Start(ctx, "configStoreEtcd.GetConfigsByService")
	defer span.End()

	cfgs := make(map[string]model.Entries)

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	res, err := cStore.cli.KV.Get(kvCtx, makeIdPrefix(id), clientv3.WithPrefix())
	cancel()
	if err != nil {
		span.RecordError(err)
		return cfgs, err
	}

	return cStore.decodeConfigsFromKvs(res.Kvs, ctx)
}

func (cStore configStoreEtcd) GetLatestConfigByService(ctx context.Context, id string) (map[string]model.Entries, error) {
	ctx, span := cStore.trace.Start(ctx, "configStoreEtcd.GetLatestConfigByService")
	defer span.End()

	cfgs := make(map[string]model.Entries)

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	opts := clientv3.WithLastRev()
	opts = append(opts, clientv3.WithPrefix())
	res, err := cStore.cli.KV.Get(kvCtx, makeIdPrefix(id), opts...)
	cancel()
	if err != nil {
		span.RecordError(err)
		return cfgs, err
	}

	return cStore.decodeConfigsFromKvs(res.Kvs, ctx)
}

func (cStore configStoreEtcd) decodeConfigsFromKvs(kvs []*mvccpb.KeyValue, ctx context.Context) (map[string]model.Entries, error) {
	ctx, span := cStore.trace.Start(ctx, "configStoreEtcd.decodeConfigsFromKvs")
	defer span.End()

	cfgs := make(map[string]model.Entries)
	for _, kv := range kvs {
		cfg := make(map[string]string)
		data := kv.Value
		err := json.Unmarshal(data, &cfg)
		if err != nil {
			span.RecordError(err)
			return cfgs, err
		}

		ver := getVersionFromKey(string(kv.Key))
		cfgs[string(ver)] = cfg
	}
	return cfgs, nil
}

func (cStore configStoreEtcd) getPrefixCount(ctx context.Context, id string) (int64, error) {
	_, span := cStore.trace.Start(ctx, "configStoreEtcd.getPrefixCount")
	defer span.End()

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	key := makeIdPrefix(id)
	res, err := cStore.cli.Get(kvCtx, key, clientv3.WithPrefix(), clientv3.WithCountOnly())
	cancel()
	if err != nil {
		return 0, err
	}
	return res.Count, nil
}

func (cStore configStoreEtcd) SaveVersion(ctx context.Context, cfg model.Config, id string) error {
	ctx, span := cStore.trace.Start(ctx, "configStoreEtcd.SaveVersion")
	defer span.End()

	key := makeKey(id, cfg.Version)

	jsonB, err := json.Marshal(cfg.Entries)
	if err != nil {
		span.RecordError(err)
		return err
	}

	c, err := cStore.getPrefixCount(ctx, id)
	if err != nil {
		return err
	}
	if c == 0 {
		return ErrorNotFound
	}

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	op := clientv3.OpPut(key, string(jsonB))
	res, err := cStore.cli.Txn(kvCtx).If(clientv3.Compare(clientv3.Version(key), "=", 0)).Then(op).Commit()
	cancel()
	if err != nil {
		span.RecordError(err)
		return err
	}

	if !res.Succeeded {
		return KeyAlreadyExistsError
	}

	return nil
}

func (cStore configStoreEtcd) DeleteConfig(ctx context.Context, id, ver string) (map[string]string, error) {
	_, span := cStore.trace.Start(ctx, "configStoreEtcd.DeleteConfig")
	defer span.End()

	entries := make(map[string]string)
	key := makeKey(id, ver)

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	res, err := cStore.cli.KV.Delete(kvCtx, key, clientv3.WithPrevKV())
	cancel()
	if err != nil {
		span.RecordError(err)
		return entries, err
	}

	kvs := res.PrevKvs
	if res.Deleted > 0 {
		kv := kvs[0]
		data := kv.Value
		err = json.Unmarshal(data, &entries)
		if err != nil {
			span.RecordError(err)
			return entries, err
		}
		return entries, nil
	}

	return entries, ErrorNotFound
}

func (cStore configStoreEtcd) DeleteConfigsWithPrefix(ctx context.Context, id string) (map[string]model.Entries, error) {
	ctx, span := cStore.trace.Start(ctx, "configStoreEtcd.DeleteConfigWithPrefix")
	defer span.End()

	cfgs := make(map[string]map[string]string)
	key := makeIdPrefix(id)

	kvCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	res, err := cStore.cli.KV.Delete(kvCtx, key, clientv3.WithPrevKV(), clientv3.WithPrefix())
	cancel()
	if err != nil {
		span.RecordError(err)
		return cfgs, err
	}
	if res.Deleted == 0 {
		return cfgs, ErrorNotFound
	}

	return cStore.decodeConfigsFromKvs(res.PrevKvs, ctx)
}
