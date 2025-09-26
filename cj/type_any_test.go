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
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "int",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: 42, JSON: "42",
			},
		},
		{
			Name: "string",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: "hello", JSON: "\"hello\"",
			},
		},
		{
			Name: "float",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: 3.14, JSON: "3.14",
			},
		},
		{
			Name: "bool",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: true, JSON: "true",
			},
		},
		{
			Name: "null",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: nil, JSON: "null", ErrUnmarshalJSONFrom: "KindMismatchError at offset 0: want non-optional value, got null",
			},
		},
		{
			Name: "null optional",
			Test: typeTestCase[*any, cj.OptionalMarshaler[*any, any, cj.Any[any]], cj.OptionalUnmarshaler[*any, any, cj.Any[any]]]{
				Value: nil, JSON: "null",
			},
		},
		{
			Name: "array",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: []any{"foo", float64(1), false}, JSON: "[\"foo\",1,false]",
			},
		},
		{
			Name: "object",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: map[string]any{"a": float64(1), "b": true}, JSON: "{\"a\":1,\"b\":true}", Options: json.Deterministic(true),
			},
		},
		{
			Name: "empty object",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: map[string]any{}, JSON: "{}",
			},
		},
		{
			Name: "empty array",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: []any{}, JSON: "[]",
			},
		},
		{
			Name: "json.Number",
			Test: typeTestCase[stdjson.Number, cj.Any[stdjson.Number], cj.Any[stdjson.Number]]{
				Value: stdjson.Number(`3.14`), JSON: "3.14",
			},
		},
		{
			Name: "json.RawMessage",
			Test: typeTestCase[stdjson.RawMessage, cj.Any[stdjson.RawMessage], cj.Any[stdjson.RawMessage]]{
				Value: stdjson.RawMessage(`{"x":1}`), JSON: "{\"x\":1}",
			},
		},
		{
			Name: "malformed JSON",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				JSON: "[1,2,", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "jsontext: unexpected EOF after offset 5",
			},
		},
		{
			Name: "deeply nested array",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: []any{[]any{[]any{1.0}}}, JSON: "[[[1]]]",
			},
		},
		{
			Name: "mixed types",
			Test: typeTestCase[any, cj.Any[any], cj.Any[any]]{
				Value: []any{1.0, "two", true, nil}, JSON: "[1,\"two\",true,null]",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("Marshal", tc.Test.TestMarshal)
			t.Run("Unmarshal", tc.Test.TestUnmarshal)
		})
	}
}
