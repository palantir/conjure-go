package cj

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/uuid"
)

// MarshalUUIDTo uses a TypeUUID to encode the value as JSON.
func MarshalUUIDTo[T ~[16]byte](value T, enc *jsontext.Encoder) error {
	return TypeUUID[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalUUIDFrom uses a TypeUUID to decode the value from JSON.
func UnmarshalUUIDFrom[T ~[16]byte](value *T, dec *jsontext.Decoder) error {
	return TypeUUID[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeUUID[T ~[16]byte] struct{}

func (t TypeUUID[T]) Compare(a, b T) int {
	// UUIDs are 16 bytes, so we can compare them directly
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}
	return 0
}

func (TypeUUID[T]) isEncoder() {}
func (TypeUUID[T]) isDecoder() {}

func (TypeUUID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(uuid.UUID(receiver).String()))
}

func (TypeUUID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	parsed, err := uuid.ParseUUID(tok.String())
	if err != nil {
		return NewSyntaxError(dec, "invalid UUID", err)
	}
	*receiver = T(parsed)
	return nil
}
