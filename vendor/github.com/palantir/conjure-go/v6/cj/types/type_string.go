// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package types

import (
	"bytes"
	"encoding"
	"fmt"
	"strings"

	"github.com/go-json-experiment/json/jsontext"
)

// String provides JSON marshaling and unmarshaling for string-like types.
// Encodes values as JSON strings, and decodes JSON strings into the underlying type.
type String[T ~string] struct{}

func (String[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(string(receiver)))
}

func (String[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	*receiver = T(tok.String())
	return nil
}

// StringerMarshaler provides JSON marshaling for types implementing fmt.Stringer.
// Encodes values as JSON strings using the result of the String() method, and supports comparison by string value.
type StringerMarshaler[T fmt.Stringer] struct{}

func (StringerMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(receiver.String()))
}

func (StringerMarshaler[T]) Compare(a, b T) int {
	return strings.Compare(a.String(), b.String())
}

// stringUnmarshaler is like encoding.TextUnmarshaler but preferred when strings
// are more performant than []byte for the final type.
type stringUnmarshaler interface {
	UnmarshalString(string) error
}

// StringUnmarshaler provides JSON unmarshaling for types implementing a custom UnmarshalString method.
// Decodes JSON strings by calling UnmarshalString on the target type.
type StringUnmarshaler[T stringUnmarshaler] struct{}

func (StringUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	return receiver.UnmarshalString(tok.String())
}

// TextMarshaler provides JSON marshaling for types implementing encoding.TextMarshaler.
// Encodes values as JSON strings using the MarshalText method, and supports comparison by text value.
type TextMarshaler[T encoding.TextMarshaler] struct{}

func (TextMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	text, err := receiver.MarshalText()
	if err != nil {
		return err
	}
	return enc.WriteToken(jsontext.String(string(text)))
}

func (t TextMarshaler[T]) Compare(a, b T) int {
	aText, errA := a.MarshalText()
	bText, errB := b.MarshalText()
	if errA != nil || errB != nil {
		return 0
	}
	return bytes.Compare(aText, bText)
}

// TextUnmarshaler provides JSON unmarshaling for types implementing encoding.TextUnmarshaler.
// Decodes JSON strings by calling UnmarshalText on the target type.
type TextUnmarshaler[T encoding.TextUnmarshaler] struct{}

func (TextUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	return receiver.UnmarshalText([]byte(tok.String()))
}
