package server

import (
	"encoding/json"
	"io"
	"kuiper/model"
)

func decodeConfigBody(r io.Reader) (model.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var config *model.Config
	if err := dec.Decode(&config); err != nil {
		return model.Config{}, err
	}
	return *config, nil
}
