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
	"github.com/palantir/pkg/bearertoken"
)

func TestBearerToken(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "bearertoken",
			Test: typeTestCase[bearertoken.Token, cj.BearerToken[bearertoken.Token], cj.BearerToken[bearertoken.Token]]{
				Value: "foo", JSON: "\"foo\"",
			},
		},
		{
			Name: "null",
			Test: typeTestCase[bearertoken.Token, cj.BearerToken[bearertoken.Token], cj.BearerToken[bearertoken.Token]]{
				JSON: "null", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at offset 4: want json string, got null",
			},
		},
		{
			Name: "invalid",
			Test: typeTestCase[bearertoken.Token, cj.BearerToken[bearertoken.Token], cj.BearerToken[bearertoken.Token]]{
				JSON: "\" \"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at offset 3: invalid character for bearer token",
			},
		},
		{
			Name: "empty",
			Test: typeTestCase[bearertoken.Token, cj.BearerToken[bearertoken.Token], cj.BearerToken[bearertoken.Token]]{
				JSON: "\"\"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at offset 2: empty bearer token",
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
