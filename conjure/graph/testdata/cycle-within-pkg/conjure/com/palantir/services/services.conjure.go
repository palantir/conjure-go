// This file was generated by Conjure and should not be manually edited.

package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/conjure-go/v6/conjure/graph/testdata/cycle-within-pkg/conjure/com/palantir/buzz"
	"github.com/palantir/conjure-go/v6/conjure/graph/testdata/cycle-within-pkg/conjure/com/palantir/foo"
	"github.com/palantir/pkg/bearertoken"
	werror "github.com/palantir/witchcraft-go-error"
)

type MyServiceClient interface {
	Endpoint1(ctx context.Context, authHeader bearertoken.Token, arg1Arg buzz.Type1) (foo.Type4, error)
	Endpoint2(ctx context.Context, authHeader bearertoken.Token, arg1Arg foo.Type1) error
}

type myServiceClient struct {
	client httpclient.Client
}

func NewMyServiceClient(client httpclient.Client) MyServiceClient {
	return &myServiceClient{client: client}
}

func (c *myServiceClient) Endpoint1(ctx context.Context, authHeader bearertoken.Token, arg1Arg buzz.Type1) (foo.Type4, error) {
	var defaultReturnVal foo.Type4
	var returnVal *foo.Type4
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Endpoint1"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/endpoint1", url.PathEscape(fmt.Sprint(arg1Arg))))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "Endpoint1 failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "Endpoint1 response cannot be nil")
	}
	return *returnVal, nil
}

func (c *myServiceClient) Endpoint2(ctx context.Context, authHeader bearertoken.Token, arg1Arg foo.Type1) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Endpoint2"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/endpoint2"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(arg1Arg))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "Endpoint2 failed")
	}
	return nil
}

type MyServiceClientWithAuth interface {
	Endpoint1(ctx context.Context, arg1Arg buzz.Type1) (foo.Type4, error)
	Endpoint2(ctx context.Context, arg1Arg foo.Type1) error
}

func NewMyServiceClientWithAuth(client MyServiceClient, authHeader bearertoken.Token) MyServiceClientWithAuth {
	return &myServiceClientWithAuth{client: client, authHeader: authHeader}
}

type myServiceClientWithAuth struct {
	client     MyServiceClient
	authHeader bearertoken.Token
}

func (c *myServiceClientWithAuth) Endpoint1(ctx context.Context, arg1Arg buzz.Type1) (foo.Type4, error) {
	return c.client.Endpoint1(ctx, c.authHeader, arg1Arg)
}

func (c *myServiceClientWithAuth) Endpoint2(ctx context.Context, arg1Arg foo.Type1) error {
	return c.client.Endpoint2(ctx, c.authHeader, arg1Arg)
}

func NewMyServiceClientWithTokenProvider(client MyServiceClient, tokenProvider httpclient.TokenProvider) MyServiceClientWithAuth {
	return &myServiceClientWithTokenProvider{client: client, tokenProvider: tokenProvider}
}

type myServiceClientWithTokenProvider struct {
	client        MyServiceClient
	tokenProvider httpclient.TokenProvider
}

func (c *myServiceClientWithTokenProvider) Endpoint1(ctx context.Context, arg1Arg buzz.Type1) (foo.Type4, error) {
	var defaultReturnVal foo.Type4
	token, err := c.tokenProvider(ctx)
	if err != nil {
		return defaultReturnVal, err
	}
	return c.client.Endpoint1(ctx, bearertoken.Token(token), arg1Arg)
}

func (c *myServiceClientWithTokenProvider) Endpoint2(ctx context.Context, arg1Arg foo.Type1) error {
	token, err := c.tokenProvider(ctx)
	if err != nil {
		return err
	}
	return c.client.Endpoint2(ctx, bearertoken.Token(token), arg1Arg)
}
