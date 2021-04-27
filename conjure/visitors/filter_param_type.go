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
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/pkg/errors"
)

type ArgumentDefinitionHeaderParam struct {
	ArgumentDefinition  spec.ArgumentDefinition
	HeaderParameterType spec.HeaderParameterType
}

type ArgumentDefinitionBodyParam struct {
	ArgumentDefinition spec.ArgumentDefinition
	BodyParameterType  spec.BodyParameterType
}

type ArgumentDefinitionPathParam struct {
	ArgumentDefinition spec.ArgumentDefinition
	PathParameterType  spec.PathParameterType
}

type ArgumentDefinitionQueryParam struct {
	ArgumentDefinition spec.ArgumentDefinition
	QueryParameterType spec.QueryParameterType
}

type ParamTypeFilterer struct {
	CurrentArgumentDefinition spec.ArgumentDefinition
	BodyParameterTypes        []ArgumentDefinitionBodyParam
	HeaderParameterTypes      []ArgumentDefinitionHeaderParam
	PathParameterTypes        []ArgumentDefinitionPathParam
	QueryParameterTypes       []ArgumentDefinitionQueryParam
}

var _ spec.ParameterTypeVisitor = &ParamTypeFilterer{}

func GetPathParams(argumentDefinitions []spec.ArgumentDefinition) ([]ArgumentDefinitionPathParam, error) {
	paramTypeFilterer, err := runFilters(argumentDefinitions)
	if err != nil {
		return nil, err
	}
	return paramTypeFilterer.PathParameterTypes, nil
}

func GetQueryParams(argumentDefinitions []spec.ArgumentDefinition) ([]ArgumentDefinitionQueryParam, error) {
	paramTypeFilterer, err := runFilters(argumentDefinitions)
	if err != nil {
		return nil, err
	}
	return paramTypeFilterer.QueryParameterTypes, nil
}

func GetBodyParams(argumentDefinitions []spec.ArgumentDefinition) ([]ArgumentDefinitionBodyParam, error) {
	paramTypeFilterer, err := runFilters(argumentDefinitions)
	if err != nil {
		return nil, err
	}
	return paramTypeFilterer.BodyParameterTypes, nil
}

func GetHeaderParams(argumentDefinitions []spec.ArgumentDefinition) ([]ArgumentDefinitionHeaderParam, error) {
	paramTypeFilterer, err := runFilters(argumentDefinitions)
	if err != nil {
		return nil, err
	}
	return paramTypeFilterer.HeaderParameterTypes, nil
}

func runFilters(argumentDefinitions []spec.ArgumentDefinition) (ParamTypeFilterer, error) {
	paramTypeFilterer := ParamTypeFilterer{}
	for _, argumentDefinition := range argumentDefinitions {
		paramTypeFilterer.CurrentArgumentDefinition = argumentDefinition
		err := argumentDefinition.ParamType.Accept(&paramTypeFilterer)
		if err != nil {
			return ParamTypeFilterer{}, err
		}
	}
	return paramTypeFilterer, nil
}

func (p *ParamTypeFilterer) VisitBody(bodyParameterType spec.BodyParameterType) error {
	p.BodyParameterTypes = append(p.BodyParameterTypes, ArgumentDefinitionBodyParam{
		ArgumentDefinition: p.CurrentArgumentDefinition,
		BodyParameterType:  bodyParameterType,
	})
	return nil
}

func (p *ParamTypeFilterer) VisitHeader(headerParameterType spec.HeaderParameterType) error {
	p.HeaderParameterTypes = append(p.HeaderParameterTypes, ArgumentDefinitionHeaderParam{
		ArgumentDefinition:  p.CurrentArgumentDefinition,
		HeaderParameterType: headerParameterType,
	})
	return nil
}

func (p *ParamTypeFilterer) VisitPath(pathParameterType spec.PathParameterType) error {
	p.PathParameterTypes = append(p.PathParameterTypes, ArgumentDefinitionPathParam{
		ArgumentDefinition: p.CurrentArgumentDefinition,
		PathParameterType:  pathParameterType,
	})
	return nil
}

func (p *ParamTypeFilterer) VisitQuery(queryParameterType spec.QueryParameterType) error {
	p.QueryParameterTypes = append(p.QueryParameterTypes, ArgumentDefinitionQueryParam{
		ArgumentDefinition: p.CurrentArgumentDefinition,
		QueryParameterType: queryParameterType,
	})
	return nil
}

func (p *ParamTypeFilterer) VisitUnknown(typeName string) error {
	return errors.New("Unknown path param type: " + typeName)
}
