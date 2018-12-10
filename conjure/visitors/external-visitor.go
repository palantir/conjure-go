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

type ExternalVisitor struct {
	externalType spec.ExternalReference
}

func NewExternalVisitor(externalType spec.ExternalReference) ConjureTypeProvider {
	return &ExternalVisitor{externalType: externalType}
}

var _ ConjureTypeProvider = &ExternalVisitor{}

func (p *ExternalVisitor) ParseType(customTypes types.CustomConjureTypes) (types.Typer, error) {
	return types.NewGoTypeFromExternalType(p.externalType), nil
}

func (p *ExternalVisitor) CollectionInitializationIfNeeded(customTypes types.CustomConjureTypes, currPkgPath string, pkgAliases map[string]string) (*expression.CallExpression, error) {
	return nil, nil
}

func (p *ExternalVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return false
}
