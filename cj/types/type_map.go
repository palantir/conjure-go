package types

import (
	"cmp"
	"maps"
	"slices"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// OrderedMapMarshaler provides JSON marshaling for maps with ordered keys.
// Encodes maps as JSON objects, sorting keys using Go's cmp.Ordered rules, and delegates encoding of keys and values.
type OrderedMapMarshaler[K cmp.Ordered, V any, KEY cj.TypeEncoder[K], VAL cj.TypeEncoder[V]] struct{}

func (OrderedMapMarshaler[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
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
	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}
	return nil
}

// SortedMapMarshaler provides JSON marshaling for maps using a custom key comparison function.
// Encodes maps as JSON objects, sorting keys using cj.MapKeyEncoder's Compare method from KEY,
// and delegates encoding of keys and values.
//
// Types compatible with OrderedMapMarshaler should likely use that unless non-standard sorting is required.
type SortedMapMarshaler[K comparable, V any, KEY cj.MapKeyEncoder[K], VAL cj.TypeEncoder[V]] struct{}

func (SortedMapMarshaler[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for _, k := range slices.SortedFunc(maps.Keys(receiver), (*new(KEY)).Compare) {
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

// UnsortedMapMarshaler provides JSON marshaling for maps without key ordering guarantees.
// Encodes maps as JSON objects, iterating in Go map order. Delegates encoding of keys and values.
// Output order is not stable across runs; use only when order does not matter and performance
// very sensitive to the sorting cost.
type UnsortedMapMarshaler[K comparable, V any, KEY cj.TypeEncoder[K], VAL cj.TypeEncoder[V]] struct{}

func (UnsortedMapMarshaler[K, V, KEY, VAL]) MarshalJSONTo(receiver map[K]V, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for k, v := range receiver {
		if err := (*new(KEY)).MarshalJSONTo(k, enc); err != nil {
			return err
		}
		if err := (*new(VAL)).MarshalJSONTo(v, enc); err != nil {
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
type MapUnmarshaler[K comparable, V any, KEY cj.TypeDecoder[K], VAL cj.TypeDecoder[V]] struct{}

func (MapUnmarshaler[K, V, KEY, VAL]) UnmarshalJSONFrom(receiver *map[K]V, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if *receiver == nil {
		*receiver = make(map[K]V)
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
