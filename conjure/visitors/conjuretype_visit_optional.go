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
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/goastwriter/expression"
)

type optionalVisitor struct {
	optionalType spec.OptionalType
}

func newOptionalVisitor(optionalType spec.OptionalType) ConjureTypeProvider {
	return &optionalVisitor{optionalType: optionalType}
}

func (p *optionalVisitor) ParseType(info types.PkgInfo) (types.Typer, error) {
	nestedTypeProvider, err := NewConjureTypeProvider(p.optionalType.ItemType)
	if err != nil {
		return nil, err
	}
	typer, err := nestedTypeProvider.ParseType(info)
	if err != nil {
		return nil, err
	}
	return types.NewOptionalType(typer), nil
}

func (p *optionalVisitor) CollectionInitializationIfNeeded(types.PkgInfo) (*expression.CallExpression, error) {
	return nil, nil
}

func (p *optionalVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	if typeCheck == IsOptional {
		return true
	}
	inner, err := NewConjureTypeProvider(p.optionalType.ItemType)
	if err != nil {
		return false
	}
	return inner.IsSpecificType(typeCheck)
}
