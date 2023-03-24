package model

type Config struct {
	Service string            `json:"service"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}
