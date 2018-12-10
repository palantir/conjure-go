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
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/expression"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
)

type PrimitiveVisitor struct {
	primitiveType spec.PrimitiveType
}

func NewPrimitiveVisitor(primitiveType spec.PrimitiveType) ConjureTypeProvider {
	return &PrimitiveVisitor{primitiveType: primitiveType}
}

var _ ConjureTypeProvider = &PrimitiveVisitor{}

func (p *PrimitiveVisitor) ParseType(customTypes types.CustomConjureTypes) (types.Typer, error) {
	switch p.primitiveType {
	case spec.PrimitiveTypeString:
		return types.String, nil
	case spec.PrimitiveTypeInteger:
		return types.Integer, nil
	case spec.PrimitiveTypeDouble:
		return types.Double, nil
	case spec.PrimitiveTypeBoolean:
		return types.Boolean, nil
	case spec.PrimitiveTypeSafelong:
		return types.SafeLongType, nil
	case spec.PrimitiveTypeUuid:
		return types.UUIDType, nil
	case spec.PrimitiveTypeDatetime:
		return types.DateTimeType, nil
	case spec.PrimitiveTypeBinary:
		return types.BinaryType, nil
	case spec.PrimitiveTypeRid:
		return types.Rid, nil
	case spec.PrimitiveTypeAny:
		return types.Any, nil
	case spec.PrimitiveTypeBearertoken:
		return types.Bearertoken, nil
	default:
		return nil, errors.New("Unsupported primitive type " + string(p.primitiveType))
	}
}

func (p *PrimitiveVisitor) CollectionInitializationIfNeeded(customTypes types.CustomConjureTypes, currPkgPath string, pkgAliases map[string]string) (*expression.CallExpression, error) {
	switch p.primitiveType {
	case spec.PrimitiveTypeBinary:
		return &expression.CallExpression{
			Function: expression.VariableVal("make"),
			Args: []astgen.ASTExpr{
				expression.Type("[]byte"),
				expression.IntVal(0),
			},
		}, nil
	default:
		return nil, nil
	}
}

func (p *PrimitiveVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	switch p.primitiveType {
	case spec.PrimitiveTypeBinary:
		return typeCheck == IsBinary
	default:
		return false
	}
}
