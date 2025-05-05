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

	"github.com/palantir/conjure-go/v6/cj/types"
)

func TestMap(t *testing.T) {
	for name, test := range map[string]typeTest{
		"empty": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			Value: map[string]int{}, JSON: "{}",
		},
		"one": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			Value: map[string]int{"foo": 1}, JSON: "{\"foo\":1}",
		},
		"several": typeTestCase[map[string]string, types.OrderedMapMarshaler[map[string]string, string, string, types.String[string], types.String[string]], types.MapUnmarshaler[map[string]string, string, string, types.String[string], types.String[string]]]{
			Value: map[string]string{"a": "x", "b": "y"}, JSON: "{\"a\":\"x\",\"b\":\"y\"}",
		},
		"nested": typeTestCase[map[string][]int, types.OrderedMapMarshaler[map[string][]int, string, []int, types.String[string], types.ListMarshaler[[]int, int, types.Int[int]]], types.MapUnmarshaler[map[string][]int, string, []int, types.String[string], types.ListUnmarshaler[[]int, int, types.Int[int]]]]{
			Value: map[string][]int{"nums": {1, 2, 3}}, JSON: "{\"nums\":[1,2,3]}",
		},
		"null_marshal": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "{}", SkipTestUnmarshal: true, Value: map[string]int(nil),
		},
		"null_unmarshal": typeTestCase[map[string]int, types.OrderedMapMarshaler[map[string]int, string, int, types.String[string], types.Int[int]], types.MapUnmarshaler[map[string]int, string, int, types.String[string], types.Int[int]]]{
			JSON: "null", SkipTestMarshal: true, Value: map[string]int{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
