// This file was generated by Conjure and should not be manually edited.

package buzz

import (
	"strings"
)

type Type1 struct {
	val Type1_Value
}

type Type1_Value string

const (
	Type1_value1  Type1_Value = "value1"
	Type1_value2  Type1_Value = "value2"
	Type1_UNKNOWN Type1_Value = "UNKNOWN"
)

// Type1_Values returns all known variants of Type1.
func Type1_Values() []Type1_Value {
	return []Type1_Value{Type1_value1, Type1_value2}
}

func New_Type1(value Type1_Value) Type1 {
	return Type1{val: value}
}

// IsUnknown returns false for all known variants of Type1 and true otherwise.
func (e Type1) IsUnknown() bool {
	switch e.val {
	case Type1_value1, Type1_value2:
		return false
	}
	return true
}

func (e Type1) Value() Type1_Value {
	if e.IsUnknown() {
		return Type1_UNKNOWN
	}
	return e.val
}

func (e Type1) String() string {
	return string(e.val)
}

func (e Type1) MarshalText() ([]byte, error) {
	return []byte(e.val), nil
}

func (e *Type1) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	default:
		*e = New_Type1(Type1_Value(v))
	case "value1":
		*e = New_Type1(Type1_value1)
	case "value2":
		*e = New_Type1(Type1_value2)
	}
	return nil
}