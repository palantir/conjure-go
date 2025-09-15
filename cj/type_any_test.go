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

package cj_test

import (
	stdjson "encoding/json"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/palantir/conjure-go/v6/cj"
)

func TestAny(t *testing.T) {
	for name, test := range map[string]typeTest{
		"int": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: 42, JSON: "42",
		},
		"string": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: "hello", JSON: "\"hello\"",
		},
		"float": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: 3.14, JSON: "3.14",
		},
		"bool": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: true, JSON: "true",
		},
		"null": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: nil, JSON: "null", ErrUnmarshalJSONFrom: "KindMismatchError at 0: want non-optional value, got null",
		},
		"null optional": typeTestCase[*any, cj.OptionalMarshaler[*any, any, cj.Any[any]], cj.OptionalUnmarshaler[*any, any, cj.Any[any]]]{
			Value: nil, JSON: "null",
		},
		"array": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: []any{"foo", float64(1), false}, JSON: "[\"foo\",1,false]",
		},
		"object": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: map[string]any{"a": float64(1), "b": true}, JSON: "{\"a\":1,\"b\":true}", Options: json.Deterministic(true),
		},
		"empty object": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: map[string]any{}, JSON: "{}",
		},
		"empty array": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: []any{}, JSON: "[]",
		},
		"json.Number": typeTestCase[stdjson.Number, cj.Any[stdjson.Number], cj.Any[stdjson.Number]]{
			Value: stdjson.Number(`3.14`), JSON: "3.14",
		},
		"json.RawMessage": typeTestCase[stdjson.RawMessage, cj.Any[stdjson.RawMessage], cj.Any[stdjson.RawMessage]]{
			Value: stdjson.RawMessage(`{"x":1}`), JSON: "{\"x\":1}",
		},
		"malformed JSON": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			JSON: "[1,2,", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "jsontext: unexpected EOF after offset 5",
		},
		"deeply nested array": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: []any{[]any{[]any{1.0}}}, JSON: "[[[1]]]",
		},
		"mixed types": typeTestCase[any, cj.Any[any], cj.Any[any]]{
			Value: []any{1.0, "two", true, nil}, JSON: "[1,\"two\",true,null]",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
