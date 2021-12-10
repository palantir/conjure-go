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
	"net/http/httptest"
	"testing"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/queryparam/api"
	"github.com/palantir/pkg/uuid"
	_ "github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryParamClient(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()
	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))

	for _, test := range []struct {
		Name string
		Args api.Response
		Err  string
	}{
		{
			Name: "default",
			Args: api.Response{
				Input:                "Input",
				Reps:                 7,
				Optional:             stringPtr("Optional"),
				ListParam:            []int{1, 2, 3},
				LastParam:            stringPtr("Optional"),
				AliasString:          "AliasString",
				AliasAliasString:     "AliasAliasString",
				OptionalAliasString:  api.OptionalAliasString{Value: (*api.AliasString)(stringPtr("OptionalAliasString"))},
				AliasInteger:         31,
				AliasAliasInteger:    529,
				OptionalAliasInteger: api.OptionalAliasInteger{Value: (*api.AliasInteger)(intPtr(4))},
				Uuid:                 uuid.NewUUID(),
				SetUuid:              []uuid.UUID{uuid.NewUUID()},
				SetAliasUuid:         []api.AliasUuid{api.AliasUuid(uuid.NewUUID())},
				AliasUuid:            api.AliasUuid(uuid.NewUUID()),
				AliasAliasUuid:       api.AliasAliasUuid(uuid.NewUUID()),
				OptionalAliasUuid:    api.OptionalAliasUuid{Value: (*api.AliasUuid)(uuidPtr(uuid.NewUUID()))},
				Enum:                 api.New_Enum(api.Enum_VAL1),
				AliasOptionalEnum:    api.AliasOptionalEnum{Value: enumPtr(api.New_Enum(api.Enum_VAL1))},
				AliasEnum:            api.AliasEnum(api.New_Enum(api.Enum_VAL2)),
				OptionalAliasEnum:    api.OptionalAliasEnum{Value: (*api.AliasEnum)(enumPtr(api.New_Enum(api.Enum_VAL2)))},
				ListAliasEnum:        []api.AliasEnum{api.AliasEnum(api.New_Enum(api.Enum_VAL2))},
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("EchoQuery", func(t *testing.T) {
				in := test.Args
				queryResp, err := client.EchoQuery(ctx,
					in.Input,
					in.Reps,
					in.Optional,
					in.ListParam,
					in.LastParam,
					in.AliasString,
					in.AliasAliasString,
					in.OptionalAliasString,
					in.AliasInteger,
					in.AliasAliasInteger,
					in.OptionalAliasInteger,
					in.Uuid,
					in.SetUuid,
					in.SetAliasUuid,
					in.AliasUuid,
					in.AliasAliasUuid,
					in.OptionalAliasUuid,
					in.Enum,
					in.AliasOptionalEnum,
					in.AliasEnum,
					in.OptionalAliasEnum,
					in.ListAliasEnum,
				)
				if test.Err == "" {
					require.NoError(t, err)
					assert.Equal(t, in, queryResp)
				} else {
					require.EqualError(t, err, test.Err)
				}
			})
			t.Run("EchoHeader", func(t *testing.T) {
				in := test.Args
				// Clear out collection arguments not allowed in headers
				// Expect collections get initialized by json marshaling.
				in.ListParam = []int{}
				in.SetUuid = []uuid.UUID{}
				in.SetAliasUuid = []api.AliasUuid{}
				in.ListAliasEnum = []api.AliasEnum{}

				headerResp, err := client.EchoHeader(ctx,
					in.Input,
					in.Reps,
					in.Optional,
					in.LastParam,
					in.AliasString,
					in.AliasAliasString,
					in.OptionalAliasString,
					in.AliasInteger,
					in.AliasAliasInteger,
					in.OptionalAliasInteger,
					in.Uuid,
					in.AliasUuid,
					in.AliasAliasUuid,
					in.OptionalAliasUuid,
					in.Enum,
					in.AliasOptionalEnum,
					in.AliasEnum,
					in.OptionalAliasEnum,
				)
				if test.Err == "" {
					require.NoError(t, err)
					assert.Equal(t, in, headerResp)
				} else {
					require.EqualError(t, err, test.Err)
				}
			})
		})
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

func (t *testImpl) EchoQuery(ctx context.Context,
	inputArg string,
	repsArg int,
	optionalArg *string,
	listParamArg []int,
	lastParamArg *string,
	aliasStringArg api.AliasString,
	aliasAliasStringArg api.AliasAliasString,
	optionalAliasStringArg api.OptionalAliasString,
	aliasIntegerArg api.AliasInteger,
	aliasAliasIntegerArg api.AliasAliasInteger,
	optionalAliasIntegerArg api.OptionalAliasInteger,
	uuidArg uuid.UUID,
	setUuidArg []uuid.UUID,
	setAliasUuidArg []api.AliasUuid,
	aliasUuidArg api.AliasUuid,
	aliasAliasUuidArg api.AliasAliasUuid,
	optionalAliasUuidArg api.OptionalAliasUuid,
	enumArg api.Enum,
	aliasOptionalEnumArg api.AliasOptionalEnum,
	aliasEnumArg api.AliasEnum,
	optionalAliasEnumArg api.OptionalAliasEnum,
	listAliasEnumArg api.ListAliasEnum,
) (api.Response, error) {
	return api.Response{
		Input:                inputArg,
		Reps:                 repsArg,
		Optional:             optionalArg,
		ListParam:            listParamArg,
		LastParam:            lastParamArg,
		AliasString:          aliasStringArg,
		AliasAliasString:     aliasAliasStringArg,
		OptionalAliasString:  optionalAliasStringArg,
		AliasInteger:         aliasIntegerArg,
		AliasAliasInteger:    aliasAliasIntegerArg,
		OptionalAliasInteger: optionalAliasIntegerArg,
		Uuid:                 uuidArg,
		SetUuid:              setUuidArg,
		SetAliasUuid:         setAliasUuidArg,
		AliasUuid:            aliasUuidArg,
		AliasAliasUuid:       aliasAliasUuidArg,
		OptionalAliasUuid:    optionalAliasUuidArg,
		Enum:                 enumArg,
		AliasOptionalEnum:    aliasOptionalEnumArg,
		AliasEnum:            aliasEnumArg,
		OptionalAliasEnum:    optionalAliasEnumArg,
		ListAliasEnum:        listAliasEnumArg,
	}, nil
}

func (t *testImpl) EchoHeader(ctx context.Context,
	inputArg string,
	repsArg int,
	optionalArg *string,
	lastParamArg *string,
	aliasStringArg api.AliasString,
	aliasAliasStringArg api.AliasAliasString,
	optionalAliasStringArg api.OptionalAliasString,
	aliasIntegerArg api.AliasInteger,
	aliasAliasIntegerArg api.AliasAliasInteger,
	optionalAliasIntegerArg api.OptionalAliasInteger,
	uuidArg uuid.UUID,
	aliasUuidArg api.AliasUuid,
	aliasAliasUuidArg api.AliasAliasUuid,
	optionalAliasUuidArg api.OptionalAliasUuid,
	enumArg api.Enum,
	aliasOptionalEnumArg api.AliasOptionalEnum,
	aliasEnumArg api.AliasEnum,
	optionalAliasEnumArg api.OptionalAliasEnum,
) (api.Response, error) {
	return api.Response{
		Input:                inputArg,
		Reps:                 repsArg,
		Optional:             optionalArg,
		LastParam:            lastParamArg,
		AliasString:          aliasStringArg,
		AliasAliasString:     aliasAliasStringArg,
		OptionalAliasString:  optionalAliasStringArg,
		AliasInteger:         aliasIntegerArg,
		AliasAliasInteger:    aliasAliasIntegerArg,
		OptionalAliasInteger: optionalAliasIntegerArg,
		Uuid:                 uuidArg,
		AliasUuid:            aliasUuidArg,
		AliasAliasUuid:       aliasAliasUuidArg,
		OptionalAliasUuid:    optionalAliasUuidArg,
		Enum:                 enumArg,
		AliasOptionalEnum:    aliasOptionalEnumArg,
		AliasEnum:            aliasEnumArg,
		OptionalAliasEnum:    optionalAliasEnumArg,
	}, nil
}

func newHTTPClient(t *testing.T, url string) httpclient.Client {
	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{url}),
	)
	require.NoError(t, err)
	return httpClient
}

func stringPtr(s string) *string     { return &s }
func intPtr(i int) *int              { return &i }
func uuidPtr(u uuid.UUID) *uuid.UUID { return &u }
func enumPtr(e api.Enum) *api.Enum   { return &e }
