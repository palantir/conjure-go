package cj

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/rid"
)

// MarshalRIDTo uses a TypeRID to encode the value as JSON.
func MarshalRIDTo[T rid.ResourceIdentifier](value T, enc *jsontext.Encoder) error {
	return TypeRID[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalRIDFrom uses a TypeRID to decode the value from JSON.
func UnmarshalRIDFrom[T rid.ResourceIdentifier](value *T, dec *jsontext.Decoder) error {
	return TypeRID[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeRID[T rid.ResourceIdentifier] struct{}

func (TypeRID[T]) isEncoder() {}
func (TypeRID[T]) isDecoder() {}

func (TypeRID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(rid.ResourceIdentifier(receiver).String()))
}

func (TypeRID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	parsed, err := rid.ParseRID(tok.String())
	if err != nil {
		return NewSyntaxError(dec, "invalid resource identifier", err)
	}
	*receiver = T(parsed)
	return nil
}
