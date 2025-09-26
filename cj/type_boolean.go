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

package cj

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// Boolean provides JSON marshaling and unmarshaling for Go bool-like types.
// Encodes values as JSON true/false, and decodes JSON booleans into the underlying type.
type Boolean[T ~bool] struct{}

func (Boolean[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver {
		return enc.WriteToken(jsontext.True)
	}
	return enc.WriteToken(jsontext.False)
}

func (Boolean[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	switch tok.Kind() {
	case 't':
		*receiver = true
	case 'f':
		*receiver = false
	default:
		return NewKindMismatchError(dec, tok.Kind(), "json boolean")
	}
	return nil
}

// BooleanMapKey provides JSON marshaling for bool-like types used as map keys.
// Encodes bool keys as the JSON strings "true" or "false" (not as booleans), to comply with JSON map key requirements.
type BooleanMapKey[T ~bool] struct{}

func (BooleanMapKey[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver {
		return enc.WriteValue(jsontext.Value(`"true"`))
	}
	return enc.WriteValue(jsontext.Value(`"false"`))
}

func (BooleanMapKey[T]) Compare(a, b T) int {
	if a == b {
		return 0
	}
	if a {
		return 1
	}
	return -1
}

func (BooleanMapKey[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	if kind := tok.Kind(); kind != '"' {
		return NewKindMismatchError(dec, kind, "json string")
	}
	b, err := strconv.ParseBool(tok.String())
	if err != nil {
		return NewInvalidValueError(dec, "invalid boolean", err)
	}
	*receiver = T(b)
	return nil
}
