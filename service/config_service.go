package service

import (
	"github.com/c12s/kuiper/model"
	"github.com/c12s/kuiper/repository"
)

type ConfigService struct {
	repo repository.ConfigRepostory
}

func New(repo repository.ConfigRepostory) *ConfigService {
	return &ConfigService{
		repo,
	}
}

func (cs ConfigService) GetGroupConfigs(id string, version string, labels []model.Label) (map[string]string, error) {
	return cs.repo.GetGroupConfigs(id, version, labels)
}

func (cs ConfigService) CreateNewGroup(group model.Group) (model.Response, error) {
	return cs.repo.CreateNewGroup(group)
}

func (cs ConfigService) CreateNewGroupVersion(id string, group model.Group) (model.Response, error) {
	return cs.repo.CreateNewGroupVersion(id, group)
}

func (cs ConfigService) UpdateGroupVersion(uuid string, version string, configs map[string]string) {

}

func (cs ConfigService) DeleteGroup(uuid string) error {
	return nil
}

func (cs ConfigService) DeleteGroupVersion(uuid string, version string) error {
	return nil
}
