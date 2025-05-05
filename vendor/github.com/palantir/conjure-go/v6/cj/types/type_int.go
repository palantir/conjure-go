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
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// Int provides JSON marshaling and unmarshaling for integer types (signed).
// Encodes values as JSON numbers, and decodes JSON numbers into the underlying type.
type Int[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (Int[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.Int(int64(receiver)))
}

func (Int[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '0' {
		return cj.NewKindMismatchError(dec, kind, "json int")
	}
	*receiver = T(tok.Int())
	return nil
}

// IntMapKey provides JSON marshaling and unmarshaling for signed integer types used as map keys.
// Encodes integer keys as JSON strings to comply with JSON map key requirements.
type IntMapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (IntMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(strconv.FormatInt(int64(receiver), 10)))
}

func (IntMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(tok.String(), 10, 64)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid int", err)
	}
	*receiver = T(i)
	return nil
}

// Uint provides JSON marshaling and unmarshaling for unsigned integer types.
// Encodes values as JSON numbers, and decodes JSON numbers into the underlying type.
type Uint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (Uint[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.Uint(uint64(receiver)))
}

func (Uint[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '0' {
		return cj.NewKindMismatchError(dec, kind, "json uint")
	}
	*receiver = T(tok.Uint())
	return nil
}

// UintMapKey provides JSON marshaling and unmarshaling for unsigned integer types used as map keys.
// Encodes unsigned integer keys as JSON strings to comply with JSON map key requirements.
type UintMapKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (UintMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = strconv.AppendUint(out, uint64(receiver), 10)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (UintMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(tok.String(), 10, 64)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid uint", err)
	}
	*receiver = T(i)
	return nil
}
