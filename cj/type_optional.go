package cj

import (
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalOptionalTo uses a TypeOptional to encode the value as JSON.
func MarshalOptionalTo[T any, U *T, ITEM TypeEncoder[T]](value *T, enc *jsontext.Encoder) error {
	return TypeOptionalMarshaler[T, U, ITEM]{}.MarshalJSONTo(value, enc)
}

// UnmarshalOptionalFrom uses a TypeOptional to decode the value from JSON.
func UnmarshalOptionalFrom[T any, U *T, ITEM TypeDecoder[T]](value *T, dec *jsontext.Decoder) error {
	return TypeOptionalUnmarshaler[T, U, ITEM]{}.UnmarshalJSONFrom(value, dec)
}

type TypeOptionalMarshaler[T any, U *T, ITEM TypeEncoder[T]] struct{}

func (TypeOptionalMarshaler[T, U, ITEM]) isEncoder() {}

func (TypeOptionalMarshaler[T, U, ITEM]) MarshalJSONTo(receiver U, enc *jsontext.Encoder) error {
	if receiver == nil {
		return enc.WriteToken(jsontext.Null)
	}
	return (*new(ITEM)).MarshalJSONTo(*receiver, enc)
}

type TypeOptionalUnmarshaler[T any, U *T, ITEM TypeDecoder[T]] struct{}

func (TypeOptionalUnmarshaler[T, U, ITEM]) isDecoder() {}

func (TypeOptionalUnmarshaler[T, U, ITEM]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	if dec.PeekKind() == 'n' {
		// still have to consume token
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*receiver = *new(T)
		return nil
	}
	return (*new(ITEM)).UnmarshalJSONFrom(receiver, dec)
}
