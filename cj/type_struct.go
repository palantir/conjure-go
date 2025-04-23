package cj

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalStructTo uses a TypeStruct to encode the value as JSON.
func MarshalStructTo[T json.MarshalerTo](value T, enc *jsontext.Encoder) error {
	return TypeStructMarshaler[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalStructFrom uses a TypeStruct to decode the value from JSON.
func UnmarshalStructFrom[T json.UnmarshalerFrom](value T, dec *jsontext.Decoder) error {
	return TypeStructUnmarshaler[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeStructMarshaler[T json.MarshalerTo] struct{}

func (TypeStructMarshaler[T]) isEncoder() {}

func (TypeStructMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return receiver.MarshalJSONTo(enc)
}

type TypeStructUnmarshaler[T json.UnmarshalerFrom] struct{}

func (TypeStructUnmarshaler[T]) isDecoder() {}

func (TypeStructUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	return receiver.UnmarshalJSONFrom(dec)
}
