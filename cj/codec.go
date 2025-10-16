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
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Marshaling Utils //

// MarshalEncode is a variant of json.MarshalEncode that instantiates a new TypeEncoder and uses its MarshalJSONTo method.
func MarshalEncode[T any, E TypeEncoder[T]](enc *jsontext.Encoder, receiver T) error {
	return (*new(E)).MarshalJSONTo(enc, receiver)
}

type anonMarshalerTo[T any, E TypeEncoder[T]] struct {
	receiver T
}

// NewMarshalerTo constructs a new json.MarshalerTo that writes the JSON encoding of 'receiver' to 'enc'.
// Instead of using the default reflection-based behavior, it uses the provided TypeEncoder type to marshal the receiver.
func NewMarshalerTo[T any, E TypeEncoder[T]](receiver T) json.MarshalerTo {
	return anonMarshalerTo[T, E]{receiver: receiver}
}

func (f anonMarshalerTo[T, E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

func (f anonMarshalerTo[T, E]) MarshalJSONTo(enc *jsontext.Encoder) error {
	return MarshalEncode[T, E](enc, f.receiver)
}

// Unmarshaling Utils //

// UnmarshalDecode is a variant of json.UnmarshalDecode that instantiates a new TypeDecoder and uses its UnmarshalJSONFrom method.
func UnmarshalDecode[T any, D TypeDecoder[T]](dec *jsontext.Decoder, receiver *T) error {
	return (*new(D)).UnmarshalJSONFrom(dec, receiver)
}

type anonUnmarshalerFrom[T any, D TypeDecoder[T]] struct {
	receiver *T
}

// NewUnmarshalerFrom constructs a new json.UnmarshalerFrom that reads the JSON encoding of 'receiver' from 'dec'.
// Instead of using the default reflection-based behavior, it uses the provided TypeDecoder type to unmarshal the receiver.
func NewUnmarshalerFrom[T any, D TypeDecoder[T]](receiver *T) json.UnmarshalerFrom {
	return &anonUnmarshalerFrom[T, D]{receiver: receiver}
}

func (f *anonUnmarshalerFrom[T, D]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, f)
}

func (f *anonUnmarshalerFrom[T, D]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return UnmarshalDecode[T, D](dec, f.receiver)
}
