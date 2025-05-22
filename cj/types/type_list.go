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

// ListMarshaler provides JSON marshaling for slices of type T using a nested encoder ITEM.
// Encodes slices as JSON arrays, delegating encoding of each element to ITEM.
type ListMarshaler[T ~[]U, U any, ITEM cj.TypeEncoder[U]] struct{}

func (ListMarshaler[T, U, ITEM]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver == nil {
		if formatNilSliceAsNull, isSet := json.GetOption(enc.Options(), json.FormatNilSliceAsNull); formatNilSliceAsNull && isSet {
			return enc.WriteToken(jsontext.Null)
		}
	}
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	for _, item := range receiver {
		if err := (*new(ITEM)).MarshalJSONTo(enc, item); err != nil {
			return err
		}
	}
	if err := enc.WriteToken(jsontext.EndArray); err != nil {
		return err
	}
	return nil
}

// ListUnmarshaler provides JSON unmarshaling for slices of type T using a nested decoder ITEM.
// Decodes JSON arrays into Go slices, delegating decoding of each element to ITEM.
type ListUnmarshaler[T ~[]U, U any, ITEM cj.TypeDecoder[U]] struct{}

func (ListUnmarshaler[T, U, ITEM]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if *receiver == nil {
		*receiver = make(T, 0)
	} else {
		*receiver = (*receiver)[:0]
	}
	if kind := tok.Kind(); kind != '[' {
		if kind == 'n' {
			// null
			*receiver = make(T, 0)
			return nil
		}
		return cj.NewKindMismatchError(dec, kind, "list opening bracket")
	}
	for {
		if dec.PeekKind() == ']' {
			_, err := dec.ReadToken()
			return err
		}
		item := *new(U)
		if err := (*new(ITEM)).UnmarshalJSONFrom(dec, &item); err != nil {
			return err
		}
		*receiver = append(*receiver, item)
	}
}
