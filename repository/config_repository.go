package repository

import (
	"github.com/c12s/kuiper/model"
	"github.com/emirpasic/gods/lists/arraylist"
)

type ConfigRepostory interface {
	GetPreviousVersions(version model.Version) ([]model.Version, error)
	CreateNewVersion(version model.Version) (model.Version, error)
	ListVersions(input model.ListRequest) (*arraylist.List, error)
}
