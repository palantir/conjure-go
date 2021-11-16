// This file was generated by Conjure and should not be manually edited.

package types

import (
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/safeyaml"
	"github.com/palantir/pkg/uuid"
)

type AnyExample struct {
	Value interface{} `json:"value"`
}

func (o AnyExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *AnyExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type BearerTokenExample struct {
	Value bearertoken.Token `json:"value"`
}

func (o BearerTokenExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BearerTokenExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type BinaryExample struct {
	Value []byte `json:"value"`
}

func (o BinaryExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BinaryExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type BooleanExample struct {
	Value bool `json:"value"`
}

func (o BooleanExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *BooleanExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type DateTimeExample struct {
	Value datetime.DateTime `json:"value"`
}

func (o DateTimeExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *DateTimeExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type DoubleExample struct {
	Value float64 `json:"value"`
}

func (o DoubleExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *DoubleExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type EmptyObjectExample struct{}

func (o EmptyObjectExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *EmptyObjectExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type EnumFieldExample struct {
	Enum EnumExample `json:"enum"`
}

func (o EnumFieldExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *EnumFieldExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type IntegerExample struct {
	Value int `json:"value"`
}

func (o IntegerExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *IntegerExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type KebabCaseObjectExample struct {
	KebabCasedField int `json:"kebab-cased-field"`
}

func (o KebabCaseObjectExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *KebabCaseObjectExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type ListExample struct {
	Value []string `json:"value"`
}

func (o ListExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]string, 0)
	}
	type ListExampleAlias ListExample
	return safejson.Marshal(ListExampleAlias(o))
}

func (o *ListExample) UnmarshalJSON(data []byte) error {
	type ListExampleAlias ListExample
	var rawListExample ListExampleAlias
	if err := safejson.Unmarshal(data, &rawListExample); err != nil {
		return err
	}
	if rawListExample.Value == nil {
		rawListExample.Value = make([]string, 0)
	}
	*o = ListExample(rawListExample)
	return nil
}

func (o ListExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ListExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type LongFieldNameOptionalExample struct {
	SomeLongName *string `json:"someLongName"`
}

func (o LongFieldNameOptionalExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *LongFieldNameOptionalExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type MapExample struct {
	Value map[string]string `json:"value"`
}

func (o MapExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make(map[string]string, 0)
	}
	type MapExampleAlias MapExample
	return safejson.Marshal(MapExampleAlias(o))
}

func (o *MapExample) UnmarshalJSON(data []byte) error {
	type MapExampleAlias MapExample
	var rawMapExample MapExampleAlias
	if err := safejson.Unmarshal(data, &rawMapExample); err != nil {
		return err
	}
	if rawMapExample.Value == nil {
		rawMapExample.Value = make(map[string]string, 0)
	}
	*o = MapExample(rawMapExample)
	return nil
}

func (o MapExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *MapExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type ObjectExample struct {
	String       string             `json:"string"`
	Integer      int                `json:"integer"`
	DoubleValue  float64            `json:"doubleValue"`
	OptionalItem *string            `json:"optionalItem"`
	Items        []string           `json:"items"`
	Set          []string           `json:"set"`
	Map          map[string]string  `json:"map"`
	Alias        StringAliasExample `json:"alias"`
}

func (o ObjectExample) MarshalJSON() ([]byte, error) {
	if o.Items == nil {
		o.Items = make([]string, 0)
	}
	if o.Set == nil {
		o.Set = make([]string, 0)
	}
	if o.Map == nil {
		o.Map = make(map[string]string, 0)
	}
	type ObjectExampleAlias ObjectExample
	return safejson.Marshal(ObjectExampleAlias(o))
}

func (o *ObjectExample) UnmarshalJSON(data []byte) error {
	type ObjectExampleAlias ObjectExample
	var rawObjectExample ObjectExampleAlias
	if err := safejson.Unmarshal(data, &rawObjectExample); err != nil {
		return err
	}
	if rawObjectExample.Items == nil {
		rawObjectExample.Items = make([]string, 0)
	}
	if rawObjectExample.Set == nil {
		rawObjectExample.Set = make([]string, 0)
	}
	if rawObjectExample.Map == nil {
		rawObjectExample.Map = make(map[string]string, 0)
	}
	*o = ObjectExample(rawObjectExample)
	return nil
}

func (o ObjectExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *ObjectExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type OptionalBooleanExample struct {
	Value *bool `json:"value"`
}

func (o OptionalBooleanExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *OptionalBooleanExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type OptionalExample struct {
	Value *string `json:"value"`
}

func (o OptionalExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *OptionalExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type OptionalIntegerExample struct {
	Value *int `json:"value"`
}

func (o OptionalIntegerExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *OptionalIntegerExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type RidExample struct {
	Value rid.ResourceIdentifier `json:"value"`
}

func (o RidExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *RidExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type SafeLongExample struct {
	Value safelong.SafeLong `json:"value"`
}

func (o SafeLongExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *SafeLongExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type SetDoubleExample struct {
	Value []float64 `json:"value"`
}

func (o SetDoubleExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]float64, 0)
	}
	type SetDoubleExampleAlias SetDoubleExample
	return safejson.Marshal(SetDoubleExampleAlias(o))
}

func (o *SetDoubleExample) UnmarshalJSON(data []byte) error {
	type SetDoubleExampleAlias SetDoubleExample
	var rawSetDoubleExample SetDoubleExampleAlias
	if err := safejson.Unmarshal(data, &rawSetDoubleExample); err != nil {
		return err
	}
	if rawSetDoubleExample.Value == nil {
		rawSetDoubleExample.Value = make([]float64, 0)
	}
	*o = SetDoubleExample(rawSetDoubleExample)
	return nil
}

func (o SetDoubleExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *SetDoubleExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type SetStringExample struct {
	Value []string `json:"value"`
}

func (o SetStringExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]string, 0)
	}
	type SetStringExampleAlias SetStringExample
	return safejson.Marshal(SetStringExampleAlias(o))
}

func (o *SetStringExample) UnmarshalJSON(data []byte) error {
	type SetStringExampleAlias SetStringExample
	var rawSetStringExample SetStringExampleAlias
	if err := safejson.Unmarshal(data, &rawSetStringExample); err != nil {
		return err
	}
	if rawSetStringExample.Value == nil {
		rawSetStringExample.Value = make([]string, 0)
	}
	*o = SetStringExample(rawSetStringExample)
	return nil
}

func (o SetStringExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *SetStringExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type SnakeCaseObjectExample struct {
	SnakeCasedField int `json:"snake_cased_field"`
}

func (o SnakeCaseObjectExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *SnakeCaseObjectExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type StringExample struct {
	Value string `json:"value"`
}

func (o StringExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *StringExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}

type UuidExample struct {
	Value uuid.UUID `json:"value"`
}

func (o UuidExample) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

func (o *UuidExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
