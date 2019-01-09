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

package errors

import (
	"bytes"
	"net/http"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/codecs"
)

// ErrorFromResponse extract serializable error from the given response.
//
// TODO This function is subject to change.
func ErrorFromResponse(response *http.Response) (Error, error) {
	var unmarshalled SerializableError
	if err := codecs.JSON.Decode(response.Body, &unmarshalled); err != nil {
		return nil, err
	}
	errorType, err := NewErrorType(unmarshalled.ErrorCode, unmarshalled.ErrorName)
	if err != nil {
		return nil, err
	}

	gErr := &genericError{
		errorType:       errorType,
		errorInstanceID: unmarshalled.ErrorInstanceID,
	}

	// best effort; on failure we continue without params
	_ = codecs.JSON.Decode(bytes.NewReader(unmarshalled.Parameters), &gErr.parameterizer)

	return gErr, nil
}

// WriteErrorResponse writes error to the response writer.
//
// TODO This function is subject to change.
func WriteErrorResponse(w http.ResponseWriter, e Error) {
	se := SerializableError{
		ErrorCode:       e.Code(),
		ErrorName:       e.Name(),
		ErrorInstanceID: e.InstanceID(),
	}
	// If the parameters fail to marshal, we will send the rest without params.
	// The other fields are primitives that should always successfully marshal.
	se.Parameters, _ = codecs.JSON.Marshal(e.Parameters())

	w.Header().Add("Content-Type", codecs.JSON.ContentType())
	w.WriteHeader(e.Code().StatusCode())
	_ = codecs.JSON.Encode(w, e) // There is nothing we can do on write failure.
}
