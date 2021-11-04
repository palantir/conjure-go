package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/palantir/conjure-go/v6/integration_test/testgenerated/server/api"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
	"github.com/palantir/witchcraft-go-server/v2/wrouter/whttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeMarker(t *testing.T) {
	router := wrouter.New(whttprouter.New())
	err := api.RegisterRoutesTestService(router, testServerImpl{})
	require.NoError(t, err)
	called := false
	router.AddRouteHandlerMiddleware(func(rw http.ResponseWriter, r *http.Request, reqVals wrouter.RequestVals, next wrouter.RouteRequestHandler) {
		if reqVals.Spec.PathTemplate == "/safe/{myPathParam1}/{myPathParam2}" {
			called = true

			pathPerms := reqVals.ParamPerms.PathParamPerms()
			assert.True(t, pathPerms.Safe("myPathParam1"))
			assert.False(t, pathPerms.Safe("myPathParam2"))

			headerPerms := reqVals.ParamPerms.HeaderParamPerms()
			assert.True(t, headerPerms.Safe("X-My-Header1-Abc"))
			assert.False(t, headerPerms.Safe("X-My-Header2"))

			queryPerms := reqVals.ParamPerms.QueryParamPerms()
			assert.True(t, queryPerms.Safe("query1"))
			assert.True(t, queryPerms.Safe("myQueryParam2"))
			assert.False(t, queryPerms.Safe("myQueryParam3"))
		}
		next(rw, r, reqVals)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	long2 := safelong.SafeLong(2)
	str := "abc"
	id := uuid.NewUUID()
	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))
	err = client.PostSafeParams(context.Background(),
		"password",
		"myPathParam1Arg",
		true,
		api.CustomObject{Data: []byte("hello world!")},
		"myQueryParam1Arg",
		"myQueryParam2Arg",
		1,
		&long2,
		&str,
		2,
		&id)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestEchoOptionalObject(t *testing.T) {
	wlog.SetDefaultLoggerProvider(wlog.NewJSONMarshalLoggerProvider())
	router := wrouter.New(whttprouter.New())
	err := api.RegisterRoutesTestService(router, testServerImpl{})
	require.NoError(t, err)
	server := httptest.NewServer(router)
	defer server.Close()
	client := api.NewTestServiceClient(newHTTPClient(t, server.URL))

	t.Run("HTTP client", func(t *testing.T) {
		t.Run("nonempty", func(t *testing.T) {
			obj := &api.CustomObject{Data: []byte("hello world")}
			objJSON, err := json.Marshal(obj)
			require.NoError(t, err)
			resp, err := http.Post(server.URL+"/echoCustomObject", "application/json", bytes.NewReader(objJSON))
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			respJSON, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.JSONEq(t, string(objJSON), string(respJSON))
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := http.Post(server.URL+"/echoCustomObject", "application/json", nil)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, http.StatusNoContent, resp.StatusCode)
			require.Equal(t, http.NoBody, resp.Body)
		})
	})
	t.Run("CGR client", func(t *testing.T) {
		t.Run("nonempty", func(t *testing.T) {
			obj := &api.CustomObject{Data: []byte("hello world")}
			resp, err := client.EchoCustomObject(context.Background(), obj)
			require.NoError(t, err)
			require.Equal(t, obj, resp)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := client.EchoCustomObject(context.Background(), nil)
			require.NoError(t, err)
			require.Equal(t, (*api.CustomObject)(nil), resp)
		})
	})
}

type testServerImpl struct{}

func (t testServerImpl) PostSafeParams(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg api.CustomObject, myQueryParam1Arg string, myQueryParam2Arg string,
	myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) error {
	switch {
	case authHeader == "":
		return errors.New("empty authHeader")
	case myPathParam1Arg == "":
		return errors.New("empty myPathParam1Arg")
	case myPathParam2Arg == false:
		return errors.New("empty myPathParam2Arg")
	case reflect.ValueOf(myBodyParamArg).IsZero():
		return errors.New("empty myBodyParamArg")
	case myQueryParam1Arg == "":
		return errors.New("empty myQueryParam1Arg")
	case myQueryParam2Arg == "":
		return errors.New("empty myQueryParam2Arg")
	case myQueryParam3Arg == 0:
		return errors.New("empty myQueryParam3Arg")
	case myQueryParam4Arg == nil:
		return errors.New("empty myQueryParam4Arg")
	case myQueryParam5Arg == nil:
		return errors.New("empty myQueryParam5Arg")
	case myHeaderParam1Arg == 0:
		return errors.New("empty myHeaderParam1Arg")
	case myHeaderParam2Arg == nil:
		return errors.New("empty myHeaderParam2Arg")
	}
	return nil
}

func (t testServerImpl) Echo(ctx context.Context, cookieToken bearertoken.Token) error {
	panic("implement me")
}

func (t testServerImpl) EchoStrings(ctx context.Context, bodyArg []string) ([]string, error) {
	panic("implement me")
}

func (t testServerImpl) EchoCustomObject(ctx context.Context, bodyArg *api.CustomObject) (*api.CustomObject, error) {
	return bodyArg, nil
}

func (t testServerImpl) GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error {
	panic("implement me")
}

func (t testServerImpl) GetPathParamAlias(ctx context.Context, authHeader bearertoken.Token, myPathParamArg api.StringAlias) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamList(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListBoolean(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []bool) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListDouble(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []float64) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []int) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListRid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []rid.ResourceIdentifier) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListSafeLong(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []safelong.SafeLong) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error {
	panic("implement me")
}

func (t testServerImpl) QueryParamListUuid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []uuid.UUID) error {
	panic("implement me")
}

func (t testServerImpl) PostPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg api.CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg api.OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (api.CustomObject, error) {
	panic("implement me")
}

func (t testServerImpl) Bytes(ctx context.Context) (api.CustomObject, error) {
	panic("implement me")
}

func (t testServerImpl) GetBinary(ctx context.Context) (io.ReadCloser, error) {
	panic("implement me")
}

func (t testServerImpl) PostBinary(ctx context.Context, myBytesArg io.ReadCloser) (io.ReadCloser, error) {
	panic("implement me")
}

func (t testServerImpl) PutBinary(ctx context.Context, myBytesArg io.ReadCloser) error {
	panic("implement me")
}

func (t testServerImpl) GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error) {
	panic("implement me")
}

func (t testServerImpl) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong) error {
	panic("implement me")
}
