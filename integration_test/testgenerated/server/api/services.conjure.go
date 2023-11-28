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
	EchoStrings(ctx context.Context, bodyArg []string) ([]string, error)
	EchoCustomObject(ctx context.Context, bodyArg *CustomObject) (*CustomObject, error)
	EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error)
	EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error)
	GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error
	GetPathParamAlias(ctx context.Context, authHeader bearertoken.Token, myPathParamArg StringAlias) error
	QueryParamList(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error
	QueryParamListBoolean(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []bool) error
	QueryParamListDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) error
	QueryParamSetDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) ([]datetime.DateTime, error)
	QueryParamListDouble(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []float64) error
	QueryParamListInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []int) error
	QueryParamListRid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []rid.ResourceIdentifier) error
	QueryParamListSafeLong(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []safelong.SafeLong) error
	QueryParamListString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error
	QueryParamListUuid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []uuid.UUID) error
	QueryParamExternalString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg string) error
	QueryParamExternalInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg int) error
	PathParamExternalString(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string) error
	PathParamExternalInteger(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg int) error
	PostPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error)
	PostSafeParams(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *SafeUuid) error
	Bytes(ctx context.Context) (CustomObject, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) error
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
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
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("PALANTIR_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "echo failed")
	}
	return nil
}

func (c *testServiceClient) EchoStrings(ctx context.Context, bodyArg []string) ([]string, error) {
	var returnVal []string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoStrings"))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
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
	requestParams = append(requestParams, httpclient.WithPathf("/echoCustomObject"))
	if bodyArg != nil {
		requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	}
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "echoCustomObject failed")
	}
	return returnVal, nil
}

func (c *testServiceClient) EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error) {
	var defaultReturnVal OptionalIntegerAlias
	var returnVal OptionalIntegerAlias
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoOptionalAlias"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/alias"))
	if bodyArg.Value != nil {
		requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	}
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "echoOptionalAlias failed")
	}
	return returnVal, nil
}

func (c *testServiceClient) EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error) {
	var defaultReturnVal OptionalListAlias
	var returnVal OptionalListAlias
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoOptionalListAlias"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/list-alias"))
	if bodyArg.Value != nil {
		requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	}
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "echoOptionalListAlias failed")
	}
	return returnVal, nil
}

func (c *testServiceClient) GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetPathParam"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/string/%s", url.PathEscape(fmt.Sprint(myPathParamArg))))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "getPathParam failed")
	}
	return nil
}

func (c *testServiceClient) GetPathParamAlias(ctx context.Context, authHeader bearertoken.Token, myPathParamArg StringAlias) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetPathParamAlias"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/alias/%s", url.PathEscape(fmt.Sprint(myPathParamArg))))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "getPathParamAlias failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamList(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamList"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/pathNew"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamList failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListBoolean(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []bool) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListBoolean"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/booleanListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListBoolean failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListDateTime"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/dateTimeListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListDateTime failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamSetDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) ([]datetime.DateTime, error) {
	var returnVal []datetime.DateTime
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamSetDateTime"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/dateTimeSetQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "queryParamSetDateTime failed")
	}
	if returnVal == nil {
		return nil, werror.ErrorWithContextParams(ctx, "queryParamSetDateTime response cannot be nil")
	}
	return returnVal, nil
}

func (c *testServiceClient) QueryParamListDouble(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []float64) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListDouble"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/doubleListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListDouble failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []int) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListInteger"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/intListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListInteger failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListRid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []rid.ResourceIdentifier) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListRid"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/ridListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListRid failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListSafeLong(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []safelong.SafeLong) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListSafeLong"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/safeLongListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListSafeLong failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListString"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/stringListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListString failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamListUuid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []uuid.UUID) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListUuid"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/uuidListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamListUuid failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamExternalString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamExternalString"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/externalStringQueryVar"))
	queryParams := make(url.Values)
	queryParams.Set("myQueryParam1", fmt.Sprint(myQueryParam1Arg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamExternalString failed")
	}
	return nil
}

func (c *testServiceClient) QueryParamExternalInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg int) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamExternalInteger"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/externalIntegerQueryVar"))
	queryParams := make(url.Values)
	queryParams.Set("myQueryParam1", fmt.Sprint(myQueryParam1Arg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "queryParamExternalInteger failed")
	}
	return nil
}

func (c *testServiceClient) PathParamExternalString(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PathParamExternalString"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/externalStringPath/%s", url.PathEscape(fmt.Sprint(myPathParam1Arg))))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "pathParamExternalString failed")
	}
	return nil
}

func (c *testServiceClient) PathParamExternalInteger(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg int) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PathParamExternalInteger"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/externalIntegerPath/%s", url.PathEscape(fmt.Sprint(myPathParam1Arg))))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "pathParamExternalInteger failed")
	}
	return nil
}

func (c *testServiceClient) PostPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error) {
	var defaultReturnVal CustomObject
	var returnVal *CustomObject
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PostPathParam"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/%s/%s", url.PathEscape(fmt.Sprint(myPathParam1Arg)), url.PathEscape(fmt.Sprint(myPathParam2Arg))))
	requestParams = append(requestParams, httpclient.WithJSONRequest(myBodyParamArg))
	requestParams = append(requestParams, httpclient.WithHeader("X-My-Header1-Abc", fmt.Sprint(myHeaderParam1Arg)))
	if myHeaderParam2Arg != nil {
		requestParams = append(requestParams, httpclient.WithHeader("X-My-Header2", fmt.Sprint(*myHeaderParam2Arg)))
	}
	queryParams := make(url.Values)
	queryParams.Set("query1", fmt.Sprint(myQueryParam1Arg))
	queryParams.Set("myQueryParam2", fmt.Sprint(myQueryParam2Arg))
	queryParams.Set("myQueryParam3", fmt.Sprint(myQueryParam3Arg))
	if myQueryParam4Arg != nil {
		queryParams.Set("myQueryParam4", fmt.Sprint(*myQueryParam4Arg))
	}
	if myQueryParam5Arg != nil {
		queryParams.Set("myQueryParam5", fmt.Sprint(*myQueryParam5Arg))
	}
	if myQueryParam6Arg.Value != nil {
		queryParams.Set("myQueryParam6", fmt.Sprint(*myQueryParam6Arg.Value))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "postPathParam failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "postPathParam response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) PostSafeParams(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *SafeUuid) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PostSafeParams"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/safe/%s/%s", url.PathEscape(fmt.Sprint(myPathParam1Arg)), url.PathEscape(fmt.Sprint(myPathParam2Arg))))
	requestParams = append(requestParams, httpclient.WithJSONRequest(myBodyParamArg))
	requestParams = append(requestParams, httpclient.WithHeader("X-My-Header1-Abc", fmt.Sprint(myHeaderParam1Arg)))
	if myHeaderParam2Arg != nil {
		requestParams = append(requestParams, httpclient.WithHeader("X-My-Header2", fmt.Sprint(*myHeaderParam2Arg)))
	}
	queryParams := make(url.Values)
	queryParams.Set("query1", fmt.Sprint(myQueryParam1Arg))
	queryParams.Set("myQueryParam2", fmt.Sprint(myQueryParam2Arg))
	queryParams.Set("myQueryParam3", fmt.Sprint(myQueryParam3Arg))
	if myQueryParam4Arg != nil {
		queryParams.Set("myQueryParam4", fmt.Sprint(*myQueryParam4Arg))
	}
	if myQueryParam5Arg != nil {
		queryParams.Set("myQueryParam5", fmt.Sprint(*myQueryParam5Arg))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "postSafeParams failed")
	}
	return nil
}

func (c *testServiceClient) Bytes(ctx context.Context) (CustomObject, error) {
	var defaultReturnVal CustomObject
	var returnVal *CustomObject
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Bytes"))
	requestParams = append(requestParams, httpclient.WithPathf("/bytes"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Get(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "bytes failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "bytes response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) GetBinary(ctx context.Context) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBinary"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Get(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "getBinary failed")
	}
	return resp.Body, nil
}

func (c *testServiceClient) PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PostBinary"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(myBytesArg))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Post(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "postBinary failed")
	}
	return resp.Body, nil
}

func (c *testServiceClient) PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PutBinary"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(myBytesArg))
	if _, err := c.client.Put(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "putBinary failed")
	}
	return nil
}

func (c *testServiceClient) GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetOptionalBinary"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Get(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "getOptionalBinary failed")
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	return &resp.Body, nil
}

func (c *testServiceClient) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Chan"))
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
	if _, err := c.client.Post(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "chan failed")
	}
	return nil
}

type TestServiceClientWithAuth interface {
	Echo(ctx context.Context) error
	EchoStrings(ctx context.Context, bodyArg []string) ([]string, error)
	EchoCustomObject(ctx context.Context, bodyArg *CustomObject) (*CustomObject, error)
	EchoOptionalAlias(ctx context.Context, bodyArg OptionalIntegerAlias) (OptionalIntegerAlias, error)
	EchoOptionalListAlias(ctx context.Context, bodyArg OptionalListAlias) (OptionalListAlias, error)
	GetPathParam(ctx context.Context, myPathParamArg string) error
	GetPathParamAlias(ctx context.Context, myPathParamArg StringAlias) error
	QueryParamList(ctx context.Context, myQueryParam1Arg []string) error
	QueryParamListBoolean(ctx context.Context, myQueryParam1Arg []bool) error
	QueryParamListDateTime(ctx context.Context, myQueryParam1Arg []datetime.DateTime) error
	QueryParamSetDateTime(ctx context.Context, myQueryParam1Arg []datetime.DateTime) ([]datetime.DateTime, error)
	QueryParamListDouble(ctx context.Context, myQueryParam1Arg []float64) error
	QueryParamListInteger(ctx context.Context, myQueryParam1Arg []int) error
	QueryParamListRid(ctx context.Context, myQueryParam1Arg []rid.ResourceIdentifier) error
	QueryParamListSafeLong(ctx context.Context, myQueryParam1Arg []safelong.SafeLong) error
	QueryParamListString(ctx context.Context, myQueryParam1Arg []string) error
	QueryParamListUuid(ctx context.Context, myQueryParam1Arg []uuid.UUID) error
	QueryParamExternalString(ctx context.Context, myQueryParam1Arg string) error
	QueryParamExternalInteger(ctx context.Context, myQueryParam1Arg int) error
	PathParamExternalString(ctx context.Context, myPathParam1Arg string) error
	PathParamExternalInteger(ctx context.Context, myPathParam1Arg int) error
	PostPathParam(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error)
	PostSafeParams(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *SafeUuid) error
	Bytes(ctx context.Context) (CustomObject, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) error
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
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

func (c *testServiceClientWithAuth) GetPathParamAlias(ctx context.Context, myPathParamArg StringAlias) error {
	return c.client.GetPathParamAlias(ctx, c.authHeader, myPathParamArg)
}

func (c *testServiceClientWithAuth) QueryParamList(ctx context.Context, myQueryParam1Arg []string) error {
	return c.client.QueryParamList(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListBoolean(ctx context.Context, myQueryParam1Arg []bool) error {
	return c.client.QueryParamListBoolean(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListDateTime(ctx context.Context, myQueryParam1Arg []datetime.DateTime) error {
	return c.client.QueryParamListDateTime(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamSetDateTime(ctx context.Context, myQueryParam1Arg []datetime.DateTime) ([]datetime.DateTime, error) {
	return c.client.QueryParamSetDateTime(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListDouble(ctx context.Context, myQueryParam1Arg []float64) error {
	return c.client.QueryParamListDouble(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListInteger(ctx context.Context, myQueryParam1Arg []int) error {
	return c.client.QueryParamListInteger(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListRid(ctx context.Context, myQueryParam1Arg []rid.ResourceIdentifier) error {
	return c.client.QueryParamListRid(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListSafeLong(ctx context.Context, myQueryParam1Arg []safelong.SafeLong) error {
	return c.client.QueryParamListSafeLong(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListString(ctx context.Context, myQueryParam1Arg []string) error {
	return c.client.QueryParamListString(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListUuid(ctx context.Context, myQueryParam1Arg []uuid.UUID) error {
	return c.client.QueryParamListUuid(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamExternalString(ctx context.Context, myQueryParam1Arg string) error {
	return c.client.QueryParamExternalString(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamExternalInteger(ctx context.Context, myQueryParam1Arg int) error {
	return c.client.QueryParamExternalInteger(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) PathParamExternalString(ctx context.Context, myPathParam1Arg string) error {
	return c.client.PathParamExternalString(ctx, c.authHeader, myPathParam1Arg)
}

func (c *testServiceClientWithAuth) PathParamExternalInteger(ctx context.Context, myPathParam1Arg int) error {
	return c.client.PathParamExternalInteger(ctx, c.authHeader, myPathParam1Arg)
}

func (c *testServiceClientWithAuth) PostPathParam(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error) {
	return c.client.PostPathParam(ctx, c.authHeader, myPathParam1Arg, myPathParam2Arg, myBodyParamArg, myQueryParam1Arg, myQueryParam2Arg, myQueryParam3Arg, myQueryParam4Arg, myQueryParam5Arg, myQueryParam6Arg, myHeaderParam1Arg, myHeaderParam2Arg)
}

func (c *testServiceClientWithAuth) PostSafeParams(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *SafeUuid) error {
	return c.client.PostSafeParams(ctx, c.authHeader, myPathParam1Arg, myPathParam2Arg, myBodyParamArg, myQueryParam1Arg, myQueryParam2Arg, myQueryParam3Arg, myQueryParam4Arg, myQueryParam5Arg, myHeaderParam1Arg, myHeaderParam2Arg)
}

func (c *testServiceClientWithAuth) Bytes(ctx context.Context) (CustomObject, error) {
	return c.client.Bytes(ctx)
}

func (c *testServiceClientWithAuth) GetBinary(ctx context.Context) (io.ReadCloser, error) {
	return c.client.GetBinary(ctx)
}

func (c *testServiceClientWithAuth) PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (io.ReadCloser, error) {
	return c.client.PostBinary(ctx, myBytesArg)
}

func (c *testServiceClientWithAuth) PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) error {
	return c.client.PutBinary(ctx, myBytesArg)
}

func (c *testServiceClientWithAuth) GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error) {
	return c.client.GetOptionalBinary(ctx)
}

func (c *testServiceClientWithAuth) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error {
	return c.client.Chan(ctx, varArg, importArg, typeArg, returnArg, httpArg, jsonArg, reqArg, rwArg)
}
