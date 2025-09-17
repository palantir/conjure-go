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
	"testing"

	"github.com/palantir/conjure-go/v6/cj"
)

func TestList(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "empty",
			Test: typeTestCase[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{}, JSON: "[]",
			},
		},
		{
			Name: "one",
			Test: typeTestCase[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{42}, JSON: "[42]",
			},
		},
		{
			Name: "several",
			Test: typeTestCase[[]string, cj.ListMarshaler[[]string, string, cj.String[string]], cj.ListUnmarshaler[[]string, string, cj.String[string]]]{
				Value: []string{"a", "b", "c"}, JSON: "[\"a\",\"b\",\"c\"]",
			},
		},
		{
			Name: "nested",
			Test: typeTestCase[[][]bool, cj.ListMarshaler[[][]bool, []bool, cj.ListMarshaler[[]bool, bool, cj.Boolean[bool]]], cj.ListUnmarshaler[[][]bool, []bool, cj.ListUnmarshaler[[]bool, bool, cj.Boolean[bool]]]]{
				Value: [][]bool{{true, false}, {}}, JSON: "[[true,false],[]]",
			},
		},
		{
			Name: "null_marshal",
			Test: typeTestCase[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]{
				JSON: "[]", SkipTestUnmarshal: true, Value: []int(nil),
			},
		},
		{
			Name: "null_unmarshal",
			Test: typeTestCase[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]{
				JSON: "null", SkipTestMarshal: true, Value: []int{},
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
