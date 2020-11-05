// This file was generated by Conjure and should not be manually edited.

package api

import (
	"regexp"
	"strings"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	werror "github.com/palantir/witchcraft-go-error"
	wparams "github.com/palantir/witchcraft-go-params"
)

var enumValuePattern = regexp.MustCompile("^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$")

type Days struct {
	val Days_Value
}

type Days_Value string

const (
	Days_FRIDAY   Days_Value = "FRIDAY"
	Days_SATURDAY Days_Value = "SATURDAY"
	Days_UNKNOWN  Days_Value = "UNKNOWN"
)

// Days_Values returns all known variants of Days.
func Days_Values() []Days_Value {
	return []Days_Value{Days_FRIDAY, Days_SATURDAY}
}

func New_Days(value Days_Value) Days {
	return Days{val: value}
}

// IsUnknown returns false for all known variants of Days and true otherwise.
func (e Days) IsUnknown() bool {
	switch e.val {
	case Days_FRIDAY, Days_SATURDAY:
		return false
	}
	return true
}

func (e Days) Value() Days_Value {
	if e.IsUnknown() {
		return Days_UNKNOWN
	}
	return e.val
}

func (e Days) String() string {
	return string(e.val)
}

func (e Days) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *Days) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		if !enumValuePattern.MatchString(v) {
			return werror.Convert(errors.NewInvalidArgument(wparams.NewSafeAndUnsafeParamStorer(map[string]interface{}{"enumType": "Days", "message": "enum value must match pattern ^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"}, map[string]interface{}{"enumValue": string(data)})))
		}
		*e = New_Days(Days_Value(v))
	case "FRIDAY":
		*e = New_Days(Days_FRIDAY)
	case "SATURDAY":
		*e = New_Days(Days_SATURDAY)
	}
	return nil
}

type Enum struct {
	val Enum_Value
}

type Enum_Value string

const (
	Enum_VALUE      Enum_Value = "VALUE"
	Enum_VALUES     Enum_Value = "VALUES"
	Enum_VALUES_1   Enum_Value = "VALUES_1"
	Enum_VALUES_1_1 Enum_Value = "VALUES_1_1"
	Enum_VALUE1     Enum_Value = "VALUE1"
	// Docs for an enum value
	Enum_VALUE2  Enum_Value = "VALUE2"
	Enum_UNKNOWN Enum_Value = "UNKNOWN"
)

// Enum_Values returns all known variants of Enum.
func Enum_Values() []Enum_Value {
	return []Enum_Value{Enum_VALUE, Enum_VALUES, Enum_VALUES_1, Enum_VALUES_1_1, Enum_VALUE1, Enum_VALUE2}
}

func New_Enum(value Enum_Value) Enum {
	return Enum{val: value}
}

// IsUnknown returns false for all known variants of Enum and true otherwise.
func (e Enum) IsUnknown() bool {
	switch e.val {
	case Enum_VALUE, Enum_VALUES, Enum_VALUES_1, Enum_VALUES_1_1, Enum_VALUE1, Enum_VALUE2:
		return false
	}
	return true
}

func (e Enum) Value() Enum_Value {
	if e.IsUnknown() {
		return Enum_UNKNOWN
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
	case "VALUE":
		*e = New_Enum(Enum_VALUE)
	case "VALUES":
		*e = New_Enum(Enum_VALUES)
	case "VALUES_1":
		*e = New_Enum(Enum_VALUES_1)
	case "VALUES_1_1":
		*e = New_Enum(Enum_VALUES_1_1)
	case "VALUE1":
		*e = New_Enum(Enum_VALUE1)
	case "VALUE2":
		*e = New_Enum(Enum_VALUE2)
	}
	return nil
}
