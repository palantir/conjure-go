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

type SetVisitor struct {
	setType spec.SetType
}

func NewSetVisitor(setType spec.SetType) ConjureTypeProvider {
	return &SetVisitor{setType: setType}
}

var _ ConjureTypeProvider = &SetVisitor{}

func (p *SetVisitor) ParseType(customTypes types.CustomConjureTypes) (types.Typer, error) {
	nestedTypeProvider, err := NewConjureTypeProvider(p.setType.ItemType)
	if err != nil {
		return nil, err
	}
	typer, err := nestedTypeProvider.ParseType(customTypes)
	if err != nil {
		return nil, err
	}
	return types.NewSetType(typer), nil
}

func (p *SetVisitor) CollectionInitializationIfNeeded(customTypes types.CustomConjureTypes, currPkgPath string, pkgAliases map[string]string) (*expression.CallExpression, error) {
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

func (p *SetVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return typeCheck == IsSet
}
