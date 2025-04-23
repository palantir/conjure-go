package cj

import (
	"cmp"
	"maps"
	"slices"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalOrderedMapTo uses a TypeOrderedMap to encode the value as JSON.
func MarshalOrderedMapTo[K cmp.Ordered, V any, KEY TypeEncoder[K], VAL TypeEncoder[V]](value map[K]V, enc *jsontext.Encoder) error {
	return TypeOrderedMapMarshaler[K, V, KEY, VAL]{}.MarshalJSONTo(value, enc)
}

// MarshalSortedMapTo uses a TypeSortedMap to encode the value as JSON.
func MarshalSortedMapTo[K comparable, V any, KEY MapKeyEncoder[K], VAL TypeEncoder[V]](value map[K]V, enc *jsontext.Encoder) error {
	return TypeSortedMapMarshaler[K, V, KEY, VAL]{}.MarshalJSONTo(value, enc)
}

// UnmarshalMapFrom uses a TypeSortedMap to decode the value from JSON.
func UnmarshalMapFrom[K comparable, V any, KEY TypeDecoder[K], VAL TypeDecoder[V]](value *map[K]V, dec *jsontext.Decoder) error {
	return TypeMapUnmarshaler[K, V, KEY, VAL]{}.UnmarshalJSONFrom(value, dec)
}

type TypeOrderedMapMarshaler[K cmp.Ordered, V any, KEY TypeEncoder[K], VAL TypeEncoder[V]] struct{}

func (TypeOrderedMapMarshaler[K, V, KEY, VAL]) isEncoder() {}

func (TypeOrderedMapMarshaler[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for _, k := range slices.Sorted(maps.Keys(receiver)) {
		err := (*new(KEY)).MarshalJSONTo(k, enc)
		if err != nil {
			return err
		}
		err = (*new(VAL)).MarshalJSONTo(receiver[k], enc)
		if err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

type TypeSortedMapMarshaler[K comparable, V any, KEY MapKeyEncoder[K], VAL TypeEncoder[V]] struct{}

func (TypeSortedMapMarshaler[K, V, KEY, VAL]) isEncoder() {}

func (TypeSortedMapMarshaler[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	keys := make([]K, 0, len(receiver))
	for k := range receiver {
		keys = append(keys, k)
	}

	for _, k := range slices.SortedFunc(maps.Keys(receiver), (*new(KEY)).Compare) {
		if err := (*new(KEY)).MarshalJSONTo(k, enc); err != nil {
			return err
		}
		if err := (*new(VAL)).MarshalJSONTo(receiver[k], enc); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

type TypeMapUnmarshaler[K comparable, V any, KEY TypeDecoder[K], VAL TypeDecoder[V]] struct{}

func (TypeMapUnmarshaler[K, V, KEY, VAL]) isDecoder() {}

func (TypeMapUnmarshaler[K, V, KEY, VAL]) UnmarshalJSONFrom(receiver *map[K]V, dec *jsontext.Decoder) error {
	panic("not implemented")
}
