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

func TestOptional(t *testing.T) {
	type optStr = *string
	type optInt = *int

	for name, test := range map[string]typeTest{
		"nil string": typeTestCase[optStr, cj.OptionalMarshaler[optStr, string, cj.String[string]], cj.OptionalUnmarshaler[optStr, string, cj.String[string]]]{
			Value: nil, JSON: "null",
		},
		"some string": typeTestCase[optStr, cj.OptionalMarshaler[optStr, string, cj.String[string]], cj.OptionalUnmarshaler[optStr, string, cj.String[string]]]{
			Value: mustPtr("foo"), JSON: "\"foo\"",
		},
		"nil int": typeTestCase[optInt, cj.OptionalMarshaler[optInt, int, cj.Int32[int]], cj.OptionalUnmarshaler[optInt, int, cj.Int32[int]]]{
			Value: nil, JSON: "null",
		},
		"some int": typeTestCase[optInt, cj.OptionalMarshaler[optInt, int, cj.Int32[int]], cj.OptionalUnmarshaler[optInt, int, cj.Int32[int]]]{
			Value: mustPtr(123), JSON: "123",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

func mustPtr[T any](v T) *T { return &v }
