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
	"reflect"
	"slices"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// SetMarshaler provides json marshaling for sets of type T using a nested encoder ITEM.
// It writes a json array, delegating encoding of each element to ITEM.
//
// Duplicate items in the receiver slice (as determined by the == operator) are skipped.
// The emitted JSON list's elements will otherwise be in the same order as the original.
//
// If the receiver is nil and json.FormatNilSliceAsNull is true, the JSON null value is written.
type SetMarshaler[T ~[]U, U comparable, ITEM TypeEncoder[U]] struct{}

func (SetMarshaler[T, U, ITEM]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver == nil && getOptionOrFalse(enc.Options(), json.FormatNilSliceAsNull) {
		if err := enc.WriteToken(jsontext.Null); err != nil {
			return WrapEncodeError(enc, "", err)
		}
		return nil
	}
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return WrapEncodeError(enc, "", err)
	}
	for i, item := range receiver {
		if slices.Contains(receiver[0:i], item) {
			// duplicate item
			continue
		}
		if err := (*new(ITEM)).MarshalJSONTo(enc, item); err != nil {
			return err
		}
	}
	if err := enc.WriteToken(jsontext.EndArray); err != nil {
		return WrapEncodeError(enc, "", err)
	}
	return nil
}

// SetUnmarshaler provides json unmarshaling for sets of type T using a nested decoder ITEM.
// It reads a json array, delegating decoding of each element to ITEM.
//
// Duplicate items in the receiver slice (as determined by the == operator) result in DuplicateSetItemError.
//
// The receiver slice is allocated even if the input JSON value is 'null'.
type SetUnmarshaler[T ~[]U, U comparable, ITEM TypeDecoder[U]] struct{}

func (SetUnmarshaler[T, U, ITEM]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	if *receiver == nil {
		*receiver = make(T, 0)
	} else {
		*receiver = (*receiver)[:0]
	}
	kind := tok.Kind()
	if kind == 'n' {
		// null
		return nil
	}
	if kind != '[' {
		return NewKindMismatchError(dec, kind, "list opening bracket")
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
		if slices.Contains(*receiver, item) {
			return NewDuplicateSetItemError(dec, reflect.TypeOf(item).Name(), len(*receiver))
		}
		*receiver = append(*receiver, item)
	}
}
