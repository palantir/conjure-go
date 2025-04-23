package cj

import (
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalListTo uses a TypeList to encode the value as JSON.
func MarshalListTo[T any, ITEM TypeArshaler[T]](value []T, enc *jsontext.Encoder) error {
	return TypeListMarshaler[T, ITEM]{}.MarshalJSONTo(value, enc)
}

// UnmarshalListFrom uses a TypeList to decode the value from JSON.
func UnmarshalListFrom[T any, ITEM TypeArshaler[T]](value *[]T, dec *jsontext.Decoder) error {
	return TypeListUnmarshaler[T, ITEM]{}.UnmarshalJSONFrom(value, dec)
}

type TypeListMarshaler[T any, ITEM TypeEncoder[T]] struct{}

func (TypeListMarshaler[T, ITEM]) isEncoder() {}
func (TypeListMarshaler[T, ITEM]) isDecoder() {}

func (TypeListMarshaler[T, ITEM]) MarshalJSONTo(receiver []T, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	for _, item := range receiver {
		err := (*new(ITEM)).MarshalJSONTo(item, enc)
		if err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndArray)
}

type TypeListUnmarshaler[T any, ITEM TypeDecoder[T]] struct{}

func (TypeListUnmarshaler[T, ITEM]) isDecoder() {}

func (TypeListUnmarshaler[T, ITEM]) UnmarshalJSONFrom(receiver *[]T, dec *jsontext.Decoder) error {
	if dec.PeekKind() == 'n' {
		// still have to consume token
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*receiver = make([]T, 0)
		return nil
	}
	if _, err := readTokenOfKind(dec, '['); err != nil {
		return err
	}
	if *receiver == nil {
		*receiver = make([]T, 0)
	} else {
		*receiver = (*receiver)[:0]
	}
	for {
		item := *new(T)
		if err := (*new(ITEM)).UnmarshalJSONFrom(&item, dec); err != nil {
			return err
		}
		*receiver = append(*receiver, item)

		if dec.PeekKind() == ']' {
			break
		}
	}
	if _, err := readTokenOfKind(dec, ']'); err != nil {
		return err
	}
	return nil
}
