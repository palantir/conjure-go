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

// ClientDecoder implements conjure-go-runtime's codecs.Decoder interface.
// It ignores unknown struct members.
type ClientDecoder[T any, D TypeDecoder[T]] struct{ codecBase }

func (ClientDecoder[T, D]) Decode(r io.Reader, v any) error {
	return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r), v.(*T))
}

func (ClientDecoder[T, D]) Unmarshal(data []byte, v any) error {
	return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewBuffer(data)), v.(*T))
}

// ClientEncoder implements conjure-go-runtime's codecs.Encoder interface.
type ClientEncoder[T any, E TypeEncoder[T]] struct{ codecBase }

func (ClientEncoder[T, E]) Encode(w io.Writer, v any) error {
	return (*new(E)).MarshalJSONTo(jsontext.NewEncoder(w), v.(T))
}

func (ClientEncoder[T, E]) Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := (*new(E)).MarshalJSONTo(jsontext.NewEncoder(buf), v.(T))
	return buf.Bytes(), err
}

// ServerDecoder implements conjure-go-runtime's codecs.Decoder interface.
// It rejects unknown struct members.
type ServerDecoder[T any, D TypeDecoder[T]] struct{ codecBase }

func (ServerDecoder[T, D]) Decode(r io.Reader, v any) error {
	return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r, json.RejectUnknownMembers(true)), v.(*T))
}

func (ServerDecoder[T, D]) Unmarshal(data []byte, v any) error {
	return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewBuffer(data), json.RejectUnknownMembers(true)), v.(*T))
}

// ServerEncoder implements conjure-go-runtime's codecs.Encoder interface.
type ServerEncoder[T any, E TypeEncoder[T]] struct{ codecBase }

func (ServerEncoder[T, E]) Encode(w io.Writer, v any) error {
	return (*new(E)).MarshalJSONTo(jsontext.NewEncoder(w), v.(T))
}

func (ServerEncoder[T, E]) Marshal(v any) ([]byte, error) {
	// TODO: pool & reuse buffers
	buf := bytes.NewBuffer(nil)
	err := (*new(E)).MarshalJSONTo(jsontext.NewEncoder(buf), v.(T))
	return buf.Bytes(), err
}

type codecBase struct{}

func (codecBase) Accept() string {
	return "application/json"
}

func (codecBase) ContentType() string {
	return "application/json"
}

func (codecBase) Encode(w io.Writer, v any) error {
	if marshaler, ok := v.(json.MarshalerTo); ok {
		return marshaler.MarshalJSONTo(jsontext.NewEncoder(w))
	}
	return json.MarshalWrite(w, v)
}

func (codecBase) Marshal(v any) ([]byte, error) {
	if marshaler, ok := v.(json.MarshalerTo); ok {
		buf := bytes.NewBuffer(nil)
		err := marshaler.MarshalJSONTo(jsontext.NewEncoder(buf))
		return buf.Bytes(), err
	}
	return json.Marshal(v)
}

func DecodeJSON[T any, D TypeDecoder[T]](r io.Reader, opts ...json.Options) (T, error) {
	var v T
	err := (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r, opts...), &v)
	return v, err
}

func EncodeJSON[T any, E TypeEncoder[T]](v T, opts ...json.Options) ([]byte, error) {
	// TODO: pool & reuse buffers
	buf := bytes.NewBuffer(nil)
	err := (*new(E)).MarshalJSONTo(jsontext.NewEncoder(buf, opts...), v)
	return buf.Bytes(), err
}
