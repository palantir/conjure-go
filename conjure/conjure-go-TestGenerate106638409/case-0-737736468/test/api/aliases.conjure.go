// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/conjure-go/v6/conjure/conjure-go-TestGenerate106638409/case-0-737736468/foundry/catalog/api/datasets"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/safeyaml"
)

type ExampleAlias string
type LongAlias safelong.SafeLong

func (a LongAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(safelong.SafeLong(a))
}

func (a *LongAlias) UnmarshalJSON(data []byte) error {
	var rawLongAlias safelong.SafeLong
	if err := safejson.Unmarshal(data, &rawLongAlias); err != nil {
		return err
	}
	*a = LongAlias(rawLongAlias)
	return nil
}

func (a LongAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *LongAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type Status int
type ObjectAlias datasets.TestType

func (a ObjectAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(datasets.TestType(a))
}

func (a *ObjectAlias) UnmarshalJSON(data []byte) error {
	var rawObjectAlias datasets.TestType
	if err := safejson.Unmarshal(data, &rawObjectAlias); err != nil {
		return err
	}
	*a = ObjectAlias(rawObjectAlias)
	return nil
}

func (a ObjectAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ObjectAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type MapAlias map[string]Status

func (a MapAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[string]Status(a))
}

func (a *MapAlias) UnmarshalJSON(data []byte) error {
	var rawMapAlias map[string]Status
	if err := safejson.Unmarshal(data, &rawMapAlias); err != nil {
		return err
	}
	*a = MapAlias(rawMapAlias)
	return nil
}

func (a MapAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type AliasAlias Status

func (a AliasAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(Status(a))
}

func (a *AliasAlias) UnmarshalJSON(data []byte) error {
	var rawAliasAlias Status
	if err := safejson.Unmarshal(data, &rawAliasAlias); err != nil {
		return err
	}
	*a = AliasAlias(rawAliasAlias)
	return nil
}

func (a AliasAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *AliasAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
