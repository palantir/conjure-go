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
	"github.com/palantir/pkg/uuid"
)

func TestString(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "empty",
			Test: typeTestCase[string, cj.String[string], cj.String[string]]{
				Value: "", JSON: "\"\"",
			},
		},
		{
			Name: "ascii",
			Test: typeTestCase[string, cj.String[string], cj.String[string]]{
				Value: "hello", JSON: "\"hello\"",
			},
		},
		{
			Name: "unicode",
			Test: typeTestCase[string, cj.String[string], cj.String[string]]{
				Value: "héllo 世界", JSON: "\"héllo 世界\"",
			},
		},
		{
			Name: "escaped",
			Test: typeTestCase[string, cj.String[string], cj.String[string]]{
				Value: "foo\nbar", JSON: "\"foo\\nbar\"",
			},
		},
		{
			Name: "null",
			Test: typeTestCase[string, cj.String[string], cj.String[string]]{
				JSON: "null", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "KindMismatchError at 4: want json string, got null",
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

func TestStringer(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "basic",
			Test: typeTestCase[stringerWithUnmarshal, cj.StringerMarshaler[stringerWithUnmarshal], cj.StringUnmarshaler[*stringerWithUnmarshal]]{
				Value: stringerWithUnmarshal("hello"), JSON: "\"hello\"",
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

type stringerWithUnmarshal string

func (d stringerWithUnmarshal) String() string { return string(d) }
func (d *stringerWithUnmarshal) UnmarshalString(s string) error {
	*d = stringerWithUnmarshal(s)
	return nil
}

func TestText(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "text",
			Test: typeTestCase[uuid.UUID, cj.TextMarshaler[uuid.UUID], cj.TextUnmarshaler[*uuid.UUID]]{
				Value: must(uuid.ParseUUID("10101010-1010-1010-1010-101010101010")), JSON: "\"10101010-1010-1010-1010-101010101010\"",
			},
		},
		{
			Name: "map",
			Test: typeTestCase[map[uuid.UUID]string, cj.ComparableMapMarshaler[map[uuid.UUID]string, uuid.UUID, string, cj.TextMarshaler[uuid.UUID], cj.String[string]], cj.MapUnmarshaler[map[uuid.UUID]string, uuid.UUID, string, cj.TextUnmarshaler[*uuid.UUID], cj.String[string]]]{
				Value: map[uuid.UUID]string{
					must(uuid.ParseUUID("00101010-1010-1010-1010-101010101010")): "foo",
					must(uuid.ParseUUID("00202020-2020-2020-2020-202020202020")): "bar",
				},
				JSON: `{"00101010-1010-1010-1010-101010101010":"foo","00202020-2020-2020-2020-202020202020":"bar"}`,
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
