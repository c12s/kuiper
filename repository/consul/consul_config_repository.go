package consul

import (
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
	return make(map[string]string), nil
}

func (ccr ConsulConfigRepository) CreateNewGroup(group model.Group) (model.Response, error) {
	return model.Response{}, nil
}

func (ccr ConsulConfigRepository) CreateNewGroupVersion(id string, group model.Group) (model.Response, error) {
	return model.Response{}, nil
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
