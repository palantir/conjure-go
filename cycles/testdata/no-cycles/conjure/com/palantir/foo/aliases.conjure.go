// This file was generated by Conjure and should not be manually edited.

package foo

import (
	"github.com/palantir/conjure-go/v6/cycles/testdata/no-cycles/conjure/com/palantir/fizz"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
)

type Type2 map[fizz.Type1]Type4

func (a Type2) MarshalJSON() ([]byte, error) {
	rawType2 := map[fizz.Type1]Type4(a)
	if rawType2 == nil {
		rawType2 = make(map[fizz.Type1]Type4, 0)
	}
	return safejson.Marshal(rawType2)
}

func (a *Type2) UnmarshalJSON(data []byte) error {
	var rawType2 map[fizz.Type1]Type4
	if err := safejson.Unmarshal(data, &rawType2); err != nil {
		return err
	}
	*a = Type2(rawType2)
	return nil
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
