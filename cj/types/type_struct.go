package types

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// StructMarshaler provides JSON marshaling for types that implement json.MarshalerTo.
// Delegates marshaling to the type's MarshalJSONTo method.
type StructMarshaler[T json.MarshalerTo] struct{}

func (StructMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return receiver.MarshalJSONTo(enc)
}

// StructUnmarshaler provides JSON unmarshaling for types that implement json.UnmarshalerFrom.
// Delegates unmarshaling to the type's UnmarshalJSONFrom method.
type StructUnmarshaler[T json.UnmarshalerFrom] struct{}

func (StructUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	return receiver.UnmarshalJSONFrom(dec)
}
