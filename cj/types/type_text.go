package types

import (
	"bytes"
	"encoding"
	"github.com/palantir/conjure-go/v6/cj"

	"github.com/go-json-experiment/json/jsontext"
)

type TextMarshaler[T encoding.TextMarshaler] struct{}

func (TextMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	text, err := receiver.MarshalText()
	if err != nil {
		return err
	}
	return enc.WriteToken(jsontext.String(string(text)))
}

func (t TextMarshaler[T]) Compare(a, b T) int {
	aText, errA := a.MarshalText()
	bText, errB := b.MarshalText()
	if errA != nil || errB != nil {
		return 0
	}
	return bytes.Compare(aText, bText)
}

type TextUnmarshaler[T encoding.TextUnmarshaler] struct{}

func (TextUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '"' {
		return cj.NewKindMismatchError(dec, tok.Kind(), "text")
	}
	return receiver.UnmarshalText([]byte(tok.String()))
}
