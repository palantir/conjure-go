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
	werror.Werror
	cjError() // sealed interface
}

// SyntaxError is an error that occurs when parsing a json string.
type SyntaxError struct {
	baseErr
}

// NewSyntaxError returns a new SyntaxError.
func NewSyntaxError(dec *jsontext.Decoder, message string) SyntaxError {
	return SyntaxError{baseErr: newDecodeErr(dec, "SyntaxError", message, nil)}
}

// WrapSyntaxError returns a new SyntaxError with a cause.
func WrapSyntaxError(dec *jsontext.Decoder, message string, cause error) SyntaxError {
	return SyntaxError{baseErr: newDecodeErr(dec, "SyntaxError", message, cause)}
}

// EncodeError is an error that occurs when encoding a json string.
type EncodeError struct {
	baseErr
}

// NewEncodeError returns a new EncodeError.
func NewEncodeError(enc *jsontext.Encoder, message string) EncodeError {
	return EncodeError{baseErr: newEncodeErr(enc, "EncodeError", message, nil)}
}

// WrapEncodeError returns a new EncodeError with a cause.
func WrapEncodeError(enc *jsontext.Encoder, message string, cause error) EncodeError {
	return EncodeError{baseErr: newEncodeErr(enc, "EncodeError", message, cause)}
}

// KindMismatchError occurs when a decoded value is not of the expected kind.
type KindMismatchError struct {
	baseErr
}

// NewKindMismatchError returns a new KindMismatchError for when the decoder encounters a kind mismatch.
func NewKindMismatchError(dec *jsontext.Decoder, got jsontext.Kind, want string) KindMismatchError {
	return KindMismatchError{
		baseErr: newDecodeErr(dec, "KindMismatchError", fmt.Sprintf("want %s, got %s", want, got.String()), nil),
	}
}

// InvalidValueError occurs when a decoded value is the correct type but otherwise not valid.
type InvalidValueError struct {
	baseErr
}

// NewInvalidValueError returns a new InvalidValueError.
func NewInvalidValueError(dec *jsontext.Decoder, message string, err error) InvalidValueError {
	return InvalidValueError{baseErr: newDecodeErr(dec, "InvalidValueError", message, err)}
}

// UnmarshalFieldError occurs when a struct field cannot be decoded.
type UnmarshalFieldError struct {
	baseErr
}

// NewUnmarshalFieldError returns a new UnmarshalFieldError.
func NewUnmarshalFieldError(dec *jsontext.Decoder, fieldDescriptor string, cause error) UnmarshalFieldError {
	if cause != nil {
		if _, ok := cause.(Error); ok {
			return UnmarshalFieldError{baseErr: newDecodeErr(dec, "", fieldDescriptor, cause)}
		}
	}
	return UnmarshalFieldError{baseErr: newDecodeErr(dec, "UnmarshalFieldError", fieldDescriptor, cause)}
}

// MissingFieldsError occurs when a struct is missing required fields.
type MissingFieldsError struct {
	baseErr
}

// NewMissingFieldsError returns a new MissingFieldsError.
func NewMissingFieldsError(dec *jsontext.Decoder, typeName string, fields []string) MissingFieldsError {
	return MissingFieldsError{
		baseErr: newDecodeErr(dec, "MissingFieldsError", fmt.Sprintf("type %s missing required fields: %v", typeName, fields), nil),
	}
}

// UnknownFieldsError occurs when a struct has unknown fields.
type UnknownFieldsError struct {
	baseErr
}

// NewUnknownFieldsError returns a new UnknownFieldsError.
func NewUnknownFieldsError(dec *jsontext.Decoder, typeName string, fields []string) UnknownFieldsError {
	return UnknownFieldsError{
		baseErr: newDecodeErr(dec, "UnknownFieldsError", fmt.Sprintf("type %s has unknown fields: %v", typeName, fields), nil),
	}
}

// DuplicateFieldKeyError occurs when a struct has duplicate fields.
type DuplicateFieldKeyError struct {
	baseErr
}

// NewDuplicateFieldKeyError returns a new DuplicateFieldKeyError.
func NewDuplicateFieldKeyError(dec *jsontext.Decoder, fieldDescriptor string) DuplicateFieldKeyError {
	return DuplicateFieldKeyError{
		baseErr: newDecodeErr(dec, "DuplicateFieldKeyError", fmt.Sprintf("field %s duplicated", fieldDescriptor), nil),
	}
}

// DuplicateMapKeyError occurs when a map has duplicate keys.
type DuplicateMapKeyError struct {
	baseErr
}

// NewDuplicateMapKeyError returns a new DuplicateMapKeyError.
func NewDuplicateMapKeyError(dec *jsontext.Decoder, typeName string) DuplicateMapKeyError {
	return DuplicateMapKeyError{
		baseErr: newDecodeErr(dec, "DuplicateMapKeyError", fmt.Sprintf("type %s has duplicate map keys", typeName), nil),
	}
}

// DuplicateSetItemError occurs when a set has duplicate items.
type DuplicateSetItemError struct {
	baseErr
}

// NewDuplicateSetItemError returns a new DuplicateSetItemError.
func NewDuplicateSetItemError(dec *jsontext.Decoder, typeName string, index int) DuplicateSetItemError {
	return DuplicateSetItemError{
		baseErr: newDecodeErr(dec, "DuplicateSetItemError", fmt.Sprintf("type %s has a duplicate set item at index %d", typeName, index), nil),
	}
}

type baseErr struct {
	message string
	index   int64
	pointer jsontext.Pointer // JSON pointer to the error location
	cause   error
	stack   werror.StackTrace
}

func newDecodeErr(dec *jsontext.Decoder, prefix, msg string, cause error) baseErr {
	return baseErr{
		message: errString(prefix, msg, dec.InputOffset()),
		index:   dec.InputOffset(),
		pointer: dec.StackPointer(),
		cause:   cause,
		stack:   werror.NewStackTraceWithSkip(1),
	}
}

func newEncodeErr(enc *jsontext.Encoder, prefix, msg string, cause error) baseErr {
	return baseErr{
		message: errString(prefix, msg, enc.OutputOffset()),
		index:   enc.OutputOffset(),
		pointer: enc.StackPointer(),
		cause:   cause,
		stack:   werror.NewStackTraceWithSkip(1),
	}
}

func (e baseErr) Error() string {
	if e.cause == nil {
		return e.message
	}
	return e.message + ": " + e.cause.Error()
}

func (e baseErr) Message() string {
	return e.message
}

func (e baseErr) Format(state fmt.State, verb rune) {
	werror.Format(e, e.SafeParams(), state, verb)
}

func (e baseErr) StackTrace() werror.StackTrace { return e.stack }
func (e baseErr) Cause() error                  { return e.cause }
func (e baseErr) Unwrap() error                 { return e.cause }

func (e baseErr) SafeParams() map[string]interface{} {
	return map[string]interface{}{}
}

func (e baseErr) UnsafeParams() map[string]interface{} {
	return map[string]interface{}{}
}

func errString(prefix, msg string, index int64) string {
	sb := new(strings.Builder)
	if prefix != "" {
		sb.WriteString(fmt.Sprintf("%s at offset %d", prefix, index))
	}
	if msg != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(msg)
	}
	return sb.String()
}

func (e baseErr) cjError() {}
