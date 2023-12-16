package etcd

import (
	"fmt"

	"github.com/c12s/kuiper/model"
)

const (
	// 						  /namespace/appName/type/configID/v1/
	VersionKeyFormat string = "/%s/%s/%s/%s/%s"
	//								/namespace/appName/type/configID
	ConfigurationKeyFormat string = "/%s/%s/%s/%s"
)

func buildConfigurationVersionKey(version model.Version) (key string) {

	return fmt.Sprintf(VersionKeyFormat,
		version.Namespace,
		"app",
		version.ConfigWrapper.Type,
		version.ConfigurationID,
		version.Tag,
	)
}

func buildPrefixesFromListInput(listInput model.ListRequest) (prefix, from, to string) {

	prefix = fmt.Sprintf(
		ConfigurationKeyFormat,
		listInput.Namespace,
		listInput.AppName,
		listInput.Type,
		listInput.ID,
	)

	if listInput.FromVersion != "" {
		from = fmt.Sprintf(
			VersionKeyFormat,
			listInput.Namespace,
			listInput.AppName,
			listInput.Type,
			listInput.ID,
			listInput.FromVersion,
		)
	}

	if listInput.ToVersion != "" {
		to = fmt.Sprintf(
			VersionKeyFormat,
			listInput.Namespace,
			listInput.AppName,
			listInput.Type,
			listInput.ID,
			listInput.ToVersion,
		)
	}

	return
}

func buildConfigurationKey(version model.Version) (key string) {

	return fmt.Sprintf(ConfigurationKeyFormat,
		version.Namespace,
		"app",
		version.ConfigWrapper.Type,
		version.ConfigurationID,
	)
}
