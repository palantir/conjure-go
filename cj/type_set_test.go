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
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/uuid"
)

func TestSet(t *testing.T) {
	tests := []struct {
		Name string
		Test typeTest
	}{
		{
			Name: "empty",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{}, JSON: "[]",
			},
		},
		{
			Name: "one",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{42}, JSON: "[42]",
			},
		},
		{
			Name: "several",
			Test: typeTestCase[[]string, cj.SetMarshaler[[]string, string, cj.String[string]], cj.SetUnmarshaler[[]string, string, cj.String[string]]]{
				Value: []string{"a", "b", "c"}, JSON: "[\"a\",\"b\",\"c\"]",
			},
		},
		{
			Name: "null_marshal",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				JSON: "[]", SkipTestUnmarshal: true, Value: []int(nil),
			},
		},
		{
			Name: "null_unmarshal",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				JSON: "null", SkipTestMarshal: true, Value: []int{},
			},
		},
		{
			Name: "comparable",
			Test: typeTestCase[[]uuid.UUID, cj.SetMarshaler[[]uuid.UUID, uuid.UUID, cj.UUID[uuid.UUID]], cj.SetUnmarshaler[[]uuid.UUID, uuid.UUID, cj.UUID[uuid.UUID]]]{
				Value: []uuid.UUID{must(uuid.ParseUUID("10101010-1010-1010-1010-101010101010")), must(uuid.ParseUUID("20202020-2020-2020-2020-202020202020"))}, JSON: "[\"10101010-1010-1010-1010-101010101010\",\"20202020-2020-2020-2020-202020202020\"]",
			},
		},
		{
			Name: "marshal_ordered_dedupe",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{5, 4, 3, 2, 1, 2, 3, 4, 5}, JSON: "[5,4,3,2,1]", SkipTestUnmarshal: true,
			},
		},
		{
			Name: "marshal_dupes",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{42, 42, 1}, JSON: "[42,1]", SkipTestUnmarshal: true,
			},
		},
		{
			Name: "ordered",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{5, 4, 3, 2, 1}, JSON: "[5,4,3,2,1]",
			},
		},
		{
			Name: "unmarshal_dupes",
			Test: typeTestCase[[]int, cj.SetMarshaler[[]int, int, cj.Int32[int]], cj.SetUnmarshaler[[]int, int, cj.Int32[int]]]{
				Value: []int{1, 42}, JSON: "[42, 42, 1]", SkipTestMarshal: true, ErrUnmarshalJSONFrom: "DuplicateSetItemError at 7: type int has a duplicate set item at index 1",
			},
		},
		{
			Name: "set_of_structs",
			Test: typeTestCase[[]testStruct, cj.SetMarshaler[[]testStruct, testStruct, cj.StructMarshaler[testStruct]], cj.SetUnmarshaler[[]testStruct, testStruct, cj.StructUnmarshaler[*testStruct]]]{
				Value: []testStruct{{Name: "c", Age: 3}, {Name: "b", Age: 2}, {Name: "a", Age: 1}},
				JSON:  "[{\"name\":\"c\",\"age\":3},{\"name\":\"b\",\"age\":2},{\"name\":\"a\",\"age\":1}]",
			},
		},
		{
			Name: "set_of_structs_dedupes",
			Test: typeTestCase[[]testStruct, cj.SetMarshaler[[]testStruct, testStruct, cj.StructMarshaler[testStruct]], cj.SetUnmarshaler[[]testStruct, testStruct, cj.StructUnmarshaler[*testStruct]]]{
				Value:             []testStruct{{Name: "c", Age: 3}, {Name: "b", Age: 2}, {Name: "a", Age: 1}, {Name: "b", Age: 2}, {Name: "b", Age: 2}},
				JSON:              "[{\"name\":\"c\",\"age\":3},{\"name\":\"b\",\"age\":2},{\"name\":\"a\",\"age\":1}]",
				SkipTestUnmarshal: true,
			},
		},
		{
			Name: "set_of_structs_errors",
			Test: typeTestCase[[]testStruct, cj.SetMarshaler[[]testStruct, testStruct, cj.StructMarshaler[testStruct]], cj.SetUnmarshaler[[]testStruct, testStruct, cj.StructUnmarshaler[*testStruct]]]{
				JSON:                 "[{\"name\":\"c\",\"age\":3},{\"name\":\"b\",\"age\":2},{\"name\":\"a\",\"age\":1},{\"name\":\"b\",\"age\":2},{\"name\":\"b\",\"age\":2}]",
				ErrUnmarshalJSONFrom: "DuplicateSetItemError at 84: type testStruct has a duplicate set item at index 3",
				SkipTestMarshal:      true,
			},
		},
		{
			Name: "datetimes",
			Test: typeTestCase[[]datetime.DateTime, cj.SetMarshaler[[]datetime.DateTime, datetime.DateTime, cj.StringerMarshaler[datetime.DateTime]], cj.SetUnmarshaler[[]datetime.DateTime, datetime.DateTime, cj.TextUnmarshaler[*datetime.DateTime]]]{
				JSON:  "[\"2025-05-12T19:26:00Z\",\"2001-01-01T19:26:00Z\",\"0001-01-01T00:00:00Z\"]",
				Value: []datetime.DateTime{must(datetime.ParseDateTime("2025-05-12T19:26:00Z")), must(datetime.ParseDateTime("2001-01-01T19:26:00Z")), must(datetime.ParseDateTime("0001-01-01T00:00:00Z"))},
			},
		},
		{
			Name: "datetimes_errors",
			Test: typeTestCase[[]datetime.DateTime, cj.SetMarshaler[[]datetime.DateTime, datetime.DateTime, cj.StringerMarshaler[datetime.DateTime]], cj.SetUnmarshaler[[]datetime.DateTime, datetime.DateTime, cj.TextUnmarshaler[*datetime.DateTime]]]{
				JSON:                 "[\"2025-05-12T19:26:00Z\",\"2001-01-01T19:26:00Z\",\"0001-01-01T00:00:00Z\",\"0001-01-01T00:00:00.00Z\"]",
				Value:                []datetime.DateTime{must(datetime.ParseDateTime("2025-05-12T19:26:00Z")), must(datetime.ParseDateTime("2001-01-01T19:26:00Z")), must(datetime.ParseDateTime("0001-01-01T00:00:00Z"))},
				ErrUnmarshalJSONFrom: "DuplicateSetItemError at 95: type DateTime has a duplicate set item at index 3",
				SkipTestMarshal:      true,
			},
		},
		{
			Name: "datetimes_dedupes",
			Test: typeTestCase[[]datetime.DateTime, cj.SetMarshaler[[]datetime.DateTime, datetime.DateTime, cj.StringerMarshaler[datetime.DateTime]], cj.SetUnmarshaler[[]datetime.DateTime, datetime.DateTime, cj.TextUnmarshaler[*datetime.DateTime]]]{
				JSON: "[\"2025-05-12T19:26:00Z\",\"2001-01-01T19:26:00Z\",\"0001-01-01T00:00:00Z\"]",
				Value: []datetime.DateTime{
					must(datetime.ParseDateTime("2025-05-12T19:26:00Z")),
					must(datetime.ParseDateTime("2001-01-01T19:26:00Z")),
					must(datetime.ParseDateTime("2001-01-01T19:26:00.000Z")),
					must(datetime.ParseDateTime("0001-01-01T00:00:00Z")),
					must(datetime.ParseDateTime("0001-01-01T00:00:00Z")),
				},
				SkipTestUnmarshal: true,
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
