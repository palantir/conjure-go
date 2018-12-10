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
	"github.com/palantir/goastwriter/expression"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
)

type OptionalVisitor struct {
	optionalType spec.OptionalType
}

func NewOptionalVisitor(optionalType spec.OptionalType) ConjureTypeProvider {
	return &OptionalVisitor{optionalType: optionalType}
}

var _ ConjureTypeProvider = &OptionalVisitor{}

func (p *OptionalVisitor) ParseType(customTypes types.CustomConjureTypes) (types.Typer, error) {
	nestedTypeProvider, err := NewConjureTypeProvider(p.optionalType.ItemType)
	if err != nil {
		return nil, err
	}
	typer, err := nestedTypeProvider.ParseType(customTypes)
	if err != nil {
		return nil, err
	}
	return types.NewOptionalType(typer), nil
}

func (p *OptionalVisitor) CollectionInitializationIfNeeded(customTypes types.CustomConjureTypes, currPkgPath string, pkgAliases map[string]string) (*expression.CallExpression, error) {
	return nil, nil
}

func (p *OptionalVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return typeCheck == IsOptional
}
