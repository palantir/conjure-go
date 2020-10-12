// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/rid"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/wrouter"
)

type TestService interface {
	Echo(ctx context.Context, cookieToken bearertoken.Token) error
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
	PostPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) (CustomObject, error)
	PostSafeParams(ctx context.Context, authHeader bearertoken.Token, myPathParam1Arg string, myPathParam2Arg bool, myBodyParamArg CustomObject, myQueryParam1Arg string, myQueryParam2Arg string, myQueryParam3Arg float64, myQueryParam4Arg *safelong.SafeLong, myQueryParam5Arg *string, myHeaderParam1Arg safelong.SafeLong, myHeaderParam2Arg *uuid.UUID) error
	Bytes(ctx context.Context) (CustomObject, error)
	GetBinary(ctx context.Context) (io.ReadCloser, error)
	PostBinary(ctx context.Context, myBytesArg io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg io.ReadCloser) error
	// An endpoint that uses go keywords
	Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong) error
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Get("Echo", "/echo", httpserver.NewConjureHandler(handler.HandleEcho)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Echo"))
	}
	if err := resource.Get("GetPathParam", "/path/string/{myPathParam}", httpserver.NewConjureHandler(handler.HandleGetPathParam)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "GetPathParam"))
	}
	if err := resource.Get("GetPathParamAlias", "/path/alias/{myPathParam}", httpserver.NewConjureHandler(handler.HandleGetPathParamAlias)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "GetPathParamAlias"))
	}
	if err := resource.Get("QueryParamList", "/pathNew", httpserver.NewConjureHandler(handler.HandleQueryParamList)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamList"))
	}
	if err := resource.Get("QueryParamListBoolean", "/booleanListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListBoolean)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListBoolean"))
	}
	if err := resource.Get("QueryParamListDateTime", "/dateTimeListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListDateTime)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListDateTime"))
	}
	if err := resource.Get("QueryParamListDouble", "/doubleListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListDouble)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListDouble"))
	}
	if err := resource.Get("QueryParamListInteger", "/intListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListInteger)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListInteger"))
	}
	if err := resource.Get("QueryParamListRid", "/ridListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListRid)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListRid"))
	}
	if err := resource.Get("QueryParamListSafeLong", "/safeLongListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListSafeLong)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListSafeLong"))
	}
	if err := resource.Get("QueryParamListString", "/stringListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListString)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListString"))
	}
	if err := resource.Get("QueryParamListUuid", "/uuidListQueryVar", httpserver.NewConjureHandler(handler.HandleQueryParamListUuid)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "QueryParamListUuid"))
	}
	if err := resource.Post("PostPathParam", "/path/{myPathParam1}/{myPathParam2}", httpserver.NewConjureHandler(handler.HandlePostPathParam)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "PostPathParam"))
	}
	if err := resource.Post("PostSafeParams", "/safe/{myPathParam1}/{myPathParam2}", httpserver.NewConjureHandler(handler.HandlePostSafeParams), wrouter.SafePathParams("myPathParam1"), wrouter.SafeHeaderParams("X-My-Header1-Abc"), wrouter.SafeQueryParams("query1", "myQueryParam2")); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "PostSafeParams"))
	}
	if err := resource.Get("Bytes", "/bytes", httpserver.NewConjureHandler(handler.HandleBytes)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Bytes"))
	}
	if err := resource.Get("GetBinary", "/binary", httpserver.NewConjureHandler(handler.HandleGetBinary)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "GetBinary"))
	}
	if err := resource.Post("PostBinary", "/binary", httpserver.NewConjureHandler(handler.HandlePostBinary)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "PostBinary"))
	}
	if err := resource.Put("PutBinary", "/binary", httpserver.NewConjureHandler(handler.HandlePutBinary)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "PutBinary"))
	}
	if err := resource.Post("Chan", "/chan/{var}", httpserver.NewConjureHandler(handler.HandleChan)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Chan"))
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleEcho(rw http.ResponseWriter, req *http.Request) error {
	authCookie, err := req.Cookie("PALANTIR_TOKEN")
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	cookieToken := bearertoken.Token(authCookie.Value)
	return t.impl.Echo(req.Context(), cookieToken)
}

func (t *testServiceHandler) HandleGetPathParam(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam, ok := pathParams["myPathParam"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "myPathParam"))
	}
	return t.impl.GetPathParam(req.Context(), bearertoken.Token(authHeader), myPathParam)
}

func (t *testServiceHandler) HandleGetPathParamAlias(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParamStr, ok := pathParams["myPathParam"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "myPathParam"))
	}
	var myPathParam StringAlias
	myPathParamQuote := strconv.Quote(myPathParamStr)
	if err := safejson.Unmarshal([]byte(myPathParamQuote), &myPathParam); err != nil {
		return werror.Wrap(err, "failed to unmarshal argument", werror.SafeParam("argName", "myPathParam"), werror.SafeParam("argType", "StringAlias"))
	}
	return t.impl.GetPathParamAlias(req.Context(), bearertoken.Token(authHeader), myPathParam)
}

func (t *testServiceHandler) HandleQueryParamList(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	myQueryParam1 := req.URL.Query()["myQueryParam1"]
	return t.impl.QueryParamList(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListBoolean(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []bool
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListBoolean(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListDateTime(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []datetime.DateTime
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := datetime.ParseDateTime(v)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListDateTime(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListDouble(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []float64
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListDouble(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListInteger(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []int
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListInteger(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListRid(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []rid.ResourceIdentifier
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := rid.ParseRID(v)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListRid(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListSafeLong(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []safelong.SafeLong
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := safelong.ParseSafeLong(v)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListSafeLong(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListString(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	myQueryParam1 := req.URL.Query()["myQueryParam1"]
	return t.impl.QueryParamListString(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListUuid(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	var myQueryParam1 []uuid.UUID
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := uuid.ParseUUID(v)
		if err != nil {
			return err
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListUuid(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandlePostPathParam(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam1, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "myPathParam1"))
	}
	myPathParam2Str, ok := pathParams["myPathParam2"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "myPathParam2"))
	}
	myPathParam2, err := strconv.ParseBool(myPathParam2Str)
	if err != nil {
		return err
	}
	myQueryParam1 := req.URL.Query().Get("query1")
	myQueryParam2 := req.URL.Query().Get("myQueryParam2")
	myQueryParam3, err := strconv.ParseFloat(req.URL.Query().Get("myQueryParam3"), 64)
	if err != nil {
		return err
	}
	var myQueryParam4 *safelong.SafeLong
	if myQueryParam4Str := req.URL.Query().Get("myQueryParam4"); myQueryParam4Str != "" {
		myQueryParam4Internal, err := safelong.ParseSafeLong(myQueryParam4Str)
		if err != nil {
			return err
		}
		myQueryParam4 = &myQueryParam4Internal
	}
	var myQueryParam5 *string
	if myQueryParam5Str := req.URL.Query().Get("myQueryParam5"); myQueryParam5Str != "" {
		myQueryParam5Internal := myQueryParam5Str
		myQueryParam5 = &myQueryParam5Internal
	}
	myHeaderParam1, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header1-Abc"))
	if err != nil {
		return err
	}
	var myHeaderParam2 *uuid.UUID
	if myHeaderParam2Str := req.Header.Get("X-My-Header2"); myHeaderParam2Str != "" {
		myHeaderParam2Internal, err := uuid.ParseUUID(myHeaderParam2Str)
		if err != nil {
			return err
		}
		myHeaderParam2 = &myHeaderParam2Internal
	}
	var myBodyParam CustomObject
	if err := codecs.JSON.Decode(req.Body, &myBodyParam); err != nil {
		return errors.NewWrappedError(errors.NewInvalidArgument(), err)
	}
	respArg, err := t.impl.PostPathParam(req.Context(), bearertoken.Token(authHeader), myPathParam1, myPathParam2, myBodyParam, myQueryParam1, myQueryParam2, myQueryParam3, myQueryParam4, myQueryParam5, myHeaderParam1, myHeaderParam2)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandlePostSafeParams(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.NewWrappedError(errors.NewPermissionDenied(), err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam1, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "myPathParam1"))
	}
	myPathParam2Str, ok := pathParams["myPathParam2"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "myPathParam2"))
	}
	myPathParam2, err := strconv.ParseBool(myPathParam2Str)
	if err != nil {
		return err
	}
	myQueryParam1 := req.URL.Query().Get("query1")
	myQueryParam2 := req.URL.Query().Get("myQueryParam2")
	myQueryParam3, err := strconv.ParseFloat(req.URL.Query().Get("myQueryParam3"), 64)
	if err != nil {
		return err
	}
	var myQueryParam4 *safelong.SafeLong
	if myQueryParam4Str := req.URL.Query().Get("myQueryParam4"); myQueryParam4Str != "" {
		myQueryParam4Internal, err := safelong.ParseSafeLong(myQueryParam4Str)
		if err != nil {
			return err
		}
		myQueryParam4 = &myQueryParam4Internal
	}
	var myQueryParam5 *string
	if myQueryParam5Str := req.URL.Query().Get("myQueryParam5"); myQueryParam5Str != "" {
		myQueryParam5Internal := myQueryParam5Str
		myQueryParam5 = &myQueryParam5Internal
	}
	myHeaderParam1, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header1-Abc"))
	if err != nil {
		return err
	}
	var myHeaderParam2 *uuid.UUID
	if myHeaderParam2Str := req.Header.Get("X-My-Header2"); myHeaderParam2Str != "" {
		myHeaderParam2Internal, err := uuid.ParseUUID(myHeaderParam2Str)
		if err != nil {
			return err
		}
		myHeaderParam2 = &myHeaderParam2Internal
	}
	var myBodyParam CustomObject
	if err := codecs.JSON.Decode(req.Body, &myBodyParam); err != nil {
		return errors.NewWrappedError(errors.NewInvalidArgument(), err)
	}
	return t.impl.PostSafeParams(req.Context(), bearertoken.Token(authHeader), myPathParam1, myPathParam2, myBodyParam, myQueryParam1, myQueryParam2, myQueryParam3, myQueryParam4, myQueryParam5, myHeaderParam1, myHeaderParam2)
}

func (t *testServiceHandler) HandleBytes(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.Bytes(req.Context())
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleGetBinary(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.GetBinary(req.Context())
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, respArg)
}

func (t *testServiceHandler) HandlePostBinary(rw http.ResponseWriter, req *http.Request) error {
	myBytes := req.Body
	respArg, err := t.impl.PostBinary(req.Context(), myBytes)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, respArg)
}

func (t *testServiceHandler) HandlePutBinary(rw http.ResponseWriter, req *http.Request) error {
	myBytes := req.Body
	return t.impl.PutBinary(req.Context(), myBytes)
}

func (t *testServiceHandler) HandleChan(rw http.ResponseWriter, req *http.Request) error {
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	var_, ok := pathParams["var"]
	if !ok {
		return werror.Wrap(errors.NewInvalidArgument(), "path param not present", werror.SafeParam("pathParamName", "var"))
	}
	type_ := req.URL.Query().Get("type")
	return_, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header2"))
	if err != nil {
		return err
	}
	var import_ map[string]string
	if err := codecs.JSON.Decode(req.Body, &import_); err != nil {
		return errors.NewWrappedError(errors.NewInvalidArgument(), err)
	}
	return t.impl.Chan(req.Context(), var_, import_, type_, return_)
}
