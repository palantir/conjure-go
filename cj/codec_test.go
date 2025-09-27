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

	"github.com/palantir/conjure-go/v6/cj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientEncoder(t *testing.T) {
	encoder := cj.ClientEncoder[string, cj.String[string]]{}

	t.Run("ContentType", func(t *testing.T) {
		assert.Equal(t, "application/json", encoder.ContentType())
	})

	t.Run("Encode_value", func(t *testing.T) {
		var buf bytes.Buffer
		err := encoder.Encode(&buf, "hello")
		require.NoError(t, err)
		assert.Equal(t, "\"hello\"", buf.String())
	})

	t.Run("Encode_pointer", func(t *testing.T) {
		var buf bytes.Buffer
		value := "world"
		err := encoder.Encode(&buf, &value)
		require.NoError(t, err)
		assert.Equal(t, "\"world\"", buf.String())
	})

	t.Run("Encode_nil_pointer", func(t *testing.T) {
		var buf bytes.Buffer
		var nilPtr *string
		err := encoder.Encode(&buf, nilPtr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encode source should not be nil")
	})

	t.Run("Encode_wrong_type", func(t *testing.T) {
		var buf bytes.Buffer
		err := encoder.Encode(&buf, 123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encode source is incompatible with encoder")
	})

	t.Run("Marshal_value", func(t *testing.T) {
		data, err := encoder.Marshal("test")
		require.NoError(t, err)
		assert.Equal(t, `"test"`, string(data))
	})

	t.Run("Marshal_pointer", func(t *testing.T) {
		value := "test"
		data, err := encoder.Marshal(&value)
		require.NoError(t, err)
		assert.Equal(t, `"test"`, string(data))
	})

	t.Run("Marshal_wrong_type", func(t *testing.T) {
		_, err := encoder.Marshal(42)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encode source is incompatible with encoder")
	})
}

func TestClientDecoder(t *testing.T) {
	decoder := cj.ClientDecoder[string, cj.String[string]]{}

	t.Run("Accept", func(t *testing.T) {
		assert.Equal(t, "application/json", decoder.Accept())
	})

	t.Run("Decode_to_value", func(t *testing.T) {
		var result string
		err := decoder.Decode(strings.NewReader(`"hello"`), &result)
		require.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("Decode_to_pointer", func(t *testing.T) {
		var resultPtr *string
		err := decoder.Decode(strings.NewReader(`"world"`), &resultPtr)
		require.NoError(t, err)
		require.NotNil(t, resultPtr)
		assert.Equal(t, "world", *resultPtr)
	})

	t.Run("Decode_nil_target", func(t *testing.T) {
		err := decoder.Decode(strings.NewReader(`"test"`), nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode target cannot be nil")
	})

	t.Run("Decode_wrong_type", func(t *testing.T) {
		var wrongType int
		err := decoder.Decode(strings.NewReader(`"test"`), &wrongType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode target is incompatible with decoder")
	})

	t.Run("Unmarshal_value", func(t *testing.T) {
		var result string
		err := decoder.Unmarshal([]byte(`"test"`), &result)
		require.NoError(t, err)
		assert.Equal(t, "test", result)
	})

	t.Run("Unmarshal_pointer", func(t *testing.T) {
		var resultPtr *string
		err := decoder.Unmarshal([]byte(`"test"`), &resultPtr)
		require.NoError(t, err)
		require.NotNil(t, resultPtr)
		assert.Equal(t, "test", *resultPtr)
	})

	t.Run("Unmarshal_wrong_type", func(t *testing.T) {
		var wrongType float64
		err := decoder.Unmarshal([]byte(`"test"`), &wrongType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode target is incompatible with decoder")
	})
}

func TestClientCodecIntegration(t *testing.T) {
	encoder := cj.ClientEncoder[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]]{}
	decoder := cj.ClientDecoder[[]int, cj.ListUnmarshaler[[]int, int, cj.Int32[int]]]{}

	original := []int{1, 2, 3, 42}

	// Encode
	data, err := encoder.Marshal(original)
	require.NoError(t, err)
	assert.Equal(t, "[1,2,3,42]", string(data))

	// Decode
	var result []int
	err = decoder.Unmarshal(data, &result)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}
