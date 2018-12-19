// This file was generated by Conjure and should not be manually edited.

package verification

import (
	"encoding/json"

	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
)

type BearerTokenExample struct {
	Value bearertoken.Bearertoken `json:"value" yaml:"value,omitempty"`
}

type BinaryExample struct {
	Value []byte `json:"value" yaml:"value,omitempty"`
}

func (o BinaryExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]byte, 0)
	}
	type BinaryExampleAlias BinaryExample
	return json.Marshal(BinaryExampleAlias(o))
}

func (o *BinaryExample) UnmarshalJSON(data []byte) error {
	type BinaryExampleAlias BinaryExample
	var rawBinaryExample BinaryExampleAlias
	if err := json.Unmarshal(data, &rawBinaryExample); err != nil {
		return err
	}
	if rawBinaryExample.Value == nil {
		rawBinaryExample.Value = make([]byte, 0)
	}
	*o = BinaryExample(rawBinaryExample)
	return nil
}

func (o BinaryExample) MarshalYAML() (interface{}, error) {
	if o.Value == nil {
		o.Value = make([]byte, 0)
	}
	type BinaryExampleAlias BinaryExample
	return BinaryExampleAlias(o), nil
}

func (o *BinaryExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type BinaryExampleAlias BinaryExample
	var rawBinaryExample BinaryExampleAlias
	if err := unmarshal(&rawBinaryExample); err != nil {
		return err
	}
	if rawBinaryExample.Value == nil {
		rawBinaryExample.Value = make([]byte, 0)
	}
	*o = BinaryExample(rawBinaryExample)
	return nil
}

type BooleanExample struct {
	Value bool `json:"value" yaml:"value,omitempty"`
}

type DateTimeExample struct {
	Value datetime.DateTime `json:"value" yaml:"value,omitempty"`
}

type DoubleExample struct {
	Value float64 `json:"value" yaml:"value,omitempty"`
}

type IntegerExample struct {
	Value int `json:"value" yaml:"value,omitempty"`
}

type RidExample struct {
	Value rid.ResourceIdentifier `json:"value" yaml:"value,omitempty"`
}

type SafeLongExample struct {
	Value safelong.SafeLong `json:"value" yaml:"value,omitempty"`
}

type StringExample struct {
	Value string `json:"value" yaml:"value,omitempty"`
}

type UuidExample struct {
	Value uuid.UUID `json:"value" yaml:"value,omitempty"`
}

type AnyExample struct {
	Value interface{} `json:"value" yaml:"value,omitempty"`
}

type ListExample struct {
	Value []string `json:"value" yaml:"value,omitempty"`
}

func (o ListExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]string, 0)
	}
	type ListExampleAlias ListExample
	return json.Marshal(ListExampleAlias(o))
}

func (o *ListExample) UnmarshalJSON(data []byte) error {
	type ListExampleAlias ListExample
	var rawListExample ListExampleAlias
	if err := json.Unmarshal(data, &rawListExample); err != nil {
		return err
	}
	if rawListExample.Value == nil {
		rawListExample.Value = make([]string, 0)
	}
	*o = ListExample(rawListExample)
	return nil
}

func (o ListExample) MarshalYAML() (interface{}, error) {
	if o.Value == nil {
		o.Value = make([]string, 0)
	}
	type ListExampleAlias ListExample
	return ListExampleAlias(o), nil
}

func (o *ListExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ListExampleAlias ListExample
	var rawListExample ListExampleAlias
	if err := unmarshal(&rawListExample); err != nil {
		return err
	}
	if rawListExample.Value == nil {
		rawListExample.Value = make([]string, 0)
	}
	*o = ListExample(rawListExample)
	return nil
}

type SetStringExample struct {
	Value []string `json:"value" yaml:"value,omitempty"`
}

func (o SetStringExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]string, 0)
	}
	type SetStringExampleAlias SetStringExample
	return json.Marshal(SetStringExampleAlias(o))
}

func (o *SetStringExample) UnmarshalJSON(data []byte) error {
	type SetStringExampleAlias SetStringExample
	var rawSetStringExample SetStringExampleAlias
	if err := json.Unmarshal(data, &rawSetStringExample); err != nil {
		return err
	}
	if rawSetStringExample.Value == nil {
		rawSetStringExample.Value = make([]string, 0)
	}
	*o = SetStringExample(rawSetStringExample)
	return nil
}

func (o SetStringExample) MarshalYAML() (interface{}, error) {
	if o.Value == nil {
		o.Value = make([]string, 0)
	}
	type SetStringExampleAlias SetStringExample
	return SetStringExampleAlias(o), nil
}

func (o *SetStringExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type SetStringExampleAlias SetStringExample
	var rawSetStringExample SetStringExampleAlias
	if err := unmarshal(&rawSetStringExample); err != nil {
		return err
	}
	if rawSetStringExample.Value == nil {
		rawSetStringExample.Value = make([]string, 0)
	}
	*o = SetStringExample(rawSetStringExample)
	return nil
}

type SetDoubleExample struct {
	Value []float64 `json:"value" yaml:"value,omitempty"`
}

func (o SetDoubleExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make([]float64, 0)
	}
	type SetDoubleExampleAlias SetDoubleExample
	return json.Marshal(SetDoubleExampleAlias(o))
}

func (o *SetDoubleExample) UnmarshalJSON(data []byte) error {
	type SetDoubleExampleAlias SetDoubleExample
	var rawSetDoubleExample SetDoubleExampleAlias
	if err := json.Unmarshal(data, &rawSetDoubleExample); err != nil {
		return err
	}
	if rawSetDoubleExample.Value == nil {
		rawSetDoubleExample.Value = make([]float64, 0)
	}
	*o = SetDoubleExample(rawSetDoubleExample)
	return nil
}

func (o SetDoubleExample) MarshalYAML() (interface{}, error) {
	if o.Value == nil {
		o.Value = make([]float64, 0)
	}
	type SetDoubleExampleAlias SetDoubleExample
	return SetDoubleExampleAlias(o), nil
}

func (o *SetDoubleExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type SetDoubleExampleAlias SetDoubleExample
	var rawSetDoubleExample SetDoubleExampleAlias
	if err := unmarshal(&rawSetDoubleExample); err != nil {
		return err
	}
	if rawSetDoubleExample.Value == nil {
		rawSetDoubleExample.Value = make([]float64, 0)
	}
	*o = SetDoubleExample(rawSetDoubleExample)
	return nil
}

type MapExample struct {
	Value map[string]string `json:"value" yaml:"value,omitempty"`
}

func (o MapExample) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		o.Value = make(map[string]string, 0)
	}
	type MapExampleAlias MapExample
	return json.Marshal(MapExampleAlias(o))
}

func (o *MapExample) UnmarshalJSON(data []byte) error {
	type MapExampleAlias MapExample
	var rawMapExample MapExampleAlias
	if err := json.Unmarshal(data, &rawMapExample); err != nil {
		return err
	}
	if rawMapExample.Value == nil {
		rawMapExample.Value = make(map[string]string, 0)
	}
	*o = MapExample(rawMapExample)
	return nil
}

func (o MapExample) MarshalYAML() (interface{}, error) {
	if o.Value == nil {
		o.Value = make(map[string]string, 0)
	}
	type MapExampleAlias MapExample
	return MapExampleAlias(o), nil
}

func (o *MapExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type MapExampleAlias MapExample
	var rawMapExample MapExampleAlias
	if err := unmarshal(&rawMapExample); err != nil {
		return err
	}
	if rawMapExample.Value == nil {
		rawMapExample.Value = make(map[string]string, 0)
	}
	*o = MapExample(rawMapExample)
	return nil
}

type OptionalExample struct {
	Value *string `json:"value" yaml:"value,omitempty"`
}

type LongOptionalExample struct {
	SomeLongName *string `json:"someLongName" yaml:"someLongName,omitempty"`
}

type EnumFieldExample struct {
	Enum EnumExample `json:"enum" yaml:"enum,omitempty"`
}

type EmptyObjectExample struct {
}

type ObjectExample struct {
	// docs for string field
	String string `json:"string" yaml:"string,omitempty" conjure-docs:"docs for string field"`
	// docs for integer field
	Integer int `json:"integer" yaml:"integer,omitempty" conjure-docs:"docs for integer field"`
	// docs for doubleValue field
	DoubleValue float64 `json:"doubleValue" yaml:"doubleValue,omitempty" conjure-docs:"docs for doubleValue field"`
	// docs for optionalItem field
	OptionalItem *string `json:"optionalItem" yaml:"optionalItem,omitempty" conjure-docs:"docs for optionalItem field"`
	// docs for items field
	Items []string `json:"items" yaml:"items,omitempty" conjure-docs:"docs for items field"`
	// docs for set field
	Set []string `json:"set" yaml:"set,omitempty" conjure-docs:"docs for set field"`
	// docs for map field
	Map map[string]string `json:"map" yaml:"map,omitempty" conjure-docs:"docs for map field"`
	// docs for alias field
	Alias StringAliasExample `json:"alias" yaml:"alias,omitempty" conjure-docs:"docs for alias field"`
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
	return json.Marshal(ObjectExampleAlias(o))
}

func (o *ObjectExample) UnmarshalJSON(data []byte) error {
	type ObjectExampleAlias ObjectExample
	var rawObjectExample ObjectExampleAlias
	if err := json.Unmarshal(data, &rawObjectExample); err != nil {
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
	return ObjectExampleAlias(o), nil
}

func (o *ObjectExample) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ObjectExampleAlias ObjectExample
	var rawObjectExample ObjectExampleAlias
	if err := unmarshal(&rawObjectExample); err != nil {
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

type KebabCaseObjectExample struct {
	KebabCasedField int `json:"kebab-cased-field" yaml:"kebab-cased-field,omitempty"`
}

type SnakeCaseObjectExample struct {
	SnakeCasedField int `json:"snake_cased_field" yaml:"snake_cased_field,omitempty"`
}

type TestCases struct {
	Client ClientTestCases `json:"client" yaml:"client,omitempty"`
}

type ClientTestCases struct {
	AutoDeserialize         map[EndpointName]PositiveAndNegativeTestCases `json:"autoDeserialize" yaml:"autoDeserialize,omitempty"`
	SingleHeaderService     map[EndpointName][]string                     `json:"singleHeaderService" yaml:"singleHeaderService,omitempty"`
	SinglePathParamService  map[EndpointName][]string                     `json:"singlePathParamService" yaml:"singlePathParamService,omitempty"`
	SingleQueryParamService map[EndpointName][]string                     `json:"singleQueryParamService" yaml:"singleQueryParamService,omitempty"`
}

func (o ClientTestCases) MarshalJSON() ([]byte, error) {
	if o.AutoDeserialize == nil {
		o.AutoDeserialize = make(map[EndpointName]PositiveAndNegativeTestCases, 0)
	}
	if o.SingleHeaderService == nil {
		o.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if o.SinglePathParamService == nil {
		o.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if o.SingleQueryParamService == nil {
		o.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	type ClientTestCasesAlias ClientTestCases
	return json.Marshal(ClientTestCasesAlias(o))
}

func (o *ClientTestCases) UnmarshalJSON(data []byte) error {
	type ClientTestCasesAlias ClientTestCases
	var rawClientTestCases ClientTestCasesAlias
	if err := json.Unmarshal(data, &rawClientTestCases); err != nil {
		return err
	}
	if rawClientTestCases.AutoDeserialize == nil {
		rawClientTestCases.AutoDeserialize = make(map[EndpointName]PositiveAndNegativeTestCases, 0)
	}
	if rawClientTestCases.SingleHeaderService == nil {
		rawClientTestCases.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if rawClientTestCases.SinglePathParamService == nil {
		rawClientTestCases.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if rawClientTestCases.SingleQueryParamService == nil {
		rawClientTestCases.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	*o = ClientTestCases(rawClientTestCases)
	return nil
}

func (o ClientTestCases) MarshalYAML() (interface{}, error) {
	if o.AutoDeserialize == nil {
		o.AutoDeserialize = make(map[EndpointName]PositiveAndNegativeTestCases, 0)
	}
	if o.SingleHeaderService == nil {
		o.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if o.SinglePathParamService == nil {
		o.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if o.SingleQueryParamService == nil {
		o.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	type ClientTestCasesAlias ClientTestCases
	return ClientTestCasesAlias(o), nil
}

func (o *ClientTestCases) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ClientTestCasesAlias ClientTestCases
	var rawClientTestCases ClientTestCasesAlias
	if err := unmarshal(&rawClientTestCases); err != nil {
		return err
	}
	if rawClientTestCases.AutoDeserialize == nil {
		rawClientTestCases.AutoDeserialize = make(map[EndpointName]PositiveAndNegativeTestCases, 0)
	}
	if rawClientTestCases.SingleHeaderService == nil {
		rawClientTestCases.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if rawClientTestCases.SinglePathParamService == nil {
		rawClientTestCases.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if rawClientTestCases.SingleQueryParamService == nil {
		rawClientTestCases.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	*o = ClientTestCases(rawClientTestCases)
	return nil
}

type PositiveAndNegativeTestCases struct {
	Positive []string `json:"positive" yaml:"positive,omitempty"`
	Negative []string `json:"negative" yaml:"negative,omitempty"`
}

func (o PositiveAndNegativeTestCases) MarshalJSON() ([]byte, error) {
	if o.Positive == nil {
		o.Positive = make([]string, 0)
	}
	if o.Negative == nil {
		o.Negative = make([]string, 0)
	}
	type PositiveAndNegativeTestCasesAlias PositiveAndNegativeTestCases
	return json.Marshal(PositiveAndNegativeTestCasesAlias(o))
}

func (o *PositiveAndNegativeTestCases) UnmarshalJSON(data []byte) error {
	type PositiveAndNegativeTestCasesAlias PositiveAndNegativeTestCases
	var rawPositiveAndNegativeTestCases PositiveAndNegativeTestCasesAlias
	if err := json.Unmarshal(data, &rawPositiveAndNegativeTestCases); err != nil {
		return err
	}
	if rawPositiveAndNegativeTestCases.Positive == nil {
		rawPositiveAndNegativeTestCases.Positive = make([]string, 0)
	}
	if rawPositiveAndNegativeTestCases.Negative == nil {
		rawPositiveAndNegativeTestCases.Negative = make([]string, 0)
	}
	*o = PositiveAndNegativeTestCases(rawPositiveAndNegativeTestCases)
	return nil
}

func (o PositiveAndNegativeTestCases) MarshalYAML() (interface{}, error) {
	if o.Positive == nil {
		o.Positive = make([]string, 0)
	}
	if o.Negative == nil {
		o.Negative = make([]string, 0)
	}
	type PositiveAndNegativeTestCasesAlias PositiveAndNegativeTestCases
	return PositiveAndNegativeTestCasesAlias(o), nil
}

func (o *PositiveAndNegativeTestCases) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type PositiveAndNegativeTestCasesAlias PositiveAndNegativeTestCases
	var rawPositiveAndNegativeTestCases PositiveAndNegativeTestCasesAlias
	if err := unmarshal(&rawPositiveAndNegativeTestCases); err != nil {
		return err
	}
	if rawPositiveAndNegativeTestCases.Positive == nil {
		rawPositiveAndNegativeTestCases.Positive = make([]string, 0)
	}
	if rawPositiveAndNegativeTestCases.Negative == nil {
		rawPositiveAndNegativeTestCases.Negative = make([]string, 0)
	}
	*o = PositiveAndNegativeTestCases(rawPositiveAndNegativeTestCases)
	return nil
}

type IgnoredTestCases struct {
	Client IgnoredClientTestCases `json:"client" yaml:"client,omitempty"`
}

type IgnoredClientTestCases struct {
	AutoDeserialize         map[EndpointName][]string `json:"autoDeserialize" yaml:"autoDeserialize,omitempty"`
	SingleHeaderService     map[EndpointName][]string `json:"singleHeaderService" yaml:"singleHeaderService,omitempty"`
	SinglePathParamService  map[EndpointName][]string `json:"singlePathParamService" yaml:"singlePathParamService,omitempty"`
	SingleQueryParamService map[EndpointName][]string `json:"singleQueryParamService" yaml:"singleQueryParamService,omitempty"`
}

func (o IgnoredClientTestCases) MarshalJSON() ([]byte, error) {
	if o.AutoDeserialize == nil {
		o.AutoDeserialize = make(map[EndpointName][]string, 0)
	}
	if o.SingleHeaderService == nil {
		o.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if o.SinglePathParamService == nil {
		o.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if o.SingleQueryParamService == nil {
		o.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	type IgnoredClientTestCasesAlias IgnoredClientTestCases
	return json.Marshal(IgnoredClientTestCasesAlias(o))
}

func (o *IgnoredClientTestCases) UnmarshalJSON(data []byte) error {
	type IgnoredClientTestCasesAlias IgnoredClientTestCases
	var rawIgnoredClientTestCases IgnoredClientTestCasesAlias
	if err := json.Unmarshal(data, &rawIgnoredClientTestCases); err != nil {
		return err
	}
	if rawIgnoredClientTestCases.AutoDeserialize == nil {
		rawIgnoredClientTestCases.AutoDeserialize = make(map[EndpointName][]string, 0)
	}
	if rawIgnoredClientTestCases.SingleHeaderService == nil {
		rawIgnoredClientTestCases.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if rawIgnoredClientTestCases.SinglePathParamService == nil {
		rawIgnoredClientTestCases.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if rawIgnoredClientTestCases.SingleQueryParamService == nil {
		rawIgnoredClientTestCases.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	*o = IgnoredClientTestCases(rawIgnoredClientTestCases)
	return nil
}

func (o IgnoredClientTestCases) MarshalYAML() (interface{}, error) {
	if o.AutoDeserialize == nil {
		o.AutoDeserialize = make(map[EndpointName][]string, 0)
	}
	if o.SingleHeaderService == nil {
		o.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if o.SinglePathParamService == nil {
		o.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if o.SingleQueryParamService == nil {
		o.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	type IgnoredClientTestCasesAlias IgnoredClientTestCases
	return IgnoredClientTestCasesAlias(o), nil
}

func (o *IgnoredClientTestCases) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type IgnoredClientTestCasesAlias IgnoredClientTestCases
	var rawIgnoredClientTestCases IgnoredClientTestCasesAlias
	if err := unmarshal(&rawIgnoredClientTestCases); err != nil {
		return err
	}
	if rawIgnoredClientTestCases.AutoDeserialize == nil {
		rawIgnoredClientTestCases.AutoDeserialize = make(map[EndpointName][]string, 0)
	}
	if rawIgnoredClientTestCases.SingleHeaderService == nil {
		rawIgnoredClientTestCases.SingleHeaderService = make(map[EndpointName][]string, 0)
	}
	if rawIgnoredClientTestCases.SinglePathParamService == nil {
		rawIgnoredClientTestCases.SinglePathParamService = make(map[EndpointName][]string, 0)
	}
	if rawIgnoredClientTestCases.SingleQueryParamService == nil {
		rawIgnoredClientTestCases.SingleQueryParamService = make(map[EndpointName][]string, 0)
	}
	*o = IgnoredClientTestCases(rawIgnoredClientTestCases)
	return nil
}
