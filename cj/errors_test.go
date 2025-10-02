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
	"errors"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyntaxError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`invalid`))

	err := cj.NewSyntaxError(dec, "bad syntax")
	assert.EqualError(t, err, "SyntaxError at offset 0: bad syntax")

	// Verify it implements the Error interface
	var cjErr cj.Error = err
	assert.NotNil(t, cjErr)
}

func TestSyntaxErrorWithCause(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`malformed`))
	cause := errors.New("underlying error")

	err := cj.WrapSyntaxError(dec, "wrapped syntax error", cause)
	assert.EqualError(t, err, "SyntaxError at offset 0: wrapped syntax error: underlying error")

	assert.Equal(t, cause, err.Cause())
	assert.Equal(t, cause, err.Unwrap())
}

func TestKindMismatchError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`123`))

	err := cj.NewKindMismatchError(dec, '0', "string")
	assert.EqualError(t, err, "KindMismatchError at offset 0: want string, got number")
}

func TestInvalidValueError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`"invalid"`))
	cause := errors.New("value validation failed")

	err := cj.NewInvalidValueError(dec, "invalid bearer token", cause)
	assert.EqualError(t, err, "InvalidValueError at offset 0: invalid bearer token: value validation failed")
	assert.Equal(t, cause, err.Cause())
}

func TestUnmarshalFieldError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"field": bad}`))
	cause := errors.New("field unmarshal failed")

	err := cj.NewUnmarshalFieldError(dec, "Person.name", cause)
	assert.EqualError(t, err, "UnmarshalFieldError at offset 0: Person.name: field unmarshal failed")
	assert.Equal(t, cause, err.Cause())
}

func TestMissingFieldsError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{}`))

	err := cj.NewMissingFieldsError(dec, "Person", []string{"name", "age"})
	assert.EqualError(t, err, "MissingFieldsError at offset 0: type Person missing required fields: [name age]")
}

func TestUnknownFieldsError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"extra": "field"}`))

	err := cj.NewUnknownFieldsError(dec, "Person", []string{"extra", "unknown"})
	assert.EqualError(t, err, "UnknownFieldsError at offset 0: type Person has unknown fields: [extra unknown]")
}

func TestDuplicateFieldKeyError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"name":"John","name":"Jane"}`))

	err := cj.NewDuplicateFieldKeyError(dec, "Person.name")
	assert.EqualError(t, err, "DuplicateFieldKeyError at offset 0: field Person.name duplicated")
}

func TestDuplicateMapKeyError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"1":1,"01":2}`))

	err := cj.NewDuplicateMapKeyError(dec, "map[int]int")
	assert.EqualError(t, err, "DuplicateMapKeyError at offset 0: type map[int]int has duplicate map keys")
}

func TestErrorStackTraces(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`invalid`))

	err := cj.NewSyntaxError(dec, "test error")

	// Verify stack trace exists
	stack := err.StackTrace()
	assert.NotNil(t, stack)
	assert.NotEmpty(t, stack)
}

func TestErrorContextPreservation(t *testing.T) {
	// Test that errors preserve JSON context like byte offsets
	jsonInput := `{"field": invalid_value}`
	dec := jsontext.NewDecoder(strings.NewReader(jsonInput))

	// Advance decoder to a specific position
	_, _ = dec.ReadToken() // Read opening brace
	_, _ = dec.ReadToken() // Read "field" key

	err := cj.NewKindMismatchError(dec, '"', "number")

	// Error should contain position information
	errStr := err.Error()
	assert.Contains(t, errStr, "at ")
	assert.Contains(t, errStr, "KindMismatchError")
}

func TestErrorIntegrationWithTypeDecoder(t *testing.T) {
	// Test that actual type decoders produce expected errors
	t.Run("wrong_type_for_string", func(t *testing.T) {
		var result string
		err := json.Unmarshal([]byte(`123`), cj.NewUnmarshalerFrom[string, cj.String[string]](&result))
		require.Error(t, err)

		// Should be a KindMismatchError
		assert.Contains(t, err.Error(), "KindMismatchError")
		assert.Contains(t, err.Error(), "want json string")
	})

	t.Run("invalid_base64", func(t *testing.T) {
		var result []byte
		err := json.Unmarshal([]byte(`"not_base64!@#"`), cj.NewUnmarshalerFrom[[]byte, cj.Binary[[]byte]](&result))
		require.Error(t, err)

		// Should be an InvalidValueError
		assert.Contains(t, err.Error(), "InvalidValueError")
	})

	t.Run("invalid_integer", func(t *testing.T) {
		var result int
		err := json.Unmarshal([]byte(`"not_a_number"`), cj.NewUnmarshalerFrom[int, cj.Int32[int]](&result))
		require.Error(t, err)

		// Should be a KindMismatchError
		assert.Contains(t, err.Error(), "KindMismatchError")
		assert.Contains(t, err.Error(), "want json int")
	})
}
