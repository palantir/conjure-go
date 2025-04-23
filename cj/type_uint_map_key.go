package cj

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalUintMapKeyTo uses a TypeUintMapKey to encode the value as JSON.
func MarshalUintMapKeyTo[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](value T, enc *jsontext.Encoder) error {
	return TypeUintMapKey[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalUintMapKeyFrom uses a TypeUintMapKey to decode the value from JSON.
func UnmarshalUintMapKeyFrom[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](value *T, dec *jsontext.Decoder) error {
	return TypeUintMapKey[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeUintMapKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (TypeUintMapKey[T]) isEncoder() {}
func (TypeUintMapKey[T]) isDecoder() {}

func (TypeUintMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = strconv.AppendUint(out, uint64(receiver), 10)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (TypeUintMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(tok.String(), 10, 64)
	if err != nil {
		return NewSyntaxError(dec, "invalid uint", err)
	}
	*receiver = T(i)
	return nil
}
