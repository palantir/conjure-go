package dj

import (
	"fmt"
)

// SyntaxError is an error that occurs when parsing a json string.
type SyntaxError struct {
	Index int
	Msg   string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("invalid json at index %d: %s", e.Index, e.Msg)
}

type TypeMismatchError struct {
	Index int
	Want  string
	Got   Type
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch at index %d: want %s got %s", e.Index, e.Want, e.Got.String())
}

type InvalidValueError struct {
	Index int
	Msg   string
	Err   error
}

func (e InvalidValueError) Cause() error  { return e.Err }
func (e InvalidValueError) Unwrap() error { return e.Err }

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("invalid value at index %d: %s: %v", e.Index, e.Msg, e.Err)
}

type UnmarshalFieldError struct {
	Index int
	Type  string
	Field string
	Err   error
}

func (e UnmarshalFieldError) Cause() error  { return e.Err }
func (e UnmarshalFieldError) Unwrap() error { return e.Err }

func (e UnmarshalFieldError) Error() string {
	return fmt.Sprintf("field %s[%q] at index %d: %v", e.Type, e.Field, e.Index, e.Err)
}

type UnmarshalMissingFieldsError struct {
	Index  int
	Type   string
	Fields []string
}

func (e UnmarshalMissingFieldsError) Error() string {
	return fmt.Sprintf("type %s at index %d missing required fields: %v", e.Type, e.Index, e.Fields)
}

type UnmarshalUnknownFieldsError struct {
	Index  int
	Type   string
	Fields []string
}

func (e UnmarshalUnknownFieldsError) Error() string {
	return fmt.Sprintf("type %s at index %d encountered %d unknown fields: %v", e.Type, e.Index, len(e.Fields), e.Fields)
}

type UnmarshalDuplicateFieldError struct {
	Index int
	Type  string
	Field string
}

func (e UnmarshalDuplicateFieldError) Error() string {
	return fmt.Sprintf("field %s[%q] duplicated at index %d", e.Type, e.Field, e.Index)
}

type UnmarshalDuplicateMapKeyError struct {
	Index int
	Type  string
}

func (e UnmarshalDuplicateMapKeyError) Error() string {
	return fmt.Sprintf("field %s encountered duplicate map key at index %d", e.Type, e.Index)
}
