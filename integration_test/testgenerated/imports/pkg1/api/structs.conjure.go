// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Struct1 struct {
	Data string `json:"data"`
}

func (o Struct1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Struct1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
