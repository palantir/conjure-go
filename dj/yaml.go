package dj

import (
	"encoding/json"
	"github.com/palantir/pkg/safeyaml"
)

// MarshalYAML marshals the given json.Marshaler to YAML.
// Used to implement yaml.Marshaler.
func MarshalYAML(marshaler json.Marshaler) (interface{}, error) {
	jsonBytes, err := marshaler.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

// UnmarshalYAML unmarshals the given json.Unmarshaler from YAML.
// Used to implement yaml.Unmarshaler.
func UnmarshalYAML(unmarshaler json.Unmarshaler, unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return unmarshaler.UnmarshalJSON(jsonBytes)
}
