package server

import (
	"errors"
	"os"
)

type Config struct {
	JaegerAddress string
	EtcdAddress   string
	NatsAddress   string
}

const (
	jaegerAddressEnv = "JAEGER_ADDRESS"
	etcdAddressEnv   = "ETCD_ADDRESS"
	natsAddressEnv   = "NATS_ADDRESS"
)

func NewConfig() (Config, error) {
	jagAddr, found := os.LookupEnv(jaegerAddressEnv)
	if !found {
		return Config{}, errors.New(jaegerAddressEnv + " environment variable not set")
	}
	etcdAddr, found := os.LookupEnv(etcdAddressEnv)
	if !found {
		return Config{}, errors.New(etcdAddressEnv + " environment variable not set")
	}
	natsAddr, found := os.LookupEnv(natsAddressEnv)
	if !found {
		return Config{}, errors.New(natsAddressEnv + " environment variable not set")
	}

	return Config{JaegerAddress: jagAddr, EtcdAddress: etcdAddr, NatsAddress: natsAddr}, nil
}
