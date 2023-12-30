package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/c12s/kuiper/model"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type ETCDRepo struct {
	Client *clientv3.Client
	Logger *zap.Logger
}

func New(log *zap.Logger) (repo *ETCDRepo, err error) {
	// dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	connectionURI := fmt.Sprintf("http://etcd:%s", dbPort)
	config := clientv3.Config{
		Endpoints:   []string{connectionURI},
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Error("error in creating etcd client.", zap.Error(err))
		panic(1)
	}
	fmt.Println("connection with db is successfully established.")

	repo = &ETCDRepo{
		Client: client,
		Logger: log,
	}

	return
}

func (etcd *ETCDRepo) ListVersions(input model.ListRequest) (results *arraylist.List, err error) {

	log := etcd.Logger.Named("[Repo:ListVersions]").With(zap.Any("input", input))
	log.Info("started")

	prefix, from, to := buildPrefixesFromListInput(input)
	listOptions := buildSearchOptions(input, prefix, &from, &to)

	searchPrefix := prefix

	if from != "" {
		searchPrefix = from
	}

	res, err := etcd.Client.Get(context.Background(), searchPrefix, listOptions...)
	if err != nil {
		etcd.Logger.Error("error in getting versions",
			zap.Error(err),
		)
	}

	results = arraylist.New()
	for _, el := range res.Kvs {
		var result model.Version
		err = json.Unmarshal(el.Value, &result)
		if err != nil {
			log.Info("Error in unmarshalling taken db element to version",
				zap.Error(err),
			)
		}

		results.Add(result)
	}

	log.Info("finished")
	return
}

func (etcd *ETCDRepo) CreateNewVersion(version model.Version) (saved model.Version, err error) {

	log := etcd.Logger.Named("[Repo:CreateNewVersion]").With(zap.Any("version", version))
	log.Info("CreateNewVersion started")

	version.CreatedAt = time.Now().Unix()
	if version.ConfigurationID == "" {
		version.ConfigurationID = uuid.NewString()
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

	bytes, err := json.Marshal(version)
	if err != nil {
		etcd.Logger.Error("error in marshalling version.",
			zap.Error(err),
		)
		return
	}

	_, err = etcd.Client.Put(context.Background(), buildConfigurationVersionKey(version), string(bytes))
	if err != nil {
		etcd.Logger.Error("error occured on putting version into db.",
			zap.Error(err),
		)
		return
	}

	saved = version
	log.Info("CreateNewVersion finished")
	return
}

func (etcd *ETCDRepo) GetPreviousVersions(version model.Version) ([]model.Version, error) {
	key := buildConfigurationKey(version)

	res, err := etcd.Client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		etcd.Logger.Info("error in getting prev version")
		return nil, err
	}

	result := make([]model.Version, 0)
	for _, element := range res.Kvs {
		newVers := model.Version{}
		err = json.Unmarshal(element.Value, &newVers)
		if err != nil {
			etcd.Logger.Info("error in unmarshalling value from kvs")
			continue
		}

		result = append(result, newVers)
	}

	return result, nil
}

func buildSearchOptions(listRequest model.ListRequest, prefix string, fromPrefix, toPrefix *string) (options []clientv3.OpOption) {
	options = make([]clientv3.OpOption, 0)

	if listRequest.FromVersion != "" && !listRequest.WithFrom {
		*fromPrefix += "\x00"
	}

	if listRequest.ToVersion != "" {
		if listRequest.WithTo {
			*toPrefix += "\x00"
		}
		options = append(options, clientv3.WithRange(*toPrefix))
		return
	}

	if listRequest.FromVersion != "" && listRequest.ToVersion == "" {
		rangeEnd := clientv3.GetPrefixRangeEnd(prefix)
		options = append(options, clientv3.WithRange(rangeEnd))
		return
	}

	options = append(options, clientv3.WithPrefix())
	return
}
