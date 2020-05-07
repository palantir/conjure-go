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
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/witchcraft-go-server/witchcraft"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v4/integration_test/internal/testutil"
	"github.com/palantir/conjure-go/v4/integration_test/testgenerated/auth/api"
)

const (
	headerAuthAccepted = "header: Authorization accepted"
	testJWT            = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0ifQ.JTD36MhcwmSuvfdCkfSYc-LHOGNA1UQ-0FKLKqdXbF4`
)

var (
	tokenProvider = httpclient.TokenProvider(func(context.Context) (string, error) {
		return testJWT, nil
	})
)

func TestBothAuthClient(t *testing.T) {
	ctx := testutil.TestContext()
	httpClient, cleanup := testutil.StartTestServer(t, func(ctx context.Context, info witchcraft.InitInfo) (cleanup func(), rErr error) {
		if err := api.RegisterRoutesBothAuthService(info.Router, bothAuthImpl{}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	defer cleanup()
	client := api.NewBothAuthServiceClient(httpClient)
	authClient := api.NewBothAuthServiceClientWithAuth(client, testJWT, testJWT)

	// test header auth calls
	resp, err := client.Default(ctx, testJWT)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
	resp, err = authClient.Default(ctx)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)

	// test invalid auth
	_, err = client.Default(ctx, "invalid token")
	assert.EqualError(t, err, "httpclient request failed: failed to unmarshal body as conjure error: json: cannot unmarshal string into Go value of type struct { Name string \"json:\\\"errorName\\\"\" }")

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
	ctx := testutil.TestContext()
	httpClient, cleanup := testutil.StartTestServer(t, func(ctx context.Context, info witchcraft.InitInfo) (cleanup func(), rErr error) {
		if err := api.RegisterRoutesBothAuthService(info.Router, bothAuthImpl{}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	defer cleanup()
	client := api.NewHeaderAuthServiceClient(httpClient)
	authClient := api.NewHeaderAuthServiceClientWithAuth(client, testJWT)
	tokenClient := api.NewHeaderAuthServiceClientWithTokenProvider(client, tokenProvider)

	// test header auth calls
	resp, err := client.Default(ctx, testJWT)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
	resp, err = authClient.Default(ctx)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
	resp, err = tokenClient.Default(ctx)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)
}

func TestCookieAuthClient(t *testing.T) {
	ctx := testutil.TestContext()
	httpClient, cleanup := testutil.StartTestServer(t, func(ctx context.Context, info witchcraft.InitInfo) (cleanup func(), rErr error) {
		if err := api.RegisterRoutesBothAuthService(info.Router, bothAuthImpl{}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	defer cleanup()
	client := api.NewCookieAuthServiceClient(httpClient)
	authClient := api.NewCookieAuthServiceClientWithAuth(client, testJWT)
	tokenClient := api.NewCookieAuthServiceClientWithTokenProvider(client, tokenProvider)

	// test cookie auth calls
	err := authClient.Cookie(ctx)
	require.NoError(t, err)
	err = client.Cookie(ctx, testJWT)
	require.NoError(t, err)
	err = tokenClient.Cookie(ctx)
	require.NoError(t, err)
}

func TestTokenProviderClient(t *testing.T) {
	ctx := testutil.TestContext()
	httpClient, cleanup := testutil.StartTestServer(t, func(ctx context.Context, info witchcraft.InitInfo) (cleanup func(), rErr error) {
		if err := api.RegisterRoutesBothAuthService(info.Router, bothAuthImpl{}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	defer cleanup()
	client := api.NewSomeHeaderAuthServiceClient(httpClient)
	tokenClient := api.NewSomeHeaderAuthServiceClientWithTokenProvider(client, tokenProvider)

	// test cookie auth calls
	resp, err := client.Default(ctx, testJWT)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)

	resp, err = tokenClient.Default(ctx)
	require.NoError(t, err)
	assert.Equal(t, headerAuthAccepted, resp)

	err = client.None(ctx)
	require.NoError(t, err)
	err = tokenClient.None(ctx)
	require.NoError(t, err)
}

type bothAuthImpl struct{}

func (bothAuthImpl) Default(ctx context.Context, authHeader bearertoken.Token) (string, error) {
	if authHeader != testJWT {
		return "", errors.NewPermissionDenied()
	}
	return headerAuthAccepted, nil
}

func (bothAuthImpl) Cookie(ctx context.Context, cookieToken bearertoken.Token) error {
	if cookieToken != testJWT {
		return errors.NewPermissionDenied()
	}
	return nil
}

func (bothAuthImpl) None(ctx context.Context) error {
	return nil
}

func (bothAuthImpl) WithArg(ctx context.Context, authHeader bearertoken.Token, argArg string) error {
	if authHeader != testJWT {
		return errors.NewPermissionDenied()
	}
	return nil
}
