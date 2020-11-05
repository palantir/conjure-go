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

package client_test

import (
	"context"
	"crypto/rand"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/witchcraft-go-server/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/client/api"
)

func TestHeaderParams(t *testing.T) {
	const (
		customKey   = "Custom-Key"
		customValue = "customValue"
	)

	called := false
	r := httprouter.New()
	r.GET("/echo", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		called = true
		assert.Equal(t, customValue, req.Header.Get(customKey))
		rw.WriteHeader(http.StatusOK)
	}))
	server := httptest.NewServer(r)
	defer server.Close()

	httpClient, err := httpclient.NewClient(
		httpclient.WithUserAgent("TestNewRequest"),
		httpclient.WithBaseURLs([]string{server.URL}),
		httpclient.WithSetHeader(customKey, customValue),
	)
	require.NoError(t, err)

	client := api.NewTestServiceClient(httpClient)
	err = client.Echo(context.Background())
	require.NoError(t, err)
	assert.True(t, called)
}

func TestPathParam(t *testing.T) {
	const (
		wantParam        = "var/conf/install.yml"
		wantEncodedParam = "var%2Fconf%2Finstall.yml"
	)

	called := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
		paramVal := strings.TrimPrefix(req.RequestURI, "/path/")
		assert.Equal(t, wantEncodedParam, paramVal)

		unescaped, err := url.PathUnescape(paramVal)
		require.NoError(t, err)
		assert.Equal(t, wantParam, unescaped)

		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	err := client.PathParam(context.Background(), wantParam)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestPathParamRid(t *testing.T) {
	const wantParam = "ri.service.instance.rtype.locator"

	called := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
		paramVal := strings.TrimPrefix(req.RequestURI, "/path/rid/")

		unescaped, err := url.PathUnescape(paramVal)
		require.NoError(t, err)
		assert.Equal(t, wantParam, unescaped)

		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ridValue := rid.MustNew("service", "instance", "rtype", "locator")
	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	err := client.PathParamRid(context.Background(), ridValue)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestPathParamRidAlias(t *testing.T) {
	const wantParam = "ri.service.instance.rtype.locator"

	called := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
		paramVal := strings.TrimPrefix(req.RequestURI, "/path/rid/alias/")

		unescaped, err := url.PathUnescape(paramVal)
		require.NoError(t, err)
		assert.Equal(t, wantParam, unescaped)

		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ridValue := api.RidAlias(rid.MustNew("service", "instance", "rtype", "locator"))
	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	err := client.PathParamRidAlias(context.Background(), ridValue)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestBytes(t *testing.T) {
	called := false
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	want := api.CustomObject{
		Data: bytes,
	}
	r := httprouter.New()
	r.GET("/bytes", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		called = true
		rest.WriteJSONResponse(rw, want, http.StatusOK)
	}))
	server := httptest.NewServer(r)
	defer server.Close()

	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	got, err := client.Bytes(context.Background())
	require.NoError(t, err)
	assert.True(t, called)

	assert.Equal(t, want, got)
}

func TestBinary(t *testing.T) {
	called := false
	want := make([]byte, 10)
	_, err := rand.Read(want)
	require.NoError(t, err)
	r := httprouter.New()
	r.GET("/binary", httprouter.Handle(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		called = true
		_, err = rw.Write(want)
		require.NoError(t, err)
	}))
	server := httptest.NewServer(r)
	defer server.Close()

	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	rc, err := client.Binary(context.Background())
	require.NoError(t, err)
	assert.True(t, called)

	got, err := ioutil.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}
