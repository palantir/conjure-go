// Copyright (c) 2023 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dj

import (
	"fmt"

	werror "github.com/palantir/witchcraft-go-error"
)

// Error is the interface that all errors in this package implement.
// This allows users to check if an error is a dj error.
type Error interface {
	djError()
}

// SyntaxError is an error that occurs when parsing a json string.
type SyntaxError struct {
	message string
	baseErr
}

// NewSyntaxError returns a new SyntaxError.
func NewSyntaxError(index int, message string) SyntaxError {
	return SyntaxError{message: message, baseErr: newStack(index, nil)}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError at %d: %s", e.index, e.message)
}

// TypeMismatchError occurs when a decoded value is not of the expected type.
type TypeMismatchError struct {
	got  Type
	want string
	baseErr
}

// NewTypeMismatchError returns a new TypeMismatchError.
func NewTypeMismatchError(index int, got Type, want string) TypeMismatchError {
	return TypeMismatchError{
		got:     got,
		want:    want,
		baseErr: newStack(index, nil),
	}
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintf("TypeMismatchError at %d: want %s got %s", e.index, e.want, e.got.String())
}

// InvalidValueError occurs when a decoded value is the correct type but otherwise not valid.
type InvalidValueError struct {
	message string
	baseErr
}

// NewInvalidValueError returns a new InvalidValueError.
func NewInvalidValueError(index int, message string, err error) InvalidValueError {
	return InvalidValueError{message: message, baseErr: newStack(index, err)}
}

func (e InvalidValueError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("InvalidValueError at %d: %s: %v", e.index, e.message, e.cause)
	}
	return fmt.Sprintf("InvalidValueError at %d: %s", e.index, e.message)
}

// UnmarshalFieldError occurs when a struct field cannot be decoded.
type UnmarshalFieldError struct {
	fieldDescriptor string
	baseErr
}

// NewUnmarshalFieldError returns a new UnmarshalFieldError.
func NewUnmarshalFieldError(index int, fieldDescriptor string, err error) UnmarshalFieldError {
	return UnmarshalFieldError{fieldDescriptor: fieldDescriptor, baseErr: newStack(index, err)}
}

func (e UnmarshalFieldError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s at %d: %s", e.fieldDescriptor, e.index, e.cause)
	}
	return fmt.Sprintf("%s at %d", e.fieldDescriptor, e.index)
}

// UnmarshalMissingFieldsError occurs when a struct is missing required fields.
type UnmarshalMissingFieldsError struct {
	typeName string
	fields   []string
	baseErr
}

// NewUnmarshalMissingFieldsError returns a new UnmarshalMissingFieldsError.
func NewUnmarshalMissingFieldsError(index int, typeName string, fields []string) UnmarshalMissingFieldsError {
	return UnmarshalMissingFieldsError{typeName: typeName, fields: fields, baseErr: newStack(index, nil)}
}

func (e UnmarshalMissingFieldsError) Error() string {
	return fmt.Sprintf("type %s at index %d missing required fields: %v", e.typeName, e.index, e.fields)
}

// UnmarshalUnknownFieldsError occurs when a struct has unknown fields.
type UnmarshalUnknownFieldsError struct {
	typeName string
	fields   []string
	baseErr
}

// NewUnmarshalUnknownFieldsError returns a new UnmarshalUnknownFieldsError.
func NewUnmarshalUnknownFieldsError(index int, typeName string, fields []string) UnmarshalUnknownFieldsError {
	return UnmarshalUnknownFieldsError{typeName: typeName, fields: fields, baseErr: newStack(index, nil)}
}

func (e UnmarshalUnknownFieldsError) Error() string {
	return fmt.Sprintf("type %s at index %d encountered %d unknown fields: %v", e.typeName, e.index, len(e.fields), e.fields)
}

// UnmarshalDuplicateFieldError occurs when a struct has duplicate fields.
type UnmarshalDuplicateFieldError struct {
	fieldDescriptor string
	baseErr
}

// NewUnmarshalDuplicateFieldError returns a new UnmarshalDuplicateFieldError.
func NewUnmarshalDuplicateFieldError(index int, fieldDescriptor string) UnmarshalDuplicateFieldError {
	return UnmarshalDuplicateFieldError{fieldDescriptor: fieldDescriptor, baseErr: newStack(index, nil)}
}

func (e UnmarshalDuplicateFieldError) Error() string {
	return fmt.Sprintf("%s duplicated at index %d", e.fieldDescriptor, e.index)
}

// UnmarshalDuplicateMapKeyError occurs when a map has duplicate keys.
type UnmarshalDuplicateMapKeyError struct {
	typeName string
	baseErr
}

// NewUnmarshalDuplicateMapKeyError returns a new UnmarshalDuplicateMapKeyError.
func NewUnmarshalDuplicateMapKeyError(index int, typeName string) UnmarshalDuplicateMapKeyError {
	return UnmarshalDuplicateMapKeyError{typeName: typeName, baseErr: newStack(index, nil)}
}

func (e UnmarshalDuplicateMapKeyError) Error() string {
	return fmt.Sprintf("%s map key duplicated at index %d", e.typeName, e.index)
}

type baseErr struct {
	index int
	cause error
	stack werror.StackTrace
}

func newStack(index int, cause error) baseErr {
	return baseErr{index: index, cause: cause, stack: werror.NewStackTraceWithSkip(2)}
}

func (e baseErr) StackTrace() werror.StackTrace { return e.stack }
func (e baseErr) Cause() error                  { return e.cause }
func (e baseErr) Unwrap() error                 { return e.cause }
