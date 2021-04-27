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
	"github.com/pkg/errors"
)

type CustomTypesVisitor struct {
	decls             types.CustomConjureTypes
	conjurePkgToGoPkg func(string) string
}

var _ spec.TypeDefinitionVisitor = &CustomTypesVisitor{}

func TypeNameToTyperName(typeName spec.TypeName) string {
	return typeName.Package + "." + typeName.Name
}

func GetCustomConjureTypes(typeDefinitions []spec.TypeDefinition, conjurePkgToGoPk func(string) string) (types.CustomConjureTypes, error) {
	newCustomTypesVisitor := &CustomTypesVisitor{
		conjurePkgToGoPkg: conjurePkgToGoPk,
		decls:             types.NewCustomConjureTypes(),
	}
	for _, singleType := range typeDefinitions {
		err := singleType.Accept(newCustomTypesVisitor)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create object")
		}
	}
	return newCustomTypesVisitor.decls, nil
}

func (c *CustomTypesVisitor) VisitAlias(aliasDefinition spec.AliasDefinition) error {
	aliasType := TypeNameToTyperName(aliasDefinition.TypeName)
	goPkg := c.getGoPackage(aliasDefinition.TypeName.Package)
	if err := c.decls.Add(aliasType, goPkg, types.NewGoType(aliasType, goPkg)); err != nil {
		return errors.Wrapf(err, "failed to create alias type %s", aliasType)
	}
	return nil
}

func (c *CustomTypesVisitor) VisitEnum(enumDefinition spec.EnumDefinition) error {
	enumType := TypeNameToTyperName(enumDefinition.TypeName)
	goPkg := c.getGoPackage(enumDefinition.TypeName.Package)
	if err := c.decls.Add(enumType, goPkg, types.NewGoType(enumType, goPkg)); err != nil {
		return errors.Wrapf(err, "failed to create enum type %s", enumType)
	}
	return nil
}

func (c *CustomTypesVisitor) VisitObject(objectDefinition spec.ObjectDefinition) error {
	objectType := TypeNameToTyperName(objectDefinition.TypeName)
	goPkg := c.getGoPackage(objectDefinition.TypeName.Package)
	if err := c.decls.Add(objectType, goPkg, types.NewGoType(objectType, goPkg)); err != nil {
		return errors.Wrapf(err, "failed to create object type %s", objectType)
	}
	return nil
}

func (c *CustomTypesVisitor) VisitUnion(unionDefinition spec.UnionDefinition) error {
	unionType := TypeNameToTyperName(unionDefinition.TypeName)
	goPkg := c.getGoPackage(unionDefinition.TypeName.Package)
	if err := c.decls.Add(unionType, goPkg, types.NewGoType(unionType, goPkg)); err != nil {
		return errors.Wrapf(err, "failed to create union type %s", unionType)
	}
	return nil
}

func (c *CustomTypesVisitor) getGoPackage(conjurePackage string) string {
	return c.conjurePkgToGoPkg(conjurePackage)
}

func (c *CustomTypesVisitor) VisitUnknown(typeName string) error {
	return errors.New("Unknown Type found " + typeName)
}
