// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type OptionalIntegerAlias struct {
	Value *int
}

func (a OptionalIntegerAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalIntegerAlias) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(int)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalIntegerAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalIntegerAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalListAlias struct {
	Value *[]string
}

func (a OptionalListAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalListAlias) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new([]string)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalListAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalListAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type StringAlias string
