package cj

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalUintTo uses a TypeUint to encode the value as JSON.
func MarshalUintTo[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](value T, enc *jsontext.Encoder) error {
	return TypeUint[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalUintFrom uses a TypeUint to decode the value from JSON.
func UnmarshalUintFrom[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](value *T, dec *jsontext.Decoder) error {
	return TypeUint[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}

func (TypeUint[T]) isEncoder() {}
func (TypeUint[T]) isDecoder() {}

func (TypeUint[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = strconv.AppendUint(out, uint64(receiver), 10)
	return enc.WriteValue(out)
}

func (TypeUint[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '0')
	if err != nil {
		return err
	}
	*receiver = T(tok.Uint())
	return nil
}
