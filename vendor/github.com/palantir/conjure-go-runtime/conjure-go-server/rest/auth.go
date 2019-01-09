// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"net/http"
	"strings"

	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/witchcraft-go-error"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/errors"
)

// ParseBearerTokenHeader parses a bearer token value out of the Authorization header. It expects a header with a key
// of 'Authorization' and a value of 'Bearer {token}'. ParseBearerTokenHeader will return the token value, or a
// PermissionDenied error if the Authorization header is missing, an empty string, or is not in the format expected.
func ParseBearerTokenHeader(req *http.Request) (bearertoken.Token, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", werror.Wrap(errors.NewPermissionDenied(), "Authorization header not found")
	}
	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) != 2 || strings.ToLower(headerSplit[0]) != "bearer" {
		return "", werror.Wrap(errors.NewPermissionDenied(), "Illegal authorization prefix, expected Bearer")
	}
	return bearertoken.Token(headerSplit[1]), nil
}

// ParseBearerTokenCookie returns the bearer token associated with this request by using the
// value of the cookieName cookie. If the cookie is not present, a PermissionDenied error is returned.
func ParseBearerTokenCookie(req *http.Request, cookieName string) (bearertoken.Token, error) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return "", errors.NewPermissionDenied()
	}
	return bearertoken.Token(cookie.Value), nil
}
