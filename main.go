package main

import (
	"github.com/c12s/kuiper/model"
	"github.com/c12s/kuiper/repository/consul"
	"github.com/c12s/kuiper/service"
)

func main() {
	configRepo, err := consul.New()

	if err != nil {
		panic(err)
	}

	configService := service.New(configRepo)

	group := model.Group{
		Configs: []model.Config{
			model.Config{
				Key:    "mysql.user",
				Value:  "asdf",
				Labels: []model.Label{},
			},
			model.Config{
				Key:    "mysql.pass",
				Value:  "123",
				Labels: []model.Label{},
			},
		},
	}

	resp, _ := configService.CreateNewGroup(group)

	println(resp.Id, resp.Version)

	stored, _ := configService.GetGroupConfigs(resp.Id, resp.Version, []model.Label{})

	for k, v := range stored {
		println(k, v)
	}
}
