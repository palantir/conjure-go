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

package types_test

import (
	"encoding/json"
	"testing"

	"github.com/palantir/conjure-go/v6/cj/types"
)

func TestAny(t *testing.T) {
	for name, test := range map[string]typeTest{
		"int": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: 42, JSON: "42",
		},
		"string": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: "hello", JSON: "\"hello\"",
		},
		"float": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: 3.14, JSON: "3.14",
		},
		"bool": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: true, JSON: "true",
		},
		"null": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: nil, JSON: "null",
		},
		"array": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: []any{"foo", float64(1), false}, JSON: "[\"foo\",1,false]",
		},
		"object": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: map[string]any{"a": float64(1), "b": true}, JSON: "{\"a\":1,\"b\":true}",
		},
		"empty object": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: map[string]any{}, JSON: "{}",
		},
		"empty array": typeTestCase[any, types.Any[any], types.Any[any]]{
			Value: []any{}, JSON: "[]",
		},
		"json.Number": typeTestCase[json.Number, types.Any[json.Number], types.Any[json.Number]]{
			Value: json.Number(`3.14`), JSON: "3.14",
		},
		"json.RawMessage": typeTestCase[json.RawMessage, types.Any[json.RawMessage], types.Any[json.RawMessage]]{
			Value: json.RawMessage(`{"x":1}`), JSON: "{\"x\":1}",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
