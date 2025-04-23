package cj

import (
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalBooleanTo uses a TypeBoolean to encode the value as JSON.
func MarshalBooleanTo[T ~bool](value T, enc *jsontext.Encoder) error {
	return TypeBoolean[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalBooleanFrom uses a TypeBoolean to decode the value from JSON.
func UnmarshalBooleanFrom[T ~bool](value *T, dec *jsontext.Decoder) error {
	return TypeBoolean[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeBoolean[T ~bool] struct{}

func (TypeBoolean[T]) isEncoder() {}
func (TypeBoolean[T]) isDecoder() {}

func (TypeBoolean[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteToken(jsontext.True)
	}
	return enc.WriteToken(jsontext.False)
}

func (TypeBoolean[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, 't', 'f')
	if err != nil {
		return err
	}
	if tok.Kind() == 't' {
		*receiver = true
	} else {
		*receiver = false
	}
	return nil
}
