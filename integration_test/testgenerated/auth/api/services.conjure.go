// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"fmt"

	httpclient "github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	bearertoken "github.com/palantir/pkg/bearertoken"
	safejson "github.com/palantir/pkg/safejson"
	werror "github.com/palantir/witchcraft-go-error"
)

type BothAuthServiceClient interface {
	Default(ctx context.Context, authHeader bearertoken.Token) (string, error)
	Cookie(ctx context.Context, cookieToken bearertoken.Token) error
	None(ctx context.Context) error
	WithArg(ctx context.Context, authHeader bearertoken.Token, argArg string) error
}

type bothAuthServiceClient struct {
	client httpclient.Client
}

func NewBothAuthServiceClient(client httpclient.Client) BothAuthServiceClient {
	return &bothAuthServiceClient{client: client}
}

func (c *bothAuthServiceClient) Default(ctx context.Context, authHeader bearertoken.Token) (string, error) {
	var defaultReturnVal string
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Default"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/default"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "default failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "default response cannot be nil")
	}
	return *returnVal, nil
}

func (c *bothAuthServiceClient) Cookie(ctx context.Context, cookieToken bearertoken.Token) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Cookie"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("P_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/cookie"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "cookie failed")
	}
	return nil
}

func (c *bothAuthServiceClient) None(ctx context.Context) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("None"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/none"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "none failed")
	}
	return nil
}

func (c *bothAuthServiceClient) WithArg(ctx context.Context, authHeader bearertoken.Token, argArg string) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("WithArg"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("POST"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/withArg"))
	requestParams = append(requestParams, httpclient.WithJSONRequest(safejson.AppendFunc(func(out []byte) ([]byte, error) {
		out = safejson.AppendQuotedString(out, argArg)
		return out, nil
	})))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "withArg failed")
	}
	return nil
}

type BothAuthServiceClientWithAuth interface {
	Default(ctx context.Context) (string, error)
	Cookie(ctx context.Context) error
	None(ctx context.Context) error
	WithArg(ctx context.Context, argArg string) error
}

func NewBothAuthServiceClientWithAuth(client BothAuthServiceClient, authHeader bearertoken.Token, cookieToken bearertoken.Token) BothAuthServiceClientWithAuth {
	return &bothAuthServiceClientWithAuth{client: client, authHeader: authHeader, cookieToken: cookieToken}
}

type bothAuthServiceClientWithAuth struct {
	client      BothAuthServiceClient
	authHeader  bearertoken.Token
	cookieToken bearertoken.Token
}

func (c *bothAuthServiceClientWithAuth) Default(ctx context.Context) (string, error) {
	return c.client.Default(ctx, c.authHeader)
}

func (c *bothAuthServiceClientWithAuth) Cookie(ctx context.Context) error {
	return c.client.Cookie(ctx, c.cookieToken)
}

func (c *bothAuthServiceClientWithAuth) None(ctx context.Context) error {
	return c.client.None(ctx)
}

func (c *bothAuthServiceClientWithAuth) WithArg(ctx context.Context, argArg string) error {
	return c.client.WithArg(ctx, c.authHeader, argArg)
}

type CookieAuthServiceClient interface {
	Cookie(ctx context.Context, cookieToken bearertoken.Token) error
}

type cookieAuthServiceClient struct {
	client httpclient.Client
}

func NewCookieAuthServiceClient(client httpclient.Client) CookieAuthServiceClient {
	return &cookieAuthServiceClient{client: client}
}

func (c *cookieAuthServiceClient) Cookie(ctx context.Context, cookieToken bearertoken.Token) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Cookie"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Cookie", fmt.Sprint("P_TOKEN=", cookieToken)))
	requestParams = append(requestParams, httpclient.WithPathf("/cookie"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "cookie failed")
	}
	return nil
}

type CookieAuthServiceClientWithAuth interface {
	Cookie(ctx context.Context) error
}

func NewCookieAuthServiceClientWithAuth(client CookieAuthServiceClient, cookieToken bearertoken.Token) CookieAuthServiceClientWithAuth {
	return &cookieAuthServiceClientWithAuth{client: client, cookieToken: cookieToken}
}

type cookieAuthServiceClientWithAuth struct {
	client      CookieAuthServiceClient
	cookieToken bearertoken.Token
}

func (c *cookieAuthServiceClientWithAuth) Cookie(ctx context.Context) error {
	return c.client.Cookie(ctx, c.cookieToken)
}

func NewCookieAuthServiceClientWithTokenProvider(client CookieAuthServiceClient, tokenProvider httpclient.TokenProvider) CookieAuthServiceClientWithAuth {
	return &cookieAuthServiceClientWithTokenProvider{client: client, tokenProvider: tokenProvider}
}

type cookieAuthServiceClientWithTokenProvider struct {
	client        CookieAuthServiceClient
	tokenProvider httpclient.TokenProvider
}

func (c *cookieAuthServiceClientWithTokenProvider) Cookie(ctx context.Context) error {
	token, err := c.tokenProvider(ctx)
	if err != nil {
		return err
	}
	return c.client.Cookie(ctx, bearertoken.Token(token))
}

type HeaderAuthServiceClient interface {
	Default(ctx context.Context, authHeader bearertoken.Token) (string, error)
}

type headerAuthServiceClient struct {
	client httpclient.Client
}

func NewHeaderAuthServiceClient(client httpclient.Client) HeaderAuthServiceClient {
	return &headerAuthServiceClient{client: client}
}

func (c *headerAuthServiceClient) Default(ctx context.Context, authHeader bearertoken.Token) (string, error) {
	var defaultReturnVal string
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Default"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/default"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "default failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "default response cannot be nil")
	}
	return *returnVal, nil
}

type HeaderAuthServiceClientWithAuth interface {
	Default(ctx context.Context) (string, error)
}

func NewHeaderAuthServiceClientWithAuth(client HeaderAuthServiceClient, authHeader bearertoken.Token) HeaderAuthServiceClientWithAuth {
	return &headerAuthServiceClientWithAuth{client: client, authHeader: authHeader}
}

type headerAuthServiceClientWithAuth struct {
	client     HeaderAuthServiceClient
	authHeader bearertoken.Token
}

func (c *headerAuthServiceClientWithAuth) Default(ctx context.Context) (string, error) {
	return c.client.Default(ctx, c.authHeader)
}

func NewHeaderAuthServiceClientWithTokenProvider(client HeaderAuthServiceClient, tokenProvider httpclient.TokenProvider) HeaderAuthServiceClientWithAuth {
	return &headerAuthServiceClientWithTokenProvider{client: client, tokenProvider: tokenProvider}
}

type headerAuthServiceClientWithTokenProvider struct {
	client        HeaderAuthServiceClient
	tokenProvider httpclient.TokenProvider
}

func (c *headerAuthServiceClientWithTokenProvider) Default(ctx context.Context) (string, error) {
	var defaultReturnVal string
	token, err := c.tokenProvider(ctx)
	if err != nil {
		return defaultReturnVal, err
	}
	return c.client.Default(ctx, bearertoken.Token(token))
}

type SomeHeaderAuthServiceClient interface {
	Default(ctx context.Context, authHeader bearertoken.Token) (string, error)
	None(ctx context.Context) error
}

type someHeaderAuthServiceClient struct {
	client httpclient.Client
}

func NewSomeHeaderAuthServiceClient(client httpclient.Client) SomeHeaderAuthServiceClient {
	return &someHeaderAuthServiceClient{client: client}
}

func (c *someHeaderAuthServiceClient) Default(ctx context.Context, authHeader bearertoken.Token) (string, error) {
	var defaultReturnVal string
	var returnVal *string
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("Default"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithHeader("Authorization", fmt.Sprint("Bearer ", authHeader)))
	requestParams = append(requestParams, httpclient.WithPathf("/default"))
	requestParams = append(requestParams, httpclient.WithJSONResponse(&returnVal))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return defaultReturnVal, werror.WrapWithContextParams(ctx, err, "default failed")
	}
	if returnVal == nil {
		return defaultReturnVal, werror.ErrorWithContextParams(ctx, "default response cannot be nil")
	}
	return *returnVal, nil
}

func (c *someHeaderAuthServiceClient) None(ctx context.Context) error {
	var requestParams []httpclient.RequestParam
	requestParams = append(requestParams, httpclient.WithRPCMethodName("None"))
	requestParams = append(requestParams, httpclient.WithRequestMethod("GET"))
	requestParams = append(requestParams, httpclient.WithPathf("/none"))
	if _, err := c.client.Do(ctx, requestParams...); err != nil {
		return werror.WrapWithContextParams(ctx, err, "none failed")
	}
	return nil
}

type SomeHeaderAuthServiceClientWithAuth interface {
	Default(ctx context.Context) (string, error)
	None(ctx context.Context) error
}

func NewSomeHeaderAuthServiceClientWithAuth(client SomeHeaderAuthServiceClient, authHeader bearertoken.Token) SomeHeaderAuthServiceClientWithAuth {
	return &someHeaderAuthServiceClientWithAuth{client: client, authHeader: authHeader}
}

type someHeaderAuthServiceClientWithAuth struct {
	client     SomeHeaderAuthServiceClient
	authHeader bearertoken.Token
}

func (c *someHeaderAuthServiceClientWithAuth) Default(ctx context.Context) (string, error) {
	return c.client.Default(ctx, c.authHeader)
}

func (c *someHeaderAuthServiceClientWithAuth) None(ctx context.Context) error {
	return c.client.None(ctx)
}

func NewSomeHeaderAuthServiceClientWithTokenProvider(client SomeHeaderAuthServiceClient, tokenProvider httpclient.TokenProvider) SomeHeaderAuthServiceClientWithAuth {
	return &someHeaderAuthServiceClientWithTokenProvider{client: client, tokenProvider: tokenProvider}
}

type someHeaderAuthServiceClientWithTokenProvider struct {
	client        SomeHeaderAuthServiceClient
	tokenProvider httpclient.TokenProvider
}

func (c *someHeaderAuthServiceClientWithTokenProvider) Default(ctx context.Context) (string, error) {
	var defaultReturnVal string
	token, err := c.tokenProvider(ctx)
	if err != nil {
		return defaultReturnVal, err
	}
	return c.client.Default(ctx, bearertoken.Token(token))
}

func (c *someHeaderAuthServiceClientWithTokenProvider) None(ctx context.Context) error {
	return c.client.None(ctx)
}
