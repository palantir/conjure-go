package dj

import (
	"fmt"
	werror "github.com/palantir/witchcraft-go-error"
)

// SyntaxError is an error that occurs when parsing a json string.
type SyntaxError struct {
	baseErr
}

// NewSyntaxError returns a new SyntaxError.
func NewSyntaxError(index int, message string, err error) SyntaxError {
	msg := fmt.Sprintf("invalid json at index %d: %s", index, message)
	if err != nil {
		msg += fmt.Sprintf(": %v", err)
	}
	return SyntaxError{baseErr: newStack(msg, err)}
}

// TypeMismatchError occurs when a decoded value is not of the expected type.
type TypeMismatchError struct {
	baseErr
}

// NewTypeMismatchError returns a new TypeMismatchError.
func NewTypeMismatchError(res Result, want string) TypeMismatchError {
	msg := fmt.Sprintf("type mismatch at index %d: want %s got %s", res.Index, want, res.Type.String())
	return TypeMismatchError{baseErr: newStack(msg, nil)}
}

// InvalidValueError occurs when a decoded value is the correct type but not valid.
type InvalidValueError struct {
	baseErr
}

// NewInvalidValueError returns a new InvalidValueError.
func NewInvalidValueError(res Result, message string, err error) InvalidValueError {
	msg := fmt.Sprintf("invalid value at index %d: %s", res.Index, message)
	if err != nil {
		msg += fmt.Sprintf(": %v", err)
	}
	return InvalidValueError{baseErr: newStack(msg, err)}
}

// UnmarshalFieldError occurs when a struct field cannot be decoded.
type UnmarshalFieldError struct {
	baseErr
}

// NewUnmarshalFieldError returns a new UnmarshalFieldError.
func NewUnmarshalFieldError(res Result, fieldDescriptor string, err error) UnmarshalFieldError {
	msg := fmt.Sprintf("%s at index %d", fieldDescriptor, res.Index)
	if err != nil {
		msg += fmt.Sprintf(": %v", err)
	}
	return UnmarshalFieldError{baseErr: newStack(msg, err)}
}

// UnmarshalMissingFieldsError occurs when a struct is missing required fields.
type UnmarshalMissingFieldsError struct {
	baseErr
}

// NewUnmarshalMissingFieldsError returns a new UnmarshalMissingFieldsError.
func NewUnmarshalMissingFieldsError(res Result, typeName string, fields []string) UnmarshalMissingFieldsError {
	msg := fmt.Sprintf("type %s at index %d missing required fields: %v", typeName, res.Index, fields)
	return UnmarshalMissingFieldsError{baseErr: newStack(msg, nil)}
}

// UnmarshalUnknownFieldsError occurs when a struct has unknown fields.
type UnmarshalUnknownFieldsError struct {
	baseErr
}

// NewUnmarshalUnknownFieldsError returns a new UnmarshalUnknownFieldsError.
func NewUnmarshalUnknownFieldsError(res Result, typeName string, fields []string) UnmarshalUnknownFieldsError {
	msg := fmt.Sprintf("type %s at index %d encountered %d unknown fields: %v", typeName, res.Index, len(fields), fields)
	return UnmarshalUnknownFieldsError{newStack(msg, nil)}
}

// UnmarshalDuplicateFieldError occurs when a struct has duplicate fields.
type UnmarshalDuplicateFieldError struct {
	baseErr
}

// NewUnmarshalDuplicateFieldError returns a new UnmarshalDuplicateFieldError.
func NewUnmarshalDuplicateFieldError(res Result, fieldDescriptor string) UnmarshalDuplicateFieldError {
	msg := fmt.Sprintf("%s duplicated at index %d", fieldDescriptor, res.Index)
	return UnmarshalDuplicateFieldError{baseErr: newStack(msg, nil)}
}

// UnmarshalDuplicateMapKeyError occurs when a map has duplicate keys.
type UnmarshalDuplicateMapKeyError struct {
	baseErr
}

// NewUnmarshalDuplicateMapKeyError returns a new UnmarshalDuplicateMapKeyError.
func NewUnmarshalDuplicateMapKeyError(res Result, typeName string) UnmarshalDuplicateMapKeyError {
	msg := fmt.Sprintf("%s map key duplicated at index %d", typeName, res.Index)
	return UnmarshalDuplicateMapKeyError{baseErr: newStack(msg, nil)}
}

type baseErr struct {
	error string
	cause error
	stack werror.StackTrace
}

func newStack(error string, cause error) baseErr {
	return baseErr{error: error, cause: cause, stack: werror.NewStackTraceWithSkip(2)}
}

func (e baseErr) Error() string                 { return e.error }
func (e baseErr) StackTrace() werror.StackTrace { return e.stack }
func (e baseErr) Cause() error                  { return e.cause }
func (e baseErr) Unwrap() error                 { return e.cause }
