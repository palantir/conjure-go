package types

import (
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
)

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
		return cj.NewSyntaxError(dec, "list expected '['")
	}
}
