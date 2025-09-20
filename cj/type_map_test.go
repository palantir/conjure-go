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
	"math"
	"testing"
	"time"

	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/datetime"
)

func TestMap(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "empty",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				Value: map[string]int{}, JSON: "{}",
			},
		},
		{
			Name: "one",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				Value: map[string]int{"foo": 1}, JSON: "{\"foo\":1}",
			},
		},
		{
			Name: "ordered",
			Test: typeTestCase[map[string]string, cj.OrderedMapMarshaler[map[string]string, string, string, cj.String[string], cj.String[string]], cj.MapUnmarshaler[map[string]string, string, string, cj.String[string], cj.String[string]]]{
				Value: map[string]string{"j": "10", "i": "9", "h": "8", "g": "7", "f": "6", "e": "5", "d": "4", "c": "3", "b": "2", "a": "1"}, JSON: "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\",\"d\":\"4\",\"e\":\"5\",\"f\":\"6\",\"g\":\"7\",\"h\":\"8\",\"i\":\"9\",\"j\":\"10\"}",
			},
		},
		{
			Name: "ordered_int_keys",
			Test: typeTestCase[map[int]int, cj.OrderedMapMarshaler[map[int]int, int, int, cj.Int32MapKey[int], cj.Int32[int]], cj.MapUnmarshaler[map[int]int, int, int, cj.Int32MapKey[int], cj.Int32[int]]]{
				Value: map[int]int{100: 100, 10: 10, 9: 9, 1: 1, 0: 0, -1: -1}, JSON: "{\"-1\":-1,\"0\":0,\"1\":1,\"9\":9,\"10\":10,\"100\":100}",
			},
		},
		{
			Name: "ordered_float_keys",
			Test: typeTestCase[map[float64]float64, cj.OrderedMapMarshaler[map[float64]float64, float64, float64, cj.FloatMapKey[float64], cj.Float[float64]], cj.MapUnmarshaler[map[float64]float64, float64, float64, cj.FloatMapKey[float64], cj.Float[float64]]]{
				Value: map[float64]float64{100: 100, 10: 10, 9: 9, 1: 1, 0: 0, -1: -1, -0.10: -0.10, -0.9: -0.9, math.Inf(1): math.Inf(1), math.Inf(-1): math.Inf(-1)},
				JSON:  "{\"-Infinity\":\"-Infinity\",\"-1\":-1,\"-0.9\":-0.9,\"-0.1\":-0.1,\"0\":0,\"1\":1,\"9\":9,\"10\":10,\"100\":100,\"Infinity\":\"Infinity\"}",
			},
		},
		{
			Name: "nested",
			Test: typeTestCase[map[string][]int, cj.OrderedMapMarshaler[map[string][]int, string, []int, cj.String[string], cj.ListMarshaler[[]int, int, cj.Int32[int]]], cj.MapUnmarshaler[map[string][]int, string, []int, cj.String[string], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]]{
				Value: map[string][]int{"nums": {1, 2, 3}}, JSON: "{\"nums\":[1,2,3]}",
			},
		},
		{
			Name: "boolean map key",
			Test: typeTestCase[map[bool]int, cj.ComparableMapMarshaler[map[bool]int, bool, int, cj.BooleanMapKey[bool], cj.Int32[int]], cj.MapUnmarshaler[map[bool]int, bool, int, cj.BooleanMapKey[bool], cj.Int32[int]]]{
				Value: map[bool]int{true: 2, false: 2}, JSON: "{\"false\":2,\"true\":2}",
			},
		},
		{
			Name: "datetime map key",
			Test: typeTestCase[map[datetime.DateTime]string, cj.ComparableMapMarshaler[map[datetime.DateTime]string, datetime.DateTime, string, cj.DateTime[datetime.DateTime], cj.String[string]], cj.MapUnmarshaler[map[datetime.DateTime]string, datetime.DateTime, string, cj.DateTime[datetime.DateTime], cj.String[string]]]{
				Value: map[datetime.DateTime]string{datetime.DateTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)): "2024-01-01T00:00:00Z", datetime.DateTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)): "2025-01-01T00:00:00Z"}, JSON: "{\"2024-01-01T00:00:00Z\":\"2024-01-01T00:00:00Z\",\"2025-01-01T00:00:00Z\":\"2025-01-01T00:00:00Z\"}",
			},
		},
		{
			Name: "null_marshal",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				JSON: "{}", SkipTestUnmarshal: true, Value: map[string]int(nil),
			},
		},
		{
			Name: "null_unmarshal",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				JSON: "null", SkipTestMarshal: true, Value: map[string]int{},
			},
		},
		{
			Name: "not an object",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				JSON: "[]", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 1: want map opening brace, got [",
			},
		},
		{
			Name: "duplicate key",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				JSON: "{\"a\":1,\"a\":2}", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "jsontext: duplicate object member name \"a\"",
			},
		},
		{
			Name: "duplicate int key",
			Test: typeTestCase[map[int]int, cj.OrderedMapMarshaler[map[int]int, int, int, cj.Int32MapKey[int], cj.Int32[int]], cj.MapUnmarshaler[map[int]int, int, int, cj.Int32MapKey[int], cj.Int32[int]]]{
				JSON: "{\"01\":1,\"1\":2}", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "DuplicateMapKeyError at 13: type *map[int]int has duplicate map keys",
			},
		},
		{
			Name: "key not string",
			Test: typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
				JSON: "{1:2}", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "jsontext: object member name must be a string after offset 1",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("Marshal", func(t *testing.T) {
				tc.Test.TestMarshal(t)
			})
			t.Run("Unmarshal", func(t *testing.T) {
				tc.Test.TestUnmarshal(t)
			})
		})
	}
}
