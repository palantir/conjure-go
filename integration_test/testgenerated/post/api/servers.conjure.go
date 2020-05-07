// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"net/http"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/rest"
	"github.com/palantir/witchcraft-go-server/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/wrouter"
)

type TestService interface {
	Echo(ctx context.Context, inputArg string) (string, error)
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Post("Echo", "/echo", rest.NewJSONHandler(handler.HandleEcho, rest.StatusCodeMapper, rest.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Echo"))
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleEcho(rw http.ResponseWriter, req *http.Request) error {
	var input string
	if err := codecs.JSON.Decode(req.Body, &input); err != nil {
		return errors.NewWrappedError(err, errors.NewInvalidArgument())
	}
	respArg, err := t.impl.Echo(req.Context(), input)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}
