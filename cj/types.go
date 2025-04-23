package cj

import (
	"github.com/go-json-experiment/json/jsontext"
)

type TypeArshaler[T any] interface {
	TypeEncoder[T]
	TypeDecoder[T]
}

type TypeEncoder[T any] interface {
	MarshalJSONTo(receiver T, enc *jsontext.Encoder) error

	isEncoder() // sealed interface
}

type TypeDecoder[T any] interface {
	UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error

	isDecoder() // sealed interface
}

type MapKeyEncoder[K comparable] interface {
	TypeEncoder[K]
	MapKeyComparator[K]
}

type MapKeyComparator[K comparable] interface {
	// Compare returns -1 if a < b, 0 if a == b, and 1 if a > b.
	// This is used to sort keys in a deterministic order.
	Compare(K, K) int
}
