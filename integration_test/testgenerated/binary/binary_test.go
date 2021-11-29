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
	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/binary/api"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const randByteLen = 10

func TestBytes(t *testing.T) {
	ctx := context.Background()
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
			Data:        randBytes,
			BinaryAlias: &binAlias,
		}
		got, err := client.Bytes(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("Binary", func(t *testing.T) {
		resp, err := client.Binary(ctx, func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(randBytes))
		})
		require.NoError(t, err)
		got, err := ioutil.ReadAll(resp)
		require.NoError(t, err)
		assert.Equal(t, randBytes, got)
	})
	t.Run("BinaryAlias", func(t *testing.T) {
		resp, err := client.BinaryAlias(ctx, func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(randBytes))
		})
		require.NoError(t, err)
		got, err := ioutil.ReadAll(resp)
		require.NoError(t, err)
		assert.Equal(t, randBytes, got)
	})
	t.Run("BinaryAliasOptional", func(t *testing.T) {
		resp, err := client.BinaryAliasOptional(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
		got, err := ioutil.ReadAll(*resp)
		require.NoError(t, err)
		assert.Len(t, got, randByteLen)
	})
	t.Run("BinaryAliasAlias", func(t *testing.T) {
		resp, err := client.BinaryAliasAlias(ctx, func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(randBytes))
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		got, err := ioutil.ReadAll(*resp)
		require.NoError(t, err)
		assert.Equal(t, randBytes, got)
	})
	t.Run("BinaryOptional", func(t *testing.T) {
		resp, err := client.BinaryOptional(ctx)
		require.NoError(t, err)
		require.Nil(t, resp)
	})
	t.Run("BinaryOptionalAlias", func(t *testing.T) {
		resp, err := client.BinaryOptionalAlias(ctx, func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(randBytes))
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		got, err := ioutil.ReadAll(*resp)
		require.NoError(t, err)
		assert.Equal(t, randBytes, got)
	})
	t.Run("BinaryOptionalAlias empty", func(t *testing.T) {
		resp, err := client.BinaryOptionalAlias(ctx, func() io.ReadCloser {
			return nil
		})
		require.NoError(t, err)
		require.Nil(t, resp)
	})
	t.Run("BinaryList", func(t *testing.T) {
		list := [][]byte{randBytes}
		got, err := client.BinaryList(ctx, list)
		require.NoError(t, err)
		assert.Equal(t, list, got)
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

func (b binaryServer) BinaryAliasOptional(ctx context.Context) (*io.ReadCloser, error) {
	randBytes := make([]byte, randByteLen)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, err
	}
	resp := ioutil.NopCloser(bytes.NewReader(randBytes))
	return &resp, nil
}

func (b binaryServer) BinaryAliasAlias(ctx context.Context, bodyArg *io.ReadCloser) (*io.ReadCloser, error) {
	if bodyArg == nil {
		return nil, nil
	}
	body, err := ioutil.ReadAll(*bodyArg)
	if err != nil {
		return nil, err
	}
	resp := ioutil.NopCloser(bytes.NewReader(body))
	return &resp, nil
}

func (b binaryServer) BinaryOptional(ctx context.Context) (*io.ReadCloser, error) {
	return nil, nil
}

func (b binaryServer) BinaryOptionalAlias(ctx context.Context, bodyArg *io.ReadCloser) (*io.ReadCloser, error) {
	if bodyArg == nil {
		return nil, nil
	}
	return bodyArg, nil
}

func (b binaryServer) BinaryList(ctx context.Context, bodyArg [][]byte) ([][]byte, error) {
	if len(bodyArg) == 0 {
		return nil, errors.NewInvalidArgument()
	}
	return bodyArg, nil
}
