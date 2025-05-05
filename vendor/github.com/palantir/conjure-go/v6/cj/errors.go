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

package cj

import (
	"fmt"
	"strings"

	"github.com/go-json-experiment/json/jsontext"
	werror "github.com/palantir/witchcraft-go-error"
)

// Error is the interface that all errors in this package implement.
// This allows users to check if an error is a cj error.
type Error interface {
	cjError()
}

// SyntaxError is an error that occurs when parsing a json string.
type SyntaxError struct {
	message string
	baseErr
}

// NewSyntaxError returns a new SyntaxError.
func NewSyntaxError(dec *jsontext.Decoder, message string) SyntaxError {
	return SyntaxError{message: message, baseErr: newDecodeErr(dec, nil)}
}

// WrapSyntaxError returns a new SyntaxError with a cause.
func WrapSyntaxError(dec *jsontext.Decoder, message string, cause error) SyntaxError {
	return SyntaxError{message: message, baseErr: newDecodeErr(dec, cause)}
}

func (e SyntaxError) Error() string {
	return e.errString("SyntaxError", e.message)
}

// KindMismatchError occurs when a decoded value is not of the expected kind.
type KindMismatchError struct {
	got  jsontext.Kind
	want string
	baseErr
}

// NewKindMismatchError returns a new KindMismatchError for when the decoder encounters a kind mismatch.
func NewKindMismatchError(dec *jsontext.Decoder, got jsontext.Kind, want string) KindMismatchError {
	return KindMismatchError{
		got:     got,
		want:    want,
		baseErr: newDecodeErr(dec, nil),
	}
}

func (e KindMismatchError) Error() string {
	return e.errString("KindMismatchError", fmt.Sprintf("want %s, got %s", e.want, e.got.String()))
}

// InvalidValueError occurs when a decoded value is the correct type but otherwise not valid.
type InvalidValueError struct {
	message string
	baseErr
}

// NewInvalidValueError returns a new InvalidValueError.
func NewInvalidValueError(dec *jsontext.Decoder, message string, err error) InvalidValueError {
	return InvalidValueError{
		message: message,
		baseErr: newDecodeErr(dec, err),
	}
}

func (e InvalidValueError) Error() string {
	return e.errString("InvalidValueError", e.message)
}

// UnmarshalFieldError occurs when a struct field cannot be decoded.
type UnmarshalFieldError struct {
	fieldDescriptor string
	baseErr
}

// NewUnmarshalFieldError returns a new UnmarshalFieldError.
func NewUnmarshalFieldError(dec *jsontext.Decoder, fieldDescriptor string, err error) UnmarshalFieldError {
	return UnmarshalFieldError{
		fieldDescriptor: fieldDescriptor,
		baseErr:         newDecodeErr(dec, err),
	}
}

func (e UnmarshalFieldError) Error() string {
	return e.errString("UnmarshalFieldError", e.fieldDescriptor)
}

// MissingRequiredFieldsError occurs when a struct is missing required fields.
type MissingRequiredFieldsError struct {
	typeName string
	fields   []string
	baseErr
}

// NewMissingRequiredFieldsError returns a new MissingRequiredFieldsError.
func NewMissingRequiredFieldsError(dec *jsontext.Decoder, typeName string, fields []string) MissingRequiredFieldsError {
	return MissingRequiredFieldsError{
		typeName: typeName,
		fields:   fields,
		baseErr:  newDecodeErr(dec, nil),
	}
}

func (e MissingRequiredFieldsError) Error() string {
	return e.errString("MissingRequiredFieldsError", fmt.Sprintf("type %s missing required fields: %v", e.typeName, e.fields))
}

// UnknownFieldsError occurs when a struct has unknown fields.
type UnknownFieldsError struct {
	typeName string
	fields   []string
	baseErr
}

// NewUnknownFieldsError returns a new UnknownFieldsError.
func NewUnknownFieldsError(dec *jsontext.Decoder, typeName string, fields []string) UnknownFieldsError {
	return UnknownFieldsError{
		typeName: typeName,
		fields:   fields,
		baseErr:  newDecodeErr(dec, nil),
	}
}

func (e UnknownFieldsError) Error() string {
	return e.errString("UnknownFieldsError", fmt.Sprintf("type %s has unknown fields: %v", e.typeName, e.fields))
}

// DuplicateFieldKeyError occurs when a struct has duplicate fields.
type DuplicateFieldKeyError struct {
	typeName        string
	fieldDescriptor string
	baseErr
}

// NewDuplicateFieldKeyError returns a new DuplicateFieldKeyError.
func NewDuplicateFieldKeyError(dec *jsontext.Decoder, typeName string, fieldDescriptor string) DuplicateFieldKeyError {
	return DuplicateFieldKeyError{
		typeName:        typeName,
		fieldDescriptor: fieldDescriptor,
		baseErr:         newDecodeErr(dec, nil),
	}
}

func (e DuplicateFieldKeyError) Error() string {
	return e.errString("DuplicateFieldKeyError", fmt.Sprintf("field %s[%q] duplicated", e.typeName, e.fieldDescriptor))
}

// DuplicateMapKeyError occurs when a map has duplicate keys.
type DuplicateMapKeyError struct {
	typeName string
	baseErr
}

// NewDuplicateMapKeyError returns a new DuplicateMapKeyError.
func NewDuplicateMapKeyError(dec *jsontext.Decoder, typeName string) DuplicateMapKeyError {
	return DuplicateMapKeyError{
		typeName: typeName,
		baseErr:  newDecodeErr(dec, nil),
	}
}

func (e DuplicateMapKeyError) Error() string {
	return e.errString("DuplicateMapKeyError", fmt.Sprintf("type %s has duplicate map keys", e.typeName))
}

type baseErr struct {
	index   int64
	pointer jsontext.Pointer // JSON pointer to the error location
	cause   error
	stack   werror.StackTrace
}

func newDecodeErr(dec *jsontext.Decoder, cause error) baseErr {
	return baseErr{
		index:   dec.InputOffset(),
		pointer: dec.StackPointer(),
		cause:   cause,
		stack:   werror.NewStackTraceWithSkip(2),
	}
}

func (e baseErr) StackTrace() werror.StackTrace { return e.stack }
func (e baseErr) Cause() error                  { return e.cause }
func (e baseErr) Unwrap() error                 { return e.cause }

func (e baseErr) errString(prefix, msg string) string {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("%s at %d", prefix, e.index))
	if msg != "" {
		sb.WriteString(": ")
		sb.WriteString(msg)
	}
	if e.cause != nil {
		sb.WriteString(": ")
		sb.WriteString(e.cause.Error())
	}
	return sb.String()
}

func (e baseErr) cjError() {}
