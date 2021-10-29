// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"io"
	"net/http"
	"strconv"

	codecs "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	errors "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	httpserver "github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	bearertoken "github.com/palantir/pkg/bearertoken"
	datetime "github.com/palantir/pkg/datetime"
	rid "github.com/palantir/pkg/rid"
	safejson "github.com/palantir/pkg/safejson"
	safelong "github.com/palantir/pkg/safelong"
	uuid "github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	wresource "github.com/palantir/witchcraft-go-server/v2/witchcraft/wresource"
	wrouter "github.com/palantir/witchcraft-go-server/v2/wrouter"
)

type TestService interface {
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
	PostBinary(ctx context.Context, myBytesArg io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg io.ReadCloser) error
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
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
	if err := resource.Get("Echo", "/echo", httpserver.NewJSONHandler(handler.HandleEcho, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add echo route")
	}
	if err := resource.Post("EchoStrings", "/echo", httpserver.NewJSONHandler(handler.HandleEchoStrings, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add echoStrings route")
	}
	if err := resource.Get("GetPathParam", "/path/string/{myPathParam}", httpserver.NewJSONHandler(handler.HandleGetPathParam, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add getPathParam route")
	}
	if err := resource.Get("GetPathParamAlias", "/path/alias/{myPathParam}", httpserver.NewJSONHandler(handler.HandleGetPathParamAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add getPathParamAlias route")
	}
	if err := resource.Get("QueryParamList", "/pathNew", httpserver.NewJSONHandler(handler.HandleQueryParamList, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamList route")
	}
	if err := resource.Get("QueryParamListBoolean", "/booleanListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListBoolean, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListBoolean route")
	}
	if err := resource.Get("QueryParamListDateTime", "/dateTimeListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListDateTime, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListDateTime route")
	}
	if err := resource.Get("QueryParamListDouble", "/doubleListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListDouble, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListDouble route")
	}
	if err := resource.Get("QueryParamListInteger", "/intListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListInteger, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListInteger route")
	}
	if err := resource.Get("QueryParamListRid", "/ridListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListRid, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListRid route")
	}
	if err := resource.Get("QueryParamListSafeLong", "/safeLongListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListSafeLong, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListSafeLong route")
	}
	if err := resource.Get("QueryParamListString", "/stringListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListString, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListString route")
	}
	if err := resource.Get("QueryParamListUuid", "/uuidListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListUuid, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add queryParamListUuid route")
	}
	if err := resource.Post("PostPathParam", "/path/{myPathParam1}/{myPathParam2}", httpserver.NewJSONHandler(handler.HandlePostPathParam, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add postPathParam route")
	}
	if err := resource.Post("PostSafeParams", "/safe/{myPathParam1}/{myPathParam2}", httpserver.NewJSONHandler(handler.HandlePostSafeParams, httpserver.StatusCodeMapper, httpserver.ErrHandler), wrouter.SafePathParams("myPathParam1"), wrouter.SafeHeaderParams("X-My-Header1-Abc"), wrouter.SafeQueryParams("query1"), wrouter.SafeQueryParams("myQueryParam2")); err != nil {
		return werror.Wrap(err, "failed to add postSafeParams route")
	}
	if err := resource.Get("Bytes", "/bytes", httpserver.NewJSONHandler(handler.HandleBytes, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add bytes route")
	}
	if err := resource.Get("GetBinary", "/binary", httpserver.NewJSONHandler(handler.HandleGetBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add getBinary route")
	}
	if err := resource.Post("PostBinary", "/binary", httpserver.NewJSONHandler(handler.HandlePostBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add postBinary route")
	}
	if err := resource.Put("PutBinary", "/binary", httpserver.NewJSONHandler(handler.HandlePutBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add putBinary route")
	}
	if err := resource.Get("GetOptionalBinary", "/optional/binary", httpserver.NewJSONHandler(handler.HandleGetOptionalBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add getOptionalBinary route")
	}
	if err := resource.Post("Chan", "/chan/{var}", httpserver.NewJSONHandler(handler.HandleChan, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add chan route")
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleEcho(rw http.ResponseWriter, req *http.Request) error {
	authCookie, err := req.Cookie("PALANTIR_TOKEN")
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	cookieToken := bearertoken.Token(authCookie.Value)
	return t.impl.Echo(req.Context(), cookieToken)
}

func (t *testServiceHandler) HandleEchoStrings(rw http.ResponseWriter, req *http.Request) error {
	var body []string
	if err := codecs.JSON.Decode(req.Body, &body); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	respArg, err := t.impl.EchoStrings(req.Context(), body)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleGetPathParam(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam, ok := pathParams["myPathParam"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam\" not present")
	}
	return t.impl.GetPathParam(req.Context(), bearertoken.Token(authHeader), myPathParam)
}

func (t *testServiceHandler) HandleGetPathParamAlias(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParamStr, ok := pathParams["myPathParam"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam\" not present")
	}
	var myPathParam StringAlias
	if err := safejson.Unmarshal([]byte(strconv.Quote(myPathParamStr)), &myPathParam); err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to unmarshal \"myPathParam\" param")
	}
	return t.impl.GetPathParamAlias(req.Context(), bearertoken.Token(authHeader), myPathParam)
}

func (t *testServiceHandler) HandleQueryParamList(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	myQueryParam1 := req.URL.Query()["myQueryParam1"]
	return t.impl.QueryParamList(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListBoolean(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []bool
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.ParseBool(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as boolean")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListBoolean(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListDateTime(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []datetime.DateTime
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := datetime.ParseDateTime(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as datetime")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListDateTime(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListDouble(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []float64
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as double")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListDouble(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListInteger(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []int
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.Atoi(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as integer")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListInteger(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListRid(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []rid.ResourceIdentifier
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := rid.ParseRID(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as rid")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListRid(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListSafeLong(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []safelong.SafeLong
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := safelong.ParseSafeLong(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as safelong")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListSafeLong(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListString(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	myQueryParam1 := req.URL.Query()["myQueryParam1"]
	return t.impl.QueryParamListString(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandleQueryParamListUuid(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1 []uuid.UUID
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := uuid.ParseUUID(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as uuid")
		}
		myQueryParam1 = append(myQueryParam1, convertedVal)
	}
	return t.impl.QueryParamListUuid(req.Context(), bearertoken.Token(authHeader), myQueryParam1)
}

func (t *testServiceHandler) HandlePostPathParam(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam1, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam1\" not present")
	}
	myPathParam2Str, ok := pathParams["myPathParam2"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam2\" not present")
	}
	myPathParam2, err := strconv.ParseBool(myPathParam2Str)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myPathParam2\" as boolean")
	}
	myQueryParam1 := req.URL.Query().Get("query1")
	myQueryParam2 := req.URL.Query().Get("myQueryParam2")
	myQueryParam3, err := strconv.ParseFloat(req.URL.Query().Get("myQueryParam3"), 64)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam3\" as double")
	}
	var myQueryParam4 *safelong.SafeLong
	if myQueryParam4Str := req.URL.Query().Get("myQueryParam4"); myQueryParam4Str != "" {
		myQueryParam4Internal, err := safelong.ParseSafeLong(myQueryParam4Str)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam4\" as safelong")
		}
		myQueryParam4 = &myQueryParam4Internal
	}
	var myQueryParam5 *string
	if myQueryParam5Str := req.URL.Query().Get("myQueryParam5"); myQueryParam5Str != "" {
		myQueryParam5Internal := myQueryParam5Str
		myQueryParam5 = &myQueryParam5Internal
	}
	var myQueryParam6 OptionalIntegerAlias
	if err := safejson.Unmarshal([]byte(strconv.Quote(req.URL.Query().Get("myQueryParam6"))), &myQueryParam6); err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to unmarshal \"myQueryParam6\" param")
	}
	myHeaderParam1, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header1-Abc"))
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam1\" as safelong")
	}
	var myHeaderParam2 *uuid.UUID
	if myHeaderParam2Str := req.Header.Get("X-My-Header2"); myHeaderParam2Str != "" {
		myHeaderParam2Internal, err := uuid.ParseUUID(myHeaderParam2Str)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam2\" as uuid")
		}
		myHeaderParam2 = &myHeaderParam2Internal
	}
	var myBodyParam CustomObject
	if err := codecs.JSON.Decode(req.Body, &myBodyParam); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	respArg, err := t.impl.PostPathParam(req.Context(), bearertoken.Token(authHeader), myPathParam1, myPathParam2, myBodyParam, myQueryParam1, myQueryParam2, myQueryParam3, myQueryParam4, myQueryParam5, myQueryParam6, myHeaderParam1, myHeaderParam2)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandlePostSafeParams(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam1, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam1\" not present")
	}
	myPathParam2Str, ok := pathParams["myPathParam2"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam2\" not present")
	}
	myPathParam2, err := strconv.ParseBool(myPathParam2Str)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myPathParam2\" as boolean")
	}
	myQueryParam1 := req.URL.Query().Get("query1")
	myQueryParam2 := req.URL.Query().Get("myQueryParam2")
	myQueryParam3, err := strconv.ParseFloat(req.URL.Query().Get("myQueryParam3"), 64)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam3\" as double")
	}
	var myQueryParam4 *safelong.SafeLong
	if myQueryParam4Str := req.URL.Query().Get("myQueryParam4"); myQueryParam4Str != "" {
		myQueryParam4Internal, err := safelong.ParseSafeLong(myQueryParam4Str)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam4\" as safelong")
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
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam1\" as safelong")
	}
	var myHeaderParam2 *uuid.UUID
	if myHeaderParam2Str := req.Header.Get("X-My-Header2"); myHeaderParam2Str != "" {
		myHeaderParam2Internal, err := uuid.ParseUUID(myHeaderParam2Str)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam2\" as uuid")
		}
		myHeaderParam2 = &myHeaderParam2Internal
	}
	var myBodyParam CustomObject
	if err := codecs.JSON.Decode(req.Body, &myBodyParam); err != nil {
		return errors.WrapWithInvalidArgument(err)
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

func (t *testServiceHandler) HandleGetOptionalBinary(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.GetOptionalBinary(req.Context())
	if err != nil {
		return err
	}
	if respArg == nil {
		rw.WriteHeader(http.StatusNoContent)
		return nil
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, *respArg)
}

func (t *testServiceHandler) HandleChan(rw http.ResponseWriter, req *http.Request) error {
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	var_, ok := pathParams["var"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"var\" not present")
	}
	type_ := req.URL.Query().Get("type")
	return_, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header2"))
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"return\" as safelong")
	}
	var import_ map[string]string
	if err := codecs.JSON.Decode(req.Body, &import_); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	return t.impl.Chan(req.Context(), var_, import_, type_, return_)
}
