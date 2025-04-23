package cj

import (
	"bytes"
	"encoding/base64"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalBinaryTo uses a TypeBinary to encode the value as JSON.
func MarshalBinaryTo[T ~[]byte](value T, enc *jsontext.Encoder) error {
	return TypeBinary[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalBinaryFrom uses a TypeBinary to decode the value from JSON.
func UnmarshalBinaryFrom[T ~[]byte](value *T, dec *jsontext.Decoder) error {
	return TypeBinary[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeBinary[T ~[]byte] struct{}

func (TypeBinary[T]) isEncoder() {}
func (TypeBinary[T]) isDecoder() {}

func (TypeBinary[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	b64out := enc.UnusedBuffer()
	b64out = append(b64out, '"')
	b64out = base64.StdEncoding.AppendEncode(b64out, receiver)
	b64out = append(b64out, '"')
	return enc.WriteValue(b64out)
}

func (TypeBinary[T]) Compare(a, b T) int {
	return bytes.Compare(a, b)
}

func (TypeBinary[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	val, err := readValueOfKind(dec, '"')
	if err != nil {
		return err
	}
	if val[len(val)-1] != '"' {
		return NewSyntaxError(dec, "expected closing quote", nil)
	}
	*receiver, err = base64.StdEncoding.AppendDecode(*receiver, val[1:len(val)-1])
	if err != nil {
		return NewSyntaxError(dec, "invalid base64", err)
	}
	return nil
}
