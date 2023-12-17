package config

import "os"

type Config struct {
	SERVICE_PORT string
	DB_HOST      string
	DB_PORT      string
}

func NewConfig() *Config {
	return &Config{
		SERVICE_PORT: os.Getenv("SERVICEPORT"),
		DB_HOST:      os.Getenv("DBHOST"),
		DB_PORT:      os.Getenv("DBPORT"),
	}
}
