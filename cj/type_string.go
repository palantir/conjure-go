package cj

import (
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalStringTo uses a TypeString to encode the value as JSON.
func MarshalStringTo[T ~string](value T, enc *jsontext.Encoder) error {
	return TypeString[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalStringFrom uses a TypeString to decode the value from JSON.
func UnmarshalStringFrom[T ~string](value *T, dec *jsontext.Decoder) error {
	return TypeString[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeString[T ~string] struct{}

func (TypeString[T]) isEncoder() {}
func (TypeString[T]) isDecoder() {}

func (TypeString[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(string(receiver)))
}

func (TypeString[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	*receiver = T(tok.String())
	return nil
}
