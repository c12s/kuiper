package consul

import (
	"fmt"
	"github.com/c12s/kuiper/model"
	"github.com/google/uuid"
)

const (
	groupConfig       = "groups/%s/%s/%s/%s/"
	groupConfigFilter = "groups/%s/%s/%s/"
)

func generateUUID() string {
	return uuid.New().String()
}

func constructGroupConfigKey(id string, version string, labels []model.Label, key string) string {
	if len(labels) == 0 {
		return fmt.Sprintf(groupConfig, id, version, "none", key)
	} else {
		labelsString := DecodeLabels(labels)

		return fmt.Sprintf(groupConfig, id, version, labelsString, key)
	}
}

func constructGroupConfigFilterKey(id string, version string, labels []model.Label) string {
	if len(labels) == 0 {
		return fmt.Sprintf(groupConfigFilter, id, version, "none")
	} else {
		labelsString := DecodeLabels(labels)

		return fmt.Sprintf(groupConfigFilter, id, version, labelsString)
	}
}
