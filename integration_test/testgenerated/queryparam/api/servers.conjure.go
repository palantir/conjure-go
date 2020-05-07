// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/wrouter"
)

type TestService interface {
	Echo(ctx context.Context, inputArg string, repsArg int, optionalArg *string, lastParamArg *string) (string, error)
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Get("Echo", "/echo", httpserver.NewJSONHandler(handler.HandleEcho, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Echo"))
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleEcho(rw http.ResponseWriter, req *http.Request) error {
	input := req.URL.Query().Get("input")
	reps, err := strconv.Atoi(req.URL.Query().Get("reps"))
	if err != nil {
		return err
	}
	var optional *string
	if optionalStr := req.URL.Query().Get("optional"); optionalStr != "" {
		optionalInternal := optionalStr
		optional = &optionalInternal
	}
	var lastParam *string
	if lastParamStr := req.URL.Query().Get("lastParam"); lastParamStr != "" {
		lastParamInternal := lastParamStr
		lastParam = &lastParamInternal
	}
	respArg, err := t.impl.Echo(req.Context(), input, reps, optional, lastParam)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}
