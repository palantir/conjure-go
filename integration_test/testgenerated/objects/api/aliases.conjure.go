// This file was generated by Conjure and should not be manually edited.

package api

import (
	"github.com/palantir/pkg/binary"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
)

type AnyAlias interface{}
type BinaryAlias []byte

func (a BinaryAlias) String() string {
	return binary.New(a).String()
}

func (a BinaryAlias) MarshalText() ([]byte, error) {
	return binary.New(a).MarshalText()
}

func (a *BinaryAlias) UnmarshalText(data []byte) error {
	rawBinaryAlias, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a = BinaryAlias(rawBinaryAlias)
	return nil
}

func (a BinaryAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *BinaryAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type ListLongAlias []interface{}
type LongAlias interface{}
type MapLongAlias map[string]interface{}
type MapStringAny map[string]interface{}
type MapStringAnyAlias map[string]AnyAlias
type MapUuidLongAlias map[uuid.UUID]interface{}

func (a MapUuidLongAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[uuid.UUID]interface{}(a))
}

func (a *MapUuidLongAlias) UnmarshalJSON(data []byte) error {
	var rawMapUuidLongAlias map[uuid.UUID]interface{}
	if err := safejson.Unmarshal(data, &rawMapUuidLongAlias); err != nil {
		return err
	}
	*a = MapUuidLongAlias(rawMapUuidLongAlias)
	return nil
}

func (a MapUuidLongAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapUuidLongAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type NestedAlias1 NestedAlias2
type NestedAlias2 NestedAlias3
type NestedAlias3 struct {
	Value *string
}

func (a NestedAlias3) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}

func (a *NestedAlias3) UnmarshalText(data []byte) error {
	rawNestedAlias3 := string(data)
	a.Value = &rawNestedAlias3
	return nil
}

func (a NestedAlias3) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *NestedAlias3) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalStructAlias struct {
	Value *Basic
}

func (a OptionalStructAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalStructAlias) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(Basic)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalStructAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalStructAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalUuidAlias struct {
	Value *uuid.UUID
}

func (a OptionalUuidAlias) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return a.Value.MarshalText()
}

func (a OptionalUuidAlias) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalUuidAlias) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(uuid.UUID)
	}
	return a.Value.UnmarshalText(data)
}

func (a OptionalUuidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalUuidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type RidAlias rid.ResourceIdentifier

func (a RidAlias) String() string {
	return rid.ResourceIdentifier(a).String()
}

func (a RidAlias) MarshalText() ([]byte, error) {
	return rid.ResourceIdentifier(a).MarshalText()
}

func (a *RidAlias) UnmarshalText(data []byte) error {
	var rawRidAlias rid.ResourceIdentifier
	if err := rawRidAlias.UnmarshalText(data); err != nil {
		return err
	}
	*a = RidAlias(rawRidAlias)
	return nil
}

func (a RidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *RidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type UuidAlias uuid.UUID

func (a UuidAlias) String() string {
	return uuid.UUID(a).String()
}

func (a UuidAlias) MarshalText() ([]byte, error) {
	return uuid.UUID(a).MarshalText()
}

func (a *UuidAlias) UnmarshalText(data []byte) error {
	var rawUuidAlias uuid.UUID
	if err := rawUuidAlias.UnmarshalText(data); err != nil {
		return err
	}
	*a = UuidAlias(rawUuidAlias)
	return nil
}

func (a UuidAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *UuidAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type UuidAlias2 Compound

func (a UuidAlias2) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(Compound(a))
}

func (a *UuidAlias2) UnmarshalJSON(data []byte) error {
	var rawUuidAlias2 Compound
	if err := safejson.Unmarshal(data, &rawUuidAlias2); err != nil {
		return err
	}
	*a = UuidAlias2(rawUuidAlias2)
	return nil
}

func (a UuidAlias2) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *UuidAlias2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
