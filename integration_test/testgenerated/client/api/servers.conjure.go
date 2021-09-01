// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"io"
	"net/http"

	codecs "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	errors "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	httpserver "github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	rid "github.com/palantir/pkg/rid"
	werror "github.com/palantir/witchcraft-go-error"
	wresource "github.com/palantir/witchcraft-go-server/v2/witchcraft/wresource"
	wrouter "github.com/palantir/witchcraft-go-server/v2/wrouter"
)

type TestService interface {
	Echo(ctx context.Context) error
	PathParam(ctx context.Context, paramArg string) error
	PathParamAlias(ctx context.Context, paramArg StringAlias) error
	PathParamRid(ctx context.Context, paramArg rid.ResourceIdentifier) error
	PathParamRidAlias(ctx context.Context, paramArg RidAlias) error
	Bytes(ctx context.Context) (CustomObject, error)
	Binary(ctx context.Context) (io.ReadCloser, error)
	MaybeBinary(ctx context.Context) (*io.ReadCloser, error)
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Get(
		"Echo",
		"/echo",
		httpserver.NewJSONHandler(handler.HandleEcho, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add echo route")
	}
	if err := resource.Get(
		"PathParam",
		"/path/{param}",
		httpserver.NewJSONHandler(handler.HandlePathParam, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add pathParam route")
	}
	if err := resource.Get(
		"PathParamAlias",
		"/path/alias/{param}",
		httpserver.NewJSONHandler(handler.HandlePathParamAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add pathParamAlias route")
	}
	if err := resource.Get(
		"PathParamRid",
		"/path/rid/{param}",
		httpserver.NewJSONHandler(handler.HandlePathParamRid, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add pathParamRid route")
	}
	if err := resource.Get(
		"PathParamRidAlias",
		"/path/rid/alias/{param}",
		httpserver.NewJSONHandler(handler.HandlePathParamRidAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add pathParamRidAlias route")
	}
	if err := resource.Get(
		"Bytes",
		"/bytes",
		httpserver.NewJSONHandler(handler.HandleBytes, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add bytes route")
	}
	if err := resource.Get(
		"Binary",
		"/binary",
		httpserver.NewJSONHandler(handler.HandleBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binary route")
	}
	if err := resource.Get(
		"MaybeBinary",
		"/optional/binary",
		httpserver.NewJSONHandler(handler.HandleMaybeBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add maybeBinary route")
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleEcho(_ http.ResponseWriter, req *http.Request) error {
	return t.impl.Echo(req.Context())
}

func (t *testServiceHandler) HandlePathParam(_ http.ResponseWriter, req *http.Request) error {
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.WrapWithContextParams(req.Context(), errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	param, ok := pathParams["param"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"param\" not present")
	}
	return t.impl.PathParam(req.Context(), param)
}

func (t *testServiceHandler) HandlePathParamAlias(_ http.ResponseWriter, req *http.Request) error {
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.WrapWithContextParams(req.Context(), errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	paramStr, ok := pathParams["param"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"param\" not present")
	}
	var param StringAlias
	if err := param.UnmarshalString(paramStr); err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "unmarshal path[\"param\"] as StringAlias")
	}
	return t.impl.PathParamAlias(req.Context(), param)
}

func (t *testServiceHandler) HandlePathParamRid(_ http.ResponseWriter, req *http.Request) error {
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.WrapWithContextParams(req.Context(), errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	paramStr, ok := pathParams["param"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"param\" not present")
	}
	param, err := rid.ParseRID(paramStr)
	if err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "unmarshal path[\"param\"] as rid")
	}
	return t.impl.PathParamRid(req.Context(), param)
}

func (t *testServiceHandler) HandlePathParamRidAlias(_ http.ResponseWriter, req *http.Request) error {
	pathParams := wrouter.PathParams(req)
	if pathParams == nil {
		return werror.WrapWithContextParams(req.Context(), errors.NewInternal(), "path params not found on request: ensure this endpoint is registered with wrouter")
	}
	paramStr, ok := pathParams["param"]
	if !ok {
		return werror.WrapWithContextParams(req.Context(), errors.NewInvalidArgument(), "path parameter \"param\" not present")
	}
	var param RidAlias
	if err := param.UnmarshalString(paramStr); err != nil {
		return werror.WrapWithContextParams(req.Context(), errors.WrapWithInvalidArgument(err), "unmarshal path[\"param\"] as RidAlias")
	}
	return t.impl.PathParamRidAlias(req.Context(), param)
}

func (t *testServiceHandler) HandleBytes(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.Bytes(req.Context())
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleBinary(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.Binary(req.Context())
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleMaybeBinary(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.MaybeBinary(req.Context())
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
