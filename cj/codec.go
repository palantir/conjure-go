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
	werror "github.com/palantir/witchcraft-go-error"
)

// ClientDecoder implements conjure-go-runtime's codecs.Decoder interface.
// It ignores unknown struct members.
type ClientDecoder[T any, D TypeDecoder[T]] struct{ codecBase }

func (ClientDecoder[T, D]) Decode(r io.Reader, v any) error {
	if v == nil {
		return werror.Error("decode target cannot be nil")
	}
	vt, ok := v.(*T)
	if !ok {
		return werror.Error("decode target is incompatible with decoder")
	}

	return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r), vt)
}

func (ClientDecoder[T, D]) Unmarshal(data []byte, v any) error {
	return ClientDecoder[T, D]{}.Decode(bytes.NewBuffer(data), v)
}

// ClientEncoder implements conjure-go-runtime's codecs.Encoder interface.
type ClientEncoder[T any, E TypeEncoder[T]] struct{ codecBase }

func (ClientEncoder[T, E]) Encode(w io.Writer, v any) error {
	vt, ok := v.(T)
	if !ok {
		return werror.Error("encode source is incompatible with encoder")
	}
	return (*new(E)).MarshalJSONTo(jsontext.NewEncoder(w), vt)
}

func (ClientEncoder[T, E]) Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := ClientEncoder[T, E]{}.Encode(buf, v)
	return buf.Bytes(), err
}

// ServerDecoder implements conjure-go-runtime's codecs.Decoder interface.
// It rejects unknown struct members.
type ServerDecoder[T any, D TypeDecoder[T]] struct{ codecBase }

func (ServerDecoder[T, D]) Decode(r io.Reader, v any) error {
	if v == nil {
		return werror.Error("decode target cannot be nil")
	}
	vt, ok := v.(*T)
	if !ok {
		return werror.Error("decode target is incompatible with decoder")
	}
	return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r, json.RejectUnknownMembers(true)), vt)
}

func (ServerDecoder[T, D]) Unmarshal(data []byte, v any) error {
	return ServerDecoder[T, D]{}.Decode(bytes.NewBuffer(data), v)
}

// ServerEncoder implements conjure-go-runtime's codecs.Encoder interface.
type ServerEncoder[T any, E TypeEncoder[T]] struct{ codecBase }

func (ServerEncoder[T, E]) Encode(w io.Writer, v any) error {
	vt, ok := v.(T)
	if !ok {
		vtp, ok := v.(*T)
		if !ok || vtp == nil {
			return werror.Error("encode source cannot be nil")
		}
		vt = *vtp
	}
	return (*new(E)).MarshalJSONTo(jsontext.NewEncoder(w), vt)
}

func (ServerEncoder[T, E]) Marshal(v any) ([]byte, error) {
	// TODO: pool & reuse buffers
	buf := bytes.NewBuffer(nil)
	err := ServerEncoder[T, E]{}.Encode(buf, v)
	return buf.Bytes(), err
}

type codecBase struct{}

func (codecBase) Accept() string {
	return "application/json"
}

func (codecBase) ContentType() string {
	return "application/json"
}
