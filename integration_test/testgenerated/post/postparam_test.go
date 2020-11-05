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

package post_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-server/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/post/api"
)

func TestPostClient(t *testing.T) {
	r := httprouter.New()
	r.POST("/echo", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var param string
		bodyBytes, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)
		err = json.Unmarshal(bodyBytes, &param)
		require.NoError(t, err)

		rest.WriteJSONResponse(rw, param, http.StatusOK)
	}))
	server := httptest.NewServer(r)
	defer server.Close()

	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	resp, err := client.Echo(context.Background(), "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", resp)
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}
