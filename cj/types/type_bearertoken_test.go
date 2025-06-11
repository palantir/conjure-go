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
	"github.com/palantir/pkg/bearertoken"
)

func TestBearerToken(t *testing.T) {
	for name, test := range map[string]typeTest{
		"bearertoken": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			Value: "foo", JSON: "\"foo\"",
		},
		"null": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			JSON: "null", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 4: want json string, got null",
		},
		"invalid": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			JSON: "\" \"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at 3: invalid character for bearer token",
		},
		"empty": typeTestCase[bearertoken.Token, types.BearerToken[bearertoken.Token], types.BearerToken[bearertoken.Token]]{
			JSON: "\"\"", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "InvalidValueError at 2: empty bearer token",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
