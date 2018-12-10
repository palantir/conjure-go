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

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
)

type CreateTypeProvider struct {
	conjureTypeProvider ConjureTypeProvider
}

var _ spec.TypeVisitor = &CreateTypeProvider{}

func NewConjureTypeProvider(rawType spec.Type) (ConjureTypeProvider, error) {
	createTypeProvider := CreateTypeProvider{}
	err := rawType.Accept(&createTypeProvider)
	if err != nil {
		return nil, err
	}
	return createTypeProvider.conjureTypeProvider, nil
}

func NewConjureTypeProviderTyper(rawType spec.Type, customTypes types.CustomConjureTypes) (types.Typer, error) {
	createTypeProvider := CreateTypeProvider{}
	err := rawType.Accept(&createTypeProvider)
	if err != nil {
		return nil, err
	}
	typer, err := createTypeProvider.conjureTypeProvider.ParseType(customTypes)
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

func (c *CreateTypeProvider) VisitPrimitive(v spec.PrimitiveType) error {
	c.conjureTypeProvider = NewPrimitiveVisitor(v)
	return nil
}

func (c *CreateTypeProvider) VisitSet(v spec.SetType) error {
	c.conjureTypeProvider = NewSetVisitor(v)
	return nil
}

func (c *CreateTypeProvider) VisitList(v spec.ListType) error {
	c.conjureTypeProvider = NewListVisitor(v)
	return nil
}

func (c *CreateTypeProvider) VisitOptional(v spec.OptionalType) error {
	c.conjureTypeProvider = NewOptionalVisitor(v)
	return nil
}

func (c *CreateTypeProvider) VisitMap(v spec.MapType) error {
	c.conjureTypeProvider = NewMapVisitor(v)
	return nil
}
func (c *CreateTypeProvider) VisitReference(v spec.TypeName) error {
	c.conjureTypeProvider = NewReferenceVisitor(v)
	return nil
}
func (c *CreateTypeProvider) VisitExternal(v spec.ExternalReference) error {
	c.conjureTypeProvider = NewExternalVisitor(v)
	return nil
}
func (c *CreateTypeProvider) VisitUnknown(v string) error {
	return errors.New("Unsupported Type Visit Unknown " + v)
}
