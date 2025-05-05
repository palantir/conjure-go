package types

import (
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// Boolean provides JSON marshaling and unmarshaling for Go bool-like types.
// Encodes values as JSON true/false, and decodes JSON booleans into the underlying type.
type Boolean[T ~bool] struct{}

func (Boolean[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteToken(jsontext.True)
	}
	return enc.WriteToken(jsontext.False)
}

func (Boolean[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch tok.Kind() {
	case 't':
		*receiver = true
	case 'f':
		*receiver = false
	default:
		return cj.NewKindMismatchError(dec, tok.Kind(), "json boolean")
	}
	return nil
}

// BooleanMapKey provides JSON marshaling for bool-like types used as map keys.
// Encodes bool keys as the JSON strings "true" or "false" (not as booleans), to comply with JSON map key requirements.
type BooleanMapKey[T ~bool] struct{}

func (BooleanMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	if receiver {
		return enc.WriteValue(jsontext.Value(`"true"`))
	}
	return enc.WriteValue(jsontext.Value(`"false"`))
}

func (BooleanMapKey[T]) Compare(a, b T) int {
	if a == b {
		return 0
	}
	if a {
		return 1
	}
	return -1
}

func (BooleanMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	b, err := strconv.ParseBool(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid boolean", err)
	}
	*receiver = T(b)
	return nil
}
