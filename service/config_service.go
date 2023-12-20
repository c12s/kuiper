package service

import (
	"fmt"
	"sort"

	"slices"

	"github.com/c12s/kuiper/model"
	"github.com/c12s/kuiper/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	DefaultAppName = "app"
)

type ConfigService struct {
	repo   repository.ConfigRepostory
	logger *zap.Logger
}

func New(repo repository.ConfigRepostory, logger *zap.Logger) *ConfigService {
	return &ConfigService{
		repo,
		logger,
	}
}

func (cs ConfigService) CreateNewVersion(version model.Version) (model.Version, error) {

	if version.Namespace == "" {
		version.Namespace = "namespace"
	}

	if version.AppName == "" {
		version.AppName = "app"
	}

	if version.ConfigurationID != "" {
		versions, err := cs.repo.GetPreviousVersions(version)
		if err != nil {
			return model.Version{}, err
		}

		if version.Type == "group" {
			group := version.Config.(model.Group)
			for index, config := range group.Configs {
				if config.ID == "" {
					config.ID = uuid.NewString()
					group.Configs[index] = config
				}
			}

			version.Config = group
		}

		if len(versions) > 0 {
			sort.Slice(versions, func(i, j int) bool {
				return versions[i].CreatedAt > versions[j].CreatedAt
			})

			latest := versions[0]
			if latest.Tag == version.Tag {
				return model.Version{}, fmt.Errorf("this version tag already exists, please set another version tag")
			}
			version.Diff = buildDiffMap(&version, &latest)
		}

	}

	return cs.repo.CreateNewVersion(version)
}

func (cs ConfigService) ListVersions(input model.ListRequest) ([]model.Version, error) {
	if input.AppName == "" {
		input.AppName = "app"
	}

	versions, err := cs.repo.ListVersions(input)
	if err != nil {
		return nil, err
	}

	if input.SortType == model.SortTypeTimestamp {
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].CreatedAt < versions[j].CreatedAt
		})
	} else if input.SortType == model.SortTypeLexically {
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].Tag < versions[j].Tag
		})
	}

	return versions, nil
}

func buildDiffMap(new, latest *model.Version) (diff model.Diffs) {
	diff = make(model.Diffs, 0)
	switch new.Type {
	case "config":
		configNew := new.Config.(model.Config)
		configLatest := latest.Config.(model.Config)
		diff = buildConfigLabelsDiff(configNew, configLatest)

	case "group":
		groupNew := new.Config.(model.Group)
		groupLatest := latest.Config.(model.Group)
		for i, config := range groupNew.Configs {

			index := slices.IndexFunc(groupLatest.Configs, func(c model.GroupConfig) bool {
				return c.ID == config.ID
			})

			if index == -1 {
				//addition of config in group
				newDiff := model.Addition{
					DiffCommon: model.DiffCommon{
						Type: "addition",
					},
					Key:   config.ID,
					Value: config.Labels,
				}
				diff = append(diff, newDiff)

			} else {
				// diff in GroupConfig like a added, deleted or replaced labels
				configNew := config.Config
				configOld := groupLatest.Configs[index].Config
				config.Diff = buildConfigLabelsDiff(configNew, configOld)
				groupNew.Configs[i] = config
				new.Config = groupNew
			}

		}

		for _, config := range groupLatest.Configs {
			contains := slices.ContainsFunc(groupNew.Configs, func(c model.GroupConfig) bool {
				return config.ID == c.ID
			})

			if !contains {
				//deletion of config in group
				newDiff := model.Deletion{
					DiffCommon: model.DiffCommon{
						Type: "deletion",
					},
					Key:   config.ID,
					Value: config.Labels,
				}
				diff = append(diff, newDiff)
			}
		}
	}

	return
}

func buildConfigLabelsDiff(new, latest model.Config) (diff model.Diffs) {
	diff = make(model.Diffs, 0)

	labelsNew := new.Labels
	labelsLatest := latest.Labels

	for key, labelN := range labelsNew {
		labelL, ok := labelsLatest[key]

		var newDiff model.IDiff
		if !ok {
			newDiff = model.Addition{
				DiffCommon: model.DiffCommon{Type: string(model.DiffTypeAddition)},
				Value:      labelN,
				Key:        key,
			}

			diff = append(diff, newDiff)
		} else if ok && labelN != labelL {
			newDiff = model.Replace{
				DiffCommon: model.DiffCommon{Type: string(model.DiffTypeReplace)},
				Key:        key,
				New:        labelN,
				Old:        labelL,
			}

			diff = append(diff, newDiff)
		}

	}

	for key, labelL := range labelsLatest {
		if _, ok := labelsNew[key]; !ok {
			delDiff := model.Deletion{
				DiffCommon: model.DiffCommon{Type: string(model.DiffTypeDeletion)},
				Key:        key,
				Value:      labelL,
			}

			diff = append(diff, delDiff)
		}
	}

	return
}
