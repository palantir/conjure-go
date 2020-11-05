// This file was generated by Conjure and should not be manually edited.

package types

import (
	"regexp"
	"strings"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	werror "github.com/palantir/witchcraft-go-error"
	wparams "github.com/palantir/witchcraft-go-params"
)

var enumValuePattern = regexp.MustCompile("^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$")

type Enum struct {
	val Enum_Value
}

type Enum_Value string

const (
	Enum_One     Enum_Value = "ONE"
	Enum_Two     Enum_Value = "TWO"
	Enum_Unknown Enum_Value = "UNKNOWN"
)

// Enum_Values returns all known variants of Enum.
func Enum_Values() []Enum_Value {
	return []Enum_Value{Enum_One, Enum_Two}
}

func New_Enum(value Enum_Value) Enum {
	return Enum{val: value}
}

// IsUnknown returns false for all known variants of Enum and true otherwise.
func (e Enum) IsUnknown() bool {
	switch e.val {
	case Enum_One, Enum_Two:
		return false
	}
	return true
}

func (e Enum) Value() Enum_Value {
	if e.IsUnknown() {
		return Enum_Unknown
	}
	return e.val
}

func (e Enum) String() string {
	return string(e.val)
}

func (e Enum) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *Enum) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Enum", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = New_Enum(Enum_Value(v))
	case "ONE":
		*e = New_Enum(Enum_One)
	case "TWO":
		*e = New_Enum(Enum_Two)
	}
	return nil
}

type EnumExample struct {
	val EnumExample_Value
}

type EnumExample_Value string

const (
	EnumExample_One        EnumExample_Value = "ONE"
	EnumExample_Two        EnumExample_Value = "TWO"
	EnumExample_OneHundred EnumExample_Value = "ONE_HUNDRED"
	EnumExample_Unknown    EnumExample_Value = "UNKNOWN"
)

// EnumExample_Values returns all known variants of EnumExample.
func EnumExample_Values() []EnumExample_Value {
	return []EnumExample_Value{EnumExample_One, EnumExample_Two, EnumExample_OneHundred}
}

func New_EnumExample(value EnumExample_Value) EnumExample {
	return EnumExample{val: value}
}

// IsUnknown returns false for all known variants of EnumExample and true otherwise.
func (e EnumExample) IsUnknown() bool {
	switch e.val {
	case EnumExample_One, EnumExample_Two, EnumExample_OneHundred:
		return false
	}
	return true
}

func (e EnumExample) Value() EnumExample_Value {
	if e.IsUnknown() {
		return EnumExample_Unknown
	}
	return e.val
}

func (e EnumExample) String() string {
	return string(e.val)
}

func (e EnumExample) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *EnumExample) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "EnumExample", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = New_EnumExample(EnumExample_Value(v))
	case "ONE":
		*e = New_EnumExample(EnumExample_One)
	case "TWO":
		*e = New_EnumExample(EnumExample_Two)
	case "ONE_HUNDRED":
		*e = New_EnumExample(EnumExample_OneHundred)
	}
	return nil
}
