package main

import (
	"github.com/c12s/kuiper/repository/consul"
	"github.com/c12s/kuiper/service"
)

func main() {
	configRepo, err := consul.New()

	if err != nil {
		panic(err)
	}

	configService := service.New(configRepo)

	configService.DeleteGroup("")
}
