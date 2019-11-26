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

	"github.com/palantir/conjure-go/v4/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v4/conjure/types"
)

type setVisitor struct {
	setType spec.SetType
}

func newSetVisitor(setType spec.SetType) ConjureTypeProvider {
	return &setVisitor{setType: setType}
}

func (p *setVisitor) ParseType(info types.PkgInfo) (types.Typer, error) {
	nestedTypeProvider, err := NewConjureTypeProvider(p.setType.ItemType)
	if err != nil {
		return nil, err
	}
	typer, err := nestedTypeProvider.ParseType(info)
	if err != nil {
		return nil, err
	}
	return types.NewSetType(typer), nil
}

func (p *setVisitor) CollectionInitializationIfNeeded(info types.PkgInfo) (*expression.CallExpression, error) {
	typer, err := p.ParseType(info)
	if err != nil {
		return nil, err
	}
	return expression.NewCallExpression(expression.MakeBuiltIn, expression.Type(typer.GoType(info)), expression.IntVal(0)), nil
}

func (p *setVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return typeCheck == IsSet
}
