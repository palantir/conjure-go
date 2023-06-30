// This file was generated by Conjure and should not be manually edited.

package bar

import (
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Type1 []Type2
type Type2 struct {
	Value *int
}

func (a Type2) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *Type2) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(int)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a Type2) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *Type2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}