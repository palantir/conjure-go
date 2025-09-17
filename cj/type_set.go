package cj

import (
	"reflect"
	"slices"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// SetMarshaler provides json marshaling for sets of type T using a nested encoder ITEM.
// The emitted JSON list's elements will be unique, but otherwise in the same order as the original.
type SetMarshaler[T ~[]U, U comparable, ITEM TypeEncoder[U]] struct{}

func (SetMarshaler[T, U, ITEM]) MarshalJSONTo(enc *jsontext.Encoder, receiver T) error {
	if receiver == nil && getOptionOrFalse(enc.Options(), json.FormatNilSliceAsNull) {
		if err := enc.WriteToken(jsontext.Null); err != nil {
			return WrapEncodeError(enc, "", err)
		}
		return nil
	}
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return WrapEncodeError(enc, "", err)
	}
	for i, item := range receiver {
		if slices.Contains(receiver[0:i], item) {
			// duplicate item
			continue
		}
		if err := (*new(ITEM)).MarshalJSONTo(enc, item); err != nil {
			return err
		}
	}
	if err := enc.WriteToken(jsontext.EndArray); err != nil {
		return WrapEncodeError(enc, "", err)
	}
	return nil
}

type SetUnmarshaler[T ~[]U, U comparable, ITEM TypeDecoder[U]] struct{}

func (SetUnmarshaler[T, U, ITEM]) UnmarshalJSONFrom(dec *jsontext.Decoder, receiver *T) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return WrapSyntaxError(dec, "", err)
	}
	if *receiver == nil {
		*receiver = make(T, 0)
	} else {
		*receiver = (*receiver)[:0]
	}
	kind := tok.Kind()
	if kind == 'n' {
		// null
		return nil
	}
	if kind != '[' {
		return NewKindMismatchError(dec, kind, "list opening bracket")
	}
	for {
		if dec.PeekKind() == ']' {
			_, err := dec.ReadToken()
			return err
		}
		item := *new(U)
		if err := (*new(ITEM)).UnmarshalJSONFrom(dec, &item); err != nil {
			return err
		}
		if slices.Contains(*receiver, item) {
			return NewDuplicateSetItemError(dec, reflect.TypeOf(item).Name(), len(*receiver))
		}
		*receiver = append(*receiver, item)
	}
}
