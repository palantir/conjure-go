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
	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
)

func VisitConjureDefinition(def spec.ConjureDefinition, visitor ConjureDefinitionVisitor) error {
	for _, typeDef := range def.Types {
		if err := typeDef.Accept(visitor); err != nil {
			return err
		}
	}
	for _, errorDef := range def.Errors {
		if err := visitor.VisitError(errorDef); err != nil {
			return err
		}
	}
	for _, svcDef := range def.Services {
		if err := visitor.VisitService(svcDef); err != nil {
			return err
		}
	}
	return nil
}

type ConjureDefinitionVisitor interface {
	spec.TypeDefinitionVisitor
	VisitError(v spec.ErrorDefinition) error
	VisitService(v spec.ServiceDefinition) error
}
