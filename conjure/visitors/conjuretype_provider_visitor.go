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
)

func NewConjureTypeProvider(rawType spec.Type) (ConjureTypeProvider, error) {
	var provider ConjureTypeProvider
	err := rawType.AcceptFuncs(
		func(v spec.PrimitiveType) error {
			provider = newPrimitiveVisitor(v)
			return nil
		},
		func(v spec.OptionalType) error {
			provider = newOptionalVisitor(v)
			return nil
		},
		func(v spec.ListType) error {
			provider = newListVisitor(v)
			return nil
		},
		func(v spec.SetType) error {
			provider = newSetVisitor(v)
			return nil
		},
		func(v spec.MapType) error {
			provider = newMapVisitor(v)
			return nil
		},
		func(v spec.TypeName) error {
			provider = newReferenceVisitor(v)
			return nil
		},
		func(v spec.ExternalReference) error {
			provider = newExternalVisitor(v)
			return nil
		},
		rawType.ErrorOnUnknown)
	return provider, err
}

func NewConjureTypeProviderTyper(rawType spec.Type, info types.PkgInfo) (types.Typer, error) {
	provider, err := NewConjureTypeProvider(rawType)
	if err != nil {
		return nil, err
	}
	typer, err := provider.ParseType(info)
	if err != nil {
		return nil, err
	}
	return typer, nil
}

func IsSpecificConjureType(rawType spec.Type, typeCheck TypeCheck) (bool, error) {
	conjureTypeProvider, err := NewConjureTypeProvider(rawType)
	if err != nil {
		return false, err
	}
	return conjureTypeProvider.IsSpecificType(typeCheck), nil
}
