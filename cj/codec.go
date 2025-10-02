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

// Marshaling Utils //

// MarshalEncode is a variant of json.MarshalEncode that instantiates a new TypeEncoder and uses its MarshalJSONTo method.
func MarshalEncode[T any, E TypeEncoder[T]](enc *jsontext.Encoder, receiver T) error {
	return (*new(E)).MarshalJSONTo(enc, receiver)
}

// Unmarshaling Utils //

// UnmarshalDecode is a variant of json.UnmarshalDecode that instantiates a new TypeDecoder and uses its UnmarshalJSONFrom method.
func UnmarshalDecode[T any, D TypeDecoder[T]](dec *jsontext.Decoder, receiver *T) error {
	return (*new(D)).UnmarshalJSONFrom(dec, receiver)
}

// ClientDecoder implements conjure-go-runtime's codecs.Decoder interface.
// It ignores unknown struct members.
type ClientDecoder[T any, D TypeDecoder[T]] struct{}

func (ClientDecoder[T, D]) Decode(r io.Reader, v any) error {
	if v == nil {
		return werror.Error("decode target cannot be nil")
	}
	switch vt := v.(type) {
	case *T:
		return json.UnmarshalRead(r, NewUnmarshalerFrom[T, D](vt))
	case **T:
		*vt = new(T)
		return json.UnmarshalRead(r, NewUnmarshalerFrom[T, D](*vt))
	default:
		return werror.Error("decode target is incompatible with decoder")
	}
}

func (ClientDecoder[T, D]) Unmarshal(data []byte, v any) error {
	return ClientDecoder[T, D]{}.Decode(bytes.NewBuffer(data), v)
}

func (ClientDecoder[T, D]) Accept() string {
	return "application/json"
}

// ClientEncoder implements conjure-go-runtime's codecs.Encoder interface.
type ClientEncoder[T any, E TypeEncoder[T]] struct{}

func (ClientEncoder[T, E]) Encode(w io.Writer, v any) error {
	switch vt := v.(type) {
	case T:
		return json.MarshalWrite(w, NewMarshalerTo[T, E](vt))
	case *T:
		if vt == nil {
			return werror.Error("encode source should not be nil")
		}
		return json.MarshalWrite(w, NewMarshalerTo[T, E](*vt))
	default:
		return werror.Error("encode source is incompatible with encoder")
	}
}

func (ClientEncoder[T, E]) Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := ClientEncoder[T, E]{}.Encode(buf, v)
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), err
}

func (ClientEncoder[T, E]) ContentType() string {
	return "application/json"
}
