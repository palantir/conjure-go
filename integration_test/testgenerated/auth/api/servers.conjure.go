// This file was generated by Conjure and should not be manually edited.

package api

import (
	"net/http"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/codecs"
	"github.com/palantir/conjure-go-runtime/conjure-go-server/rest"
	"github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/witchcraft/wresource"
	"github.com/palantir/witchcraft-go-server/wrouter"
)

// RegisterRoutesBothAuthService registers handlers for the BothAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesBothAuthService(router wrouter.Router, impl BothAuthService) error {
	handler := bothAuthServiceHandler{impl: impl}
	resource := wresource.New("bothauthservice", router)
	if err := resource.Get("Default", "/default", rest.HandlerFunc(handler.HandleDefault)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Default"))
	}
	if err := resource.Get("Cookie", "/cookie", rest.HandlerFunc(handler.HandleCookie)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Cookie"))
	}
	if err := resource.Get("None", "/none", rest.HandlerFunc(handler.HandleNone)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "None"))
	}
	if err := resource.Post("WithArg", "/withArg", rest.HandlerFunc(handler.HandleWithArg)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "WithArg"))
	}
	return nil
}

type bothAuthServiceHandler struct {
	impl BothAuthService
}

func (b *bothAuthServiceHandler) HandleDefault(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := rest.ParseBearerTokenHeader(req)
	if err != nil {
		return err
	}
	respArg, err := b.impl.Default(req.Context(), authHeader)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

func (b *bothAuthServiceHandler) HandleCookie(rw http.ResponseWriter, req *http.Request) error {
	cookieToken, err := rest.ParseBearerTokenCookie(req, "P_TOKEN")
	if err != nil {
		return err
	}
	return b.impl.Cookie(req.Context(), cookieToken)
}

func (b *bothAuthServiceHandler) HandleNone(rw http.ResponseWriter, req *http.Request) error {
	return b.impl.None(req.Context())
}

func (b *bothAuthServiceHandler) HandleWithArg(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := rest.ParseBearerTokenHeader(req)
	if err != nil {
		return err
	}
	var arg string
	if err := codecs.JSON.Decode(req.Body, &arg); err != nil {
		return werror.Wrap(err, "failed to unmarshal request body", werror.SafeParam("bodyParamName", "arg"), werror.SafeParam("bodyParamType", "string"))
	}
	return b.impl.WithArg(req.Context(), authHeader, arg)
}

// RegisterRoutesHeaderAuthService registers handlers for the HeaderAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesHeaderAuthService(router wrouter.Router, impl HeaderAuthService) error {
	handler := headerAuthServiceHandler{impl: impl}
	resource := wresource.New("headerauthservice", router)
	if err := resource.Get("Default", "/default", rest.HandlerFunc(handler.HandleDefault)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Default"))
	}
	return nil
}

type headerAuthServiceHandler struct {
	impl HeaderAuthService
}

func (h *headerAuthServiceHandler) HandleDefault(rw http.ResponseWriter, req *http.Request) error {
	authHeader, err := rest.ParseBearerTokenHeader(req)
	if err != nil {
		return err
	}
	respArg, err := h.impl.Default(req.Context(), authHeader)
	if err != nil {
		return err
	}
	rw.Header().Add("Content-Type", codecs.JSON.ContentType())
	return codecs.JSON.Encode(rw, respArg)
}

// RegisterRoutesCookieAuthService registers handlers for the CookieAuthService endpoints with a witchcraft wrouter.
// This should typically be called in a witchcraft server's InitFunc.
// impl provides an implementation of each endpoint, which can assume the request parameters have been parsed
// in accordance with the Conjure specification.
func RegisterRoutesCookieAuthService(router wrouter.Router, impl CookieAuthService) error {
	handler := cookieAuthServiceHandler{impl: impl}
	resource := wresource.New("cookieauthservice", router)
	if err := resource.Get("Cookie", "/cookie", rest.HandlerFunc(handler.HandleCookie)); err != nil {
		return werror.Wrap(err, "failed to add route", werror.SafeParam("routeName", "Cookie"))
	}
	return nil
}

type cookieAuthServiceHandler struct {
	impl CookieAuthService
}

func (c *cookieAuthServiceHandler) HandleCookie(rw http.ResponseWriter, req *http.Request) error {
	cookieToken, err := rest.ParseBearerTokenCookie(req, "P_TOKEN")
	if err != nil {
		return err
	}
	return c.impl.Cookie(req.Context(), cookieToken)
}