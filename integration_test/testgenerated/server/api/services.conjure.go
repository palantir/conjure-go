// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	httpclient "github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	codecs "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	bearertoken "github.com/palantir/pkg/bearertoken"
	datetime "github.com/palantir/pkg/datetime"
	rid "github.com/palantir/pkg/rid"
	safejson "github.com/palantir/pkg/safejson"
	safelong "github.com/palantir/pkg/safelong"
	uuid "github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	gjson "github.com/tidwall/gjson"
)

type TestServiceClient interface {
	Echo(ctx context.Context, cookieToken bearertoken.Token) error
	EchoStrings(ctx context.Context, bodyArg []string) ([]string, error)
	GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error
	GetPathParamAlias(ctx context.Context, authHeader bearertoken.Token, myPathParamArg StringAlias) error
	QueryParamList(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error
	QueryParamListBoolean(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []bool) error
	QueryParamListDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) error
	QueryParamListDouble(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []float64) error
	QueryParamListInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []int) error
	QueryParamListRid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []rid.ResourceIdentifier) error
	QueryParamListSafeLong(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []safelong.SafeLong) error
	QueryParamListString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) error
	QueryParamListUuid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []uuid.UUID) error
	PostPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error)
	PostSafeParams(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) error
	Bytes(ctx context.Context) (CustomObject, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) error
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
	// An endpoint that uses go keywords
	Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong) error
}

type testServiceClient struct {
	client httpclient.Client
}

func NewTestServiceClient(client httpclient.Client) TestServiceClient {
	return &testServiceClient{client: client}
}

func (c *testServiceClient) Echo(ctx context.Context, cookieToken bearertoken.Token) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Echo"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("PALANTIR_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "echo failed")
		return
	}
	return
}

func (c *testServiceClient) EchoStrings(ctx context.Context, bodyArg []string) (returnVal []string, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("EchoStrings"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/echo"))
	requestParams = append(requestParams, httpclient.WithRequestAppendFunc(codecs.JSON.ContentType(), func(out []byte) ([]byte, error) {
		out = append(out, '[')
		for i := range bodyArg {
			out = safejson.AppendQuotedString(out, bodyArg[i])
			if i < len(bodyArg)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
		if !gjson.ValidBytes(out) {
			return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
		}
		return out, nil
	}))
	requestParams = append(requestParams, httpclient.WithResponseUnmarshalFunc(codecs.JSON.Accept(), func(data []byte) ([]byte, error) {
		ctx := context.TODO()
		if !gjson.ValidBytes(data) {
			return werror.ErrorWithContextParams(ctx, "invalid JSON for list<string>")
		}
		value := gjson.ParseBytes(data)
		var err error
		if !value.IsArray() {
			err = werror.ErrorWithContextParams(ctx, "list<string> expected JSON array")
			return err
		}
		value.ForEach(func(_, value gjson.Result) bool {
			var listElement string
			if value.Type != gjson.String {
				err = werror.ErrorWithContextParams(ctx, "list<string> list element expected JSON string")
				return false
			}
			listElement = value.Str
			returnVal = append(returnVal, listElement)
			return err == nil
		})
		if err != nil {
			return err
		}
	}))
	requestParams = append(requestParams, httpclient.WithRequiredResponse())
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "echoStrings failed")
		return
	}
	return
}

func (c *testServiceClient) GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetPathParam"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/string/%s", url.PathEscape(fmt.Sprint(myPathParamArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "getPathParam failed")
		return
	}
	return
}

func (c *testServiceClient) GetPathParamAlias(ctx context.Context, authHeader bearertoken.Token, myPathParamArg StringAlias) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetPathParamAlias"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/alias/%s", url.PathEscape(fmt.Sprint(myPathParamArg))))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "getPathParamAlias failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamList(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamList"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/pathNew"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamList failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListBoolean(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []bool) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListBoolean"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/booleanListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListBoolean failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListDateTime(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []datetime.DateTime) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListDateTime"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/dateTimeListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListDateTime failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListDouble(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []float64) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListDouble"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/doubleListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListDouble failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListInteger(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []int) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListInteger"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/intListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListInteger failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListRid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []rid.ResourceIdentifier) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListRid"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/ridListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListRid failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListSafeLong(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []safelong.SafeLong) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListSafeLong"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/safeLongListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListSafeLong failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListString(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []string) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListString"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/stringListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListString failed")
		return
	}
	return
}

func (c *testServiceClient) QueryParamListUuid(ctx context.Context, authHeader bearertoken.Token, myQueryParam1Arg []uuid.UUID) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("QueryParamListUuid"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/uuidListQueryVar"))
	queryParams := make(url.Values)
	for _, v := range myQueryParam1Arg {
		queryParams.Add("myQueryParam1", fmt.Sprint(v))
	}
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "queryParamListUuid failed")
		return
	}
	return
}

func (c *testServiceClient) PostPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (returnVal CustomObject, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PostPathParam"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/path/%s/%s", url.PathEscape(fmt.Sprint(myPathParam1Arg)), url.PathEscape(fmt.Sprint(myPathParam2Arg))))
	requestParams = append(requestParams, httpclient.WithRequestAppendFunc(codecs.JSON.ContentType(), myBodyParamArg.AppendJSON))
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
	requestParams = append(requestParams, httpclient.WithResponseUnmarshalFunc(codecs.JSON.Accept(), returnVal.UnmarshalJSON))
	requestParams = append(requestParams, httpclient.WithRequiredResponse())
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "postPathParam failed")
		return
	}
	return
}

func (c *testServiceClient) PostSafeParams(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PostSafeParams"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/safe/%s/%s", url.PathEscape(fmt.Sprint(myPathParam1Arg)), url.PathEscape(fmt.Sprint(myPathParam2Arg))))
	requestParams = append(requestParams, httpclient.WithRequestAppendFunc(codecs.JSON.ContentType(), myBodyParamArg.AppendJSON))
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
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "postSafeParams failed")
		return
	}
	return
}

func (c *testServiceClient) Bytes(ctx context.Context) (returnVal CustomObject, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Bytes"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/bytes"))
	requestParams = append(requestParams, httpclient.WithResponseUnmarshalFunc(codecs.JSON.Accept(), returnVal.UnmarshalJSON))
	requestParams = append(requestParams, httpclient.WithRequiredResponse())
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "bytes failed")
		return
	}
	return
}

func (c *testServiceClient) GetBinary(ctx context.Context) (returnVal io.ReadCloser, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "getBinary failed")
		return
	}
	returnVal = resp.Body
	return
}

func (c *testServiceClient) PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (returnVal io.ReadCloser, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PostBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(myBytesArg))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "postBinary failed")
		return
	}
	returnVal = resp.Body
	return
}

func (c *testServiceClient) PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("PutBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("PUT"))
	requestParams = append(requestParams, httpclient.WithPathf("/binary"))
	requestParams = append(requestParams, httpclient.WithRawRequestBodyProvider(myBytesArg))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "putBinary failed")
		return
	}
	return
}

func (c *testServiceClient) GetOptionalBinary(ctx context.Context) (returnVal *io.ReadCloser, returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("GetOptionalBinary"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/optional/binary"))
	requestParams = append(requestParams, httpclient.WithRawResponseBody())
	resp, err := c.client.Do(ctx, requestParams...)
	if err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "getOptionalBinary failed")
		return
	}
	if resp.StatusCode == http.StatusNoContent {
		return
	}
	returnVal = &resp.Body
	return
}

func (c *testServiceClient) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong) (returnErr error) {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Chan"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithPathf("/chan/%s", url.PathEscape(fmt.Sprint(varArg))))
	requestParams = append(requestParams, httpclient.WithRequestAppendFunc(codecs.JSON.ContentType(), func(out []byte) ([]byte, error) {
		out = append(out, '{')
		{
			var i int
			for k, v := range importArg {
				{
					out = safejson.AppendQuotedString(out, k)
				}
				out = append(out, ':')
				{
					out = safejson.AppendQuotedString(out, v)
				}
				i++
				if i < len(importArg) {
					out = append(out, ',')
				}
			}
		}
		out = append(out, '}')
		if !gjson.ValidBytes(out) {
			return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
		}
		return out, nil
	}))
	requestParams = append(requestParams, httpclient.WithHeader("X-My-Header2", fmt.Sprint(returnArg)))
	queryParams := make(url.Values)
	queryParams.Set("type", fmt.Sprint(typeArg))
	requestParams = append(requestParams, httpclient.WithQueryValues(queryParams))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		returnErr = werror.WrapWithContextParams(ctx, err, "chan failed")
		return
	}
	return
}

type TestServiceClientWithAuth interface {
	Echo(ctx context.Context) error
	EchoStrings(ctx context.Context, bodyArg []string) ([]string, error)
	GetPathParam(ctx context.Context, myPathParamArg string) error
	GetPathParamAlias(ctx context.Context, myPathParamArg StringAlias) error
	QueryParamList(ctx context.Context, myQueryParam1Arg []string) error
	QueryParamListBoolean(ctx context.Context, myQueryParam1Arg []bool) error
	QueryParamListDateTime(ctx context.Context, myQueryParam1Arg []datetime.DateTime) error
	QueryParamListDouble(ctx context.Context, myQueryParam1Arg []float64) error
	QueryParamListInteger(ctx context.Context, myQueryParam1Arg []int) error
	QueryParamListRid(ctx context.Context, myQueryParam1Arg []rid.ResourceIdentifier) error
	QueryParamListSafeLong(ctx context.Context, myQueryParam1Arg []safelong.SafeLong) error
	QueryParamListString(ctx context.Context, myQueryParam1Arg []string) error
	QueryParamListUuid(ctx context.Context, myQueryParam1Arg []uuid.UUID) error
	PostPathParam(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error)
	PostSafeParams(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) error
	Bytes(ctx context.Context) (CustomObject, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) error
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
	// An endpoint that uses go keywords
	Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong) error
}

func NewTestServiceClientWithAuth(client TestServiceClient, authHeader bearertoken.Token, cookieToken bearertoken.Token) TestServiceClientWithAuth {
	return &testServiceClientWithAuth{client: client, authHeader: authHeader, cookieToken: cookieToken}
}

type testServiceClientWithAuth struct {
	client      TestServiceClient
	authHeader  bearertoken.Token
	cookieToken bearertoken.Token
}

func (c *testServiceClientWithAuth) Echo(ctx context.Context) (returnErr error) {
	return c.client.Echo(ctx, c.cookieToken)
}

func (c *testServiceClientWithAuth) EchoStrings(ctx context.Context, bodyArg []string) (returnVal []string, returnErr error) {
	return c.client.EchoStrings(ctx, bodyArg)
}

func (c *testServiceClientWithAuth) GetPathParam(ctx context.Context, myPathParamArg string) (returnErr error) {
	return c.client.GetPathParam(ctx, c.authHeader, myPathParamArg)
}

func (c *testServiceClientWithAuth) GetPathParamAlias(ctx context.Context, myPathParamArg StringAlias) (returnErr error) {
	return c.client.GetPathParamAlias(ctx, c.authHeader, myPathParamArg)
}

func (c *testServiceClientWithAuth) QueryParamList(ctx context.Context, myQueryParam1Arg []string) (returnErr error) {
	return c.client.QueryParamList(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListBoolean(ctx context.Context, myQueryParam1Arg []bool) (returnErr error) {
	return c.client.QueryParamListBoolean(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListDateTime(ctx context.Context, myQueryParam1Arg []datetime.DateTime) (returnErr error) {
	return c.client.QueryParamListDateTime(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListDouble(ctx context.Context, myQueryParam1Arg []float64) (returnErr error) {
	return c.client.QueryParamListDouble(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListInteger(ctx context.Context, myQueryParam1Arg []int) (returnErr error) {
	return c.client.QueryParamListInteger(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListRid(ctx context.Context, myQueryParam1Arg []rid.ResourceIdentifier) (returnErr error) {
	return c.client.QueryParamListRid(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListSafeLong(ctx context.Context, myQueryParam1Arg []safelong.SafeLong) (returnErr error) {
	return c.client.QueryParamListSafeLong(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListString(ctx context.Context, myQueryParam1Arg []string) (returnErr error) {
	return c.client.QueryParamListString(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) QueryParamListUuid(ctx context.Context, myQueryParam1Arg []uuid.UUID) (returnErr error) {
	return c.client.QueryParamListUuid(ctx, c.authHeader, myQueryParam1Arg)
}

func (c *testServiceClientWithAuth) PostPathParam(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myQueryParam6Arg OptionalIntegerAlias, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (returnVal CustomObject, returnErr error) {
	return c.client.PostPathParam(ctx, c.authHeader, myPathParam1Arg, myPathParam2Arg, myBodyParamArg, myQueryParam1Arg, myQueryParam2Arg, myQueryParam3Arg, myQueryParam4Arg, myQueryParam5Arg, myQueryParam6Arg, myHeaderParam1Arg, myHeaderParam2Arg)
}

func (c *testServiceClientWithAuth) PostSafeParams(ctx context.Context, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (returnErr error) {
	return c.client.PostSafeParams(ctx, c.authHeader, myPathParam1Arg, myPathParam2Arg, myBodyParamArg, myQueryParam1Arg, myQueryParam2Arg, myQueryParam3Arg, myQueryParam4Arg, myQueryParam5Arg, myHeaderParam1Arg, myHeaderParam2Arg)
}

func (c *testServiceClientWithAuth) Bytes(ctx context.Context) (returnVal CustomObject, returnErr error) {
	return c.client.Bytes(ctx)
}

func (c *testServiceClientWithAuth) GetBinary(ctx context.Context) (returnVal io.ReadCloser, returnErr error) {
	return c.client.GetBinary(ctx)
}

func (c *testServiceClientWithAuth) PostBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (returnVal io.ReadCloser, returnErr error) {
	return c.client.PostBinary(ctx, myBytesArg)
}

func (c *testServiceClientWithAuth) PutBinary(ctx context.Context, myBytesArg func() io.ReadCloser) (returnErr error) {
	return c.client.PutBinary(ctx, myBytesArg)
}

func (c *testServiceClientWithAuth) GetOptionalBinary(ctx context.Context) (returnVal *io.ReadCloser, returnErr error) {
	return c.client.GetOptionalBinary(ctx)
}

func (c *testServiceClientWithAuth) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong) (returnErr error) {
	return c.client.Chan(ctx, varArg, importArg, typeArg, returnArg)
}
