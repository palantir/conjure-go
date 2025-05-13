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
	"github.com/palantir/pkg/safelong"
)

// Int32 provides JSON marshaling and unmarshaling for integer types (signed).
// Encodes values as JSON numbers, and decodes JSON numbers into the underlying type.
type Int32[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (Int32[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return enc.WriteToken(jsontext.Int(int64(receiver)))
}

func (Int32[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '0' {
		return cj.NewKindMismatchError(dec, kind, "json int")
	}
	num, err := strconv.ParseInt(tok.String(), 10, 32)
	if err != nil {
		return cj.NewInvalidValueError(dec, "invalid int32", err)
	}
	*receiver = T(num)
	return nil
}

// Int32MapKey provides JSON marshaling and unmarshaling for signed integer types used as map keys.
// Encodes integer keys as JSON strings to comply with JSON map key requirements.
type Int32MapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (Int32MapKey[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return enc.WriteToken(jsontext.String(strconv.FormatInt(int64(receiver), 10)))
}

func (Int32MapKey[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(tok.String(), 10, 32)
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid int32", err)
	}
	*receiver = T(i)
	return nil
}

// SafeLong provides JSON marshaling and unmarshaling for integer types (signed).
// Encodes values as JSON numbers, and decodes JSON numbers into the underlying type.
type SafeLong[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (SafeLong[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return enc.WriteToken(jsontext.Int(int64(receiver)))
}

func (SafeLong[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '0' {
		return cj.NewKindMismatchError(dec, kind, "json int")
	}
	num, err := safelong.ParseSafeLong(tok.String())
	if err != nil {
		return cj.NewInvalidValueError(dec, "invalid int32", err)
	}
	*receiver = T(num)
	return nil
}

// SafeLongMapKey provides JSON marshaling and unmarshaling for signed integer types used as map keys.
// Encodes integer keys as JSON strings to comply with JSON map key requirements.
type SafeLongMapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (SafeLongMapKey[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return enc.WriteToken(jsontext.String(strconv.FormatInt(int64(receiver), 10)))
}

func (SafeLongMapKey[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	i, err := safelong.ParseSafeLong(tok.String())
	if err != nil {
		return cj.NewInvalidValueError(dec, "", err)
	}
	*receiver = T(i)
	return nil
}
