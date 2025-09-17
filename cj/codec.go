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

// Marshal is a variant of json.Marshal that instantiates a new TypeEncoder and uses its MarshalJSONTo method.
func Marshal[T any, E TypeEncoder[T]](receiver T, opts ...json.Options) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := MarshalWrite[T, E](buf, receiver, opts...)
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), err
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

// ClientDecoder implements conjure-go-runtime's codecs.Decoder interface.
// It ignores unknown struct members.
type ClientDecoder[T any, D TypeDecoder[T]] struct{}

func (ClientDecoder[T, D]) Decode(r io.Reader, v any) error {
	if v == nil {
		return werror.Error("decode target cannot be nil")
	}
	switch vt := v.(type) {
	case *T:
		return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r), vt)
	case **T:
		*vt = new(T)
		return (*new(D)).UnmarshalJSONFrom(jsontext.NewDecoder(r), *vt)
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
		return (*new(E)).MarshalJSONTo(jsontext.NewEncoder(w), vt)
	case *T:
		if vt == nil {
			return werror.Error("encode source should not be nil")
		}
		return (*new(E)).MarshalJSONTo(jsontext.NewEncoder(w), *vt)
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

// VisitObjectFields is a helper for use in UnmarshalJSONFrom implementations that reads the opening and closing braces
// and calls visitField for each key and value in the object.
func VisitObjectFields(dec *jsontext.Decoder, visitField func(key string, dec *jsontext.Decoder) error) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '{' {
		return NewKindMismatchError(dec, kind, "object opening brace")
	}
	for {
		key, err := dec.ReadToken()
		if err != nil {
			return err
		}
		kind := key.Kind()
		if kind == '}' {
			return nil // End of object
		}
		if kind != '"' {
			return NewKindMismatchError(dec, kind, "object closing brace or next key")
		}
		if err := visitField(key.String(), dec); err != nil {
			return err
		}
	}
}

func getOptionOrTrue(options jsontext.Options, setter func(bool) jsontext.Options) bool {
	value, ok := json.GetOption(options, setter)
	if !ok {
		return true
	}
	return value
}

func getOptionOrFalse(options jsontext.Options, setter func(bool) jsontext.Options) bool {
	value, ok := json.GetOption(options, setter)
	if !ok {
		return false
	}
	return value
}
