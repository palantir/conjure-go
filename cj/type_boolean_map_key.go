package cj

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalBooleanMapKeyTo uses a TypeBooleanMapKey to encode the value as JSON.
func MarshalBooleanMapKeyTo[T ~bool](value T, enc *jsontext.Encoder) error {
	return TypeBooleanMapKey[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalBooleanMapKeyFrom uses a TypeBooleanMapKey to decode the value from JSON.
func UnmarshalBooleanMapKeyFrom[T ~bool](value *T, dec *jsontext.Decoder) error {
	return TypeBooleanMapKey[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeBooleanMapKey[T ~bool] struct{}

func (TypeBooleanMapKey[T]) isEncoder() {}
func (TypeBooleanMapKey[T]) isDecoder() {}

func (TypeBooleanMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteValue(jsontext.Value(`"true"`))
	}
	return enc.WriteValue(jsontext.Value(`"false"`))
}

func (TypeBooleanMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	b, err := strconv.ParseBool(tok.String())
	if err != nil {
		return NewSyntaxError(dec, "invalid boolean", err)
	}
	*receiver = T(b)
	return nil
}
