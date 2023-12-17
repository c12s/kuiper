package repository

import "github.com/c12s/kuiper/model"

type ConfigRepostory interface {
	// GetGroupConfigs(id string, version string, labels []model.Label) (map[string]string, error)
	// CreateNewGroup(group model.Group) (model.Response, error)
	GetPreviousVersions(version model.Version) ([]model.Version, error)
	CreateNewVersion(version model.Version) (model.Version, error)
	ListVersions(input model.ListRequest) ([]model.Version, error)
	// CreateNewGroupVersion(id string, group model.Group) (model.Response, error)
	// UpdateGroupVersion(id string, version string, group model.Group) error
	// DeleteGroup(id string) error
	// DeleteGroupVersion(id string, version string) error
}
