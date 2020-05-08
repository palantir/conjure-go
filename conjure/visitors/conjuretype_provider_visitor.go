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
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v5/conjure/types"
)

type conjureTypeVisitor struct {
	conjureTypeProvider ConjureTypeProvider
}

var _ spec.TypeVisitor = &conjureTypeVisitor{}

func NewConjureTypeProvider(rawType spec.Type) (ConjureTypeProvider, error) {
	createTypeProvider := conjureTypeVisitor{}
	if err := rawType.Accept(&createTypeProvider); err != nil {
		return nil, err
	}
	return createTypeProvider.conjureTypeProvider, nil
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

func (c *conjureTypeVisitor) VisitPrimitive(v spec.PrimitiveType) error {
	c.conjureTypeProvider = newPrimitiveVisitor(v)
	return nil
}

func (c *conjureTypeVisitor) VisitSet(v spec.SetType) error {
	c.conjureTypeProvider = newSetVisitor(v)
	return nil
}

func (c *conjureTypeVisitor) VisitList(v spec.ListType) error {
	c.conjureTypeProvider = newListVisitor(v)
	return nil
}

func (c *conjureTypeVisitor) VisitOptional(v spec.OptionalType) error {
	c.conjureTypeProvider = newOptionalVisitor(v)
	return nil
}

func (c *conjureTypeVisitor) VisitMap(v spec.MapType) error {
	c.conjureTypeProvider = newMapVisitor(v)
	return nil
}
func (c *conjureTypeVisitor) VisitReference(v spec.TypeName) error {
	c.conjureTypeProvider = newReferenceVisitor(v)
	return nil
}
func (c *conjureTypeVisitor) VisitExternal(v spec.ExternalReference) error {
	c.conjureTypeProvider = newExternalVisitor(v)
	return nil
}
func (c *conjureTypeVisitor) VisitUnknown(v string) error {
	return errors.New("Unsupported Type Visit Unknown " + v)
}
