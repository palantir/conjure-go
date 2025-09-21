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
	"github.com/go-json-experiment/json/jsontext"
)

// VisitObjectFields is a helper for use in UnmarshalJSONFrom implementations that reads the opening and closing braces
// and calls visitField for each key and value in the object.
func VisitObjectFields(dec *jsontext.Decoder, visitField func(key string, dec *jsontext.Decoder) error) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if kind := tok.Kind(); kind != '{' {
		return NewKindMismatchError(dec, kind, "object opening brace")
	}
	for {
		key, err := dec.ReadToken()
		if err != nil {
			return err
		}
		kind := key.Kind()
		if kind == '}' {
			return nil // End of object
		}
		if kind != '"' {
			return NewKindMismatchError(dec, kind, "object closing brace or next key")
		}
		if err := visitField(key.String(), dec); err != nil {
			return err
		}
		// continue to next key
	}
}
