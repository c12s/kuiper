package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/c12s/kuiper/internal/domain"
	"go.etcd.io/etcd/clientv3"
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

	_, err = s.client.KV.Put(ctx, key, value)
	if err != nil {
		return domain.NewError(domain.ErrTypeDb, err.Error())
	}
	return nil
}

func (s StandaloneConfigEtcdStore) Get(ctx context.Context, Org, name, version string) (*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{
		Org:     Org,
		Name:    name,
		Version: version,
	}.Key()
	resp, err := s.client.KV.Get(ctx, key)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	if resp.Count == 0 {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("standalone config (Org: %s, name: %s, version: %s) not found", Org, name, version))
	}

	dao, err := NewStandaloneConfigDAO(resp.Kvs[0].Value)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
	}

	paramSet := domain.NewParamSet(dao.Name, dao.ParamSet)
	return domain.NewStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, resp.Kvs[0].ModRevision, *paramSet), nil
}

// todo: allow only for past day/week, enable compaction for older revisions
func (s StandaloneConfigEtcdStore) GetHistory(ctx context.Context, Org, name, version string) ([]*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{
		Org:     Org,
		Name:    name,
		Version: version,
	}.Key()
	resp, err := s.client.KV.Get(ctx, key)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	if resp.Count == 0 {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("standalone config (Org: %s, name: %s, version: %s) not found", Org, name, version))
	}

	latestKeyRev := resp.Kvs[0].ModRevision
	firstKeyRev := resp.Kvs[0].CreateRevision
	currVersion := resp.Kvs[0].Version

	configs := make([]*domain.StandaloneConfig, 0, latestKeyRev-firstKeyRev)
	for rev := latestKeyRev; rev >= firstKeyRev; rev-- {
		resp, err := s.client.KV.Get(ctx, key, clientv3.WithRev(rev))
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.Count == 0 {
			continue
		}
		if resp.Kvs[0].Version == currVersion {
			continue
		}
		dao, err := NewStandaloneConfigDAO(resp.Kvs[0].Value)
		if err != nil {
			log.Println(err)
			continue
		}
		paramSet := domain.NewParamSet(dao.Name, dao.ParamSet)
		configs = append(configs, domain.NewStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, rev, *paramSet))
	}

	return configs, nil
}

func (s StandaloneConfigEtcdStore) List(ctx context.Context) ([]*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{}.KeyPrefixAll()
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
		configs = append(configs, domain.NewStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, resp.Kvs[0].ModRevision, *paramSet))
	}

	return configs, nil
}

func (s StandaloneConfigEtcdStore) Delete(ctx context.Context, Org, name, version string) (*domain.StandaloneConfig, *domain.Error) {
	key := StandaloneConfigDAO{
		Org:     Org,
		Name:    name,
		Version: version,
	}.Key()
	resp, err := s.client.KV.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	if len(resp.PrevKvs) == 0 {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("standalone config (Org: %s, name: %s, version: %s) not found", Org, name, version))
	}

	dao, err := NewStandaloneConfigDAO(resp.PrevKvs[0].Value)
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeMarshalSS, err.Error())
	}

	paramSet := domain.NewParamSet(dao.Name, dao.ParamSet)
	return domain.NewStandaloneConfig(domain.Org(dao.Org), dao.Version, dao.CreatedAt, resp.PrevKvs[0].ModRevision, *paramSet), nil
}

func (s StandaloneConfigEtcdStore) AddToNamespaces(ctx context.Context, config *domain.StandaloneConfig) *domain.Error {
	for _, namespace := range config.Namespaces() {
		dao := NamespaceDAO{
			Namespace: string(namespace),
			Org:       string(config.Org()),
			Name:      config.Name(),
			Version:   config.Version(),
			Revision:  config.Revision(),
		}

		key := dao.Key()
		value, err := dao.Marshal()
		if err != nil {
			return domain.NewError(domain.ErrTypeMarshalSS, err.Error())
		}

		_, err = s.client.KV.Put(ctx, key, value)
		if err != nil {
			return domain.NewError(domain.ErrTypeDb, err.Error())
		}
	}
	return nil
}

func (s StandaloneConfigEtcdStore) AddToNodes(ctx context.Context, config *domain.StandaloneConfig) *domain.Error {
	for _, node := range config.Nodes() {
		dao := NodeDAO{
			Node:     string(node),
			Org:      string(config.Org()),
			Name:     config.Name(),
			Version:  config.Version(),
			Revision: config.Revision(),
		}

		key := dao.Key()
		value, err := dao.Marshal()
		if err != nil {
			return domain.NewError(domain.ErrTypeMarshalSS, err.Error())
		}

		_, err = s.client.KV.Put(ctx, key, value)
		if err != nil {
			return domain.NewError(domain.ErrTypeDb, err.Error())
		}
	}
	return nil
}

func (s StandaloneConfigEtcdStore) ListNamespace(ctx context.Context, namespace domain.Namespace, org domain.Org) ([]*domain.StandaloneConfig, *domain.Error) {
	key := NamespaceDAO{
		Namespace: string(namespace),
		Org:       string(org),
	}.KeyPrefixByNamespace()

	resp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	configs := make([]*domain.StandaloneConfig, 0, resp.Count)
	for _, kv := range resp.Kvs {
		dao, err := NewNamespaceDAO(kv.Value)
		if err != nil {
			log.Println(err)
			continue
		}

		configKey := StandaloneConfigDAO{
			Org:     dao.Org,
			Name:    dao.Name,
			Version: dao.Version,
		}.Key()
		resp, err := s.client.KV.Get(ctx, configKey, clientv3.WithRev(dao.Revision))
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeDb, err.Error())
		}
		if resp.Count == 0 {
			log.Printf("standalone config (Org: %s, name: %s, version: %s, revision: %d) not found", dao.Org, dao.Name, dao.Version, dao.Revision)
			continue
		}

		configDao, err := NewStandaloneConfigDAO(resp.Kvs[0].Value)
		if err != nil {
			log.Println(err)
			continue
		}

		paramSet := domain.NewParamSet(dao.Name, configDao.ParamSet)
		configs = append(configs, domain.NewStandaloneConfig(domain.Org(configDao.Org), configDao.Version, configDao.CreatedAt, resp.Kvs[0].ModRevision, *paramSet))
	}

	return configs, nil
}

func (s StandaloneConfigEtcdStore) ListNode(ctx context.Context, node domain.Node, org domain.Org) ([]*domain.StandaloneConfig, *domain.Error) {
	key := NodeDAO{
		Node: string(node),
		Org:  string(org),
	}.KeyPrefixByNode()

	resp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	configs := make([]*domain.StandaloneConfig, 0, resp.Count)
	for _, kv := range resp.Kvs {
		dao, err := NewNodeDAO(kv.Value)
		if err != nil {
			log.Println(err)
			continue
		}

		configKey := StandaloneConfigDAO{
			Org:     dao.Org,
			Name:    dao.Name,
			Version: dao.Version,
		}.Key()
		resp, err := s.client.KV.Get(ctx, configKey, clientv3.WithRev(dao.Revision))
		if err != nil {
			return nil, domain.NewError(domain.ErrTypeDb, err.Error())
		}
		if resp.Count == 0 {
			log.Printf("standalone config (Org: %s, name: %s, version: %s, revision: %d) not found", dao.Org, dao.Name, dao.Version, dao.Revision)
			continue
		}

		configDao, err := NewStandaloneConfigDAO(resp.Kvs[0].Value)
		if err != nil {
			log.Println(err)
			continue
		}

		paramSet := domain.NewParamSet(dao.Name, configDao.ParamSet)
		configs = append(configs, domain.NewStandaloneConfig(domain.Org(configDao.Org), configDao.Version, configDao.CreatedAt, resp.Kvs[0].ModRevision, *paramSet))
	}

	return configs, nil
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
	return "standalone/"
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

type NamespaceDAO struct {
	Namespace string
	Org       string
	Name      string
	Version   string
	Revision  int64
}

func (dao NamespaceDAO) Key() string {
	return fmt.Sprintf("namespaces/%s/standalone/%s/%s/%s", dao.Namespace, dao.Org, dao.Name, dao.Version)
}

func (dao NamespaceDAO) KeyPrefixByNamespace() string {
	return fmt.Sprintf("namespaces/%s/standalone/%s/", dao.Namespace, dao.Org)
}

func (dao NamespaceDAO) Marshal() (string, error) {
	jsonBytes, err := json.Marshal(dao)
	return string(jsonBytes), err
}

func NewNamespaceDAO(marshalled []byte) (NamespaceDAO, error) {
	dao := &NamespaceDAO{}
	err := json.Unmarshal(marshalled, dao)
	if err != nil {
		return NamespaceDAO{}, err
	}
	return *dao, nil
}

type NodeDAO struct {
	Node     string
	Org      string
	Name     string
	Version  string
	Revision int64
}

func (dao NodeDAO) Key() string {
	return fmt.Sprintf("nodes/%s/standalone/%s/%s/%s", dao.Node, dao.Org, dao.Name, dao.Version)
}

func (dao NodeDAO) KeyPrefixByNode() string {
	return fmt.Sprintf("nodes/%s/standalone/%s/", dao.Node, dao.Org)
}

func (dao NodeDAO) Marshal() (string, error) {
	jsonBytes, err := json.Marshal(dao)
	return string(jsonBytes), err
}

func NewNodeDAO(marshalled []byte) (NodeDAO, error) {
	dao := &NodeDAO{}
	err := json.Unmarshal(marshalled, dao)
	if err != nil {
		return NodeDAO{}, err
	}
	return *dao, nil
}
