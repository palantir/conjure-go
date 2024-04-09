// Copyright (c) 2023 Palantir Technologies. All rights reserved.
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

package dj

import (
	"encoding/json"

	"github.com/palantir/pkg/safeyaml"
)

// MarshalYAML marshals the given json.Marshaler to YAML.
// Used to implement yaml.Marshaler, for example:
//
//	func (f foo) MarshalYAML() (interface{}, error) {
//	  return dj.MarshalYAML(f)
//	}
func MarshalYAML(marshaler json.Marshaler) (interface{}, error) {
	jsonBytes, err := marshaler.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}

// UnmarshalYAML unmarshals the given json.Unmarshaler from YAML.
// Used to implement yaml.Unmarshaler, for example:
//
//	func (f *foo) UnmarshalYAML(unmarshal func(interface{}) error) error {
//	  return dj.UnmarshalYAML(f, unmarshal)
//	}
func UnmarshalYAML(unmarshaler json.Unmarshaler, unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return unmarshaler.UnmarshalJSON(jsonBytes)
}
