package cj

import (
	"math"
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalFloatTo uses a TypeFloat to encode the value as JSON.
func MarshalFloatTo[T ~float64](value T, enc *jsontext.Encoder) error {
	return TypeFloat[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalFloatFrom uses a TypeFloat to decode the value from JSON.
func UnmarshalFloatFrom[T ~float64](value *T, dec *jsontext.Decoder) error {
	return TypeFloat[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeFloat[T ~float64] struct{}

func (TypeFloat[T]) isEncoder() {}
func (TypeFloat[T]) isDecoder() {}

func (TypeFloat[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	n := float64(receiver)
	switch {
	case math.Float64bits(n) == 0:
		out = append(out, '0')
	case math.IsNaN(n):
		out = append(out, "\"NaN\""...)
	case math.IsInf(n, +1):
		out = append(out, "\"Infinity\""...)
	case math.IsInf(n, -1):
		out = append(out, "\"-Infinity\""...)
	default:
		out = strconv.AppendFloat(out, n, 'f', -1, 64)
	}
	return enc.WriteValue(out)
}

func (TypeFloat[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Float())
	return nil
}
