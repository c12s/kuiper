package model

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type IConfig interface {
	GetType() ConfigType
}

type ConfigWrapper struct {
	Type   ConfigType `json:"type"`
	Config IConfig    `json:"config"`
}

type Config struct {
	Labels map[string]string `json:"labels"`
}

func (c Config) GetType() ConfigType {
	return ConfigTypeConfig
}

type Group struct {
	Configs []GroupConfig `json:"configs"`
}

type GroupConfig struct {
	ID string `json:"id"`
	Config
	Diff Diffs `json:"diff"`
}

func (g Group) GetType() ConfigType {
	return ConfigTypeGroup
}

type Version struct {
	Namespace       string `json:"namespace"`
	CreatorUsername string `json:"creatorUsername"`
	AppName         string `json:"appName"`
	VersionTag      string `json:"versionTag"`
	ConfigurationID string `json:"configurationID"`
	CreatedAt       int64  `json:"createdAt"`
	Weight          int64  `json:"weight"`
	ConfigWrapper
	Diff Diffs `json:"diff,omitempty"`
}

func (v *Version) UnmarshalJSON(bytes []byte) (err error) {

	params := struct {
		Namespace       string         `json:"namespace"`
		AppName         string         `json:"appName"`
		CreatorUsername string         `json:"creatorUsername"`
		VersionTag      string         `json:"versionTag"`
		CreatedAt       int64          `json:"createdAt"`
		Weight          int64          `json:"weight"`
		ConfigurationID string         `json:"configurationID"`
		Type            string         `json:"type"`
		ConfigMap       map[string]any `json:"config"`
		Diff            Diffs          `json:"diff"`
	}{}

	err = json.Unmarshal(bytes, &params)
	if err != nil {
		fmt.Printf("Error in unmarshalling type from payload because of: %s", err)
		return
	}

	v.Namespace = params.Namespace
	v.AppName = params.AppName
	v.CreatorUsername = params.CreatorUsername
	v.VersionTag = params.VersionTag
	v.CreatedAt = params.CreatedAt
	v.Weight = params.Weight
	v.Type = ConfigType(params.Type)
	v.ConfigurationID = params.ConfigurationID
	v.Diff = params.Diff

	configBytes, _ := json.Marshal(params.ConfigMap)

	switch params.Type {
	case string(ConfigTypeConfig):
		config := Config{}
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			fmt.Printf("an error occured on unmarshalling %s caused by: %s", ConfigTypeConfig, err)
			return
		}
		v.Config = config

	case string(ConfigTypeGroup):
		group := Group{}
		err = json.Unmarshal(configBytes, &group)
		if err != nil {
			fmt.Printf("an error occured on unmarshalling %s caused by: %s", ConfigTypeGroup, err)
			return
		}
		v.Config = group

	default:
		fmt.Printf("unknown config type: %s", params.Type)
	}
	return
}

type Diffs []IDiff

func (diffs *Diffs) UnmarshalJSON(bytes []byte) (err error) {

	diffsSliceOfMap := make([]map[string]any, 0)

	err = json.Unmarshal(bytes, &diffsSliceOfMap)
	if err != nil {
		fmt.Println("error in unmarshalling diffs to sliceOfMap")
		return
	}

	for _, diffMap := range diffsSliceOfMap {
		diffType := diffMap["type"]

		diffBytes, _ := json.Marshal(diffMap)
		switch diffType {
		case string(DiffTypeAddition):
			addition := Addition{}
			json.Unmarshal(diffBytes, &addition)
			*diffs = append(*diffs, addition)
		case string(DiffTypeReplace):
			replace := Replace{}
			json.Unmarshal(diffBytes, &replace)
			*diffs = append(*diffs, replace)
		case string(DiffTypeDeletion):
			deletion := Deletion{}
			json.Unmarshal(diffBytes, &deletion)
			*diffs = append(*diffs, deletion)
		}
	}

	return nil
}

type IDiff interface {
	Diff()
}

type DiffType string

const (
	DiffTypeAddition DiffType = "addition"
	DiffTypeReplace  DiffType = "replace"
	DiffTypeDeletion DiffType = "deletion"
)

func GetDiffTypeValues() []string {
	return []string{
		string(DiffTypeAddition),
		string(DiffTypeReplace),
		string(DiffTypeDeletion),
	}
}

func (dt *DiffType) IsValid() bool {
	if dt != nil && slices.Contains(GetDiffTypeValues(), string(*dt)) {
		return true
	}

	return false
}

type DiffCommon struct {
	Type string `json:"type"`
}

type Addition struct {
	DiffCommon
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (Addition) Diff() {}

type Replace struct {
	DiffCommon
	Key string `json:"key"`
	New string `json:"new"`
	Old string `json:"old"`
}

func (Replace) Diff() {}

type Deletion struct {
	DiffCommon
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (Deletion) Diff() {}

type ConfigType string

const (
	ConfigTypeConfig ConfigType = "config"
	ConfigTypeGroup  ConfigType = "group"
)

func GetConfigTypeValues() []string {
	return []string{string(ConfigTypeConfig), string(ConfigTypeGroup)}
}

func (ct *ConfigType) IsValid() bool {

	if ct != nil && slices.Contains(GetConfigTypeValues(), string(*ct)) {
		return true
	}

	return false
}

type ListRequest struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Namespace   string `json:"namespace"`
	AppName     string `json:"appName"`
	FromVersion string `json:"fromVersion"`
	WithFrom    bool   `json:"withFrom"`
	ToVersion   string `json:"toVersion"`
	WithTo      bool   `json:"withTo"`
	SortType    SortType
}

type SortType string

var (
	SortTypeLexically SortType = "lexically"
	SortTypeTimestamp SortType = "timestamp"
)

func ParseListRequest(
	typeRaw,
	idRaw,
	namespaceRaw,
	appNameRaw,
	fromVersionRaw,
	withFromRaw,
	toVersionRaw,
	withToRaw,
	sortTypeRaw string,
) (listRequest ListRequest, err error) {

	if typeRaw != "" {
		listRequest.Type = typeRaw
	} else {
		err = fmt.Errorf("type cannot be empty. applicable values: config, group")
		return
	}

	if idRaw != "" {
		listRequest.ID = idRaw
	} else {
		err = fmt.Errorf("uuid cannot be empty, please provide valid uuid")
		return
	}

	if namespaceRaw != "" {
		listRequest.Namespace = namespaceRaw
	}

	if appNameRaw != "" {
		listRequest.AppName = appNameRaw
	}

	if fromVersionRaw != "" {
		listRequest.FromVersion = fromVersionRaw
	}

	if withFromRaw != "" {
		listRequest.WithFrom, _ = strconv.ParseBool(withFromRaw)
	}

	if toVersionRaw != "" {
		listRequest.ToVersion = toVersionRaw
	}

	if withToRaw != "" {
		listRequest.WithTo, _ = strconv.ParseBool(withToRaw)
	}

	listRequest.SortType = SortTypeLexically
	if sortTypeRaw != "" {
		if strings.EqualFold(sortTypeRaw, string(SortTypeTimestamp)) {
			listRequest.SortType = SortTypeTimestamp
		}
	}

	return
}
