package types

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type StructMarshaler[T json.MarshalerTo] struct{}

func (StructMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return receiver.MarshalJSONTo(enc)
}

type StructUnmarshaler[T json.UnmarshalerFrom] struct{}

func (StructUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	return receiver.UnmarshalJSONFrom(dec)
}
