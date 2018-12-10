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
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/expression"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
)

type MapVisitor struct {
	mapType spec.MapType
}

func NewMapVisitor(mapType spec.MapType) ConjureTypeProvider {
	return &MapVisitor{mapType: mapType}
}

var _ ConjureTypeProvider = &MapVisitor{}

func (p *MapVisitor) ParseType(customTypes types.CustomConjureTypes) (types.Typer, error) {
	keyTypeProvider, err := getTyper(p.mapType.KeyType, customTypes)
	if err != nil {
		return nil, err
	}
	valueTypeProvider, err := getTyper(p.mapType.ValueType, customTypes)
	if err != nil {
		return nil, err
	}
	return types.NewMapType(keyTypeProvider, valueTypeProvider), nil
}

func getTyper(typeFromSpec spec.Type, customTypes types.CustomConjureTypes) (types.Typer, error) {
	typeProvider, err := NewConjureTypeProvider(typeFromSpec)
	if err != nil {
		return nil, err
	}
	typer, err := typeProvider.ParseType(customTypes)
	if err != nil {
		return nil, err
	}
	return typer, nil
}

func (p *MapVisitor) CollectionInitializationIfNeeded(customTypes types.CustomConjureTypes, currPkgPath string, pkgAliases map[string]string) (*expression.CallExpression, error) {
	typer, err := p.ParseType(customTypes)
	if err != nil {
		return nil, err
	}
	goType := typer.GoType(currPkgPath, pkgAliases)
	return &expression.CallExpression{
		Function: expression.VariableVal("make"),
		Args: []astgen.ASTExpr{
			expression.Type(goType),
			expression.IntVal(0),
		},
	}, nil
}

func (p *MapVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return typeCheck == IsMap
}
