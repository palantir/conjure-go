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
	collector := ParamTypeFilterer{}
	for i := range argumentDefinitions {
		argDef := argumentDefinitions[i]
		paramType := argDef.ParamType
		if err := argDef.ParamType.AcceptFuncs(
			func(bodyParameterType spec.BodyParameterType) error {
				collector.BodyParameterTypes = append(collector.BodyParameterTypes, ArgumentDefinitionBodyParam{
					ArgumentDefinition: argDef,
					BodyParameterType:  bodyParameterType,
				})
				return nil
			},
			func(headerParameterType spec.HeaderParameterType) error {
				collector.HeaderParameterTypes = append(collector.HeaderParameterTypes, ArgumentDefinitionHeaderParam{
					ArgumentDefinition:  argDef,
					HeaderParameterType: headerParameterType,
				})
				return nil
			},
			func(pathParameterType spec.PathParameterType) error {
				collector.PathParameterTypes = append(collector.PathParameterTypes, ArgumentDefinitionPathParam{
					ArgumentDefinition: argDef,
					PathParameterType:  pathParameterType,
				})
				return nil
			},
			func(queryParameterType spec.QueryParameterType) error {
				collector.QueryParameterTypes = append(collector.QueryParameterTypes, ArgumentDefinitionQueryParam{
					ArgumentDefinition: argDef,
					QueryParameterType: queryParameterType,
				})
				return nil
			},
			paramType.ErrorOnUnknown,
		); err != nil {
			return ParamTypeFilterer{}, err
		}
	}
	return collector, nil
}
