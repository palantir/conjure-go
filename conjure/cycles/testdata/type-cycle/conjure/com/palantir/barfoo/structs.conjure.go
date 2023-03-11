// This file was generated by Conjure and should not be manually edited.

package barfoo

import (
	"github.com/palantir/conjure-go/v6/conjure/cycles/testdata/type-cycle/conjure/com/palantir/buzz"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Type4 struct {
	Field1 buzz.Type1 `json:"field1"`
	Field2 Type31     `json:"field2"`
}

func (o Type4) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Type4) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Type3 struct {
	Field1 buzz.Type1 `json:"field1"`
	Field2 Type4      `json:"field2"`
}

func (o Type3) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Type3) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
