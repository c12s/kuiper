package model

type Entries = map[string]string

type Config struct {
	Service string  `json:"service"`
	Version string  `json:"version"`
	Entries Entries `json:"entries"`
}
