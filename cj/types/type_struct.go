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
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// StructMarshaler provides JSON marshaling for types that implement json.MarshalerTo.
// Delegates marshaling to the type's MarshalJSONTo method.
type StructMarshaler[T json.MarshalerTo] struct{}

func (StructMarshaler[T]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	return receiver.MarshalJSONTo(enc)
}

// StructUnmarshaler provides JSON unmarshaling for types that implement json.UnmarshalerFrom.
// Delegates unmarshaling to the type's UnmarshalJSONFrom method.
type StructUnmarshaler[T json.UnmarshalerFrom] struct{}

func (StructUnmarshaler[T]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver T) error {
	return receiver.UnmarshalJSONFrom(dec)
}

// VisitObjectFields is a helper for use in UnmarshalJSONFrom implementations that reads the opening and closing braces
// and calls visitField for each key and value in the object.
func VisitObjectFields(dec *jsontext.Decoder, visitField func(key string, dec *jsontext.Decoder) error) error {
	if tok, err := dec.ReadToken(); err != nil {
		return err
	} else if kind := tok.Kind(); kind != '{' {
		return cj.NewKindMismatchError(dec, kind, "object opening brace")
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
			return cj.NewKindMismatchError(dec, kind, "object closing brace or next key")
		}
		if err := visitField(key.String(), dec); err != nil {
			return err
		}
	}
}
