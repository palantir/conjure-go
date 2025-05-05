// Copyright (c) 2025 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	stdjson "encoding/json"

	"github.com/go-json-experiment/json"
)

var DefaultOptions = json.JoinOptions(
	// Marshal Options
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

func UnmarshalAs[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}

func MarshalAs[T any](value T) ([]byte, error) {
	return json.Marshal(value)
}
