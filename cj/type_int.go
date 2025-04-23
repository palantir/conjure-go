package cj

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalIntTo uses a TypeInt to encode the value as JSON.
func MarshalIntTo[T ~int | ~int8 | ~int16 | ~int32 | ~int64](value T, enc *jsontext.Encoder) error {
	return TypeInt[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalIntFrom uses a TypeInt to decode the value from JSON.
func UnmarshalIntFrom[T ~int | ~int8 | ~int16 | ~int32 | ~int64](value *T, dec *jsontext.Decoder) error {
	return TypeInt[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}

func (TypeInt[T]) isEncoder() {}
func (TypeInt[T]) isDecoder() {}

func (TypeInt[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = strconv.AppendInt(out, int64(receiver), 10)
	return enc.WriteValue(out)
}

func (TypeInt[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Int())
	return nil
}
