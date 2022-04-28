// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
)

type TestServiceClient interface {
	Echo(ctx context.Context, cookieToken bearertoken.Token) error
	// These are some endpoint docs
	EchoStrings(ctx context.Context, bodyArg []string) ([]string, error)
	EchoCustomObject(ctx context.Context, bodyArg *CustomObject) (*CustomObject, error)
	EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error)
	EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error)
	GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error
	GetListBoolean(ctx context.Context, myQueryParam1Arg []bool) ([]bool, error)
	PutMapStringString(ctx context.Context, myParamArg map[string]string) (map[string]string, error)
	PutMapStringAny(ctx context.Context, myParamArg map[string]interface{}) (map[string]interface{}, error)
	GetDateTime(ctx context.Context, myParamArg datetime.DateTime) (datetime.DateTime, error)
	GetDouble(ctx context.Context, myParamArg float64) (float64, error)
	GetRid(ctx context.Context, myParamArg rid.ResourceIdentifier) (rid.ResourceIdentifier, error)
	GetSafeLong(ctx context.Context, myParamArg safelong.SafeLong) (safelong.SafeLong, error)
	GetUuid(ctx context.Context, myParamArg uuid.UUID) (uuid.UUID, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
	// An endpoint that uses reserved flag names
	GetReserved(ctx context.Context, confArg string, bearertokenArg string) error
	// An endpoint that uses go keywords
	Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) Echo(ctx context.Context, cookieToken bearertoken.Token) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Echo"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("PALANTIR_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "echo failed")
	}
	return nil
}

func (c *testServiceClient) EchoStrings(ctx context.Context, bodyArg []string) ([]string, error) {
	var returnVal []string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoStrings"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "echoStrings failed")
	}
	if returnVal == nil {
		return nil, werror.ErrorWithContextParams(ctx, "echoStrings response cannot be nil")
	}
	return returnVal, nil
}

func (c *testServiceClient) EchoCustomObject(ctx context.Context, bodyArg *CustomObject) (*CustomObject, error) {
	var returnVal *CustomObject
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoCustomObject"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/echoCustomObject"))
	if bodyArg != nil {
		requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	}
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "echoCustomObject failed")
	}
	return returnVal, nil
}

func (c *testServiceClient) EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error) {
	var defaultReturnVal OptionalIntegerAlias
	var returnVal OptionalIntegerAlias
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoOptionalAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/alias"))
	if bodyArg.Value != nil {
		requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	}
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "echoOptionalAlias failed")
	}
	return returnVal, nil
}

func (c *testServiceClient) EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error) {
	var defaultReturnVal OptionalListAlias
	var returnVal OptionalListAlias
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoOptionalListAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/list-alias"))
	if bodyArg.Value != nil {
		requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	}
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "echoOptionalListAlias failed")
	}
	return returnVal, nil
}

func (c *testServiceClient) GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetPathParam"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/string/%s", url.PathEscape(fmt.Sprint(myPathParamArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "getPathParam failed")
	}
	return nil
}

func (c *testServiceClient) GetListBoolean(ctx context.Context, myQueryParam1Arg []bool) ([]bool, error) {
	var returnVal []bool
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetListBoolean"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/booleanListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "getListBoolean failed")
	}
	if returnVal == nil {
		return nil, werror.ErrorWithContextParams(ctx, "getListBoolean response cannot be nil")
	}
	return returnVal, nil
}

func (c *testServiceClient) PutMapStringString(ctx context.Context, myParamArg map[string]string) (map[string]string, error) {
	var returnVal map[string]string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PutMapStringString"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("PUT"))
	requestParams = append(requestParams, httpclient.WithPathf("/mapStringString"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(myParamArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "putMapStringString failed")
	}
	if returnVal == nil {
		return nil, werror.ErrorWithContextParams(ctx, "putMapStringString response cannot be nil")
	}
	return returnVal, nil
}

func (c *testServiceClient) PutMapStringAny(ctx context.Context, myParamArg map[string]interface{}) (map[string]interface{}, error) {
	var returnVal map[string]interface{}
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PutMapStringAny"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("PUT"))
	requestParams = append(requestParams, httpclient.WithPathf("/mapStringAny"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(myParamArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "putMapStringAny failed")
	}
	if returnVal == nil {
		return nil, werror.ErrorWithContextParams(ctx, "putMapStringAny response cannot be nil")
	}
	return returnVal, nil
}

func (c *testServiceClient) GetDateTime(ctx context.Context, myParamArg datetime.DateTime) (datetime.DateTime, error) {
	var defaultReturnVal datetime.DateTime
	var returnVal *datetime.DateTime
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetDateTime"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/getDateTime"))
	queryParams := make(url.Values)
	queryParams.Set("myParam", fmt.Sprint(myParamArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "getDateTime failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "getDateTime response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) GetDouble(ctx context.Context, myParamArg float64) (float64, error) {
	var defaultReturnVal float64
	var returnVal *float64
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetDouble"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/getDouble"))
	queryParams := make(url.Values)
	queryParams.Set("myParam", fmt.Sprint(myParamArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "getDouble failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "getDouble response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) GetRid(ctx context.Context, myParamArg rid.ResourceIdentifier) (rid.ResourceIdentifier, error) {
	var defaultReturnVal rid.ResourceIdentifier
	var returnVal *rid.ResourceIdentifier
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetRid"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/getRid"))
	queryParams := make(url.Values)
	queryParams.Set("myParam", fmt.Sprint(myParamArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "getRid failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "getRid response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) GetSafeLong(ctx context.Context, myParamArg safelong.SafeLong) (safelong.SafeLong, error) {
	var defaultReturnVal safelong.SafeLong
	var returnVal *safelong.SafeLong
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetSafeLong"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/getSafeLong"))
	queryParams := make(url.Values)
	queryParams.Set("myParam", fmt.Sprint(myParamArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "getSafeLong failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "getSafeLong response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) GetUuid(ctx context.Context, myParamArg uuid.UUID) (uuid.UUID, error) {
	var defaultReturnVal uuid.UUID
	var returnVal *uuid.UUID
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetUuid"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/getUuid"))
	queryParams := make(url.Values)
	queryParams.Set("myParam", fmt.Sprint(myParamArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "getUuid failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "getUuid response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) GetBinary(ctx context.Context) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "getBinary failed")
	}
	return resp.Body, nil
}

func (c *testServiceClient) GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetOptionalBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "getOptionalBinary failed")
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	return &resp.Body, nil
}

func (c *testServiceClient) GetReserved(ctx context.Context, confArg string, bearertokenArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetReserved"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/getReserved"))
	queryParams := make(url.Values)
	queryParams.Set("conf", fmt.Sprint(confArg))
	queryParams.Set("bearertoken", fmt.Sprint(bearertokenArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "getReserved failed")
	}
	return nil
}

func (c *testServiceClient) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Chan"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/chan/%s", url.PathEscape(fmt.Sprint(varArg))))
	requestParams = append(requestParams, httpclient.WithJSONRequest(importArg))
	requestParams = append(requestParams, httpclient.WithHeader("X-My-Header2", fmt.Sprint(returnArg)))
	queryParams := make(url.Values)
	queryParams.Set("type", fmt.Sprint(typeArg))
	queryParams.Set("http", fmt.Sprint(httpArg))
	queryParams.Set("json", fmt.Sprint(jsonArg))
	queryParams.Set("req", fmt.Sprint(reqArg))
	queryParams.Set("rw", fmt.Sprint(rwArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "chan failed")
	}
	return nil
}

type TestServiceClientWithAuth interface {
	Echo(ctx context.Context) error
	// These are some endpoint docs
	EchoStrings(ctx context.Context, bodyArg []string) ([]string, error)
	EchoCustomObject(ctx context.Context, bodyArg *CustomObject) (*CustomObject, error)
	EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error)
	EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error)
	GetPathParam(ctx context.Context, myPathParamArg string) error
	GetListBoolean(ctx context.Context, myQueryParam1Arg []bool) ([]bool, error)
	PutMapStringString(ctx context.Context, myParamArg map[string]string) (map[string]string, error)
	PutMapStringAny(ctx context.Context, myParamArg map[string]interface{}) (map[string]interface{}, error)
	GetDateTime(ctx context.Context, myParamArg datetime.DateTime) (datetime.DateTime, error)
	GetDouble(ctx context.Context, myParamArg float64) (float64, error)
	GetRid(ctx context.Context, myParamArg rid.ResourceIdentifier) (rid.ResourceIdentifier, error)
	GetSafeLong(ctx context.Context, myParamArg safelong.SafeLong) (safelong.SafeLong, error)
	GetUuid(ctx context.Context, myParamArg uuid.UUID) (uuid.UUID, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
	// An endpoint that uses reserved flag names
	GetReserved(ctx context.Context, confArg string, bearertokenArg string) error
	// An endpoint that uses go keywords
	Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error
}

func NewTestServiceClientWithAuth(client TestServiceClient, authHeader bearertoken.Token, cookieToken bearertoken.Token) TestServiceClientWithAuth {
	return &testServiceClientWithAuth{client: client, authHeader: authHeader, cookieToken: cookieToken}
}

type testServiceClientWithAuth struct {
	client      TestServiceClient
	authHeader  bearertoken.Token
	cookieToken bearertoken.Token
}

func (c *testServiceClientWithAuth) Echo(ctx context.Context) error {
	return c.client.Echo(ctx, c.cookieToken)
}

func (c *testServiceClientWithAuth) EchoStrings(ctx context.Context, bodyArg []string) ([]string, error) {
	return c.client.EchoStrings(ctx, bodyArg)
}

func (c *testServiceClientWithAuth) EchoCustomObject(ctx context.Context, bodyArg *CustomObject) (*CustomObject, error) {
	return c.client.EchoCustomObject(ctx, bodyArg)
}

func (c *testServiceClientWithAuth) EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error) {
	return c.client.EchoOptionalAlias(ctx, bodyArg)
}

func (c *testServiceClientWithAuth) EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error) {
	return c.client.EchoOptionalListAlias(ctx, bodyArg)
}

func (c *testServiceClientWithAuth) GetPathParam(ctx context.Context, myPathParamArg string) error {
	return c.client.GetPathParam(ctx, c.authHeader, myPathParamArg)
}

func (c *testServiceClientWithAuth) GetListBoolean(ctx context.Context, myQueryParam1Arg []bool) ([]bool, error) {
	return c.client.GetListBoolean(ctx, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) PutMapStringString(ctx context.Context, myParamArg map[string]string) (map[string]string, error) {
	return c.client.PutMapStringString(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) PutMapStringAny(ctx context.Context, myParamArg map[string]interface{}) (map[string]interface{}, error) {
	return c.client.PutMapStringAny(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) GetDateTime(ctx context.Context, myParamArg datetime.DateTime) (datetime.DateTime, error) {
	return c.client.GetDateTime(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) GetDouble(ctx context.Context, myParamArg float64) (float64, error) {
	return c.client.GetDouble(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) GetRid(ctx context.Context, myParamArg rid.ResourceIdentifier) (rid.ResourceIdentifier, error) {
	return c.client.GetRid(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) GetSafeLong(ctx context.Context, myParamArg safelong.SafeLong) (safelong.SafeLong, error) {
	return c.client.GetSafeLong(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) GetUuid(ctx context.Context, myParamArg uuid.UUID) (uuid.UUID, error) {
	return c.client.GetUuid(ctx, myParamArg)
}

func (c *testServiceClientWithAuth) GetBinary(ctx context.Context) (io.ReadCloser, error) {
	return c.client.GetBinary(ctx)
}

func (c *testServiceClientWithAuth) GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error) {
	return c.client.GetOptionalBinary(ctx)
}

func (c *testServiceClientWithAuth) GetReserved(ctx context.Context, confArg string, bearertokenArg string) error {
	return c.client.GetReserved(ctx, confArg, bearertokenArg)
}

func (c *testServiceClientWithAuth) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error {
	return c.client.Chan(ctx, varArg, importArg, typeArg, returnArg, httpArg, jsonArg, reqArg, rwArg)
}
