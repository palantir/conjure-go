package cj

import (
	"math"
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalFloatMapKeyTo uses a TypeFloatMapKey to encode the value as JSON.
func MarshalFloatMapKeyTo[T ~float64](value T, enc *jsontext.Encoder) error {
	return TypeFloatMapKey[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalFloatMapKeyFrom uses a TypeFloatMapKey to decode the value from JSON.
func UnmarshalFloatMapKeyFrom[T ~float64](value *T, dec *jsontext.Decoder) error {
	return TypeFloatMapKey[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeFloatMapKey[T ~float64] struct{}

func (TypeFloatMapKey[T]) isEncoder() {}
func (TypeFloatMapKey[T]) isDecoder() {}

func (TypeFloatMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	n := float64(receiver)
	switch {
	case math.Float64bits(n) == 0:
		out = append(out, "\"0\""...)
	case math.IsNaN(n):
		out = append(out, "\"NaN\""...)
	case math.IsInf(n, +1):
		out = append(out, "\"Infinity\""...)
	case math.IsInf(n, -1):
		out = append(out, "\"-Infinity\""...)
	default:
		out = append(out, '"')
		out = strconv.AppendFloat(out, n, 'f', -1, 64)
		out = append(out, '"')
	}
	return enc.WriteValue(out)
}

func (TypeFloatMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(tok.String(), 64)
	if err != nil {
		return NewSyntaxError(dec, "invalid float", err)
	}
	*receiver = T(f)
	return nil
}
