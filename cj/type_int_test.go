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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/palantir/conjure-go/v6/cj"
	"github.com/palantir/pkg/safelong"
)

func TestInt(t *testing.T) {
	for name, test := range map[string]typeTest{
		"zero": typeTestCase[int, cj.Int32[int], cj.Int32[int]]{
			Value: 0, JSON: "0",
		},
		"positive": typeTestCase[int, cj.Int32[int], cj.Int32[int]]{
			Value: 42, JSON: "42",
		},
		"negative": typeTestCase[int, cj.Int32[int], cj.Int32[int]]{
			Value: -7, JSON: "-7",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

func TestSafeLong(t *testing.T) {
	const maxSafe = 9007199254740991
	const minSafe = -9007199254740991
	for name, test := range map[string]typeTest{
		"zero": typeTestCase[int64, cj.SafeLong[int64], cj.SafeLong[int64]]{
			Value: 0, JSON: "0",
		},
		"positive": typeTestCase[int64, cj.SafeLong[int64], cj.SafeLong[int64]]{
			Value: maxSafe, JSON: "9007199254740991",
		},
		"negative": typeTestCase[int64, cj.SafeLong[int64], cj.SafeLong[int64]]{
			Value: minSafe, JSON: "-9007199254740991",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

func TestMapOfInt(t *testing.T) {
	for name, test := range map[string]typeTest{
		"empty": typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
			Value: map[string]int{}, JSON: "{}",
		},
		"simple": typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
			Value: map[string]int{"a": 1, "b": -2}, JSON: "{\"a\":1,\"b\":-2}",
		},
		"null": typeTestCase[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]], cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]]{
			JSON: "null", SkipTestMarshal: true, Value: map[string]int{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

func TestMapOfSafeLong(t *testing.T) {
	type SL = safelong.SafeLong
	for name, test := range map[string]typeTest{
		"empty": typeTestCase[map[string]SL, cj.OrderedMapMarshaler[map[string]SL, string, SL, cj.String[string], cj.SafeLong[SL]], cj.MapUnmarshaler[map[string]SL, string, SL, cj.String[string], cj.SafeLong[SL]]]{
			Value: map[string]SL{}, JSON: "{}",
		},
		"simple": typeTestCase[map[string]SL, cj.OrderedMapMarshaler[map[string]SL, string, SL, cj.String[string], cj.SafeLong[SL]], cj.MapUnmarshaler[map[string]SL, string, SL, cj.String[string], cj.SafeLong[SL]]]{
			Value: map[string]SL{"a": 42, "b": -42}, JSON: "{\"a\":42,\"b\":-42}",
		},
		"null": typeTestCase[map[string]SL, cj.OrderedMapMarshaler[map[string]SL, string, SL, cj.String[string], cj.SafeLong[SL]], cj.MapUnmarshaler[map[string]SL, string, SL, cj.String[string], cj.SafeLong[SL]]]{
			JSON: "null", SkipTestMarshal: true, Value: map[string]SL{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

func jsonEqual(a, b string) bool {
	var o1, o2 any
	if err := json.Unmarshal([]byte(a), &o1); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &o2); err != nil {
		return false
	}
	return reflect.DeepEqual(o1, o2)
}

func TestMapIntKeySafeLongValue(t *testing.T) {
	type SL = safelong.SafeLong
	for name, test := range map[string]typeTest{
		"simple": typeTestCase[map[int]SL, cj.OrderedMapMarshaler[map[int]SL, int, SL, cj.Int32MapKey[int], cj.SafeLong[SL]], cj.MapUnmarshaler[map[int]SL, int, SL, cj.Int32MapKey[int], cj.SafeLong[SL]]]{
			Value: map[int]SL{1: 100, -2: -200}, JSON: "{\"-2\":-200,\"1\":100}",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}

func TestMapSafeLongKeyIntValue(t *testing.T) {
	type SL = safelong.SafeLong
	for name, test := range map[string]typeTest{
		"simple": typeTestCase[map[SL]int, cj.OrderedMapMarshaler[map[SL]int, SL, int, cj.SafeLongMapKey[SL], cj.Int32[int]], cj.MapUnmarshaler[map[SL]int, SL, int, cj.SafeLongMapKey[SL], cj.Int32[int]]]{
			Value: map[SL]int{100: 1, -200: -2}, JSON: "{\"-200\":-2,\"100\":1}",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Marshal", test.TestMarshal)
			t.Run("Unmarshal", test.TestUnmarshal)
		})
	}
}
