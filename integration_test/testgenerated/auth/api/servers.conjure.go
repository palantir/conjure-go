// This file was generated by Conjure and should not be manually edited.

package api

import (
	"context"
	"net/http"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-server/httpserver"
	"github.com/palantir/pkg/bearertoken"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/v2/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/v2/wrouter"
)

type BothAuthService interface {
	Default(ctx context.Context, authHeader bearertoken.Token) (string, error)
	Cookie(ctx context.Context, cookieToken bearertoken.Token) error
	None(ctx context.Context) error
	WithArg(ctx context.Context, authHeader bearertoken.Token, argArg string) error
}

// RegisterRoutesBothAuthService registers handlers for the BothAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesBothAuthService(router wrouter.Router, impl BothAuthService) error {
	handler := bothAuthServiceHandler{impl: impl}
	resource := wresource.New("bothauthservice", router)
	if err := resource.Get("Default", "/default", httpserver.NewJSONHandler(handler.HandleDefault, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add default route")
	}
	if err := resource.Get("Cookie", "/cookie", httpserver.NewJSONHandler(handler.HandleCookie, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add cookie route")
	}
	if err := resource.Get("None", "/none", httpserver.NewJSONHandler(handler.HandleNone, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add none route")
	}
	if err := resource.Post("WithArg", "/withArg", httpserver.NewJSONHandler(handler.HandleWithArg, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add withArg route")
	}
	return nil
}

type bothAuthServiceHandler struct {
	impl BothAuthService
}

func (b *bothAuthServiceHandler) HandleDefault(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	respArg, err := b.impl.Default(req.Context(), bearertoken.Token(authHeader))
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (b *bothAuthServiceHandler) HandleCookie(rw http.ResponseWriter, req *http.Request) error {
	authCookie, err := req.Cookie("P_TOKEN")
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	cookieToken := bearertoken.Token(authCookie.Value)
	if err := b.impl.Cookie(req.Context(), cookieToken); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (b *bothAuthServiceHandler) HandleNone(rw http.ResponseWriter, req *http.Request) error {
	if err := b.impl.None(req.Context()); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (b *bothAuthServiceHandler) HandleWithArg(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	var arg string
	if err := codecs.JSON.Decode(req.Body, &arg); err != nil {
		return errors.WrapWithInvalidArgument(err)
	}
	if err := b.impl.WithArg(req.Context(), bearertoken.Token(authHeader), arg); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

type CookieAuthService interface {
	Cookie(ctx context.Context, cookieToken bearertoken.Token) error
}

// RegisterRoutesCookieAuthService registers handlers for the CookieAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesCookieAuthService(router wrouter.Router, impl CookieAuthService) error {
	handler := cookieAuthServiceHandler{impl: impl}
	resource := wresource.New("cookieauthservice", router)
	if err := resource.Get("Cookie", "/cookie", httpserver.NewJSONHandler(handler.HandleCookie, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add cookie route")
	}
	return nil
}

type cookieAuthServiceHandler struct {
	impl CookieAuthService
}

func (c *cookieAuthServiceHandler) HandleCookie(rw http.ResponseWriter, req *http.Request) error {
	authCookie, err := req.Cookie("P_TOKEN")
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	cookieToken := bearertoken.Token(authCookie.Value)
	if err := c.impl.Cookie(req.Context(), cookieToken); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

type HeaderAuthService interface {
	Default(ctx context.Context, authHeader bearertoken.Token) (string, error)
}

// RegisterRoutesHeaderAuthService registers handlers for the HeaderAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesHeaderAuthService(router wrouter.Router, impl HeaderAuthService) error {
	handler := headerAuthServiceHandler{impl: impl}
	resource := wresource.New("headerauthservice", router)
	if err := resource.Get("Default", "/default", httpserver.NewJSONHandler(handler.HandleDefault, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add default route")
	}
	return nil
}

type headerAuthServiceHandler struct {
	impl HeaderAuthService
}

func (h *headerAuthServiceHandler) HandleDefault(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	respArg, err := h.impl.Default(req.Context(), bearertoken.Token(authHeader))
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

type SomeHeaderAuthService interface {
	Default(ctx context.Context, authHeader bearertoken.Token) (string, error)
	None(ctx context.Context) error
}

// RegisterRoutesSomeHeaderAuthService registers handlers for the SomeHeaderAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesSomeHeaderAuthService(router wrouter.Router, impl SomeHeaderAuthService) error {
	handler := someHeaderAuthServiceHandler{impl: impl}
	resource := wresource.New("someheaderauthservice", router)
	if err := resource.Get("Default", "/default", httpserver.NewJSONHandler(handler.HandleDefault, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add default route")
	}
	if err := resource.Get("None", "/none", httpserver.NewJSONHandler(handler.HandleNone, httpserver.StatusCodeMapper, httpserver.ErrHandler)); err != nil {
		return werror.Wrap(err, "failed to add none route")
	}
	return nil
}

type someHeaderAuthServiceHandler struct {
	impl SomeHeaderAuthService
}

func (s *someHeaderAuthServiceHandler) HandleDefault(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := httpserver.ParseBearerTokenHeader(req)
	if err != nil {
		return errors.WrapWithPermissionDenied(err)
	}
	respArg, err := s.impl.Default(req.Context(), bearertoken.Token(authHeader))
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (s *someHeaderAuthServiceHandler) HandleNone(rw http.ResponseWriter, req *http.Request) error {
	if err := s.impl.None(req.Context()); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusNoContent)
	return nil
}
