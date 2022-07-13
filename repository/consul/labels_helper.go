package consul

import (
	"fmt"
	"github.com/c12s/kuiper/model"
	"sort"
	"strings"
)

func DecodeLabels(labels []model.Label) string {
	keys := make([]string, 0, len(labels))
	pairs := make([]string, 0, len(labels))

	labelsMap := make(map[string]string)

	for _, l := range labels {
		keys = append(keys, l.Key)
		labelsMap[l.Key] = l.Value
	}

	sort.Strings(keys)

	for _, k := range keys {
		val := fmt.Sprintf("%s=%s", k, labelsMap[k])
		pairs = append(pairs, val)
	}

	return strings.Join(pairs[:], "&")
}
