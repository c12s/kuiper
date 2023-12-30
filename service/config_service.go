package service

import (
	"fmt"
	"sort"

	"slices"

	"github.com/c12s/kuiper/errors"
	"github.com/c12s/kuiper/model"
	"github.com/c12s/kuiper/model/response"
	"github.com/c12s/kuiper/repository"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	DefaultAppName   = "app"
	DefaultNamespace = "namespace"
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
		version.Namespace = DefaultNamespace
	}

	if version.AppName == "" {
		version.AppName = DefaultAppName
	}

	if version.VersionTag == "" {
		return model.Version{}, fmt.Errorf(errors.VersionTagIsRequired)
	}

	if version.ConfigurationID != "" {
		versions, err := cs.repo.GetPreviousVersions(version)
		if err != nil {
			return model.Version{}, err
		}

		if version.Type == model.ConfigTypeGroup {
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
			indexOfVersion := slices.IndexFunc(versions, func(v model.Version) bool {
				return v.VersionTag == version.VersionTag
			})

			if indexOfVersion != -1 {
				return model.Version{}, fmt.Errorf(errors.VersionAlreadyExist)
			}

			sort.Slice(versions, func(i, j int) bool {
				return versions[i].CreatedAt > versions[j].CreatedAt
			})

			latest := versions[0]
			version.Weight = latest.Weight + 1
			version.Diff = buildDiffMap(&version, &latest)
		}

	} else {
		version.Weight = 1
	}

	return cs.repo.CreateNewVersion(version)
}

func (cs ConfigService) ListVersions(input model.ListRequest) (*arraylist.List, error) {

	if input.Namespace == "" {
		input.Namespace = DefaultNamespace
	}

	if input.AppName == "" {
		input.AppName = DefaultAppName
	}

	versions, err := cs.repo.ListVersions(input)
	if err != nil {
		return nil, err
	}

	sortVersionsListByWeight(versions)

	return versions, nil
}

func (cs ConfigService) ListVersionsDiff(input model.ListRequest) (response *arraylist.List, err error) {

	response = arraylist.New()

	if input.Namespace == "" {
		input.Namespace = DefaultNamespace
	}

	if input.AppName == "" {
		input.AppName = DefaultAppName
	}

	versions, err := cs.repo.ListVersions(input)
	if err != nil {
		return nil, err
	}

	sortVersionsListByWeight(versions)

	if input.Type == string(model.ConfigTypeConfig) {
		response = buildListDiffsConfig(versions)
	} else if input.Type == string(model.ConfigTypeGroup) {
		response = buildListDiffsGroup(versions)
	}

	return
}

func buildDiffMap(new, latest *model.Version) (diff model.Diffs) {
	diff = make(model.Diffs, 0)
	switch new.Type {
	case model.ConfigTypeConfig:
		configNew := new.Config.(model.Config)
		configLatest := latest.Config.(model.Config)
		diff = buildConfigLabelsDiff(configNew, configLatest)

	case model.ConfigTypeGroup:
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
						Type: string(model.DiffTypeAddition),
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
						Type: string(model.DiffTypeDeletion),
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

func sortVersionsListByWeight(versions *arraylist.List) {
	for i := 0; i < versions.Size()-1; i++ {

		for j := 0; j < versions.Size()-i-1; j++ {
			iElement, _ := versions.Get(j)
			iVersion := iElement.(model.Version)
			jElement, _ := versions.Get(j + 1)
			jVersion := jElement.(model.Version)
			if iVersion.Weight > jVersion.Weight {
				versions.Swap(j, j+1)
			}
		}
	}
}

func buildListDiffsConfig(versions *arraylist.List) (built *arraylist.List) {
	built = arraylist.New()

	for i := 0; i < versions.Size(); i++ {
		element, _ := versions.Get(i)
		version := element.(model.Version)

		newDiff := response.ConfigDiff{
			VersionTag: version.VersionTag,
			Diffs:      version.Diff,
		}

		built.Add(newDiff)
	}

	return built
}

func buildListDiffsGroup(versions *arraylist.List) (built *arraylist.List) {
	built = arraylist.New()

	for i := 0; i < versions.Size(); i++ {
		element, _ := versions.Get(i)
		version := element.(model.Version)

		newDiff := response.GroupDiff{
			VersionTag: version.VersionTag,
			GroupDiffs: version.Diff,
		}

		groupObject := version.Config.(model.Group)

		configsDiff := make([]response.GroupConfigDiff, 0)
		for _, config := range groupObject.Configs {
			if config.Diff == nil {
				continue
			}
			newConfigDiff := response.GroupConfigDiff{
				ConfigID: config.ID,
				Diffs:    config.Diff,
			}
			configsDiff = append(configsDiff, newConfigDiff)
		}

		newDiff.GroupConfigsDiff = configsDiff

		built.Add(newDiff)
	}

	return built
}
