// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"

	httpclient "github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	safejson "github.com/palantir/pkg/safejson"
	werror "github.com/palantir/witchcraft-go-error"
)

type TestServiceClient interface {
	BinaryAlias(ctx context.Context, bodyArg func() io.ReadCloser) (io.ReadCloser, error)
	BinaryAliasOptional(ctx context.Context) (*io.ReadCloser, error)
	BinaryAliasAlias(ctx context.Context, bodyArg func() io.ReadCloser) (*io.ReadCloser, error)
	Binary(ctx context.Context, bodyArg func() io.ReadCloser) (io.ReadCloser, error)
	BinaryOptional(ctx context.Context) (*io.ReadCloser, error)
	BinaryList(ctx context.Context, bodyArg [][]byte) ([][]byte, error)
	Bytes(ctx context.Context, bodyArg CustomObject) (CustomObject, error)
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) BinaryAlias(ctx context.Context, bodyArg func() io.ReadCloser) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("BinaryAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binaryAlias"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(bodyArg))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binaryAlias failed")
	}
	return resp.Body, nil
}

func (c *testServiceClient) BinaryAliasOptional(ctx context.Context) (*io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("BinaryAliasOptional"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binaryAliasOptional"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binaryAliasOptional failed")
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	return &resp.Body, nil
}

func (c *testServiceClient) BinaryAliasAlias(ctx context.Context, bodyArg func() io.ReadCloser) (*io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("BinaryAliasAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binaryAliasAlias"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(bodyArg))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binaryAliasAlias failed")
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	return &resp.Body, nil
}

func (c *testServiceClient) Binary(ctx context.Context, bodyArg func() io.ReadCloser) (io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Binary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(bodyArg))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binary failed")
	}
	return resp.Body, nil
}

func (c *testServiceClient) BinaryOptional(ctx context.Context) (*io.ReadCloser, error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("BinaryOptional"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binaryOptional"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binaryOptional failed")
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	return &resp.Body, nil
}

func (c *testServiceClient) BinaryList(ctx context.Context, bodyArg [][]byte) ([][]byte, error) {
	var returnVal [][]byte
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("BinaryList"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binaryList"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(safejson.AppendFunc(func(out []byte) ([]byte, error) {
		out = append(out, '[')
		{
			for i := range bodyArg {
				out = append(out, '"')
				if len(bodyArg[i]) > 0 {
					b64out := make([]byte, 0, base64.StdEncoding.EncodedLen(len(bodyArg[i])))
					base64.StdEncoding.Encode(b64out, bodyArg[i])
					out = append(out, b64out...)
				}
				out = append(out, '"')
				if i < len(bodyArg)-1 {
					out = append(out, ',')
				}
			}
		}
		out = append(out, ']')
		return out, nil
	})))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "binaryList failed")
	}
	if returnVal == nil {
		return nil, werror.ErrorWithContextParams(ctx, "binaryList response cannot be nil")
	}
	return returnVal, nil
}

func (c *testServiceClient) Bytes(ctx context.Context, bodyArg CustomObject) (CustomObject, error) {
	var defaultReturnVal CustomObject
	var returnVal *CustomObject
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Bytes"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/bytes"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(bodyArg))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "bytes failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "bytes response cannot be nil")
	}
	return *returnVal, nil
}