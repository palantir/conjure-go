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

package testutil

import (
	"net/http/httptest"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-router/wrouter"
	"github.com/palantir/witchcraft-go-router/wrouter/whttprouter"
	"github.com/stretchr/testify/require"
)

func StartTestServer(t *testing.T, register func(router wrouter.Router) error) (httpclient.Client, func()) {
	router := wrouter.New(whttprouter.New())
	require.NoError(t, register(router))
	server := httptest.NewServer(router)
	client, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{server.URL}),
	)
	require.NoError(t, err)
	return client, server.Close
}
