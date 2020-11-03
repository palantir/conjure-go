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
	val EnumValue
}

type EnumValue string

const (
	EnumOne     EnumValue = "ONE"
	EnumTwo     EnumValue = "TWO"
	EnumUnknown EnumValue = "UNKNOWN"
)

// Enum_Values returns all known variants of Enum.
func Enum_Values() []EnumValue {
	return []EnumValue{EnumOne, EnumTwo}
}

func NewEnum(value EnumValue) Enum {
	return Enum{val: value}
}

// IsUnknown returns false for all known variants of Enum and true otherwise.
func (e Enum) IsUnknown() bool {
	switch e.val {
	case EnumOne, EnumTwo:
		return false
	}
	return true
}

func (e Enum) Value() EnumValue {
	if e.IsUnknown() {
		return EnumUnknown
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
		*e = NewEnum(EnumValue(v))
	case "ONE":
		*e = NewEnum(EnumOne)
	case "TWO":
		*e = NewEnum(EnumTwo)
	}
	return nil
}

type EnumExample struct {
	val EnumExampleValue
}

type EnumExampleValue string

const (
	EnumExampleOne        EnumExampleValue = "ONE"
	EnumExampleTwo        EnumExampleValue = "TWO"
	EnumExampleOneHundred EnumExampleValue = "ONE_HUNDRED"
	EnumExampleUnknown    EnumExampleValue = "UNKNOWN"
)

// EnumExample_Values returns all known variants of EnumExample.
func EnumExample_Values() []EnumExampleValue {
	return []EnumExampleValue{EnumExampleOne, EnumExampleTwo, EnumExampleOneHundred}
}

func NewEnumExample(value EnumExampleValue) EnumExample {
	return EnumExample{val: value}
}

// IsUnknown returns false for all known variants of EnumExample and true otherwise.
func (e EnumExample) IsUnknown() bool {
	switch e.val {
	case EnumExampleOne, EnumExampleTwo, EnumExampleOneHundred:
		return false
	}
	return true
}

func (e EnumExample) Value() EnumExampleValue {
	if e.IsUnknown() {
		return EnumExampleUnknown
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
		*e = NewEnumExample(EnumExampleValue(v))
	case "ONE":
		*e = NewEnumExample(EnumExampleOne)
	case "TWO":
		*e = NewEnumExample(EnumExampleTwo)
	case "ONE_HUNDRED":
		*e = NewEnumExample(EnumExampleOneHundred)
	}
	return nil
}
