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
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	werror "github.com/palantir/witchcraft-go-error"
	wparams "github.com/palantir/witchcraft-go-params"
	"github.com/palantir/witchcraft-go-server/wrouter"
	"github.com/palantir/witchcraft-go-server/wrouter/whttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/queryparam/api"
)

func TestQueryParamClient(t *testing.T) {
	server := createTestServer()
	defer server.Close()
	optionalStr := "optional"
	listArg := []int{1, 2, 3}

	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	resp, err := client.Echo(context.Background(), "hello", 3, nil, listArg, &optionalStr)
	require.NoError(t, err)
	assert.Equal(t, "hello hello hello", resp)

	_, err = client.Echo(context.Background(), "hello", -3, &optionalStr, listArg, nil)
	if assert.Error(t, err) {
		cerr := werror.RootCause(err).(errors.Error)
		assert.Equal(t, "reps must be non-negative, was -3", cerr.UnsafeParams()["message"])
		assert.Equal(t, errors.InvalidArgument, cerr.Code())
	}
}

func createTestServer() *httptest.Server {
	r := wrouter.New(whttprouter.New())
	if err := api.RegisterRoutesTestService(r, &testImpl{}); err != nil {
		panic(err)
	}
	server := httptest.NewServer(r)
	return server
}

type testImpl struct{}

func (t *testImpl) Echo(ctx context.Context, inputArg string, repsArg int, optionalArg *string, listParamArg []int, lastParamArg *string) (string, error) {
	if repsArg < 0 {
		return "", errors.NewInvalidArgument(wparams.NewSafeParamStorer(map[string]interface{}{"message": fmt.Sprintf("reps must be non-negative, was %d", repsArg)}))
	}
	if len(listParamArg) < 3 {
		return "", errors.NewInvalidArgument(wparams.NewSafeParamStorer(map[string]interface{}{"message": "listParamArg must have 3 elements"}))
	}
	var parts []string
	for i := 0; i < repsArg; i++ {
		parts = append(parts, inputArg)
	}
	return strings.Join(parts, " "), nil
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}
