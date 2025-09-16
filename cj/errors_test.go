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

	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/cj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyntaxError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`invalid`))
	
	err := cj.NewSyntaxError(dec, "bad syntax")
	
	assert.Contains(t, err.Error(), "SyntaxError")
	assert.Contains(t, err.Error(), "bad syntax")
	assert.Contains(t, err.Error(), "at 0")
	
	// Verify it implements the Error interface
	var cjErr cj.Error = err
	assert.NotNil(t, cjErr)
}

func TestSyntaxErrorWithCause(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`malformed`))
	cause := errors.New("underlying error")
	
	err := cj.WrapSyntaxError(dec, "wrapped syntax error", cause)
	
	assert.Contains(t, err.Error(), "SyntaxError")
	assert.Contains(t, err.Error(), "wrapped syntax error")
	assert.Contains(t, err.Error(), "underlying error")
	assert.Equal(t, cause, err.Cause())
	assert.Equal(t, cause, err.Unwrap())
}

func TestKindMismatchError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`123`))
	
	err := cj.NewKindMismatchError(dec, '0', "string")
	
	assert.Contains(t, err.Error(), "KindMismatchError")
	assert.Contains(t, err.Error(), "want string, got number")
	assert.Contains(t, err.Error(), "at 0")
}

func TestInvalidValueError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`"invalid"`))
	cause := errors.New("value validation failed")
	
	err := cj.NewInvalidValueError(dec, "invalid bearer token", cause)
	
	assert.Contains(t, err.Error(), "InvalidValueError")
	assert.Contains(t, err.Error(), "invalid bearer token")
	assert.Contains(t, err.Error(), "value validation failed")
	assert.Equal(t, cause, err.Cause())
}

func TestUnmarshalFieldError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"field": bad}`))
	cause := errors.New("field unmarshal failed")
	
	err := cj.NewUnmarshalFieldError(dec, "Person.name", cause)
	
	assert.Contains(t, err.Error(), "UnmarshalFieldError")
	assert.Contains(t, err.Error(), "Person.name")
	assert.Contains(t, err.Error(), "field unmarshal failed")
	assert.Equal(t, cause, err.Cause())
}

func TestMissingFieldsError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{}`))
	
	err := cj.NewMissingFieldsError(dec, "Person", []string{"name", "age"})
	
	assert.Contains(t, err.Error(), "MissingFieldsError")
	assert.Contains(t, err.Error(), "type Person missing required fields: [name age]")
	assert.Contains(t, err.Error(), "at 0")
}

func TestUnknownFieldsError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"extra": "field"}`))
	
	err := cj.NewUnknownFieldsError(dec, "Person", []string{"extra", "unknown"})
	
	assert.Contains(t, err.Error(), "UnknownFieldsError")
	assert.Contains(t, err.Error(), "type Person has unknown fields: [extra unknown]")
	assert.Contains(t, err.Error(), "at 0")
}

func TestDuplicateFieldKeyError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"name":"John","name":"Jane"}`))
	
	err := cj.NewDuplicateFieldKeyError(dec, "Person.name")
	
	assert.Contains(t, err.Error(), "DuplicateFieldKeyError")
	assert.Contains(t, err.Error(), "field Person.name duplicated")
	assert.Contains(t, err.Error(), "at 0")
}

func TestDuplicateMapKeyError(t *testing.T) {
	dec := jsontext.NewDecoder(strings.NewReader(`{"1":1,"01":2}`))
	
	err := cj.NewDuplicateMapKeyError(dec, "map[int]int")
	
	assert.Contains(t, err.Error(), "DuplicateMapKeyError")
	assert.Contains(t, err.Error(), "type map[int]int has duplicate map keys")
	assert.Contains(t, err.Error(), "at 0")
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
		err := cj.Unmarshal[string, cj.String[string]]([]byte("123"), &result)
		require.Error(t, err)
		
		// Should be a KindMismatchError
		assert.Contains(t, err.Error(), "KindMismatchError")
		assert.Contains(t, err.Error(), "want json string")
	})

	t.Run("invalid_base64", func(t *testing.T) {
		var result []byte
		err := cj.Unmarshal[[]byte, cj.Binary[[]byte]]([]byte(`"not_base64!@#"`), &result)
		require.Error(t, err)
		
		// Should be an InvalidValueError
		assert.Contains(t, err.Error(), "InvalidValueError")
	})

	t.Run("invalid_integer", func(t *testing.T) {
		var result int
		err := cj.Unmarshal[int, cj.Int32[int]]([]byte(`"not_a_number"`), &result)
		require.Error(t, err)
		
		// Should be a KindMismatchError
		assert.Contains(t, err.Error(), "KindMismatchError")
		assert.Contains(t, err.Error(), "want json int")
	})
}