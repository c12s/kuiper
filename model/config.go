package model

type Config struct {
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}
