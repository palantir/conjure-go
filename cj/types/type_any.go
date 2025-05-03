package types

import (
	stdjson "encoding/json"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Any[T any] struct{}

func (Any[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, receiver, json.WithMarshalers(json.MarshalFunc(marshalJSONNumber)))
}

func (Any[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, receiver)
}

func marshalJSONNumber(number stdjson.Number) ([]byte, error) {
	return []byte(number), nil
}
