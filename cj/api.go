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
	"io"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// TypeEncoder is implemented by types that can encode a Go value of type T to JSON using the provided jsontext.Encoder.
// Implementations for each Conjure type (e.g., Boolean, Integer, String, List, Map, etc.) are found in the types/ package.
// Each implementation ensures correct marshaling of the corresponding Go type to the appropriate JSON representation.
// Implementations' zero values must be valid for use by container encoders.
type TypeEncoder[T any] interface {
	// MarshalJSONTo writes the JSON encoding of 'receiver' to 'enc'.
	MarshalJSONTo(enc *jsontext.Encoder, receiver T) error
}

// TypeDecoder is implemented by types that can decode JSON into a Go value of type T using the provided jsontext.Decoder.
// Implementations in the types/ package handle unmarshaling for each supported Conjure type, including type validation and error handling.
// Implementations' zero values must be valid for use by container decoders.
type TypeDecoder[T any] interface {
	// UnmarshalJSONFrom reads JSON from 'dec' and stores the result in 'receiver'.
	UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error
}

// MapKeyEncoder is implemented by types that can be used as map keys in Conjure types
// but do not implement cmp.Ordered. The encoder's Compare method is used to sort map keys in a deterministic order.
// TypeEncoder implementations for comparable types (numbers, strings, etc) should not implement Compare.
type MapKeyEncoder[K comparable] interface {
	TypeEncoder[K]

	// Compare returns -1 if a < b, 0 if a == b, and 1 if a > b.
	// This is used to sort keys in a deterministic order.
	Compare(K, K) int
}

// Marshal is a variant of json.Marshal that instantiates a new TypeEncoder and uses its MarshalJSONTo method.
func Marshal[T any, E TypeEncoder[T]](receiver T, opts ...json.Options) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := MarshalWrite[T, E](buf, receiver, opts...)
	return buf.Bytes(), err
}

// MarshalWrite is a variant of json.MarshalWrite that instantiates a new TypeEncoder and uses its MarshalJSONTo method.
func MarshalWrite[T any, E TypeEncoder[T]](out io.Writer, receiver T, opts ...json.Options) error {
	return MarshalEncode[T, E](jsontext.NewEncoder(out, opts...), receiver)
}

// MarshalEncode is a variant of json.MarshalEncode that instantiates a new TypeEncoder and uses its MarshalJSONTo method.
func MarshalEncode[T any, E TypeEncoder[T]](enc *jsontext.Encoder, receiver T) error {
	return (*new(E)).MarshalJSONTo(enc, receiver)
}

// Unmarshal is a variant of json.Unmarshal that instantiates a new TypeDecoder and uses its UnmarshalJSONFrom method.
func Unmarshal[T any, D TypeDecoder[T]](data []byte, receiver *T, opts ...json.Options) error {
	return UnmarshalRead[T, D](bytes.NewReader(data), receiver, opts...)
}

// UnmarshalRead is a variant of json.UnmarshalRead that instantiates a new TypeDecoder and uses its UnmarshalJSONFrom method.
func UnmarshalRead[T any, D TypeDecoder[T]](in io.Reader, receiver *T, opts ...json.Options) error {
	return UnmarshalDecode[T, D](jsontext.NewDecoder(in, opts...), receiver)
}

// UnmarshalDecode is a variant of json.UnmarshalDecode that instantiates a new TypeDecoder and uses its UnmarshalJSONFrom method.
func UnmarshalDecode[T any, D TypeDecoder[T]](dec *jsontext.Decoder, receiver *T) error {
	return (*new(D)).UnmarshalJSONFrom(dec, receiver)
}
