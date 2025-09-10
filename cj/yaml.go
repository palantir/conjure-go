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

package cj

import (
	"bytes"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/pkg/safeyaml"
)

// MarshalYAML is meant as a dependency for types implementing yaml.Marshaler.
// It encodes the object to JSON then convert that JSON to YAML.
func MarshalYAML[T json.MarshalerTo](receiver T, opts ...json.Options) (any, error) {
	buf := bytes.NewBuffer(nil)
	enc := jsontext.NewEncoder(buf, opts...)
	if err := receiver.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(buf.Bytes())
}

// UnmarshalYAML is meant as a dependency for types implementing yaml.Unmarshaler.
// It converts the YAML to JSON then unmarshals that JSON into the receiver.
func UnmarshalYAML[T json.UnmarshalerFrom](receiver T, unmarshal func(any) error, opts ...json.Options) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	dec := jsontext.NewDecoder(bytes.NewReader(jsonBytes), opts...)
	return receiver.UnmarshalJSONFrom(dec)
}
