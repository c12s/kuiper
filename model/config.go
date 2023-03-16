package model

type Config struct {
	ID      string            `json:"id"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}
