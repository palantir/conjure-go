package cj

import (
	"github.com/go-json-experiment/json/jsontext"
)

// TypeEncoder is implemented by types that can encode a Go value of type T to JSON using the provided jsontext.Encoder.
// Implementations for each Conjure type (e.g., Boolean, Integer, String, List, Map, etc.) are found in the types/ package.
// Each implementation ensures correct marshaling of the corresponding Go type to the appropriate JSON representation.
type TypeEncoder[T any] interface {
	// MarshalJSONTo writes the JSON encoding of 'receiver' to 'enc'.
	MarshalJSONTo(receiver T, enc *jsontext.Encoder) error
}

// TypeDecoder is implemented by types that can decode JSON into a Go value of type T using the provided jsontext.Decoder.
// Implementations in the types/ package handle unmarshaling for each supported Conjure type, including type validation and error handling.
type TypeDecoder[T any] interface {
	// UnmarshalJSONFrom reads JSON from 'dec' and stores the result in 'receiver'.
	UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error
}

// MapKeyEncoder is implemented by types that can be used as map keys in Conjure types
// but do not implement cmp.Ordered. The encoder's Compare method is used to sort map keys in a deterministic order.
// TypeEncoder implementations for comparable types (numbers, strings, etc) should not implement Compare.
type MapKeyEncoder[K comparable] interface {
	TypeEncoder[K]

	// Compare returns -1 if a < b, 0 if a == b, and 1 if a > b.
	// This is used to sort keys in a deterministic order.
	Compare(K, K) int
}
