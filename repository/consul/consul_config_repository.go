package consul

import (
	"errors"
	"fmt"
	"github.com/c12s/kuiper/model"
	"github.com/hashicorp/consul/api"
	"os"
)

type ConsulConfigRepository struct {
	cli *api.Client
}

func New() (*ConsulConfigRepository, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)

	if err != nil {
		return nil, err
	}

	ccr := ConsulConfigRepository{
		cli: client,
	}

	return &ccr, nil
}

func (ccr ConsulConfigRepository) GetGroupConfigs(id string, version string, labels []model.Label) (map[string]string, error) {
	kv := ccr.cli.KV()

	configKey := constructGroupConfigFilterKey(id, version, labels)
	pairs, _, err := kv.List(configKey, nil)

	if err != nil {
		return nil, err
	}

	if pairs == nil {
		return nil, errors.New("Config not found")
	}

	configs := make(map[string]string)

	for _, pair := range pairs {
		configs[pair.Key] = string(pair.Value)
	}

	return configs, nil
}

func (ccr ConsulConfigRepository) CreateNewGroup(group model.Group) (model.Response, error) {
	kv := ccr.cli.KV()

	id := generateUUID()
	version := generateUUID()

	for _, c := range group.Configs {
		groupConfigKey := constructGroupConfigKey(id, version, c.Labels, c.Key)

		p := &api.KVPair{Key: groupConfigKey, Value: []byte(c.Value)}

		_, err := kv.Put(p, nil)

		if err != nil {
			return model.Response{}, err
		}
	}

	return model.Response{id, version}, nil
}

func (ccr ConsulConfigRepository) CreateNewGroupVersion(id string, group model.Group) (model.Response, error) {
	kv := ccr.cli.KV()

	version := generateUUID()

	for _, c := range group.Configs {
		groupConfigKey := constructGroupConfigKey(id, version, c.Labels, c.Key)

		p := &api.KVPair{Key: groupConfigKey, Value: []byte(c.Value)}

		_, err := kv.Put(p, nil)

		if err != nil {
			return model.Response{}, err
		}
	}

	return model.Response{id, version}, nil
}

func (ccr ConsulConfigRepository) UpdateGroupVersion(id string, version string, group model.Group) error {
	return nil
}

func (ccr ConsulConfigRepository) DeleteGroup(id string) error {
	return nil
}

func (ccr ConsulConfigRepository) DeleteGroupVersion(id string, version string) error {
	return nil
}
