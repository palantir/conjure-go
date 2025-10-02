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

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatNilSliceAsNull(t *testing.T) {
	t.Run("default_behavior_empty_array", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]](nil))
		require.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})

	t.Run("format_as_null_enabled", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]](nil), json.FormatNilSliceAsNull(true))
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("format_as_null_disabled", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]](nil), json.FormatNilSliceAsNull(false))
		require.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})

	t.Run("non_nil_slice_unaffected", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[[]int, cj.ListMarshaler[[]int, int, cj.Int32[int]]]([]int{}), json.FormatNilSliceAsNull(true))
		require.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})
}

func TestFormatNilMapAsNull(t *testing.T) {
	t.Run("default_behavior_empty_object", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](nil))
		require.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("format_as_null_enabled", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](nil), json.FormatNilMapAsNull(true))
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("format_as_null_disabled", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](nil), json.FormatNilMapAsNull(false))
		require.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("non_nil_map_unaffected", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](map[string]int{}), json.FormatNilMapAsNull(true))
		require.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})
}

func TestDeterministicOrdering(t *testing.T) {
	originalMap := map[string]int{"z": 1, "a": 2, "m": 3}

	t.Run("deterministic_by_default", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](originalMap))
		require.NoError(t, err)
		assert.Equal(t, `{"a":2,"m":3,"z":1}`, string(data))
	})

	t.Run("deterministic_explicitly_enabled", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](originalMap), json.Deterministic(true))
		require.NoError(t, err)
		assert.Equal(t, `{"a":2,"m":3,"z":1}`, string(data))
	})

	t.Run("deterministic_disabled", func(t *testing.T) {
		data, err := json.Marshal(cj.NewMarshalerTo[map[string]int, cj.OrderedMapMarshaler[map[string]int, string, int, cj.String[string], cj.Int32[int]]](originalMap), json.Deterministic(false))
		require.NoError(t, err)
		// Result should still be valid JSON, but order may vary
		// We can't assert exact order, but we can verify it parses correctly
		var result map[string]int
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		assert.Equal(t, originalMap, result)
	})
}

func TestRejectUnknownMembers(t *testing.T) {

	// We can't easily test RejectUnknownMembers with our current type system
	// since it's typically used in generated code, but we can test the concept
	// using the Any type and json options
	t.Run("strict_parsing_with_unknown_field", func(t *testing.T) {
		jsonWithUnknown := `{"name":"John","age":25,"unknown":"field"}`
		var result testStruct

		// This should work with lenient parsing (default)
		err := json.Unmarshal([]byte(jsonWithUnknown), &result)
		require.NoError(t, err)
		assert.Equal(t, "John", result.Name)
		assert.Equal(t, 25, result.Age)
	})

	t.Run("strict_parsing_rejects_unknown", func(t *testing.T) {
		jsonWithUnknown := `{"name":"John","age":25,"unknown":"field"}`
		var result testStruct

		// This should fail with strict parsing
		err := json.Unmarshal([]byte(jsonWithUnknown), &result, json.RejectUnknownMembers(true))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown")
	})

	t.Run("strict_parsing_accepts_known_fields", func(t *testing.T) {
		validJSON := `{"name":"John","age":25}`
		var result testStruct

		// This should work even with strict parsing
		err := json.Unmarshal([]byte(validJSON), &result, json.RejectUnknownMembers(true))
		require.NoError(t, err)
		assert.Equal(t, "John", result.Name)
		assert.Equal(t, 25, result.Age)
	})
}

// Create a simple struct-like type that visits object fields
type testStruct struct {
	Name string
	Age  int
}

func (t testStruct) MarshalJSON() ([]byte, error) {
	return json.Marshal(t, jsontext.AllowDuplicateNames(true))
}

func (t testStruct) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	{
		if err := enc.WriteToken(jsontext.String("name")); err != nil {
			return err
		}
		if err := cj.MarshalEncode[string, cj.String[string]](enc, t.Name); err != nil {
			return err
		}
	}
	{
		if err := enc.WriteToken(jsontext.String("age")); err != nil {
			return err
		}
		if err := cj.MarshalEncode[int, cj.Int32[int]](enc, t.Age); err != nil {
			return err
		}
	}
	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}
	return nil
}

func (t *testStruct) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *testStruct) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return cj.WrapSyntaxError(dec, "", err)
	}
	if kind := tok.Kind(); kind != '{' {
		return cj.NewKindMismatchError(dec, kind, "testStruct opening brace")
	}
	var seenName bool
	var seenAge bool
	var unknownMembers []string
	for {
		key, err := dec.ReadToken()
		if err != nil {
			return cj.WrapSyntaxError(dec, "", err)
		}
		kind := key.Kind()
		if kind == '}' {
			break
		}
		if kind != '"' {
			return cj.NewKindMismatchError(dec, kind, "ConjureDefinition closing brace or next key")
		}
		switch key.String() {
		case "name":
			if seenName {
				return cj.NewDuplicateFieldKeyError(dec, "testStruct[\"name\"]")
			}
			if err := cj.UnmarshalDecode[string, cj.String[string]](dec, &t.Name); err != nil {
				return err
			}
			seenName = true
		case "age":
			if seenAge {
				return cj.NewDuplicateFieldKeyError(dec, "testStruct[\"age\"]")
			}
			if err := cj.UnmarshalDecode[int, cj.Int32[int]](dec, &t.Age); err != nil {
				return err
			}
			seenAge = true
		default:
			unknownMembers = append(unknownMembers, key.String())
			if err := dec.SkipValue(); err != nil {
				return err
			}
		}
	}
	var missingFields []string
	if !seenName {
		missingFields = append(missingFields, "name")
	}
	if !seenAge {
		missingFields = append(missingFields, "age")
	}
	if len(missingFields) > 0 {
		return cj.NewMissingFieldsError(dec, "testStruct", missingFields)
	}
	if len(unknownMembers) > 0 {
		if strict, _ := json.GetOption(dec.Options(), json.RejectUnknownMembers); strict {
			return cj.NewUnknownFieldsError(dec, "testStruct", unknownMembers)
		}
	}
	return nil
}
