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

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
)

func ConjureDefinitionsByPackage(root spec.ConjureDefinition) (map[string]spec.ConjureDefinition, error) {
	collector := &ConjureTypePackageCollector{root: root}
	if err := VisitConjureDefinition(root, collector); err != nil {
		return nil, err
	}
	return collector.packages, nil
}

type ConjureTypePackageCollector struct {
	root     spec.ConjureDefinition
	packages map[string]spec.ConjureDefinition
}

func (c *ConjureTypePackageCollector) VisitAlias(v spec.AliasDefinition) error {
	c.putType(v.TypeName.Package, spec.NewTypeDefinitionFromAlias(v))
	return nil
}

func (c *ConjureTypePackageCollector) VisitEnum(v spec.EnumDefinition) error {
	c.putType(v.TypeName.Package, spec.NewTypeDefinitionFromEnum(v))
	return nil
}

func (c *ConjureTypePackageCollector) VisitObject(v spec.ObjectDefinition) error {
	c.putType(v.TypeName.Package, spec.NewTypeDefinitionFromObject(v))
	return nil
}

func (c *ConjureTypePackageCollector) VisitUnion(v spec.UnionDefinition) error {
	c.putType(v.TypeName.Package, spec.NewTypeDefinitionFromUnion(v))
	return nil
}

func (c *ConjureTypePackageCollector) VisitError(v spec.ErrorDefinition) error {
	packageName := v.ErrorName.Package
	c.putPackage(packageName)
	existing := c.packages[packageName]
	existing.Errors = append(existing.Errors, v)
	c.packages[packageName] = existing
	return nil
}

func (c *ConjureTypePackageCollector) VisitService(v spec.ServiceDefinition) error {
	packageName := v.ServiceName.Package
	c.putPackage(packageName)
	existing := c.packages[packageName]
	existing.Services = append(existing.Services, v)
	c.packages[packageName] = existing
	return nil
}

func (c *ConjureTypePackageCollector) putType(packageName string, def spec.TypeDefinition) {
	c.putPackage(packageName)
	existing := c.packages[packageName]
	existing.Types = append(existing.Types, def)
	c.packages[packageName] = existing
}

func (c *ConjureTypePackageCollector) putPackage(packageName string) {
	if c.packages == nil {
		c.packages = make(map[string]spec.ConjureDefinition)
	}
	if _, ok := c.packages[packageName]; !ok {
		c.packages[packageName] = spec.ConjureDefinition{Version: c.root.Version}
	}
}

func (c *ConjureTypePackageCollector) VisitUnknown(typeName string) error {
	return errors.New("Unknown Type found " + typeName)
}
