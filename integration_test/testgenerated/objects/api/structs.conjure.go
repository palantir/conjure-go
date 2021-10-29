// This file was generated by Conjure and should not be manually edited.

package api

import (
	binary "github.com/palantir/pkg/binary"
	boolean "github.com/palantir/pkg/boolean"
	safejson "github.com/palantir/pkg/safejson"
	safeyaml "github.com/palantir/pkg/safeyaml"
	uuid "github.com/palantir/pkg/uuid"
)

type AnyValue struct {
	Value interface{} `json:"value"`
}

func (o AnyValue) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *AnyValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Basic struct {
	Data string `json:"data"`
}

func (o Basic) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Basic) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type BinaryMap struct {
	Map map[binary.Binary][]byte `json:"map"`
}

func (o BinaryMap) MarshalJSON() ([]byte, error) {
	if o.Map == nil {
		o.Map = make(map[binary.Binary][]byte, 0)
	}
	type BinaryMapAlias BinaryMap
	return safejson.Marshal(BinaryMapAlias(o))
}

func (o *BinaryMap) UnmarshalJSON(data []byte) error {
	type BinaryMapAlias BinaryMap
	var rawBinaryMap BinaryMapAlias
	if err := safejson.Unmarshal(data, &rawBinaryMap); err != nil {
		return err
	}
	if rawBinaryMap.Map == nil {
		rawBinaryMap.Map = make(map[binary.Binary][]byte, 0)
	}
	*o = BinaryMap(rawBinaryMap)
	return nil
}

func (o BinaryMap) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BinaryMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type BooleanIntegerMap struct {
	Map map[boolean.Boolean]int `json:"map"`
}

func (o BooleanIntegerMap) MarshalJSON() ([]byte, error) {
	if o.Map == nil {
		o.Map = make(map[boolean.Boolean]int, 0)
	}
	type BooleanIntegerMapAlias BooleanIntegerMap
	return safejson.Marshal(BooleanIntegerMapAlias(o))
}

func (o *BooleanIntegerMap) UnmarshalJSON(data []byte) error {
	type BooleanIntegerMapAlias BooleanIntegerMap
	var rawBooleanIntegerMap BooleanIntegerMapAlias
	if err := safejson.Unmarshal(data, &rawBooleanIntegerMap); err != nil {
		return err
	}
	if rawBooleanIntegerMap.Map == nil {
		rawBooleanIntegerMap.Map = make(map[boolean.Boolean]int, 0)
	}
	*o = BooleanIntegerMap(rawBooleanIntegerMap)
	return nil
}

func (o BooleanIntegerMap) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BooleanIntegerMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Collections struct {
	MapVar   map[string][]int   `json:"mapVar"`
	ListVar  []string           `json:"listVar"`
	MultiDim [][]map[string]int `json:"multiDim"`
}

func (o Collections) MarshalJSON() ([]byte, error) {
	if o.MapVar == nil {
		o.MapVar = make(map[string][]int, 0)
	}
	if o.ListVar == nil {
		o.ListVar = make([]string, 0)
	}
	if o.MultiDim == nil {
		o.MultiDim = make([][]map[string]int, 0)
	}
	type CollectionsAlias Collections
	return safejson.Marshal(CollectionsAlias(o))
}

func (o *Collections) UnmarshalJSON(data []byte) error {
	type CollectionsAlias Collections
	var rawCollections CollectionsAlias
	if err := safejson.Unmarshal(data, &rawCollections); err != nil {
		return err
	}
	if rawCollections.MapVar == nil {
		rawCollections.MapVar = make(map[string][]int, 0)
	}
	if rawCollections.ListVar == nil {
		rawCollections.ListVar = make([]string, 0)
	}
	if rawCollections.MultiDim == nil {
		rawCollections.MultiDim = make([][]map[string]int, 0)
	}
	*o = Collections(rawCollections)
	return nil
}

func (o Collections) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Collections) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type Compound struct {
	Obj Collections `json:"obj"`
}

func (o Compound) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Compound) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type ExampleUuid struct {
	Uid uuid.UUID `json:"uid"`
}

func (o ExampleUuid) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ExampleUuid) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type MapOptional struct {
	Map map[OptionalUuidAlias]string `json:"map"`
}

func (o MapOptional) MarshalJSON() ([]byte, error) {
	if o.Map == nil {
		o.Map = make(map[OptionalUuidAlias]string, 0)
	}
	type MapOptionalAlias MapOptional
	return safejson.Marshal(MapOptionalAlias(o))
}

func (o *MapOptional) UnmarshalJSON(data []byte) error {
	type MapOptionalAlias MapOptional
	var rawMapOptional MapOptionalAlias
	if err := safejson.Unmarshal(data, &rawMapOptional); err != nil {
		return err
	}
	if rawMapOptional.Map == nil {
		rawMapOptional.Map = make(map[OptionalUuidAlias]string, 0)
	}
	*o = MapOptional(rawMapOptional)
	return nil
}

func (o MapOptional) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *MapOptional) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type OptionalFields struct {
	Opt1 *string           `json:"opt1"`
	Opt2 *string           `json:"opt2"`
	Reqd string            `json:"reqd"`
	Opt3 OptionalUuidAlias `json:"opt3"`
}

func (o OptionalFields) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *OptionalFields) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

// A type using go keywords
type Type struct {
	Type []string          `json:"type"`
	Chan map[string]string `json:"chan"`
}

func (o Type) MarshalJSON() ([]byte, error) {
	if o.Type == nil {
		o.Type = make([]string, 0)
	}
	if o.Chan == nil {
		o.Chan = make(map[string]string, 0)
	}
	type TypeAlias Type
	return safejson.Marshal(TypeAlias(o))
}

func (o *Type) UnmarshalJSON(data []byte) error {
	type TypeAlias Type
	var rawType TypeAlias
	if err := safejson.Unmarshal(data, &rawType); err != nil {
		return err
	}
	if rawType.Type == nil {
		rawType.Type = make([]string, 0)
	}
	if rawType.Chan == nil {
		rawType.Chan = make(map[string]string, 0)
	}
	*o = Type(rawType)
	return nil
}

func (o Type) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
