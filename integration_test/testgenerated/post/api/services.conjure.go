// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"

	httpclient "github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	codecs "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	safejson "github.com/palantir/pkg/safejson"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type TestServiceClient interface {
	Echo(ctx context.Context, inputArg string) (string, error)
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) Echo(ctx context.Context, inputArg string) (returnVal string, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Echo"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	requestParams = append(requestParams, httpclient.WithRequestAppendFunc(codecs.JSON.ContentType(), func(out []byte) ([]byte, error) {
		out = safejson.AppendQuotedString(out, inputArg)
		return out, nil
	}))
	requestParams = append(requestParams, httpclient.WithResponseUnmarshalFunc(codecs.JSON.Accept(), func(data []byte) ([]byte, error) {
		ctx := context.TODO()
		if !gjson.ValidBytes(data) {
			return werror.ErrorWithContextParams(ctx, "invalid JSON for string")
		}
		value := gjson.ParseBytes(data)
		var err error
		if value.Type != gjson.String {
			err = werror.ErrorWithContextParams(ctx, "string expected JSON string")
			return err
		}
		returnVal = value.Str
	}))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "echo failed")
		return
	}
	return
}
