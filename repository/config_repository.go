package repository

import "github.com/c12s/kuiper/model"

type ConfigRepostory interface {
	GetPreviousVersions(version model.Version) ([]model.Version, error)
	CreateNewVersion(version model.Version) (model.Version, error)
	ListVersions(input model.ListRequest) ([]model.Version, error)
}
