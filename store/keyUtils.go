package store

import (
	"fmt"
	"strings"
)

func makeKey(id, ver string) string {
	key := fmt.Sprintf("config/%s/%s/", id, ver)
	return key
}
func makeIdPrefix(id string) string {
	key := fmt.Sprintf("config/%s/", id)
	return key
}

func getVersionFromKey(key string) string {
	split := strings.Split(string(key), "/")
	return split[len(split)-2]
}
