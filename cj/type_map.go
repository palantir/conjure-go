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
	"cmp"
	"maps"
	"reflect"
	"slices"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// OrderedMapMarshaler provides JSON marshaling for maps with ordered keys (strings and numbers).
// Encodes maps as JSON objects, sorting keys using Go's cmp.Ordered rules, and delegates encoding of keys and values.
//
// Disable sorting with json.Deterministic(false).
// Format nil maps as 'null' with json.FormatNilMapAsNull(true).
type OrderedMapMarshaler[T ~map[K]V, K cmp.Ordered, V any, KEY TypeEncoder[K], VAL TypeEncoder[V]] struct{}

func (OrderedMapMarshaler[T, K, V, KEY, VAL]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver == nil {
		if formatNilMapAsNull, isSet := json.GetOption(enc.Options(), json.FormatNilMapAsNull); formatNilMapAsNull && isSet {
			return enc.WriteToken(jsontext.Null)
		}
	}
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if deterministic, isSet := json.GetOption(enc.Options(), json.Deterministic); len(receiver) > 1 && (!isSet || deterministic) {
		for _, k := range slices.Sorted(maps.Keys(receiver)) {
			if err := (*new(KEY)).MarshalJSONTo(enc, k); err != nil {
				return err
			}
			if err := (*new(VAL)).MarshalJSONTo(enc, receiver[k]); err != nil {
				return err
			}
		}
	} else {
		for k, v := range receiver {
			if err := (*new(KEY)).MarshalJSONTo(enc, k); err != nil {
				return err
			}
			if err := (*new(VAL)).MarshalJSONTo(enc, v); err != nil {
				return err
			}
		}
	}
	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}
	return nil
}

// ComparableMapMarshaler provides JSON marshaling for maps using a custom key comparison function.
// Encodes maps as JSON objects, sorting keys using cj.MapKeyEncoder's Compare method from KEY,
// and delegates encoding of keys and values.
//
// Types compatible with OrderedMapMarshaler should likely use that unless non-standard sorting is required.
//
// Disable sorting with json.Deterministic(false).
// Format nil maps as 'null' with json.FormatNilMapAsNull(true).
type ComparableMapMarshaler[T ~map[K]V, K comparable, V any, KEY MapKeyEncoder[K], VAL TypeEncoder[V]] struct{}

func (ComparableMapMarshaler[T, K, V, KEY, VAL]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver == nil {
		if formatNilMapAsNull, isSet := json.GetOption(enc.Options(), json.FormatNilMapAsNull); formatNilMapAsNull && isSet {
			return enc.WriteToken(jsontext.Null)
		}
	}
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if deterministic, isSet := json.GetOption(enc.Options(), json.Deterministic); len(receiver) > 1 && (!isSet || deterministic) {
		for _, k := range slices.SortedFunc(maps.Keys(receiver), (*new(KEY)).Compare) {
			if err := (*new(KEY)).MarshalJSONTo(enc, k); err != nil {
				return err
			}
			if err := (*new(VAL)).MarshalJSONTo(enc, receiver[k]); err != nil {
				return err
			}
		}
	} else {
		for k, v := range receiver {
			if err := (*new(KEY)).MarshalJSONTo(enc, k); err != nil {
				return err
			}
			if err := (*new(VAL)).MarshalJSONTo(enc, v); err != nil {
				return err
			}
		}
	}
	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}
	return nil
}

// MapUnmarshaler provides JSON unmarshaling for maps, using nested KEY and VAL decoders for keys and values.
// Decodes JSON objects into Go maps of the specified types.
type MapUnmarshaler[T ~map[K]V, K comparable, V any, KEY TypeDecoder[K], VAL TypeDecoder[V]] struct{}

func (MapUnmarshaler[T, K, V, KEY, VAL]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '{' {
		if kind == 'n' {
			// null
			*receiver = make(T)
			return nil
		}
		return NewKindMismatchError(dec, kind, "map opening brace")
	}
	if *receiver == nil {
		*receiver = make(T)
	} else {
		clear(*receiver)
	}
	for {
		if dec.PeekKind() == '}' {
			_, err := dec.ReadToken()
			return err
		}
		key := *new(K)
		if err := (*new(KEY)).UnmarshalJSONFrom(dec, &key); err != nil {
			return err
		}
		val := *new(V)
		if err := (*new(VAL)).UnmarshalJSONFrom(dec, &val); err != nil {
			return err
		}
		if _, ok := (*receiver)[key]; ok {
			return NewDuplicateMapKeyError(dec, reflect.TypeOf(receiver).String())
		}
		(*receiver)[key] = val
	}
}
