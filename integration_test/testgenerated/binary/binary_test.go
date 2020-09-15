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
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/witchcraft-go-server/wrouter"
	"github.com/palantir/witchcraft-go-server/wrouter/whttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v5/integration_test/testgenerated/binary/api"
)

const randByteLen = 10

func TestBytes(t *testing.T) {
	randBytes := make([]byte, randByteLen)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)
	server := newBinaryServer(t)
	server.Start()
	defer server.Close()
	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))

	t.Run("Bytes", func(t *testing.T) {
		binAlias := api.BinaryAlias(randBytes)
		want := api.CustomObject{
			Data: randBytes,
			BinaryAlias: &binAlias,
		}
		got, err := client.Bytes(context.Background(), want)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("Binary", func(t *testing.T) {
		resp, err := client.Binary(context.Background(), func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(randBytes))
		})
		require.NoError(t, err)
		got, err := ioutil.ReadAll(resp)
		require.NoError(t, err)
		assert.Equal(t, randBytes, got)
	})
	t.Run("BinaryAlias", func(t *testing.T) {
		resp, err := client.BinaryAlias(context.Background(), func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(randBytes))
		})
		require.NoError(t, err)
		got, err := ioutil.ReadAll(resp)
		require.NoError(t, err)
		assert.Equal(t, randBytes, got)
	})
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}

func newBinaryServer(t *testing.T) *httptest.Server {
	router := wrouter.New(whttprouter.New())
	err := api.RegisterRoutesTestService(router, &binaryServer{})
	require.NoError(t, err)
	server := httptest.NewUnstartedServer(router)
	return server
}

type binaryServer struct {
}

func (b binaryServer) BinaryAlias(ctx context.Context, bodyArg io.ReadCloser) (io.ReadCloser, error) {
	body, err := ioutil.ReadAll(bodyArg)
	if err != nil {
		return nil, err
	}
	resp := ioutil.NopCloser(bytes.NewReader(body))
	return resp, nil
}

func (b binaryServer) Binary(ctx context.Context, bodyArg io.ReadCloser) (io.ReadCloser, error) {
	body, err := ioutil.ReadAll(bodyArg)
	if err != nil {
		return nil, err
	}
	resp := ioutil.NopCloser(bytes.NewReader(body))
	return resp, nil
}

func (b binaryServer) Bytes(ctx context.Context, bodyArg api.CustomObject) (api.CustomObject, error) {
	// verify we can decode binary
	if len(bodyArg.Data) != randByteLen {
		return api.CustomObject{}, errors.NewInvalidArgument()
	}
	if len(*bodyArg.BinaryAlias) != randByteLen {
		return api.CustomObject{}, errors.NewInvalidArgument()
	}
	return bodyArg, nil
}
