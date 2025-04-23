package cj

import (
	"time"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/datetime"
)

// MarshalDateTimeTo uses a TypeDateTime to encode the value as JSON.
func MarshalDateTimeTo[T time.Time | datetime.DateTime](value T, enc *jsontext.Encoder) error {
	return TypeDateTime[T]{}.MarshalJSONTo(value, enc)
}

// UnmarshalDateTimeFrom uses a TypeDateTime to decode the value from JSON.
func UnmarshalDateTimeFrom[T time.Time | datetime.DateTime](value *T, dec *jsontext.Decoder) error {
	return TypeDateTime[T]{}.UnmarshalJSONFrom(value, dec)
}

type TypeDateTime[T time.Time | datetime.DateTime] struct{}

func (TypeDateTime[T]) isEncoder() {}
func (TypeDateTime[T]) isDecoder() {}

func (TypeDateTime[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	out := enc.UnusedBuffer()
	out = append(out, '"')
	out = time.Time(receiver).AppendFormat(out, time.RFC3339Nano)
	out = append(out, '"')
	return enc.WriteValue(out)
}

func (TypeDateTime[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	tok, err := readTokenOfKind(dec, '"')
	if err != nil {
		return err
	}
	parse, err := time.Parse(tok.String(), time.RFC3339Nano)
	if err != nil {
		return NewSyntaxError(dec, "invalid datetime", err)
	}
	*receiver = T(parse)
	return nil
}
