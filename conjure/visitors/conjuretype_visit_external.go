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

type externalVisitor struct {
	externalType spec.ExternalReference
}

func newExternalVisitor(externalType spec.ExternalReference) ConjureTypeProvider {
	return &externalVisitor{externalType: externalType}
}

func (p *externalVisitor) ParseType(info types.PkgInfo) (types.Typer, error) {
	t, err := types.NewGoTypeFromExternalType(p.externalType)
	if err == nil {
		return t, nil
	}

	nestedTypeProvider, err := NewConjureTypeProvider(p.externalType.Fallback)
	if err != nil {
		return nil, err
	}

	return nestedTypeProvider.ParseType(info)
}

func (p *externalVisitor) CollectionInitializationIfNeeded(types.PkgInfo) (*expression.CallExpression, error) {
	return nil, nil
}

func (p *externalVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	if p != nil && typeCheck == IsSafeMarker {
		ref := p.externalType.ExternalReference
		return ref.Package == "com.palantir.logsafe" && ref.Name == "Safe"
	}
	return false
}
