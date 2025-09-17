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
	"bytes"
	"encoding"
	"encoding/base64"
	"slices"
	"strings"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/binary"
)

// Binary provides JSON marshaling and unmarshaling for byte slice types (e.g., []byte).
// Encodes values as base64-encoded JSON strings using base64.StdEncoding.
// Implements comparison for equality and ordering based on byte content
// so values can be used as map keys.
type Binary[T ~[]byte] struct{}

func (Binary[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	dst := slices.Grow(enc.AvailableBuffer(), base64.StdEncoding.EncodedLen(len(receiver))+2)
	dst = append(dst, '"')
	dst = base64.StdEncoding.AppendEncode(dst, receiver)
	dst = append(dst, '"')
	return enc.WriteValue(dst)
}

func (Binary[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	if kind := dec.PeekKind(); kind != '"' {
		return NewKindMismatchError(dec, kind, "json string")
	}
	val, err := dec.ReadValue()
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	unquoted, err := jsontext.AppendUnquote(nil, val)
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	if len(unquoted) == 0 {
		*receiver = make(T, 0)
		return nil
	}
	decodedLen := base64.StdEncoding.DecodedLen(len(unquoted))
	if cap(*receiver) < decodedLen {
		*receiver = make(T, 0, decodedLen)
	} else {
		*receiver = (*receiver)[:0]
	}
	decoded, err := base64.StdEncoding.AppendDecode(*receiver, unquoted)
	if err != nil {
		return NewInvalidValueError(dec, "", err)
	}
	*receiver = decoded
	return nil
}

type BinaryMapKey[T binary.Binary] struct{}

func (BinaryMapKey[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return (String[binary.Binary]{}).MarshalJSONTo(enc, binary.Binary(receiver))
}

func (BinaryMapKey[T]) Compare(a, b T) int {
	return strings.Compare(string(a), string(b))
}

func (BinaryMapKey[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	b64 := tok.String()
	if _, err := base64.StdEncoding.DecodeString(b64); err != nil {
		return NewInvalidValueError(dec, "", err)
	}
	*receiver = T(b64)
	return nil
}

// BinaryMarshaler provides JSON marshaling for types implementing encoding.BinaryMarshaler.
// Values are marshaled as base64-encoded JSON strings using the MarshalBinary method.
type BinaryMarshaler[T encoding.BinaryMarshaler] struct{}

func (BinaryMarshaler[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	decoded, err := receiver.MarshalBinary()
	if err != nil {
		return WrapEncodeError(enc, "", err)
	}
	dst := enc.AvailableBuffer()
	dst = append(dst, '"')
	dst = base64.StdEncoding.AppendEncode(dst, decoded)
	dst = append(dst, '"')
	return enc.WriteValue(dst)
}

func (t BinaryMarshaler[T]) Compare(a, b T) int {
	aBinary, errA := a.MarshalBinary()
	bBinary, errB := b.MarshalBinary()
	if errA != nil || errB != nil {
		// If either fails, treat as equal (could log or handle differently)
		return 0
	}
	return bytes.Compare(aBinary, bBinary)
}

// BinaryUnmarshaler provides JSON unmarshaling for types implementing encoding.BinaryUnmarshaler.
// Expects base64-encoded JSON strings, which are decoded and passed to UnmarshalBinary.
type BinaryUnmarshaler[T encoding.BinaryUnmarshaler] struct{}

func (BinaryUnmarshaler[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver T) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(tok.String())
	if err != nil {
		return NewInvalidValueError(dec, "", err)
	}
	return receiver.UnmarshalBinary(decoded)
}
