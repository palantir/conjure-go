// This file was generated by Conjure and should not be manually edited.

package types

import (
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/binary"
	"github.com/palantir/pkg/boolean"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
)

type AliasString string
type BearerTokenAliasExample bearertoken.Token

func (a BearerTokenAliasExample) String() string {
	return bearertoken.Token(a).String()
}

func (a BearerTokenAliasExample) MarshalText() ([]byte, error) {
	return bearertoken.Token(a).MarshalText()
}

func (a *BearerTokenAliasExample) UnmarshalText(data []byte) error {
	var rawBearerTokenAliasExample bearertoken.Token
	if err := rawBearerTokenAliasExample.UnmarshalText(data); err != nil {
		return err
	}
	*a = BearerTokenAliasExample(rawBearerTokenAliasExample)
	return nil
}

func (a BearerTokenAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *BearerTokenAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type BinaryAliasExample []byte

func (a BinaryAliasExample) String() string {
	return binary.New(a).String()
}

func (a BinaryAliasExample) MarshalText() ([]byte, error) {
	return binary.New(a).MarshalText()
}

func (a *BinaryAliasExample) UnmarshalText(data []byte) error {
	rawBinaryAliasExample, err := binary.Binary(data).Bytes()
	if err != nil {
		return err
	}
	*a = BinaryAliasExample(rawBinaryAliasExample)
	return nil
}

func (a BinaryAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *BinaryAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type BooleanAliasExample bool
type DateTimeAliasExample datetime.DateTime

func (a DateTimeAliasExample) String() string {
	return datetime.DateTime(a).String()
}

func (a DateTimeAliasExample) MarshalText() ([]byte, error) {
	return datetime.DateTime(a).MarshalText()
}

func (a *DateTimeAliasExample) UnmarshalText(data []byte) error {
	var rawDateTimeAliasExample datetime.DateTime
	if err := rawDateTimeAliasExample.UnmarshalText(data); err != nil {
		return err
	}
	*a = DateTimeAliasExample(rawDateTimeAliasExample)
	return nil
}

func (a DateTimeAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *DateTimeAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type DoubleAliasExample float64
type IntegerAliasExample int
type ListAnyAliasExample []interface{}
type ListBearerTokenAliasExample []bearertoken.Token

func (a ListBearerTokenAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]bearertoken.Token(a))
}

func (a *ListBearerTokenAliasExample) UnmarshalJSON(data []byte) error {
	var rawListBearerTokenAliasExample []bearertoken.Token
	if err := safejson.Unmarshal(data, &rawListBearerTokenAliasExample); err != nil {
		return err
	}
	*a = ListBearerTokenAliasExample(rawListBearerTokenAliasExample)
	return nil
}

func (a ListBearerTokenAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ListBearerTokenAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type ListBinaryAliasExample [][]byte

func (a ListBinaryAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([][]byte(a))
}

func (a *ListBinaryAliasExample) UnmarshalJSON(data []byte) error {
	var rawListBinaryAliasExample [][]byte
	if err := safejson.Unmarshal(data, &rawListBinaryAliasExample); err != nil {
		return err
	}
	*a = ListBinaryAliasExample(rawListBinaryAliasExample)
	return nil
}

func (a ListBinaryAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ListBinaryAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type ListBooleanAliasExample []bool
type ListDateTimeAliasExample []datetime.DateTime

func (a ListDateTimeAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]datetime.DateTime(a))
}

func (a *ListDateTimeAliasExample) UnmarshalJSON(data []byte) error {
	var rawListDateTimeAliasExample []datetime.DateTime
	if err := safejson.Unmarshal(data, &rawListDateTimeAliasExample); err != nil {
		return err
	}
	*a = ListDateTimeAliasExample(rawListDateTimeAliasExample)
	return nil
}

func (a ListDateTimeAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ListDateTimeAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type ListDoubleAliasExample []float64
type ListIntegerAliasExample []int
type ListOptionalAnyAliasExample []*interface{}
type ListRidAliasExample []rid.ResourceIdentifier

func (a ListRidAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]rid.ResourceIdentifier(a))
}

func (a *ListRidAliasExample) UnmarshalJSON(data []byte) error {
	var rawListRidAliasExample []rid.ResourceIdentifier
	if err := safejson.Unmarshal(data, &rawListRidAliasExample); err != nil {
		return err
	}
	*a = ListRidAliasExample(rawListRidAliasExample)
	return nil
}

func (a ListRidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ListRidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type ListSafeLongAliasExample []safelong.SafeLong

func (a ListSafeLongAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]safelong.SafeLong(a))
}

func (a *ListSafeLongAliasExample) UnmarshalJSON(data []byte) error {
	var rawListSafeLongAliasExample []safelong.SafeLong
	if err := safejson.Unmarshal(data, &rawListSafeLongAliasExample); err != nil {
		return err
	}
	*a = ListSafeLongAliasExample(rawListSafeLongAliasExample)
	return nil
}

func (a ListSafeLongAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ListSafeLongAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type ListStringAliasExample []string
type ListUuidAliasExample []uuid.UUID

func (a ListUuidAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]uuid.UUID(a))
}

func (a *ListUuidAliasExample) UnmarshalJSON(data []byte) error {
	var rawListUuidAliasExample []uuid.UUID
	if err := safejson.Unmarshal(data, &rawListUuidAliasExample); err != nil {
		return err
	}
	*a = ListUuidAliasExample(rawListUuidAliasExample)
	return nil
}

func (a ListUuidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ListUuidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapBearerTokenAliasExample map[bearertoken.Token]bool

func (a MapBearerTokenAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[bearertoken.Token]bool(a))
}

func (a *MapBearerTokenAliasExample) UnmarshalJSON(data []byte) error {
	var rawMapBearerTokenAliasExample map[bearertoken.Token]bool
	if err := safejson.Unmarshal(data, &rawMapBearerTokenAliasExample); err != nil {
		return err
	}
	*a = MapBearerTokenAliasExample(rawMapBearerTokenAliasExample)
	return nil
}

func (a MapBearerTokenAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapBearerTokenAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapBinaryAliasExample map[binary.Binary]bool

func (a MapBinaryAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[binary.Binary]bool(a))
}

func (a *MapBinaryAliasExample) UnmarshalJSON(data []byte) error {
	var rawMapBinaryAliasExample map[binary.Binary]bool
	if err := safejson.Unmarshal(data, &rawMapBinaryAliasExample); err != nil {
		return err
	}
	*a = MapBinaryAliasExample(rawMapBinaryAliasExample)
	return nil
}

func (a MapBinaryAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapBinaryAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapBooleanAliasExample map[boolean.Boolean]bool
type MapDateTimeAliasExample map[datetime.DateTime]bool

func (a MapDateTimeAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[datetime.DateTime]bool(a))
}

func (a *MapDateTimeAliasExample) UnmarshalJSON(data []byte) error {
	var rawMapDateTimeAliasExample map[datetime.DateTime]bool
	if err := safejson.Unmarshal(data, &rawMapDateTimeAliasExample); err != nil {
		return err
	}
	*a = MapDateTimeAliasExample(rawMapDateTimeAliasExample)
	return nil
}

func (a MapDateTimeAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapDateTimeAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapDoubleAliasExample map[float64]bool
type MapEnumExampleAlias map[EnumExample]string

func (a MapEnumExampleAlias) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[EnumExample]string(a))
}

func (a *MapEnumExampleAlias) UnmarshalJSON(data []byte) error {
	var rawMapEnumExampleAlias map[EnumExample]string
	if err := safejson.Unmarshal(data, &rawMapEnumExampleAlias); err != nil {
		return err
	}
	*a = MapEnumExampleAlias(rawMapEnumExampleAlias)
	return nil
}

func (a MapEnumExampleAlias) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapEnumExampleAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapIntegerAliasExample map[int]bool
type MapRidAliasExample map[rid.ResourceIdentifier]bool

func (a MapRidAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[rid.ResourceIdentifier]bool(a))
}

func (a *MapRidAliasExample) UnmarshalJSON(data []byte) error {
	var rawMapRidAliasExample map[rid.ResourceIdentifier]bool
	if err := safejson.Unmarshal(data, &rawMapRidAliasExample); err != nil {
		return err
	}
	*a = MapRidAliasExample(rawMapRidAliasExample)
	return nil
}

func (a MapRidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapRidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapSafeLongAliasExample map[safelong.SafeLong]bool

func (a MapSafeLongAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[safelong.SafeLong]bool(a))
}

func (a *MapSafeLongAliasExample) UnmarshalJSON(data []byte) error {
	var rawMapSafeLongAliasExample map[safelong.SafeLong]bool
	if err := safejson.Unmarshal(data, &rawMapSafeLongAliasExample); err != nil {
		return err
	}
	*a = MapSafeLongAliasExample(rawMapSafeLongAliasExample)
	return nil
}

func (a MapSafeLongAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapSafeLongAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type MapStringAliasExample map[string]bool
type MapUuidAliasExample map[uuid.UUID]bool

func (a MapUuidAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[uuid.UUID]bool(a))
}

func (a *MapUuidAliasExample) UnmarshalJSON(data []byte) error {
	var rawMapUuidAliasExample map[uuid.UUID]bool
	if err := safejson.Unmarshal(data, &rawMapUuidAliasExample); err != nil {
		return err
	}
	*a = MapUuidAliasExample(rawMapUuidAliasExample)
	return nil
}

func (a MapUuidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *MapUuidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type OptionalAnyAliasExample struct {
	Value *interface{}
}

func (a OptionalAnyAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalAnyAliasExample) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(interface{})
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalAnyAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalAnyAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type OptionalBearerTokenAliasExample struct {
	Value *bearertoken.Token
}

func (a OptionalBearerTokenAliasExample) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return a.Value.MarshalText()
}

func (a OptionalBearerTokenAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalBearerTokenAliasExample) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(bearertoken.Token)
	}
	return a.Value.UnmarshalText(data)
}

func (a OptionalBearerTokenAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalBearerTokenAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalBooleanAliasExample struct {
	Value *bool
}

func (a OptionalBooleanAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalBooleanAliasExample) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(bool)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalBooleanAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalBooleanAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type OptionalDateTimeAliasExample struct {
	Value *datetime.DateTime
}

func (a OptionalDateTimeAliasExample) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return a.Value.MarshalText()
}

func (a OptionalDateTimeAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalDateTimeAliasExample) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(datetime.DateTime)
	}
	return a.Value.UnmarshalText(data)
}

func (a OptionalDateTimeAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalDateTimeAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalDoubleAliasExample struct {
	Value *float64
}

func (a OptionalDoubleAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalDoubleAliasExample) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(float64)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalDoubleAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalDoubleAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type OptionalIntegerAliasExample struct {
	Value *int
}

func (a OptionalIntegerAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalIntegerAliasExample) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(int)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalIntegerAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalIntegerAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type OptionalRidAliasExample struct {
	Value *rid.ResourceIdentifier
}

func (a OptionalRidAliasExample) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return a.Value.MarshalText()
}

func (a OptionalRidAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalRidAliasExample) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(rid.ResourceIdentifier)
	}
	return a.Value.UnmarshalText(data)
}

func (a OptionalRidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalRidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalSafeLongAliasExample struct {
	Value *safelong.SafeLong
}

func (a OptionalSafeLongAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalSafeLongAliasExample) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(safelong.SafeLong)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a OptionalSafeLongAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalSafeLongAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type OptionalStringAliasExample struct {
	Value *string
}

func (a OptionalStringAliasExample) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}

func (a OptionalStringAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalStringAliasExample) UnmarshalText(data []byte) error {
	rawOptionalStringAliasExample := string(data)
	a.Value = &rawOptionalStringAliasExample
	return nil
}

func (a OptionalStringAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalStringAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type OptionalUuidAliasExample struct {
	Value *uuid.UUID
}

func (a OptionalUuidAliasExample) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return a.Value.MarshalText()
}

func (a OptionalUuidAliasExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *OptionalUuidAliasExample) UnmarshalText(data []byte) error {
	if a.Value == nil {
		a.Value = new(uuid.UUID)
	}
	return a.Value.UnmarshalText(data)
}

func (a OptionalUuidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *OptionalUuidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type RawOptionalExample struct {
	Value *int
}

func (a RawOptionalExample) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return []byte("null"), nil
	}
	return safejson.Marshal(a.Value)
}

func (a *RawOptionalExample) UnmarshalJSON(data []byte) error {
	if a.Value == nil {
		a.Value = new(int)
	}
	return safejson.Unmarshal(data, a.Value)
}

func (a RawOptionalExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *RawOptionalExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type ReferenceAliasExample AnyExample

func (a ReferenceAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(AnyExample(a))
}

func (a *ReferenceAliasExample) UnmarshalJSON(data []byte) error {
	var rawReferenceAliasExample AnyExample
	if err := safejson.Unmarshal(data, &rawReferenceAliasExample); err != nil {
		return err
	}
	*a = ReferenceAliasExample(rawReferenceAliasExample)
	return nil
}

func (a ReferenceAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *ReferenceAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type RidAliasExample rid.ResourceIdentifier

func (a RidAliasExample) String() string {
	return rid.ResourceIdentifier(a).String()
}

func (a RidAliasExample) MarshalText() ([]byte, error) {
	return rid.ResourceIdentifier(a).MarshalText()
}

func (a *RidAliasExample) UnmarshalText(data []byte) error {
	var rawRidAliasExample rid.ResourceIdentifier
	if err := rawRidAliasExample.UnmarshalText(data); err != nil {
		return err
	}
	*a = RidAliasExample(rawRidAliasExample)
	return nil
}

func (a RidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *RidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}

type SafeLongAliasExample safelong.SafeLong

func (a SafeLongAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(safelong.SafeLong(a))
}

func (a *SafeLongAliasExample) UnmarshalJSON(data []byte) error {
	var rawSafeLongAliasExample safelong.SafeLong
	if err := safejson.Unmarshal(data, &rawSafeLongAliasExample); err != nil {
		return err
	}
	*a = SafeLongAliasExample(rawSafeLongAliasExample)
	return nil
}

func (a SafeLongAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SafeLongAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type SetAnyAliasExample []interface{}
type SetBearerTokenAliasExample []bearertoken.Token

func (a SetBearerTokenAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]bearertoken.Token(a))
}

func (a *SetBearerTokenAliasExample) UnmarshalJSON(data []byte) error {
	var rawSetBearerTokenAliasExample []bearertoken.Token
	if err := safejson.Unmarshal(data, &rawSetBearerTokenAliasExample); err != nil {
		return err
	}
	*a = SetBearerTokenAliasExample(rawSetBearerTokenAliasExample)
	return nil
}

func (a SetBearerTokenAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SetBearerTokenAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type SetBinaryAliasExample [][]byte

func (a SetBinaryAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([][]byte(a))
}

func (a *SetBinaryAliasExample) UnmarshalJSON(data []byte) error {
	var rawSetBinaryAliasExample [][]byte
	if err := safejson.Unmarshal(data, &rawSetBinaryAliasExample); err != nil {
		return err
	}
	*a = SetBinaryAliasExample(rawSetBinaryAliasExample)
	return nil
}

func (a SetBinaryAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SetBinaryAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type SetBooleanAliasExample []bool
type SetDateTimeAliasExample []datetime.DateTime

func (a SetDateTimeAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]datetime.DateTime(a))
}

func (a *SetDateTimeAliasExample) UnmarshalJSON(data []byte) error {
	var rawSetDateTimeAliasExample []datetime.DateTime
	if err := safejson.Unmarshal(data, &rawSetDateTimeAliasExample); err != nil {
		return err
	}
	*a = SetDateTimeAliasExample(rawSetDateTimeAliasExample)
	return nil
}

func (a SetDateTimeAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SetDateTimeAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type SetDoubleAliasExample []float64
type SetIntegerAliasExample []int
type SetOptionalAnyAliasExample []*interface{}
type SetRidAliasExample []rid.ResourceIdentifier

func (a SetRidAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]rid.ResourceIdentifier(a))
}

func (a *SetRidAliasExample) UnmarshalJSON(data []byte) error {
	var rawSetRidAliasExample []rid.ResourceIdentifier
	if err := safejson.Unmarshal(data, &rawSetRidAliasExample); err != nil {
		return err
	}
	*a = SetRidAliasExample(rawSetRidAliasExample)
	return nil
}

func (a SetRidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SetRidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type SetSafeLongAliasExample []safelong.SafeLong

func (a SetSafeLongAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]safelong.SafeLong(a))
}

func (a *SetSafeLongAliasExample) UnmarshalJSON(data []byte) error {
	var rawSetSafeLongAliasExample []safelong.SafeLong
	if err := safejson.Unmarshal(data, &rawSetSafeLongAliasExample); err != nil {
		return err
	}
	*a = SetSafeLongAliasExample(rawSetSafeLongAliasExample)
	return nil
}

func (a SetSafeLongAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SetSafeLongAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type SetStringAliasExample []string
type SetUuidAliasExample []uuid.UUID

func (a SetUuidAliasExample) MarshalJSON() ([]byte, error) {
	return safejson.Marshal([]uuid.UUID(a))
}

func (a *SetUuidAliasExample) UnmarshalJSON(data []byte) error {
	var rawSetUuidAliasExample []uuid.UUID
	if err := safejson.Unmarshal(data, &rawSetUuidAliasExample); err != nil {
		return err
	}
	*a = SetUuidAliasExample(rawSetUuidAliasExample)
	return nil
}

func (a SetUuidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *SetUuidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return a.UnmarshalJSON(jsonBytes)
}

type StringAliasExample string
type UuidAliasExample uuid.UUID

func (a UuidAliasExample) String() string {
	return uuid.UUID(a).String()
}

func (a UuidAliasExample) MarshalText() ([]byte, error) {
	return uuid.UUID(a).MarshalText()
}

func (a *UuidAliasExample) UnmarshalText(data []byte) error {
	var rawUuidAliasExample uuid.UUID
	if err := rawUuidAliasExample.UnmarshalText(data); err != nil {
		return err
	}
	*a = UuidAliasExample(rawUuidAliasExample)
	return nil
}

func (a UuidAliasExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (a *UuidAliasExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
