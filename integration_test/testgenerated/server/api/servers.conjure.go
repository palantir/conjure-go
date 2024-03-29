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
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/pkg/uuid"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/v2/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
)

type TestService interface {
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
	PostBinary(ctx context.Context, myBytesArg io.ReadCloser) (io.ReadCloser, error)
	PutBinary(ctx context.Context, myBytesArg io.ReadCloser) error
	GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error)
	// An endpoint that uses go keywords
	Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService, routerParams ...wrouter.RouteParam) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Get("Echo", "/echo", httpserver.NewJSONHandler(handler.HandleEcho, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add echo route")
	}
	if err := resource.Post("EchoStrings", "/echo", httpserver.NewJSONHandler(handler.HandleEchoStrings, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add echoStrings route")
	}
	if err := resource.Post("EchoCustomObject", "/echoCustomObject", httpserver.NewJSONHandler(handler.HandleEchoCustomObject, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add echoCustomObject route")
	}
	if err := resource.Post("EchoOptionalAlias", "/optional/alias", httpserver.NewJSONHandler(handler.HandleEchoOptionalAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add echoOptionalAlias route")
	}
	if err := resource.Post("EchoOptionalListAlias", "/optional/list-alias", httpserver.NewJSONHandler(handler.HandleEchoOptionalListAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add echoOptionalListAlias route")
	}
	if err := resource.Get("GetPathParam", "/path/string/{myPathParam}", httpserver.NewJSONHandler(handler.HandleGetPathParam, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add getPathParam route")
	}
	if err := resource.Get("GetPathParamAlias", "/path/alias/{myPathParam}", httpserver.NewJSONHandler(handler.HandleGetPathParamAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler), append(routerParams, wrouter.ForbiddenPathParams("myPathParam"))...); err != nil {
		return werror.Wrap(err, "failed to add getPathParamAlias route")
	}
	if err := resource.Get("QueryParamList", "/pathNew", httpserver.NewJSONHandler(handler.HandleQueryParamList, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamList route")
	}
	if err := resource.Get("QueryParamListBoolean", "/booleanListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListBoolean, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListBoolean route")
	}
	if err := resource.Get("QueryParamListDateTime", "/dateTimeListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListDateTime, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListDateTime route")
	}
	if err := resource.Get("QueryParamSetDateTime", "/dateTimeSetQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamSetDateTime, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamSetDateTime route")
	}
	if err := resource.Get("QueryParamListDouble", "/doubleListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListDouble, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListDouble route")
	}
	if err := resource.Get("QueryParamListInteger", "/intListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListInteger, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListInteger route")
	}
	if err := resource.Get("QueryParamListRid", "/ridListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListRid, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListRid route")
	}
	if err := resource.Get("QueryParamListSafeLong", "/safeLongListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListSafeLong, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListSafeLong route")
	}
	if err := resource.Get("QueryParamListString", "/stringListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListString, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListString route")
	}
	if err := resource.Get("QueryParamListUuid", "/uuidListQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamListUuid, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamListUuid route")
	}
	if err := resource.Get("QueryParamExternalString", "/externalStringQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamExternalString, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamExternalString route")
	}
	if err := resource.Get("QueryParamExternalInteger", "/externalIntegerQueryVar", httpserver.NewJSONHandler(handler.HandleQueryParamExternalInteger, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add queryParamExternalInteger route")
	}
	if err := resource.Post("PathParamExternalString", "/externalStringPath/{myPathParam1}", httpserver.NewJSONHandler(handler.HandlePathParamExternalString, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add pathParamExternalString route")
	}
	if err := resource.Post("PathParamExternalInteger", "/externalIntegerPath/{myPathParam1}", httpserver.NewJSONHandler(handler.HandlePathParamExternalInteger, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add pathParamExternalInteger route")
	}
	if err := resource.Post("PostPathParam", "/path/{myPathParam1}/{myPathParam2}", httpserver.NewJSONHandler(handler.HandlePostPathParam, httpserver.StatusCodeMapper, httpserver.ErrHandler), append(routerParams, wrouter.SafeQueryParams("myQueryParam6"))...); err != nil {
		return werror.Wrap(err, "failed to add postPathParam route")
	}
	if err := resource.Post("PostSafeParams", "/safe/{myPathParam1}/{myPathParam2}", httpserver.NewJSONHandler(handler.HandlePostSafeParams, httpserver.StatusCodeMapper, httpserver.ErrHandler), append(routerParams, wrouter.SafePathParams("myPathParam1"), wrouter.SafeHeaderParams("X-My-Header1-Abc"), wrouter.SafeHeaderParams("X-My-Header2"), wrouter.SafeQueryParams("query1"), wrouter.SafeQueryParams("myQueryParam2"), wrouter.ForbiddenQueryParams("myQueryParam4"))...); err != nil {
		return werror.Wrap(err, "failed to add postSafeParams route")
	}
	if err := resource.Get("Bytes", "/bytes", httpserver.NewJSONHandler(handler.HandleBytes, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add bytes route")
	}
	if err := resource.Get("GetBinary", "/binary", httpserver.NewJSONHandler(handler.HandleGetBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add getBinary route")
	}
	if err := resource.Post("PostBinary", "/binary", httpserver.NewJSONHandler(handler.HandlePostBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add postBinary route")
	}
	if err := resource.Put("PutBinary", "/binary", httpserver.NewJSONHandler(handler.HandlePutBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add putBinary route")
	}
	if err := resource.Get("GetOptionalBinary", "/optional/binary", httpserver.NewJSONHandler(handler.HandleGetOptionalBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
		return werror.Wrap(err, "failed to add getOptionalBinary route")
	}
	if err := resource.Post("Chan", "/chan/{var}", httpserver.NewJSONHandler(handler.HandleChan, httpserver.StatusCodeMapper, httpserver.ErrHandler), routerParams...); err != nil {
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
	if err := t.impl.Echo(req.Context(), cookieToken); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleEchoStrings(rw http.ResponseWriter, req *http.Request) error {
	var bodyArg []string
	if err := codecs.JSON.Decode(req.Body, &bodyArg); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	respArg, err := t.impl.EchoStrings(req.Context(), bodyArg)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleEchoCustomObject(rw http.ResponseWriter, req *http.Request) error {
	var bodyArg *CustomObject
	if req.Body != nil && req.Body != http.NoBody {
		if err := codecs.JSON.Decode(req.Body, &bodyArg); err != nil {
			return errors.WrapWithInvalidArgument(err)
		}
	}
	respArg, err := t.impl.EchoCustomObject(req.Context(), bodyArg)
	if err != nil {
		return err
	}
	if respArg == nil {
		rw.WriteHeader(http.StatusNoContent)
		return nil
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, *respArg)
}

func (t *testServiceHandler) HandleEchoOptionalAlias(rw http.ResponseWriter, req *http.Request) error {
	var bodyArg OptionalIntegerAlias
	if req.Body != nil && req.Body != http.NoBody {
		if err := codecs.JSON.Decode(req.Body, &bodyArg); err != nil {
			return errors.WrapWithInvalidArgument(err)
		}
	}
	respArg, err := t.impl.EchoOptionalAlias(req.Context(), bodyArg)
	if err != nil {
		return err
	}
	if respArg.Value == nil {
		rw.WriteHeader(http.StatusNoContent)
		return nil
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleEchoOptionalListAlias(rw http.ResponseWriter, req *http.Request) error {
	var bodyArg OptionalListAlias
	if req.Body != nil && req.Body != http.NoBody {
		if err := codecs.JSON.Decode(req.Body, &bodyArg); err != nil {
			return errors.WrapWithInvalidArgument(err)
		}
	}
	respArg, err := t.impl.EchoOptionalListAlias(req.Context(), bodyArg)
	if err != nil {
		return err
	}
	if respArg.Value == nil {
		rw.WriteHeader(http.StatusNoContent)
		return nil
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
	myPathParamArg, ok := pathParams["myPathParam"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam\" not present")
	}
	if err := t.impl.GetPathParam(req.Context(), bearertoken.Token(authHeader), myPathParamArg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
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
	myPathParamArgStr, ok := pathParams["myPathParam"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam\" not present")
	}
	myPathParamArg := StringAlias(myPathParamArgStr)
	if err := t.impl.GetPathParamAlias(req.Context(), bearertoken.Token(authHeader), myPathParamArg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamList(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	myQueryParam1Arg := req.URL.Query()["myQueryParam1"]
	if err := t.impl.QueryParamList(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListBoolean(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []bool
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.ParseBool(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as boolean")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListBoolean(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListDateTime(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []datetime.DateTime
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := datetime.ParseDateTime(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as datetime")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListDateTime(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamSetDateTime(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []datetime.DateTime
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := datetime.ParseDateTime(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as datetime")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	respArg, err := t.impl.QueryParamSetDateTime(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleQueryParamListDouble(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []float64
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as double")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListDouble(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListInteger(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []int
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := strconv.Atoi(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as integer")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListInteger(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListRid(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []rid.ResourceIdentifier
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := rid.ParseRID(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as rid")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListRid(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListSafeLong(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []safelong.SafeLong
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := safelong.ParseSafeLong(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as safelong")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListSafeLong(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListString(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	myQueryParam1Arg := req.URL.Query()["myQueryParam1"]
	if err := t.impl.QueryParamListString(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamListUuid(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var myQueryParam1Arg []uuid.UUID
	for _, v := range req.URL.Query()["myQueryParam1"] {
		convertedVal, err := uuid.ParseUUID(v)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as uuid")
		}
		myQueryParam1Arg = append(myQueryParam1Arg, convertedVal)
	}
	if err := t.impl.QueryParamListUuid(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamExternalString(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	myQueryParam1Arg := req.URL.Query().Get("myQueryParam1")
	if err := t.impl.QueryParamExternalString(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandleQueryParamExternalInteger(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	myQueryParam1Arg, err := strconv.Atoi(req.URL.Query().Get("myQueryParam1"))
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam1\" as integer")
	}
	if err := t.impl.QueryParamExternalInteger(req.Context(), bearertoken.Token(authHeader), myQueryParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandlePathParamExternalString(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam1ArgStr, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam1\" not present")
	}
	myPathParam1Arg := myPathParam1ArgStr
	if err := t.impl.PathParamExternalString(req.Context(), bearertoken.Token(authHeader), myPathParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (t *testServiceHandler) HandlePathParamExternalInteger(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.Wrap(errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	myPathParam1ArgStr, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam1\" not present")
	}
	myPathParam1Arg, err := strconv.Atoi(myPathParam1ArgStr)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myPathParam1\" as integer")
	}
	if err := t.impl.PathParamExternalInteger(req.Context(), bearertoken.Token(authHeader), myPathParam1Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
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
	myPathParam1Arg, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam1\" not present")
	}
	myPathParam2ArgStr, ok := pathParams["myPathParam2"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam2\" not present")
	}
	myPathParam2Arg, err := strconv.ParseBool(myPathParam2ArgStr)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myPathParam2\" as boolean")
	}
	myQueryParam1Arg := req.URL.Query().Get("query1")
	myQueryParam2Arg := req.URL.Query().Get("myQueryParam2")
	myQueryParam3Arg, err := strconv.ParseFloat(req.URL.Query().Get("myQueryParam3"), 64)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam3\" as double")
	}
	var myQueryParam4Arg *safelong.SafeLong
	if myQueryParam4ArgStr := req.URL.Query().Get("myQueryParam4"); myQueryParam4ArgStr != "" {
		myQueryParam4ArgInternal, err := safelong.ParseSafeLong(myQueryParam4ArgStr)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam4\" as safelong")
		}
		myQueryParam4Arg = &myQueryParam4ArgInternal
	}
	var myQueryParam5Arg *string
	if myQueryParam5ArgStr := req.URL.Query().Get("myQueryParam5"); myQueryParam5ArgStr != "" {
		myQueryParam5ArgInternal := myQueryParam5ArgStr
		myQueryParam5Arg = &myQueryParam5ArgInternal
	}
	var myQueryParam6ArgValue *int
	if myQueryParam6ArgValueStr1 := req.URL.Query().Get("myQueryParam6"); myQueryParam6ArgValueStr1 != "" {
		myQueryParam6ArgValueInternal1, err := strconv.Atoi(myQueryParam6ArgValueStr1)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam6\" as integer")
		}
		myQueryParam6ArgValue = &myQueryParam6ArgValueInternal1
	}
	myQueryParam6Arg := OptionalIntegerAlias{Value: myQueryParam6ArgValue}
	myHeaderParam1Arg, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header1-Abc"))
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam1\" as safelong")
	}
	var myHeaderParam2Arg *uuid.UUID
	if myHeaderParam2ArgStr := req.Header.Get("X-My-Header2"); myHeaderParam2ArgStr != "" {
		myHeaderParam2ArgInternal, err := uuid.ParseUUID(myHeaderParam2ArgStr)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam2\" as uuid")
		}
		myHeaderParam2Arg = &myHeaderParam2ArgInternal
	}
	var myBodyParamArg CustomObject
	if err := codecs.JSON.Decode(req.Body, &myBodyParamArg); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	respArg, err := t.impl.PostPathParam(req.Context(), bearertoken.Token(authHeader), myPathParam1Arg, myPathParam2Arg, myBodyParamArg, myQueryParam1Arg, myQueryParam2Arg, myQueryParam3Arg, myQueryParam4Arg, myQueryParam5Arg, myQueryParam6Arg, myHeaderParam1Arg, myHeaderParam2Arg)
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
	myPathParam1Arg, ok := pathParams["myPathParam1"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam1\" not present")
	}
	myPathParam2ArgStr, ok := pathParams["myPathParam2"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"myPathParam2\" not present")
	}
	myPathParam2Arg, err := strconv.ParseBool(myPathParam2ArgStr)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myPathParam2\" as boolean")
	}
	myQueryParam1Arg := req.URL.Query().Get("query1")
	myQueryParam2Arg := req.URL.Query().Get("myQueryParam2")
	myQueryParam3Arg, err := strconv.ParseFloat(req.URL.Query().Get("myQueryParam3"), 64)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam3\" as double")
	}
	var myQueryParam4Arg *safelong.SafeLong
	if myQueryParam4ArgStr := req.URL.Query().Get("myQueryParam4"); myQueryParam4ArgStr != "" {
		myQueryParam4ArgInternal, err := safelong.ParseSafeLong(myQueryParam4ArgStr)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myQueryParam4\" as safelong")
		}
		myQueryParam4Arg = &myQueryParam4ArgInternal
	}
	var myQueryParam5Arg *string
	if myQueryParam5ArgStr := req.URL.Query().Get("myQueryParam5"); myQueryParam5ArgStr != "" {
		myQueryParam5ArgInternal := myQueryParam5ArgStr
		myQueryParam5Arg = &myQueryParam5ArgInternal
	}
	myHeaderParam1Arg, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header1-Abc"))
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam1\" as safelong")
	}
	var myHeaderParam2Arg *SafeUuid
	if myHeaderParam2ArgStr := req.Header.Get("X-My-Header2"); myHeaderParam2ArgStr != "" {
		myHeaderParam2ArgInternalValue1, err := uuid.ParseUUID(myHeaderParam2ArgStr)
		if err != nil {
			return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"myHeaderParam2\" as uuid")
		}
		myHeaderParam2ArgInternal := SafeUuid(myHeaderParam2ArgInternalValue1)
		myHeaderParam2Arg = &myHeaderParam2ArgInternal
	}
	var myBodyParamArg CustomObject
	if err := codecs.JSON.Decode(req.Body, &myBodyParamArg); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	if err := t.impl.PostSafeParams(req.Context(), bearertoken.Token(authHeader), myPathParam1Arg, myPathParam2Arg, myBodyParamArg, myQueryParam1Arg, myQueryParam2Arg, myQueryParam3Arg, myQueryParam4Arg, myQueryParam5Arg, myHeaderParam1Arg, myHeaderParam2Arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
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
	myBytesArg := req.Body
	respArg, err := t.impl.PostBinary(req.Context(), myBytesArg)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, respArg)
}

func (t *testServiceHandler) HandlePutBinary(rw http.ResponseWriter, req *http.Request) error {
	myBytesArg := req.Body
	if err := t.impl.PutBinary(req.Context(), myBytesArg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
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
	varArg, ok := pathParams["var"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"var\" not present")
	}
	typeArg := req.URL.Query().Get("type")
	httpArg := req.URL.Query().Get("http")
	jsonArg := req.URL.Query().Get("json")
	reqArg := req.URL.Query().Get("req")
	rwArg := req.URL.Query().Get("rw")
	returnArg, err := safelong.ParseSafeLong(req.Header.Get("X-My-Header2"))
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "failed to parse \"return\" as safelong")
	}
	var importArg map[string]string
	if err := codecs.JSON.Decode(req.Body, &importArg); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	if err := t.impl.Chan(req.Context(), varArg, importArg, typeArg, returnArg, httpArg, jsonArg, reqArg, rwArg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}
