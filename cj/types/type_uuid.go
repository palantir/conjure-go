package types

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/uuid"
)

// UUID provides JSON marshaling and unmarshaling for uuid.UUID.
type UUID[T ~[16]byte] struct{}

func (UUID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(uuid.UUID(receiver).String()))
}

func (t UUID[T]) Compare(a, b T) int {
	// UUIDs are 16 bytes, so we can compare them directly
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}
	return 0
}

func (UUID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parsed, err := uuid.ParseUUID(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid UUID", err)
	}
	*receiver = T(parsed)
	return nil
}
