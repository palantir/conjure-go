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

type primitiveVisitor struct {
	primitiveType spec.PrimitiveType
}

func newPrimitiveVisitor(primitiveType spec.PrimitiveType) ConjureTypeProvider {
	return &primitiveVisitor{primitiveType: primitiveType}
}

func (p *primitiveVisitor) ParseType(_ types.PkgInfo) (types.Typer, error) {
	switch p.primitiveType {
	case spec.PrimitiveTypeAny:
		return types.Any, nil
	case spec.PrimitiveTypeBearertoken:
		return types.Bearertoken, nil
	case spec.PrimitiveTypeBinary:
		return types.BinaryType, nil
	case spec.PrimitiveTypeBoolean:
		return types.Boolean, nil
	case spec.PrimitiveTypeDatetime:
		return types.DateTime, nil
	case spec.PrimitiveTypeDouble:
		return types.Double, nil
	case spec.PrimitiveTypeInteger:
		return types.Integer, nil
	case spec.PrimitiveTypeRid:
		return types.RID, nil
	case spec.PrimitiveTypeSafelong:
		return types.SafeLong, nil
	case spec.PrimitiveTypeString:
		return types.String, nil
	case spec.PrimitiveTypeUuid:
		return types.UUID, nil
	default:
		return nil, errors.New("Unsupported primitive type " + string(p.primitiveType))
	}
}

func (p *primitiveVisitor) CollectionInitializationIfNeeded(_ types.PkgInfo) (*expression.CallExpression, error) {
	switch p.primitiveType {
	case spec.PrimitiveTypeBinary:
		return expression.NewCallExpression(expression.MakeBuiltIn, expression.ByteSliceType, expression.IntVal(0)), nil
	default:
		return nil, nil
	}
}

func (p *primitiveVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	switch typeCheck {
	case IsString:
		return p.primitiveType == spec.PrimitiveTypeString
	case IsBinary:
		return p.primitiveType == spec.PrimitiveTypeBinary
	case IsText:
		switch p.primitiveType {
		case spec.PrimitiveTypeBearertoken,
			spec.PrimitiveTypeDatetime,
			spec.PrimitiveTypeRid,
			spec.PrimitiveTypeString,
			spec.PrimitiveTypeUuid:
			return true
		}
	}
	return false
}
