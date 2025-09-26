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

func TestBoolean(t *testing.T) {
	type boolAlias bool
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "true",
			Test: typeTestCase[bool, cj.Boolean[bool], cj.Boolean[bool]]{
				Value: true, JSON: "true",
			},
		},
		{
			Name: "boolAlias(true)",
			Test: typeTestCase[boolAlias, cj.Boolean[boolAlias], cj.Boolean[boolAlias]]{
				Value: true, JSON: "true",
			},
		},
		{
			Name: "false",
			Test: typeTestCase[bool, cj.Boolean[bool], cj.Boolean[bool]]{
				Value: false, JSON: "false",
			},
		},
		{
			Name: "boolAlias(false)",
			Test: typeTestCase[boolAlias, cj.Boolean[boolAlias], cj.Boolean[boolAlias]]{
				Value: false, JSON: "false",
			},
		},
		{
			Name: "null",
			Test: typeTestCase[bool, cj.Boolean[bool], cj.Boolean[bool]]{
				JSON:                 "null",
				SkipTestMarshal:      true,
				ErrUnmarshalJSONFrom: "KindMismatchError at offset 4: want json boolean, got null",
			},
		},
		{
			Name: "map_keys",
			Test: typeTestCase[map[bool]bool, cj.ComparableMapMarshaler[map[bool]bool, bool, bool, cj.BooleanMapKey[bool], cj.Boolean[bool]], cj.MapUnmarshaler[map[bool]bool, bool, bool, cj.BooleanMapKey[bool], cj.Boolean[bool]]]{
				Value: map[bool]bool{true: true, false: false},
				JSON:  "{\"false\":false,\"true\":true}",
			},
		},
		{
			Name: "map_key_invalid",
			Test: typeTestCase[bool, cj.BooleanMapKey[bool], cj.BooleanMapKey[bool]]{
				JSON:                 `"invalid"`,
				SkipTestMarshal:      true,
				ErrUnmarshalJSONFrom: "InvalidValueError at offset 9: invalid boolean: strconv.ParseBool: parsing \"invalid\": invalid syntax",
			},
		},
		{
			Name: "map_key_not_string",
			Test: typeTestCase[bool, cj.BooleanMapKey[bool], cj.BooleanMapKey[bool]]{
				JSON:                 `true`,
				SkipTestMarshal:      true,
				ErrUnmarshalJSONFrom: "KindMismatchError at offset 4: want json string, got true",
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
