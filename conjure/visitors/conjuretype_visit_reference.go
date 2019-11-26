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
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/v4/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v4/conjure/types"
)

type referenceVisitor struct {
	typeName spec.TypeName
}

func newReferenceVisitor(typeName spec.TypeName) ConjureTypeProvider {
	return &referenceVisitor{typeName: typeName}
}

func (r *referenceVisitor) ParseType(info types.PkgInfo) (types.Typer, error) {
	name := TypeNameToTyperName(r.typeName)
	if custom, ok := info.CustomTypes().Get(name); ok {
		return custom, nil
	}
	return nil, errors.New("Could not find specified conjure type " + name)
}

func (r *referenceVisitor) CollectionInitializationIfNeeded(info types.PkgInfo) (*expression.CallExpression, error) {
	return nil, nil
}

func (r *referenceVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	return false
}
