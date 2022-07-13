package service

import "github.com/c12s/kuiper/repository"

type ConfigService struct {
	repo repository.ConfigRepostory
}

func New(repo repository.ConfigRepostory) *ConfigService {
	return &ConfigService{
		repo,
	}
}

func (cs ConfigService) CreateNewGroup(version string, configs map[string]string) (string, error) {
	return "", nil
}

func (cs ConfigService) CreateNewGroupVersion(uuid string, version string, configs map[string]string) {

}

func (cs ConfigService) UpdateGroupVersion(uuid string, version string, configs map[string]string) {

}

func (cs ConfigService) DeleteGroup(uuid string) error {
	return nil
}

func (cs ConfigService) DeleteGroupVersion(uuid string, version string) error {
	return nil
}
