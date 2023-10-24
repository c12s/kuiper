package configs

import (
	"os"
)

type Config struct {
	natsAddress     string
	magnetarAddress string
	oortAddress     string
	serverAddress   string
}

func (c *Config) NatsAddress() string {
	return c.natsAddress
}

func (c *Config) MagnetarAddress() string {
	return c.magnetarAddress
}

func (c *Config) OortAddress() string {
	return c.oortAddress
}

func (c *Config) ServerAddress() string {
	return c.serverAddress
}

func NewFromEnv() (*Config, error) {
	return &Config{
		natsAddress:     os.Getenv("NATS_ADDRESS"),
		magnetarAddress: os.Getenv("MAGNETAR_ADDRESS"),
		oortAddress:     os.Getenv("OORT_ADDRESS"),
		serverAddress:   os.Getenv("KUIPER_ADDRESS"),
	}, nil
}
