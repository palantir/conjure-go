package types

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// OptionalMarshaler provides JSON marshaling for optional (pointer) values of type T.
// Encodes nil pointers as JSON null, otherwise delegates encoding to ITEM.
type OptionalMarshaler[T any, ITEM cj.TypeEncoder[T]] struct{}

func (OptionalMarshaler[T, ITEM]) MarshalJSONTo(receiver *T, enc *jsontext.Encoder) error {
	if receiver == nil {
		return enc.WriteToken(jsontext.Null)
	}
	return (*new(ITEM)).MarshalJSONTo(*receiver, enc)
}

// OptionalUnmarshaler provides JSON unmarshaling for optional (pointer) values of type T.
// Decodes JSON null as nil, otherwise delegates decoding to ITEM.
type OptionalUnmarshaler[T any, ITEM cj.TypeDecoder[T]] struct{}

func (OptionalUnmarshaler[T, ITEM]) UnmarshalJSONFrom(receiver **T, dec *jsontext.Decoder) error {
	if dec.PeekKind() == 'n' {
		// still have to consume token
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*receiver = nil
		return nil
	}
	return (*new(ITEM)).UnmarshalJSONFrom(*receiver, dec)
}
