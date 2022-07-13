package consul

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

const (
	groups              = "groups/%s/%s/%s/"
	groupsNoLabels      = "groups/%s/%s/"
	groupConfig         = "groups/%s/%s/%s/%s/"
	groupConfigNoLabels = "groups/%s/%s/%s/"
)

func generateUUID() string {
	return uuid.New().String()
}

func constructGroupKey(ctx context.Context, id string, version string, labels string) string {
	if labels == "" {
		return fmt.Sprintf(groupsNoLabels, id, version)
	} else {
		return fmt.Sprintf(groups, id, version, labels)
	}
}

func constructGroupConfigKey(ctx context.Context, groupId string, configId string, version string, labels string) string {
	if labels == "" {
		return fmt.Sprintf(groupConfigNoLabels, groupId, version, configId)
	} else {
		return fmt.Sprintf(groupConfig, groupId, version, labels, configId)
	}
}
