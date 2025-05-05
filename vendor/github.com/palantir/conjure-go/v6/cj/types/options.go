package types

import (
	stdjson "encoding/json"

	"github.com/go-json-experiment/json"
)

var DefaultOptions = json.JoinOptions(
	// Marshal Options
	json.Deterministic(true),
	json.WithMarshalers(json.JoinMarshalers(
		json.MarshalFunc(marshalJSONNumber),
		json.MarshalFunc(marshalJSONRawMessage),
	)),
	// Unmarshal Options
	json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFunc[*stdjson.Number](unmarshalJSONNumber),
		json.UnmarshalFunc[*stdjson.RawMessage](unmarshalJSONRawMessage),
	)),
)

// marshalJSONNumber marshals a json.Number as-is, since this type is
// not recognized by the json v2 encoder and gets quoted as a string.
func marshalJSONNumber(number stdjson.Number) ([]byte, error) {
	return []byte(number), nil
}

// unmarshalJSONNumber unmarshals a json.Number as-is, since this type is
// not recognized by the json v2 encoder and gets quoted as a string.
func unmarshalJSONNumber(data []byte, number *stdjson.Number) error {
	*number = stdjson.Number(data)
	return nil
}

// marshalJSONRawMessage marshals a json.RawMessage as-is, since this type is
// not recognized by the json v2 encoder and gets quoted as bytes.
func marshalJSONRawMessage(rawMessage stdjson.RawMessage) ([]byte, error) {
	return rawMessage, nil
}

// unmarshalJSONRawMessage unmarshals a json.RawMessage as-is, since this type is
// not recognized by the json v2 encoder and gets quoted as bytes.
func unmarshalJSONRawMessage(data []byte, message *stdjson.RawMessage) error {
	*message = data
	return nil
}
