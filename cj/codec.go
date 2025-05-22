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

// ClientCodec implements conjure-go-runtime's Codec interface.
// It ignores unknown struct members.
var ClientCodec codecClient

// ServerCodec implements conjure-go-runtime's Codec interface.
// It rejects unknown struct members.
var ServerCodec codecServer

type codecClient struct{ codecBase }

func (codecClient) Decode(r io.Reader, v any) error {
	if unmarshaler, ok := v.(json.UnmarshalerFrom); ok {
		return unmarshaler.UnmarshalJSONFrom(jsontext.NewDecoder(r))
	}
	return json.UnmarshalRead(r, *&v)
}

func (codecClient) Unmarshal(data []byte, v any) error {
	if unmarshaler, ok := v.(json.UnmarshalerFrom); ok {
		return unmarshaler.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewBuffer(data)))
	}
	return json.Unmarshal(data, *&v)
}

type codecServer struct{ codecBase }

func (codecServer) Decode(r io.Reader, v any) error {
	if unmarshaler, ok := v.(json.UnmarshalerFrom); ok {
		return unmarshaler.UnmarshalJSONFrom(jsontext.NewDecoder(r, json.RejectUnknownMembers(true)))
	}
	return json.UnmarshalRead(r, *&v, json.RejectUnknownMembers(true))
}

func (codecServer) Unmarshal(data []byte, v any) error {
	if unmarshaler, ok := v.(json.UnmarshalerFrom); ok {
		return unmarshaler.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewBuffer(data), json.RejectUnknownMembers(true)))
	}
	return json.Unmarshal(data, *&v, json.RejectUnknownMembers(true))
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

func ReadJSONFrom[T json.UnmarshalerFrom](r io.Reader, v T, opts ...json.Options) error {
	return v.UnmarshalJSONFrom(jsontext.NewDecoder(r, opts...))
}

func WriteJSONTo[T json.MarshalerTo](w io.Writer, v T, opts ...json.Options) error {
	return v.MarshalJSONTo(jsontext.NewEncoder(w, opts...))
}
