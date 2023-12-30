package response

import "github.com/c12s/kuiper/model"

type ConfigDiff struct {
	VersionTag string      `json:"versionTag"`
	Diffs      model.Diffs `json:"diffs"`
}

type GroupDiff struct {
	VersionTag       string            `json:"versionTag"`
	GroupConfigsDiff []GroupConfigDiff `json:"groupConfigsDiff"`
	GroupDiffs       model.Diffs       `json:"groupDiffs"`
}

type GroupConfigDiff struct {
	ConfigID string      `json:"configID"`
	Diffs    model.Diffs `json:"diffs"`
}
