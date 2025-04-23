package cj

import (
	"encoding"

	"github.com/go-json-experiment/json/jsontext"
)

// MarshalTextTo uses a TypeText to encode the value as JSON.
func MarshalTextTo[T encoding.TextMarshaler](value T, enc *jsontext.Encoder) error {
	return TypeTextMarshaler[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalTextFrom uses a TypeText to decode the value from JSON.
func UnmarshalTextFrom[T encoding.TextUnmarshaler](value T, dec *jsontext.Decoder) error {
	return TypeTextUnmarshaler[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeTextMarshaler[T encoding.TextMarshaler] struct{}

func (TypeTextMarshaler[T]) isEncoder() {}

func (TypeTextMarshaler[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	text, err := receiver.MarshalText()
	if err != nil {
		return err
	}
	out := enc.UnusedBuffer()
	out, err = jsontext.AppendQuote(out, text)
	if err != nil {
		return err
	}
	return enc.WriteValue(out)
}

type TypeTextUnmarshaler[T encoding.TextUnmarshaler] struct{}

func (TypeTextUnmarshaler[T]) isDecoder() {}

func (TypeTextUnmarshaler[T]) UnmarshalJSONFrom(receiver T, dec *jsontext.Decoder) error {
	val, err := readValueOfKind(dec, '"')
	if err != nil {
		return err
	}
	unquoted, err := jsontext.AppendUnquote(nil, val)
	if err != nil {
		return err
	}
	return receiver.UnmarshalText(unquoted)
}
