// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"

	codecs "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	errors "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	httpserver "github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	safejson "github.com/palantir/pkg/safejson"
	werror "github.com/palantir/witchcraft-go-error"
	wresource "github.com/palantir/witchcraft-go-server/v2/witchcraft/wresource"
	wrouter "github.com/palantir/witchcraft-go-server/v2/wrouter"
	gjson "github.com/tidwall/gjson"
)

type TestService interface {
	BinaryAlias(ctx context.Context, bodyArg io.ReadCloser) (io.ReadCloser, error)
	BinaryAliasOptional(ctx context.Context) (*io.ReadCloser, error)
	BinaryAliasAlias(ctx context.Context, bodyArg io.ReadCloser) (*io.ReadCloser, error)
	Binary(ctx context.Context, bodyArg io.ReadCloser) (io.ReadCloser, error)
	BinaryOptional(ctx context.Context) (*io.ReadCloser, error)
	BinaryList(ctx context.Context, bodyArg [][]byte) ([][]byte, error)
	Bytes(ctx context.Context, bodyArg CustomObject) (CustomObject, error)
}

// RegisterRoutesTestService registers handlers for the TestService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesTestService(router wrouter.Router, impl TestService) error {
	handler := testServiceHandler{impl: impl}
	resource := wresource.New("testservice", router)
	if err := resource.Post(
		"BinaryAlias",
		"/binaryAlias",
		httpserver.NewJSONHandler(handler.HandleBinaryAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binaryAlias route")
	}
	if err := resource.Post(
		"BinaryAliasOptional",
		"/binaryAliasOptional",
		httpserver.NewJSONHandler(handler.HandleBinaryAliasOptional, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binaryAliasOptional route")
	}
	if err := resource.Post(
		"BinaryAliasAlias",
		"/binaryAliasAlias",
		httpserver.NewJSONHandler(handler.HandleBinaryAliasAlias, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binaryAliasAlias route")
	}
	if err := resource.Post(
		"Binary",
		"/binary",
		httpserver.NewJSONHandler(handler.HandleBinary, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binary route")
	}
	if err := resource.Post(
		"BinaryOptional",
		"/binaryOptional",
		httpserver.NewJSONHandler(handler.HandleBinaryOptional, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binaryOptional route")
	}
	if err := resource.Post(
		"BinaryList",
		"/binaryList",
		httpserver.NewJSONHandler(handler.HandleBinaryList, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add binaryList route")
	}
	if err := resource.Post(
		"Bytes",
		"/bytes",
		httpserver.NewJSONHandler(handler.HandleBytes, httpserver.StatusCodeMapper, httpserver.ErrHandler),
	); err != nil {
		return werror.WrapWithContextParams(context.TODO(), err, "failed to add bytes route")
	}
	return nil
}

type testServiceHandler struct {
	impl TestService
}

func (t *testServiceHandler) HandleBinaryAlias(rw http.ResponseWriter, req *http.Request) error {
	body := req.Body
	respArg, err := t.impl.BinaryAlias(req.Context(), body)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleBinaryAliasOptional(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.BinaryAliasOptional(req.Context())
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

func (t *testServiceHandler) HandleBinaryAliasAlias(rw http.ResponseWriter, req *http.Request) error {
	body := req.Body
	respArg, err := t.impl.BinaryAliasAlias(req.Context(), body)
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

func (t *testServiceHandler) HandleBinary(rw http.ResponseWriter, req *http.Request) error {
	body := req.Body
	respArg, err := t.impl.Binary(req.Context(), body)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.Binary.ContentType())
	return codecs.Binary.Encode(rw, respArg)
}

func (t *testServiceHandler) HandleBinaryOptional(rw http.ResponseWriter, req *http.Request) error {
	respArg, err := t.impl.BinaryOptional(req.Context())
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

func (t *testServiceHandler) HandleBinaryList(rw http.ResponseWriter, req *http.Request) error {
	var body [][]byte
	if err := codecs.JSON.Decode(req.Body, &body); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	respArg, err := t.impl.BinaryList(req.Context(), body)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, safejson.AppendFunc(func(out []byte) ([]byte, error) {
		out = append(out, '[')
		for i := range respArg {
			out = append(out, '"')
			if len(respArg[i]) > 0 {
				b64out := make([]byte, base64.StdEncoding.EncodedLen(len(respArg[i])))
				base64.StdEncoding.Encode(b64out, respArg[i])
				out = append(out, b64out...)
			}
			out = append(out, '"')
			if i < len(respArg)-1 {
				out = append(out, ',')
			}
		}
		out = append(out, ']')
		if !gjson.ValidBytes(out) {
			return nil, werror.ErrorWithContextParams(context.TODO(), "generated invalid json: please report this as a bug on github.com/palantir/conjure-go/issues")
		}
		return out, nil
	}))
}

func (t *testServiceHandler) HandleBytes(rw http.ResponseWriter, req *http.Request) error {
	var body CustomObject
	if err := codecs.JSON.Decode(req.Body, &body); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	respArg, err := t.impl.Bytes(req.Context(), body)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}
