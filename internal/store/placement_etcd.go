package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/c12s/kuiper/internal/domain"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PlacementEtcdStore struct {
	client *clientv3.Client
}

func NewPlacementEtcdStore(client *clientv3.Client) domain.PlacementStore {
	return PlacementEtcdStore{
		client: client,
	}
}

func (s PlacementEtcdStore) Place(ctx context.Context, config domain.Config, req *domain.PlacementReq) *domain.Error {
	dao := PlacementReqDAO{
		Id:         req.Id(),
		Org:        string(config.Org()),
		Name:       config.Name(),
		Version:    config.Version(),
		Node:       string(req.Node()),
		Namespace:  string(req.Namespace()),
		Status:     req.Status(),
		AcceptedAt: req.AcceptedAtUnixSec(),
		ResolvedAt: req.ResolvedAtUnixSec(),
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

func (s PlacementEtcdStore) GetPlacement(ctx context.Context, org domain.Org, name string, version string) ([]*domain.PlacementReq, *domain.Error) {
	key := PlacementReqDAO{
		Org:     string(org),
		Name:    name,
		Version: version,
	}.KeyPrefixByConfig()
	resp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, domain.NewError(domain.ErrTypeDb, err.Error())
	}

	reqs := make([]*domain.PlacementReq, 0, resp.Count)
	for _, kv := range resp.Kvs {
		dao, err := NewPlacementReqDAO(kv.Value)
		if err != nil {
			log.Println(err)
			continue
		}
		reqs = append(reqs, domain.NewPlacementReq(dao.Id, domain.Node(dao.Node), domain.Namespace(dao.Namespace), dao.Status, dao.AcceptedAt, dao.ResolvedAt))
	}

	return reqs, nil
}

type PlacementReqDAO struct {
	Id         string
	Org        string
	Name       string
	Version    string
	Node       string
	Namespace  string
	Status     domain.PlacementReqStatus
	AcceptedAt int64
	ResolvedAt int64
}

func (dao PlacementReqDAO) Key() string {
	return fmt.Sprintf("placements/standalone/%s/%s/%s/%s", dao.Org, dao.Name, dao.Version, dao.Id)
}

func (dao PlacementReqDAO) KeyPrefixByConfig() string {
	return fmt.Sprintf("placements/standalone/%s/%s/%s/", dao.Org, dao.Name, dao.Version)
}

func (dao PlacementReqDAO) Marshal() (string, error) {
	jsonBytes, err := json.Marshal(dao)
	return string(jsonBytes), err
}

func NewPlacementReqDAO(marshalled []byte) (PlacementReqDAO, error) {
	dao := &PlacementReqDAO{}
	err := json.Unmarshal(marshalled, dao)
	if err != nil {
		return PlacementReqDAO{}, err
	}
	return *dao, nil
}
