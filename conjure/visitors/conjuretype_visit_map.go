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

type mapVisitor struct {
	mapType spec.MapType
}

func newMapVisitor(mapType spec.MapType) ConjureTypeProvider {
	return &mapVisitor{mapType: mapType}
}

func (p *mapVisitor) ParseType(info types.PkgInfo) (types.Typer, error) {
	keyTypeProvider, err := getTyper(p.mapType.KeyType, info)
	if err != nil {
		return nil, err
	}
	valueTypeProvider, err := getTyper(p.mapType.ValueType, info)
	if err != nil {
		return nil, err
	}
	return types.NewMapType(keyTypeProvider, valueTypeProvider), nil
}

func getTyper(typeFromSpec spec.Type, info types.PkgInfo) (types.Typer, error) {
	typeProvider, err := NewConjureTypeProvider(typeFromSpec)
	if err != nil {
		return nil, err
	}
	typer, err := typeProvider.ParseType(info)
	if err != nil {
		return nil, err
	}
	return typer, nil
}

func (p *mapVisitor) CollectionInitializationIfNeeded(info types.PkgInfo) (*expression.CallExpression, error) {
	typer, err := p.ParseType(info)
	if err != nil {
		return nil, err
	}
	return expression.NewCallExpression(expression.MakeBuiltIn, expression.Type(typer.GoType(info)), expression.IntVal(0)), nil
}

func (p *mapVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return typeCheck == IsMap
}
