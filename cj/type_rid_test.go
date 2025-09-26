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
	"github.com/palantir/pkg/rid"
)

func TestRID(t *testing.T) {
	type ridAlias rid.ResourceIdentifier
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "empty",
			Test: typeTestCase[rid.ResourceIdentifier, cj.RID[rid.ResourceIdentifier], cj.RID[rid.ResourceIdentifier]]{
				Value: rid.ResourceIdentifier{}, JSON: "\"ri....\"", ErrUnmarshalJSONFrom: "SyntaxError at offset 8: invalid resource identifier: rid first segment (service) does not match ^[a-z][a-z0-9\\-]*$ pattern: rid third segment (type) does not match ^[a-z][a-z0-9\\-]*$ pattern: rid fourth segment (locator) does not match ^[a-zA-Z0-9\\-\\._]+$ pattern",
			},
		},
		{
			Name: "resource",
			Test: typeTestCase[rid.ResourceIdentifier, cj.RID[rid.ResourceIdentifier], cj.RID[rid.ResourceIdentifier]]{
				Value: must(rid.ParseRID("ri.example.main.foo.bar")), JSON: "\"ri.example.main.foo.bar\"",
			},
		},
		{
			Name: "alias",
			Test: typeTestCase[ridAlias, cj.RID[ridAlias], cj.RID[ridAlias]]{
				Value: ridAlias(must(rid.ParseRID("ri.example.main.foo.bar"))), JSON: "\"ri.example.main.foo.bar\"",
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
