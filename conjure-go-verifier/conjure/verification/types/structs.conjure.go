// This file was generated by Conjure and should not be manually edited.

package types

import (
	"encoding/base64"
	"math"
	"strconv"

	bearertoken "github.com/palantir/pkg/bearertoken"
	datetime "github.com/palantir/pkg/datetime"
	rid "github.com/palantir/pkg/rid"
	safejson "github.com/palantir/pkg/safejson"
	safelong "github.com/palantir/pkg/safelong"
	uuid "github.com/palantir/pkg/uuid"
)

type AnyExample struct {
	Value interface{} `json:"value"`
}

func (o AnyExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o AnyExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		if o.Value == nil {
			out = append(out, "null"...)
		} else if jsonBytes, err := safejson.Marshal(o.Value); err != nil {
			return nil, err
		} else {
			out = append(out, jsonBytes...)
		}
	}
	out = append(out, '}')
	return out, nil
}

type BearerTokenExample struct {
	Value bearertoken.Token `json:"value"`
}

func (o BearerTokenExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o BearerTokenExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = safejson.AppendQuotedString(out, o.Value.String())
	}
	out = append(out, '}')
	return out, nil
}

type BinaryExample struct {
	Value []byte `json:"value"`
}

func (o BinaryExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o BinaryExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = append(out, '"')
		if len(o.Value) > 0 {
			b64out := make([]byte, 0, base64.StdEncoding.EncodedLen(len(o.Value)))
			base64.StdEncoding.Encode(b64out, o.Value)
			out = append(out, b64out...)
		}
		out = append(out, '"')
	}
	out = append(out, '}')
	return out, nil
}

type BooleanExample struct {
	Value bool `json:"value"`
}

func (o BooleanExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o BooleanExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		if o.Value {
			out = append(out, "true"...)
		} else {
			out = append(out, "false"...)
		}
		out = append(out, "true"...)
		out = append(out, "false"...)
	}
	out = append(out, '}')
	return out, nil
}

type DateTimeExample struct {
	Value datetime.DateTime `json:"value"`
}

func (o DateTimeExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o DateTimeExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = safejson.AppendQuotedString(out, o.Value.String())
	}
	out = append(out, '}')
	return out, nil
}

type DoubleExample struct {
	Value float64 `json:"value"`
}

func (o DoubleExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o DoubleExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		switch {
		default:
			out = strconv.AppendFloat(out, o.Value, -1, 10, 64)
		case math.IsNaN(o.Value):
			out = append(out, "\"NaN\""...)
		case math.IsInf(o.Value, 1):
			out = append(out, "\"Infinity\""...)
		case math.IsInf(o.Value, -1):
			out = append(out, "\"-Infinity\""...)
		}
	}
	out = append(out, '}')
	return out, nil
}

type EmptyObjectExample struct{}

func (o EmptyObjectExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o EmptyObjectExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	out = append(out, '}')
	return out, nil
}

type EnumFieldExample struct {
	Enum EnumExample `json:"enum"`
}

func (o EnumFieldExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o EnumFieldExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"enum\":"...)
		out = safejson.AppendQuotedString(out, o.Enum.String())
	}
	out = append(out, '}')
	return out, nil
}

type IntegerExample struct {
	Value int `json:"value"`
}

func (o IntegerExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o IntegerExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = strconv.AppendInt(out, int64(o.Value), 10)
	}
	out = append(out, '}')
	return out, nil
}

type KebabCaseObjectExample struct {
	KebabCasedField int `json:"kebab-cased-field"`
}

func (o KebabCaseObjectExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o KebabCaseObjectExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"kebab-cased-field\":"...)
		out = strconv.AppendInt(out, int64(o.KebabCasedField), 10)
	}
	out = append(out, '}')
	return out, nil
}

type ListExample struct {
	Value []string `json:"value"`
}

func (o ListExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o ListExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = append(out, '[')
		{
			for i := range o.Value {
				out = safejson.AppendQuotedString(out, o.Value[i])
				if i < len(o.Value)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
	}
	out = append(out, '}')
	return out, nil
}

type LongFieldNameOptionalExample struct {
	SomeLongName *string `json:"someLongName"`
}

func (o LongFieldNameOptionalExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o LongFieldNameOptionalExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"someLongName\":"...)
		if o.SomeLongName != nil {
			optVal := *o.SomeLongName
			out = safejson.AppendQuotedString(out, optVal)
		} else {
			out = append(out, "null"...)
		}
	}
	out = append(out, '}')
	return out, nil
}

type MapExample struct {
	Value map[string]string `json:"value"`
}

func (o MapExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o MapExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.Value {
				out = safejson.AppendQuotedString(out, k)
				out = append(out, ':')
				out = safejson.AppendQuotedString(out, v)
				i++
				if i < len(o.Value) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
	}
	out = append(out, '}')
	return out, nil
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
	return o.AppendJSON(nil)
}

func (o ObjectExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"string\":"...)
		out = safejson.AppendQuotedString(out, o.String)
		out = append(out, ',')
	}
	{
		out = append(out, "\"integer\":"...)
		out = strconv.AppendInt(out, int64(o.Integer), 10)
		out = append(out, ',')
	}
	{
		out = append(out, "\"doubleValue\":"...)
		switch {
		default:
			out = strconv.AppendFloat(out, o.DoubleValue, -1, 10, 64)
		case math.IsNaN(o.DoubleValue):
			out = append(out, "\"NaN\""...)
		case math.IsInf(o.DoubleValue, 1):
			out = append(out, "\"Infinity\""...)
		case math.IsInf(o.DoubleValue, -1):
			out = append(out, "\"-Infinity\""...)
		}
		out = append(out, ',')
	}
	{
		out = append(out, "\"optionalItem\":"...)
		if o.OptionalItem != nil {
			optVal := *o.OptionalItem
			out = safejson.AppendQuotedString(out, optVal)
		} else {
			out = append(out, "null"...)
		}
		out = append(out, ',')
	}
	{
		out = append(out, "\"items\":"...)
		out = append(out, '[')
		{
			for i := range o.Items {
				out = safejson.AppendQuotedString(out, o.Items[i])
				if i < len(o.Items)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
		out = append(out, ',')
	}
	{
		out = append(out, "\"set\":"...)
		out = append(out, '[')
		{
			for i := range o.Set {
				out = safejson.AppendQuotedString(out, o.Set[i])
				if i < len(o.Set)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
		out = append(out, ',')
	}
	{
		out = append(out, "\"map\":"...)
		out = append(out, '{')
		{
			var i int
			for k, v := range o.Map {
				out = safejson.AppendQuotedString(out, k)
				out = append(out, ':')
				out = safejson.AppendQuotedString(out, v)
				i++
				if i < len(o.Map) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		out = append(out, ',')
	}
	{
		out = append(out, "\"alias\":"...)
		var err error
		out, err = o.Alias.AppendJSON(out)
		if err != nil {
			return nil, err
		}
	}
	out = append(out, '}')
	return out, nil
}

type OptionalBooleanExample struct {
	Value *bool `json:"value"`
}

func (o OptionalBooleanExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o OptionalBooleanExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		if o.Value != nil {
			optVal := *o.Value
			if optVal {
				out = append(out, "true"...)
			} else {
				out = append(out, "false"...)
			}
			out = append(out, "true"...)
			out = append(out, "false"...)
		} else {
			out = append(out, "null"...)
		}
	}
	out = append(out, '}')
	return out, nil
}

type OptionalExample struct {
	Value *string `json:"value"`
}

func (o OptionalExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o OptionalExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		if o.Value != nil {
			optVal := *o.Value
			out = safejson.AppendQuotedString(out, optVal)
		} else {
			out = append(out, "null"...)
		}
	}
	out = append(out, '}')
	return out, nil
}

type OptionalIntegerExample struct {
	Value *int `json:"value"`
}

func (o OptionalIntegerExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o OptionalIntegerExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		if o.Value != nil {
			optVal := *o.Value
			out = strconv.AppendInt(out, int64(optVal), 10)
		} else {
			out = append(out, "null"...)
		}
	}
	out = append(out, '}')
	return out, nil
}

type RidExample struct {
	Value rid.ResourceIdentifier `json:"value"`
}

func (o RidExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o RidExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = safejson.AppendQuotedString(out, o.Value.String())
	}
	out = append(out, '}')
	return out, nil
}

type SafeLongExample struct {
	Value safelong.SafeLong `json:"value"`
}

func (o SafeLongExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o SafeLongExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = strconv.AppendInt(out, int64(o.Value), 10)
	}
	out = append(out, '}')
	return out, nil
}

type SetDoubleExample struct {
	Value []float64 `json:"value"`
}

func (o SetDoubleExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o SetDoubleExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = append(out, '[')
		{
			for i := range o.Value {
				switch {
				default:
					out = strconv.AppendFloat(out, o.Value[i], -1, 10, 64)
				case math.IsNaN(o.Value[i]):
					out = append(out, "\"NaN\""...)
				case math.IsInf(o.Value[i], 1):
					out = append(out, "\"Infinity\""...)
				case math.IsInf(o.Value[i], -1):
					out = append(out, "\"-Infinity\""...)
				}
				if i < len(o.Value)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
	}
	out = append(out, '}')
	return out, nil
}

type SetStringExample struct {
	Value []string `json:"value"`
}

func (o SetStringExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o SetStringExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = append(out, '[')
		{
			for i := range o.Value {
				out = safejson.AppendQuotedString(out, o.Value[i])
				if i < len(o.Value)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
	}
	out = append(out, '}')
	return out, nil
}

type SnakeCaseObjectExample struct {
	SnakeCasedField int `json:"snake_cased_field"`
}

func (o SnakeCaseObjectExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o SnakeCaseObjectExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"snake_cased_field\":"...)
		out = strconv.AppendInt(out, int64(o.SnakeCasedField), 10)
	}
	out = append(out, '}')
	return out, nil
}

type StringExample struct {
	Value string `json:"value"`
}

func (o StringExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o StringExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = safejson.AppendQuotedString(out, o.Value)
	}
	out = append(out, '}')
	return out, nil
}

type UuidExample struct {
	Value uuid.UUID `json:"value"`
}

func (o UuidExample) MarshalJSON() ([]byte, error) {
	return o.AppendJSON(nil)
}

func (o UuidExample) AppendJSON(out []byte) ([]byte, error) {
	out = append(out, '{')
	{
		out = append(out, "\"value\":"...)
		out = safejson.AppendQuotedString(out, o.Value.String())
	}
	out = append(out, '}')
	return out, nil
}
