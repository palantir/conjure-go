// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package errors_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	wparams "github.com/palantir/witchcraft-go-params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v5/integration_test/testgenerated/errors/api"
)

var _ errors.Error = &api.MyNotFound{}
var _ json.Marshaler = &api.MyNotFound{}
var _ json.Unmarshaler = &api.MyNotFound{}
var _ wparams.ParamStorer = &api.MyNotFound{}

var testError = api.NewMyNotFound(
	api.Basic{
		Data: "some data",
	},
	[]int{1, 2, 3},
	"type",
	"something",
	nil,
)

var testJSON = fmt.Sprintf(`{
  "errorCode": "NOT_FOUND",
  "errorName": "MyNamespace:MyNotFound",
  "errorInstanceId": "%s",
  "parameters": {
    "safeArgA": {
      "data": "some data"
    },
    "safeArgB": [
      1,
      2,
      3
    ],
    "type": "type",
    "unsafeArgA": "something",
    "unsafeArgB": null
  }
}`, testError.InstanceID())

var testErrorInternal = api.NewMyInternal(
	api.Basic{
		Data: "some data",
	},
	[]int{1, 2, 3},
	"type",
	"something",
	nil,
)

var testJSONInternal = fmt.Sprintf(`{
  "errorCode": "INTERNAL",
  "errorName": "MyNamespace:MyInternal",
  "errorInstanceId": "%s",
  "parameters": {
    "safeArgA": {
      "data": "some data"
    },
    "safeArgB": [
      1,
      2,
      3
    ],
    "type": "type",
    "unsafeArgA": "something",
    "unsafeArgB": null
  }
}`, testErrorInternal.InstanceID())

func TestError_ErrorMethods(t *testing.T) {
	assert.Equal(t, errors.NotFound, testError.Code())
	assert.Equal(t, "MyNamespace:MyNotFound", testError.Name())
	assert.NotNil(t, testError.InstanceID())
	assert.Equal(t, map[string]interface{}{
		"safeArgA":   testError.SafeArgA,
		"safeArgB":   testError.SafeArgB,
		"type":       testError.Type,
		"unsafeArgA": testError.UnsafeArgA,
		"unsafeArgB": testError.UnsafeArgB,
	}, testError.Parameters())
}

func TestError_MarshalJSON(t *testing.T) {
	bytes, err := json.MarshalIndent(testError, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, testJSON, string(bytes))
}

func TestError_UnmarshalJSON(t *testing.T) {
	var myNotFound api.MyNotFound
	err := json.Unmarshal([]byte(testJSON), &myNotFound)
	assert.NoError(t, err)
	assert.Equal(t, testError, &myNotFound)
}

func TestError_SafeParams(t *testing.T) {
	safeParams := testError.SafeParams()
	for _, key := range []string{"safeArgA", "safeArgB", "type", "errorInstanceId"} {
		assert.Contains(t, safeParams, key)
	}
}

func TestError_UnsafeParams(t *testing.T) {
	unsafeParams := testError.UnsafeParams()
	for _, key := range []string{"unsafeArgA", "unsafeArgB"} {
		assert.Contains(t, unsafeParams, key)
	}
}

func TestError_Init(t *testing.T) {
	genericErr, err := errors.UnmarshalError([]byte(testJSON))
	assert.NoError(t, err)
	myNotFoundErr, ok := genericErr.(*api.MyNotFound)
	require.True(t, ok)
	assert.Equal(t, myNotFoundErr, testError)

	genericErr, err = errors.UnmarshalError([]byte(testJSONInternal))
	assert.NoError(t, err)
	myInternalErr, ok := genericErr.(*api.MyInternal)
	require.True(t, ok)
	assert.Equal(t, myInternalErr, testErrorInternal)
}
