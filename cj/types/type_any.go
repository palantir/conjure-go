package types

import (
	stdjson "encoding/json"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Any provides generic JSON marshaling and unmarshaling for any Go type T.
// It is a fallback encoder/decoder for types not otherwise handled by more specific
// implementations. Use this when you want to delegate to the default Go JSON logic,
// but still participate in the MarshalerTo/UnmarshalerFrom interfaces.
type Any[T any] struct{}

func (Any[T]) MarshalJSONTo(receiver T, enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, receiver, json.WithMarshalers(json.MarshalFunc(marshalJSONNumber)))
}

func (Any[T]) UnmarshalJSONFrom(receiver *T, dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, receiver)
}

// marshalJSONNumber marshals a stdjson.Number as-is, since this type is
// not recognized by the json v2 encoder and gets quoted as a string.
func marshalJSONNumber(number stdjson.Number) ([]byte, error) {
	return []byte(number), nil
}
