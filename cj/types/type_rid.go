package types

import (
	"strings"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/rid"
)

// RID provides JSON marshaling and unmarshaling for types based on rid.ResourceIdentifier.
// Encodes values as JSON strings using the canonical string representation of the resource identifier.
// Implements comparison based on all RID fields for use as a map key.
type RID[T rid.ResourceIdentifier] struct{}

func (RID[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(rid.ResourceIdentifier(receiver).String()))
}

func (RID[T]) Compare(a, b T) int {
	ra, rb := rid.ResourceIdentifier(a), rid.ResourceIdentifier(b)
	if cmp := strings.Compare(ra.Service, rb.Service); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Instance, rb.Instance); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Type, rb.Type); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(ra.Locator, rb.Locator); cmp != 0 {
		return cmp
	}
	return 0
}

func (RID[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readStringToken(dec)
	if err != nil {
		return err
	}
	parsed, err := rid.ParseRID(tok.String())
	if err != nil {
		return cj.WrapSyntaxError(dec, "invalid resource identifier", err)
	}
	*receiver = T(parsed)
	return nil
}
