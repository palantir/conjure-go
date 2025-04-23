package cj

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalIntMapKeyTo uses a TypeIntMapKey to encode the value as JSON.
func MarshalIntMapKeyTo[T ~int | ~int8 | ~int16 | ~int32 | ~int64](value T, enc *jsontext.Encoder) error {
	return TypeIntMapKey[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalIntMapKeyFrom uses a TypeIntMapKey to decode the value from JSON.
func UnmarshalIntMapKeyFrom[T ~int | ~int8 | ~int16 | ~int32 | ~int64](value *T, dec *jsontext.Decoder) error {
	return TypeIntMapKey[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeIntMapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (TypeIntMapKey[T]) isEncoder() {}
func (TypeIntMapKey[T]) isDecoder() {}

func (TypeIntMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = strconv.AppendInt(out, int64(receiver), 10)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (TypeIntMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(tok.String(), 10, 64)
	if err != nil {
		return NewSyntaxError(dec, "invalid int", err)
	}
	*receiver = T(i)
	return nil
}
