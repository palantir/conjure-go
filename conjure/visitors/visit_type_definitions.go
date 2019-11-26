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

	"github.com/palantir/conjure-go/v4/conjure-api/conjure/spec"
)

type ConjureTypeFilterVisitor struct {
	AliasDefinitions  []spec.AliasDefinition
	EnumDefinitions   []spec.EnumDefinition
	ObjectDefinitions []spec.ObjectDefinition
	UnionDefinitions  []spec.UnionDefinition
}

var _ spec.TypeDefinitionVisitor = &ConjureTypeFilterVisitor{}

func NewConjureTypeFilterVisitor() *ConjureTypeFilterVisitor {
	return &ConjureTypeFilterVisitor{}
}

func (c *ConjureTypeFilterVisitor) VisitAlias(aliasDefinition spec.AliasDefinition) error {
	c.AliasDefinitions = append(c.AliasDefinitions, aliasDefinition)
	return nil
}

func (c *ConjureTypeFilterVisitor) VisitEnum(enumDefinition spec.EnumDefinition) error {
	c.EnumDefinitions = append(c.EnumDefinitions, enumDefinition)
	return nil
}

func (c *ConjureTypeFilterVisitor) VisitObject(objectDefinition spec.ObjectDefinition) error {
	c.ObjectDefinitions = append(c.ObjectDefinitions, objectDefinition)
	return nil
}

func (c *ConjureTypeFilterVisitor) VisitUnion(unionDefinition spec.UnionDefinition) error {
	c.UnionDefinitions = append(c.UnionDefinitions, unionDefinition)
	return nil
}

func (c *ConjureTypeFilterVisitor) VisitUnknown(typeName string) error {
	return errors.New("Unknown Type found " + typeName)
}
