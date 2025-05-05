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
	"testing"
	"time"

	"github.com/palantir/conjure-go/v6/cj/types"
	"github.com/palantir/pkg/datetime"
)

func TestMap(t *testing.T) {
	for name, test := range map[string]typeTest{
		"empty": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			Value: map[string]int{}, JSON: "{}",
		},
		"one": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			Value: map[string]int{"foo": 1}, JSON: "{\"foo\":1}",
		},
		"ordered": typeTestCase[map[string]string, types.OrderedMapMarshaler[map[string]string, string, string, types.String[string], types.String[string]], types.MapUnmarshaler[map[string]string, string, string, types.String[string], types.String[string]]]{
			Value: map[string]string{"j": "10", "i": "9", "h": "8", "g": "7", "f": "6", "e": "5", "d": "4", "c": "3", "b": "2", "a": "1"}, JSON: "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\",\"d\":\"4\",\"e\":\"5\",\"f\":\"6\",\"g\":\"7\",\"h\":\"8\",\"i\":\"9\",\"j\":\"10\"}",
		},
		"nested": typeTestCase[map[string][]int, types.OrderedMapMarshaler[map[string][]int, string, []int, types.String[string], types.ListMarshaler[[]int, int, types.Int[int]]], types.MapUnmarshaler[map[string][]int, string, []int, types.String[string], types.ListUnmarshaler[[]int, int, types.Int[int]]]]{
			Value: map[string][]int{"nums": {1, 2, 3}}, JSON: "{\"nums\":[1,2,3]}",
		},
		"boolean map key": typeTestCase[map[bool]int, types.ComparableMapMarshaler[map[bool]int, bool, int, types.BooleanMapKey[bool], types.Int[int]], types.MapUnmarshaler[map[bool]int, bool, int, types.BooleanMapKey[bool], types.Int[int]]]{
			Value: map[bool]int{true: 2, false: 2}, JSON: "{\"false\":2,\"true\":2}",
		},
		"datetime map key": typeTestCase[map[datetime.DateTime]string, types.ComparableMapMarshaler[map[datetime.DateTime]string, datetime.DateTime, string, types.DateTime[datetime.DateTime], types.String[string]], types.MapUnmarshaler[map[datetime.DateTime]string, datetime.DateTime, string, types.DateTime[datetime.DateTime], types.String[string]]]{
			Value: map[datetime.DateTime]string{datetime.DateTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)): "2024-01-01T00:00:00Z", datetime.DateTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)): "2025-01-01T00:00:00Z"}, JSON: "{\"2024-01-01T00:00:00Z\":\"2024-01-01T00:00:00Z\",\"2025-01-01T00:00:00Z\":\"2025-01-01T00:00:00Z\"}",
		},
		"null_marshal": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "{}", SkipTestUnmarshal: true, Value: map[string]int(nil),
		},
		"null_unmarshal": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "null", SkipTestMarshal: true, Value: map[string]int{},
		},
		"not an object": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "[]", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 1: want map opening brace, got [",
		},
		"duplicate key": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "{\"a\":1,\"a\":2}", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "jsontext: duplicate object member name \"a\"",
		},
		"key not string": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "{1:2}", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "jsontext: object member name must be a string after offset 1",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", func(t *testing.T) {
				test.TestMarshal(t)
			})
			t.Run("Unmarshal", func(t *testing.T) {
				test.TestUnmarshal(t)
			})
		})
	}
}
