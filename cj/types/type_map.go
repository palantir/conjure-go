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
	"cmp"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// OrderedMapMarshaler provides JSON marshaling for maps with ordered keys (strings and numbers).
// Encodes maps as JSON objects, sorting keys using Go's cmp.Ordered rules, and delegates encoding of keys and values.
//
// Disable sorting with json.Deterministic(false).
// Format nil maps as 'null' with json.FormatNilMapAsNull(true).
type OrderedMapMarshaler[T ~map[K]V, K cmp.Ordered, V any, KEY cj.TypeEncoder[K], VAL cj.TypeEncoder[V]] struct{}

func (OrderedMapMarshaler[T, K, V, KEY, VAL]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	var keys iter.Seq[K]
	if deterministic, isSet := json.GetOption(enc.Options(), json.Deterministic); deterministic || !isSet {
		keys = slices.Values(slices.Sorted(maps.Keys(receiver)))
	} else {
		keys = maps.Keys(receiver)
	}
	return marshalMapWithKeySequence[T, K, V, KEY, VAL](receiver, keys, enc)
}

// ComparableMapMarshaler provides JSON marshaling for maps using a custom key comparison function.
// Encodes maps as JSON objects, sorting keys using cj.MapKeyEncoder's Compare method from KEY,
// and delegates encoding of keys and values.
//
// Types compatible with OrderedMapMarshaler should likely use that unless non-standard sorting is required.
//
// Disable sorting with json.Deterministic(false).
// Format nil maps as 'null' with json.FormatNilMapAsNull(true).
type ComparableMapMarshaler[T ~map[K]V, K comparable, V any, KEY cj.MapKeyEncoder[K], VAL cj.TypeEncoder[V]] struct{}

func (ComparableMapMarshaler[T, K, V, KEY, VAL]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	var keys iter.Seq[K]
	if deterministic, isSet := json.GetOption(enc.Options(), json.Deterministic); deterministic || !isSet {
		keys = slices.Values(slices.SortedFunc(maps.Keys(receiver), (*new(KEY)).Compare))
	} else {
		keys = maps.Keys(receiver)
	}
	return marshalMapWithKeySequence[T, K, V, KEY, VAL](receiver, keys, enc)
}

func marshalMapWithKeySequence[T ~map[K]V, K comparable, V any, KEY cj.TypeEncoder[K], VAL cj.TypeEncoder[V]](receiver T, keys iter.Seq[K], enc *jsontext.Encoder) error {
	if receiver == nil {
		if formatNilMapAsNull, isSet := json.GetOption(enc.Options(), json.FormatNilMapAsNull); formatNilMapAsNull && isSet {
			return enc.WriteToken(jsontext.Null)
		}
	}
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for k := range keys {
		if err := (*new(KEY)).MarshalJSONTo(k, enc); err != nil {
			return err
		}
		if err := (*new(VAL)).MarshalJSONTo(receiver[k], enc); err != nil {
			return err
		}
	}
	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}
	return nil
}

// MapUnmarshaler provides JSON unmarshaling for maps, using nested KEY and VAL decoders for keys and values.
// Decodes JSON objects into Go maps of the specified types.
type MapUnmarshaler[T ~map[K]V, K comparable, V any, KEY cj.TypeDecoder[K], VAL cj.TypeDecoder[V]] struct{}

func (MapUnmarshaler[T, K, V, KEY, VAL]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if *receiver == nil {
		*receiver = make(T)
	} else {
		clear(*receiver)
	}
	switch tok.Kind() {
	case '{':
		for {
			if dec.PeekKind() == '}' {
				_, err := dec.ReadToken()
				return err
			}
			key := *new(K)
			if err := (*new(KEY)).UnmarshalJSONFrom(&key, dec); err != nil {
				return err
			}
			if _, ok := (*receiver)[key]; ok {
				return cj.NewDuplicateMapKeyError(dec, fmt.Sprintf("%T", receiver))
			}
			val := *new(V)
			if err := (*new(VAL)).UnmarshalJSONFrom(&val, dec); err != nil {
				return err
			}
			(*receiver)[key] = val
		}
	case 'n':
		return nil
	default:
		return cj.NewKindMismatchError(dec, tok.Kind(), "map opening brace")
	}
}
