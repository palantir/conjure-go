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

func TestList(t *testing.T) {
	for name, test := range map[string]typeTest{
		"empty": typeTestCase[[]int, types.ListMarshaler[[]int, int, types.Int[int]], types.ListUnmarshaler[[]int, int, types.Int[int]]]{
			Value: []int{}, JSON: "[]",
		},
		"one": typeTestCase[[]int, types.ListMarshaler[[]int, int, types.Int[int]], types.ListUnmarshaler[[]int, int, types.Int[int]]]{
			Value: []int{42}, JSON: "[42]",
		},
		"several": typeTestCase[[]string, types.ListMarshaler[[]string, string, types.String[string]], types.ListUnmarshaler[[]string, string, types.String[string]]]{
			Value: []string{"a", "b", "c"}, JSON: "[\"a\",\"b\",\"c\"]",
		},
		"nested": typeTestCase[[][]bool, types.ListMarshaler[[][]bool, []bool, types.ListMarshaler[[]bool, bool, types.Boolean[bool]]], types.ListUnmarshaler[[][]bool, []bool, types.ListUnmarshaler[[]bool, bool, types.Boolean[bool]]]]{
			Value: [][]bool{{true, false}, {}}, JSON: "[[true,false],[]]",
		},
		"null_marshal": typeTestCase[[]int, types.ListMarshaler[[]int, int, types.Int[int]], types.ListUnmarshaler[[]int, int, types.Int[int]]]{
			JSON: "[]", SkipTestUnmarshal: true, Value: []int(nil),
		},
		"null_unmarshal": typeTestCase[[]int, types.ListMarshaler[[]int, int, types.Int[int]], types.ListUnmarshaler[[]int, int, types.Int[int]]]{
			JSON: "null", SkipTestMarshal: true, Value: []int{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
