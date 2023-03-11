// This file was generated by Conjure and should not be manually edited.

package foo1

import (
	"github.com/palantir/conjure-go/v6/conjure/cycles/testdata/pkg-cycle/conjure/com/palantir/fizz"
	"github.com/palantir/conjure-go/v6/conjure/cycles/testdata/pkg-cycle/conjure/com/palantir/foo"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Type1 struct {
	Field1 foo.Type2 `json:"field1"`
	Field2 Type3     `json:"field2"`
}

func (o Type1) MarshalJSON() ([]byte, error) {
	if o.Field1 == nil {
		o.Field1 = make(map[fizz.Type1]foo.Type4, 0)
	}
	type Type1Alias Type1
	return safejson.Marshal(Type1Alias(o))
}

func (o *Type1) UnmarshalJSON(data []byte) error {
	type Type1Alias Type1
	var rawType1 Type1Alias
	if err := safejson.Unmarshal(data, &rawType1); err != nil {
		return err
	}
	if rawType1.Field1 == nil {
		rawType1.Field1 = make(map[fizz.Type1]foo.Type4, 0)
	}
	*o = Type1(rawType1)
	return nil
}

func (o Type1) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Type1) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
