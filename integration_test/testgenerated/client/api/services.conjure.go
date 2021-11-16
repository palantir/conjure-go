// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/pkg/rid"
	werror "github.com/palantir/witchcraft-go-error"
)

type TestServiceClient interface {
	Echo(ctx context.Context) error
	PathParam(ctx context.Context, paramArg string) error
	PathParamAlias(ctx context.Context, paramArg StringAlias) error
	PathParamRid(ctx context.Context, paramArg rid.ResourceIdentifier) error
	PathParamRidAlias(ctx context.Context, paramArg RidAlias) error
	Bytes(ctx context.Context) (CustomObject, error)
	Binary(ctx context.Context) (io.ReadCloser, error)
	MaybeBinary(ctx context.Context) (*io.ReadCloser, error)
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) Echo(ctx context.Context) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Echo"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "echo failed")
	}
	return nil
}

func (c *testServiceClient) PathParam(ctx context.Context, paramArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PathParam"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/path/%s", url.PathEscape(fmt.Sprint(paramArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "pathParam failed")
	}
	return nil
}

func (c *testServiceClient) PathParamAlias(ctx context.Context, paramArg StringAlias) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PathParamAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/path/alias/%s", url.PathEscape(fmt.Sprint(paramArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "pathParamAlias failed")
	}
	return nil
}

func (c *testServiceClient) PathParamRid(ctx context.Context, paramArg rid.ResourceIdentifier) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PathParamRid"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/path/rid/%s", url.PathEscape(fmt.Sprint(paramArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "pathParamRid failed")
	}
	return nil
}

func (c *testServiceClient) PathParamRidAlias(ctx context.Context, paramArg RidAlias) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PathParamRidAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/path/rid/alias/%s", url.PathEscape(fmt.Sprint(paramArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "pathParamRidAlias failed")
	}
	return nil
}

func (c *testServiceClient) Bytes(ctx context.Context) (CustomObject, error) {
	var defaultReturnVal CustomObject
	var returnVal *CustomObject
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Bytes"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/bytes"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "bytes failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "bytes response cannot be nil")
	}
	return *returnVal, nil
}

func (c *testServiceClient) Binary(ctx context.Context) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Binary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binary failed")
	}
	return resp.Body, nil
}

func (c *testServiceClient) MaybeBinary(ctx context.Context) (*io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("MaybeBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "maybeBinary failed")
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	return &resp.Body, nil
}
