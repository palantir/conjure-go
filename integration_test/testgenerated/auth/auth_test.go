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

package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/palantir/conjure-go-runtime/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-server/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/integration_test/testgenerated/auth/api"
)

const (
	headerAuthAccepted = "header: Authorization accepted"
	headerAuthInvalid  = "header: Invalid auth"
	testJWT            = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0ifQ.JTD36MhcwmSuvfdCkfSYc-LHOGNA1UQ-0FKLKqdXbF4`
)

func TestBothAuthClient(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client := api.NewBothAuthServiceClient(newHTTPClient(t, server.URL))
	authClient := api.NewBothAuthServiceClientWithAuth(client, "Bearer "+testJWT, testJWT)

	// test header auth calls
	resp, err := client.Default(ctx, "Bearer "+testJWT)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
	resp, err = authClient.Default(ctx)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)

	// test invalid auth
	_, err = client.Default(ctx, "invalid token")
	assert.EqualError(t, err, "httpclient request failed: server returned a status >= 400")

	// test cookie auth calls
	err = client.Cookie(ctx, testJWT)
	require.NoError(t, err)
	err = authClient.Cookie(ctx)
	require.NoError(t, err)

	// test none auth calls
	err = client.None(ctx)
	require.NoError(t, err)
	err = authClient.None(ctx)
	require.NoError(t, err)
}

func TestHeaderAuthClient(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client := api.NewHeaderAuthServiceClient(newHTTPClient(t, server.URL))
	authClient := api.NewHeaderAuthServiceClientWithAuth(client, "Bearer "+testJWT)

	// test header auth calls
	resp, err := client.Default(ctx, "Bearer "+testJWT)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
	resp, err = authClient.Default(ctx)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
}

func TestCookieAuthClient(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client := api.NewCookieAuthServiceClient(newHTTPClient(t, server.URL))
	authClient := api.NewCookieAuthServiceClientWithAuth(client, testJWT)

	// test cookie auth calls
	err := client.Cookie(ctx, testJWT)
	require.NoError(t, err)
	err = authClient.Cookie(ctx)
	require.NoError(t, err)
}

func createTestServer() *httptest.Server {
	r := httprouter.New()
	r.GET("/default", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		authVal := req.Header.Get("Authorization")
		if !strings.HasPrefix(authVal, "Bearer ") {
			rest.WriteJSONResponse(rw, headerAuthInvalid, http.StatusUnauthorized)
			return
		}
		authContent := strings.TrimPrefix(authVal, "Bearer ")
		if authContent != testJWT {
			rest.WriteJSONResponse(rw, headerAuthInvalid, http.StatusUnauthorized)
			return
		}
		rest.WriteJSONResponse(rw, headerAuthAccepted, http.StatusOK)
	}))
	r.GET("/cookie", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		authVal := req.Header.Get("P_TOKEN")
		if authVal != testJWT {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}))
	r.GET("/none", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		rw.WriteHeader(http.StatusOK)
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
