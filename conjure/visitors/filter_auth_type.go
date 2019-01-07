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

package visitors

import (
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
)

type AuthTypeFilterer struct {
	headerAuthType *spec.HeaderAuthType
	cookieAuthType *spec.CookieAuthType
}

var _ spec.AuthTypeVisitor = &AuthTypeFilterer{}

func GetPossibleHeaderAuth(authType spec.AuthType) (*spec.HeaderAuthType, error) {
	authTypeFilterer := AuthTypeFilterer{}
	err := authType.Accept(&authTypeFilterer)
	return authTypeFilterer.headerAuthType, err
}

func GetPossibleCookieAuth(authType spec.AuthType) (*spec.CookieAuthType, error) {
	authTypeFilterer := AuthTypeFilterer{}
	err := authType.Accept(&authTypeFilterer)
	return authTypeFilterer.cookieAuthType, err
}

func (a *AuthTypeFilterer) VisitHeader(v spec.HeaderAuthType) error {
	a.headerAuthType = &v
	return nil
}

func (a *AuthTypeFilterer) VisitCookie(v spec.CookieAuthType) error {
	a.cookieAuthType = &v
	return nil
}

func (a *AuthTypeFilterer) VisitUnknown(typeName string) error {
	return errors.New("Unknown auth type " + typeName)
}
