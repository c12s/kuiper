package configs

import (
	"os"
)

type Config struct {
	natsAddress       string
	magnetarAddress   string
	agentQueueAddress string
	oortAddress       string
	apolloAddress     string
	serverAddress     string
	tokenKey          string
}

func (c *Config) NatsAddress() string {
	return c.natsAddress
}

func (c *Config) MagnetarAddress() string {
	return c.magnetarAddress
}

func (c *Config) AgentQueueAddress() string {
	return c.agentQueueAddress
}

func (c *Config) OortAddress() string {
	return c.oortAddress
}

func (c *Config) ApolloAddress() string {
	return c.apolloAddress
}

func (c *Config) ServerAddress() string {
	return c.serverAddress
}

func (c *Config) TokenKey() string {
	return c.tokenKey
}

func NewFromEnv() (*Config, error) {
	return &Config{
		natsAddress:       os.Getenv("NATS_ADDRESS"),
		magnetarAddress:   os.Getenv("MAGNETAR_ADDRESS"),
		agentQueueAddress: os.Getenv("AGENT_QUEUE_ADDRESS"),
		oortAddress:       os.Getenv("OORT_ADDRESS"),
		apolloAddress:     os.Getenv("APOLLO_ADDRESS"),
		serverAddress:     os.Getenv("KUIPER_ADDRESS"),
		tokenKey:          os.Getenv("SECRET_KEY"),
	}, nil
}