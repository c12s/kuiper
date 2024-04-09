package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/c12s/kuiper/internal/domain"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type StandaloneConfigEtcdStore struct {
	client *clientv3.Client
}

func NewStandaloneConfigEtcdStore(client *clientv3.Client) domain.StandaloneConfigStore {
	return StandaloneConfigEtcdStore{
		client: client,
	}
}

func (s StandaloneConfigEtcdStore) Put(ctx context.Context, config *domain.StandaloneConfig) *domain.Error {
	dao := StandaloneConfigDAO{
		Org:       string(config.Org()),
		Name:      config.Name(),
		Version:   config.Version(),
		CreatedAt: config.CreatedAtUnixSec(),
		ParamSet:  config.ParamSet(),
	}

	key := dao.Key()
	value, err := dao.Marshal()
	if err != nil {
		return domain.NewError(domain.ErrTypeMarshalSS, err.Error())
	}

	resp, err := s.client.KV.Txn(ctx).If(clientv3.CreateRevision(key)).Then(clientv3.OpPut(key, value)).Commit()
	if !resp.Succeeded {
		return domain.NewError(domain.ErrTypeVersionExists, fmt.Sprintf("standalone config (Org: %s, name: %s, version: %s) already exists", config.Org(), config.Name(), config.Version()))
	}
	if err != nil {
		return domain.NewError(domain.ErrTypeDb, err.Error())
	}
	return nil
}

func (s StandaloneConfigEtcdStore) Get(ctx context.Context, org domain.Org, name, version string) (*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{
		Org:     string(org),
		Name:    name,
		Version: version,
	}.Key()
	resp, err := s.client.KV.Get(ctx, key)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	if resp.Count == 0 {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("standalone config (Org: %s, name: %s, version: %s) not found", org, name, version))
	}

	dao, err := NewStandaloneConfigDAO(resp.Kvs[0].Value)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
	}

	paramSet := domain.NewParamSet(dao.Name, dao.ParamSet)
	return domain.InitStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, *paramSet), nil
}

func (s StandaloneConfigEtcdStore) List(ctx context.Context, org domain.Org) ([]*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{
		Org: string(org),
	}.KeyPrefixAll()
	resp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	configs := make([]*domain.StandaloneConfig, 0, resp.Count)
	for _, kv := range resp.Kvs {
		dao, err := NewStandaloneConfigDAO(kv.Value)
		if err != nil {
			log.Println(err)
			continue
		}
		paramSet := domain.NewParamSet(dao.Name, dao.ParamSet)
		configs = append(configs, domain.InitStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, *paramSet))
	}

	return configs, nil
}

func (s StandaloneConfigEtcdStore) Delete(ctx context.Context, org domain.Org, name, version string) (*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{
		Org:     string(org),
		Name:    name,
		Version: version,
	}.Key()
	resp, err := s.client.KV.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	if len(resp.PrevKvs) == 0 {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("standalone config (Org: %s, name: %s, version: %s) not found", org, name, version))
	}

	dao, err := NewStandaloneConfigDAO(resp.PrevKvs[0].Value)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
	}

	paramSet := domain.NewParamSet(dao.Name, dao.ParamSet)
	return domain.InitStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, *paramSet), nil
}

type StandaloneConfigDAO struct {
	Org       string
	Name      string
	Version   string
	CreatedAt int64
	ParamSet  map[string]string
}

func (dao StandaloneConfigDAO) Key() string {
	return fmt.Sprintf("standalone/%s/%s/%s", dao.Org, dao.Name, dao.Version)
}

func (dao StandaloneConfigDAO) KeyPrefixAll() string {
	return fmt.Sprintf("standalone/%s/", dao.Org)
}

func (dao StandaloneConfigDAO) Marshal() (string, error) {
	jsonBytes, err := json.Marshal(dao)
	return string(jsonBytes), err
}

func NewStandaloneConfigDAO(marshalled []byte) (StandaloneConfigDAO, error) {
	dao := &StandaloneConfigDAO{}
	err := json.Unmarshal(marshalled, dao)
	if err != nil {
		return StandaloneConfigDAO{}, err
	}
	return *dao, nil
}
