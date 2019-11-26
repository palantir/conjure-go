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

package queryparam_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/palantir/conjure-go-runtime/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-server/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v4/integration_test/testgenerated/queryparam/api"
)

func TestQueryParamClient(t *testing.T) {
	server := createTestServer()
	defer server.Close()
	optionalStr := "optional"

	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	resp, err := client.Echo(context.Background(), "hello", 3, nil, &optionalStr)
	require.NoError(t, err)
	assert.Equal(t, "hello hello hello", resp)

	_, err = client.Echo(context.Background(), "hello", -3, &optionalStr, nil)
	assert.EqualError(t, err, "httpclient request failed: server returned a status >= 400")
}

func createTestServer() *httptest.Server {
	r := httprouter.New()
	r.GET("/echo", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		input := req.URL.Query().Get("input")
		reps, err := strconv.Atoi(req.URL.Query().Get("reps"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if reps < 0 {
			http.Error(rw, fmt.Sprintf("reps must be non-negative, was %d", reps), http.StatusBadRequest)
			return
		}
		var parts []string
		for i := 0; i < reps; i++ {
			parts = append(parts, input)
		}
		rest.WriteJSONResponse(rw, strings.Join(parts, " "), http.StatusOK)
	}))
	server := httptest.NewServer(r)
	return server
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}
