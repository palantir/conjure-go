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
	"bytes"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Assert types implement interfaces so compiler flags incompatibilities here.
var (
	_ cj.TypeEncoder[any]    = cj.Any[any]{}
	_ cj.TypeEncoder[string] = cj.String[string]{}
	_ cj.TypeEncoder[int]    = cj.Int32[int]{}
	_ cj.TypeEncoder[int]    = cj.SafeLong[int]{}
	_ cj.TypeEncoder[[]int]  = cj.ListMarshaler[[]int, int, cj.Int32[int]]{}

	_ cj.TypeDecoder[any]    = cj.Any[any]{}
	_ cj.TypeDecoder[string] = cj.String[string]{}
	_ cj.TypeDecoder[int]    = cj.Int32[int]{}
	_ cj.TypeDecoder[int]    = cj.SafeLong[int]{}
)

func TestMarshal(t *testing.T) {
	t.Run("simple_string", func(t *testing.T) {
		data, err := cj.Marshal[string, cj.String[string]]("hello")
		require.NoError(t, err)
		assert.Equal(t, `"hello"`, string(data))
	})

	t.Run("list", func(t *testing.T) {
		data, err := cj.Marshal[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]]([]int{1, 2, 3})
		require.NoError(t, err)
		assert.Equal(t, "[1,2,3]", string(data))
	})

	t.Run("with_options", func(t *testing.T) {
		data, err := cj.Marshal[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]](
			nil, json.FormatNilSliceAsNull(true))
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})
}

func TestMarshalWrite(t *testing.T) {
	t.Run("simple_string", func(t *testing.T) {
		var buf bytes.Buffer
		err := cj.MarshalWrite[string, cj.String[string]](&buf, "world")
		require.NoError(t, err)
		assert.Equal(t, "\"world\"\n", buf.String())
	})

	t.Run("map", func(t *testing.T) {
		var buf bytes.Buffer
		value := map[string]int{"a": 1, "b": 2}
		err := cj.MarshalWrite[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](&buf, value)
		require.NoError(t, err)
		assert.Equal(t, "{\"a\":1,\"b\":2}\n", buf.String())
	})
}

func TestMarshalEncode(t *testing.T) {
	t.Run("simple_bool", func(t *testing.T) {
		var buf bytes.Buffer
		enc := jsontext.NewEncoder(&buf)
		err := cj.MarshalEncode[bool, cj.Boolean[bool]](enc, true)
		require.NoError(t, err)
		assert.Equal(t, "true\n", buf.String())
	})

	t.Run("float", func(t *testing.T) {
		var buf bytes.Buffer
		enc := jsontext.NewEncoder(&buf)
		err := cj.MarshalEncode[float64, cj.Float[float64]](enc, 123.456)
		require.NoError(t, err)
		assert.Equal(t, "123.456\n", buf.String())
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("simple_string", func(t *testing.T) {
		var result string
		err := cj.Unmarshal[string, cj.String[string]]([]byte(`"hello"`), &result)
		require.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("list", func(t *testing.T) {
		var result []int
		err := cj.Unmarshal[[]int, cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]([]byte("[1,2,3]"), &result)
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("invalid_json", func(t *testing.T) {
		var result string
		err := cj.Unmarshal[string, cj.String[string]]([]byte(`invalid`), &result)
		assert.Error(t, err)
	})
}

func TestUnmarshalRead(t *testing.T) {
	t.Run("simple_bool", func(t *testing.T) {
		var result bool
		err := cj.UnmarshalRead[bool, cj.Boolean[bool]](strings.NewReader("false"), &result)
		require.NoError(t, err)
		assert.Equal(t, false, result)
	})

	t.Run("map", func(t *testing.T) {
		var result map[string]int
		err := cj.UnmarshalRead[map[string]int, cj.MapUnmarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](
			strings.NewReader(`{"x":10,"y":20}`), &result)
		require.NoError(t, err)
		assert.Equal(t, map[string]int{"x": 10, "y": 20}, result)
	})
}

func TestUnmarshalDecode(t *testing.T) {
	t.Run("simple_float", func(t *testing.T) {
		dec := jsontext.NewDecoder(strings.NewReader("42.5"))
		var result float64
		err := cj.UnmarshalDecode[float64, cj.Float[float64]](dec, &result)
		require.NoError(t, err)
		assert.Equal(t, 42.5, result)
	})

	t.Run("optional", func(t *testing.T) {
		dec := jsontext.NewDecoder(strings.NewReader("null"))
		var result *string
		err := cj.UnmarshalDecode[*string, cj.OptionalUnmarshaler[*string, string, cj.String[string]]](dec, &result)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestRoundTrip(t *testing.T) {
	type testCase[T any, E cj.TypeEncoder[T], D cj.TypeDecoder[T]] struct {
		name  string
		value T
	}

	t.Run("string", func(t *testing.T) {
		tc := testCase[string, cj.String[string], cj.String[string]]{
			name:  "test_string",
			value: "hello world",
		}

		// Marshal
		data, err := cj.Marshal[string, cj.String[string]](tc.value)
		require.NoError(t, err)

		// Unmarshal
		var result string
		err = cj.Unmarshal[string, cj.String[string]](data, &result)
		require.NoError(t, err)
		assert.Equal(t, tc.value, result)
	})

	t.Run("complex_nested", func(t *testing.T) {
		type complexType = map[string][]int
		tc := testCase[complexType,
			cj.OrderedMapMarshaler[complexType, string, []int, cj.String[string], cj.ListMarshaler[[]int, int, cj.Int32[int]]],
			cj.MapUnmarshaler[complexType, string, []int, cj.String[string], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]]{
			name:  "nested_structure",
			value: map[string][]int{"nums": {1, 2, 3}, "empty": {}},
		}

		// Marshal
		data, err := cj.Marshal[complexType,
			cj.OrderedMapMarshaler[complexType, string, []int, cj.String[string], cj.ListMarshaler[[]int, int, cj.Int32[int]]]](tc.value)
		require.NoError(t, err)

		// Unmarshal
		var result complexType
		err = cj.Unmarshal[complexType,
			cj.MapUnmarshaler[complexType, string, []int, cj.String[string], cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]](data, &result)
		require.NoError(t, err)
		assert.Equal(t, tc.value, result)
	})
}
