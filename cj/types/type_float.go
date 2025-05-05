package types

import (
	"math"
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// Float provides JSON marshaling and unmarshaling for float64-like types.
// Encodes values as JSON numbers, and decodes JSON numbers into the underlying type.
// Special values like "NaN", "Infinity", and "-Infinity" are handled by jsontext.Float.
type Float[T ~float64] struct{}

func (Float[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.Float(float64(receiver)))
}

func (Float[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch kind := tok.Kind(); kind {
	case '0':
		*receiver = T(tok.Float())
	case '"':
		switch tok.String() {
		case "NaN":
			*receiver = T(math.NaN())
		case "Infinity":
			*receiver = T(math.Inf(+1))
		case "-Infinity":
			*receiver = T(math.Inf(-1))
		default:
			return cj.NewKindMismatchError(dec, kind, "json float")
		}
	default:
		return cj.NewKindMismatchError(dec, kind, "json float")
	}
	return nil
}

// FloatMapKey provides JSON marshaling and unmarshaling for float64-like types used as map keys.
// Encodes float keys as JSON strings to comply with JSON map key requirements, supporting
// special values like "NaN", "Infinity", and "-Infinity".
type FloatMapKey[T ~float64] struct{}

func (FloatMapKey[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	switch {
	case math.IsNaN(float64(receiver)):
		return enc.WriteToken(jsontext.String("NaN"))
	case math.IsInf(float64(receiver), +1):
		return enc.WriteToken(jsontext.String("Infinity"))
	case math.IsInf(float64(receiver), -1):
		return enc.WriteToken(jsontext.String("-Infinity"))
	default:
		return enc.WriteToken(jsontext.String(strconv.FormatFloat(float64(receiver), 'f', -1, 64)))
	}
}

func (FloatMapKey[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	switch s := tok.String(); s {
	case "NaN":
		*receiver = T(math.NaN())
	case "Infinity":
		*receiver = T(math.Inf(1))
	case "-Infinity":
		*receiver = T(math.Inf(-1))
	default:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return cj.WrapSyntaxError(dec, "invalid float", err)
		}
		*receiver = T(f)
	}
	return nil
}
