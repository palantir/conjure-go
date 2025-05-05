package types

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

// ListMarshaler provides JSON marshaling for slices of type T using a nested encoder ITEM.
// Encodes slices as JSON arrays, delegating encoding of each element to ITEM.
type ListMarshaler[T any, ITEM cj.TypeEncoder[T]] struct{}

func (ListMarshaler[T, ITEM]) MarshalJSONTo(receiver []T, enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	for _, item := range receiver {
		err := (*new(ITEM)).MarshalJSONTo(item, enc)
		if err != nil {
			return err
		}
	}
	if err := enc.WriteToken(jsontext.EndArray); err != nil {
		return err
	}
	return nil
}

// ListUnmarshaler provides JSON unmarshaling for slices of type T using a nested decoder ITEM.
// Decodes JSON arrays into Go slices, delegating decoding of each element to ITEM.
type ListUnmarshaler[T any, ITEM cj.TypeDecoder[T]] struct{}

func (ListUnmarshaler[T, ITEM]) UnmarshalJSONFrom(receiver *[]T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if *receiver == nil {
		*receiver = make([]T, 0)
	} else {
		*receiver = (*receiver)[:0]
	}
	switch tok.Kind() {
	case '[':
		for {
			if dec.PeekKind() == ']' {
				_, err := dec.ReadToken()
				return err
			}
			item := *new(T)
			if err := (*new(ITEM)).UnmarshalJSONFrom(&item, dec); err != nil {
				return err
			}
			*receiver = append(*receiver, item)
		}
	case 'n':
		return nil
	default:
		return cj.NewKindMismatchError(dec, tok.Kind(), "list opening bracket")
	}
}
