// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package types

import (
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

type EndpointArgumentType uint

const (
	_ EndpointArgumentType = iota
	PathParam
	HeaderParam
	QueryParam
	BodyParam
)

type ServiceDefinition struct {
	Docs
	Name       string
	Endpoints  []*EndpointDefinition
	conjurePkg string
	importPath string
}

func (d ServiceDefinition) HasHeaderAuth() bool {
	for _, endpoint := range d.Endpoints {
		if endpoint.HeaderAuth {
			return true
		}
	}
	return false
}

func (d ServiceDefinition) HasCookieAuth() bool {
	for _, endpoint := range d.Endpoints {
		if endpoint.CookieAuth != nil {
			return true
		}
	}
	return false
}

type EndpointDefinition struct {
	Docs
	Deprecated Docs

	EndpointName string
	HTTPMethod   spec.HttpMethod
	HTTPPath     string
	HeaderAuth   bool    // only one of HeaderAuth or CookieAuth allowed
	CookieAuth   *string // nil if no cookie auth, else value is cookie key
	Params       []*EndpointArgumentDefinition
	Returns      *Type // nil if no return
	Markers      []Type
	Tags         []string
}

type EndpointArgumentDefinition struct {
	Docs
	Name      string
	Type      Type
	ParamType EndpointArgumentType
	ParamID   string
	Markers   []Type
	Safety    *spec.LogSafety
	Tags      []string
}

func (d EndpointDefinition) PathParams() []*EndpointArgumentDefinition {
	return d.filterParams(PathParam)
}

func (d EndpointDefinition) HeaderParams() []*EndpointArgumentDefinition {
	return d.filterParams(HeaderParam)
}

func (d EndpointDefinition) QueryParams() []*EndpointArgumentDefinition {
	return d.filterParams(QueryParam)
}

func (d EndpointDefinition) BodyParam() *EndpointArgumentDefinition {
	if p := d.filterParams(BodyParam); len(p) != 0 {
		return p[0]
	}
	return nil
}

func (d EndpointDefinition) filterParams(argType EndpointArgumentType) []*EndpointArgumentDefinition {
	var p []*EndpointArgumentDefinition
	for _, param := range d.Params {
		if param.ParamType == argType {
			p = append(p, param)
		}
	}
	return p
}
