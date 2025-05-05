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

func TestBoolean(t *testing.T) {
	type boolAlias bool
	for name, test := range map[string]typeTest{
		"true": typeTestCase[bool, types.Boolean[bool], types.Boolean[bool]]{
			Value: true, JSON: "true",
		},
		"boolAlias(true)": typeTestCase[boolAlias, types.Boolean[boolAlias], types.Boolean[boolAlias]]{
			Value: true, JSON: "true",
		},
		"false": typeTestCase[bool, types.Boolean[bool], types.Boolean[bool]]{
			Value: false, JSON: "false",
		},
		"boolAlias(false)": typeTestCase[boolAlias, types.Boolean[boolAlias], types.Boolean[boolAlias]]{
			Value: false, JSON: "false",
		},
		"null": typeTestCase[bool, types.Boolean[bool], types.Boolean[bool]]{
			JSON:                 "null",
			SkipTestMarshal:      true,
			ErrUnmarshalJSONFrom: "KindMismatchError at 4: want json boolean, got null",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
