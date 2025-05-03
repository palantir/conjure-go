package cj

import (
	"bytes"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/safeyaml"
)

func YAMLV3MarshalerFromJSON(obj json.MarshalerTo) (any, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func YAMLV3UnmarshalerToJSON(obj json.UnmarshalerFrom, unmarshal func(any) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return obj.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(jsonBytes)))
}
