package cj

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalAnyTo uses a TypeAny to encode the value as JSON.
func MarshalAnyTo[T any](value T, enc *jsontext.Encoder) error {
	return TypeAny[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalAnyFrom uses a TypeAny to decode the value from JSON.
func UnmarshalAnyFrom[T any](value *T, dec *jsontext.Decoder) error {
	return TypeAny[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeAny[T any] struct{}

func (TypeAny[T]) isEncoder() {}
func (TypeAny[T]) isDecoder() {}

func (TypeAny[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, receiver)
}

func (TypeAny[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, receiver)
}
