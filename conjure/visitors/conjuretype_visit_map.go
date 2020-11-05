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

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

type mapVisitor struct {
	mapType spec.MapType
}

func newMapVisitor(mapType spec.MapType) ConjureTypeProvider {
	return &mapVisitor{mapType: mapType}
}

func (p *mapVisitor) ParseType(info types.PkgInfo) (types.Typer, error) {
	keyTypeProvider, err := NewConjureTypeProvider(p.mapType.KeyType)
	if err != nil {
		return nil, err
	}
	keyTyper, err := keyTypeProvider.ParseType(info)
	if err != nil {
		return nil, err
	}

	// Use binary.Binary for map keys since []byte is invalid in go maps.
	if keyTypeProvider.IsSpecificType(IsBinary) {
		keyTyper = types.BinaryPkg
	}
	// Use boolean.Boolean for map keys since conjure boolean keys are serialized as strings
	if keyTypeProvider.IsSpecificType(IsBoolean) {
		keyTyper = types.BooleanPkg
	}

	valueTypeProvider, err := NewConjureTypeProvider(p.mapType.ValueType)
	if err != nil {
		return nil, err
	}
	valueTyper, err := valueTypeProvider.ParseType(info)
	if err != nil {
		return nil, err
	}

	return types.NewMapType(keyTyper, valueTyper), nil
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
