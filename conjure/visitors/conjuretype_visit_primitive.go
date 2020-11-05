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

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

type primitiveVisitor struct {
	primitiveType spec.PrimitiveType
}

func newPrimitiveVisitor(primitiveType spec.PrimitiveType) ConjureTypeProvider {
	return &primitiveVisitor{primitiveType: primitiveType}
}

func (p *primitiveVisitor) ParseType(_ types.PkgInfo) (types.Typer, error) {
	switch p.primitiveType.Value() {
	case spec.PrimitiveType_ANY:
		return types.Any, nil
	case spec.PrimitiveType_BEARERTOKEN:
		return types.Bearertoken, nil
	case spec.PrimitiveType_BINARY:
		return types.BinaryType, nil
	case spec.PrimitiveType_BOOLEAN:
		return types.Boolean, nil
	case spec.PrimitiveType_DATETIME:
		return types.DateTime, nil
	case spec.PrimitiveType_DOUBLE:
		return types.Double, nil
	case spec.PrimitiveType_INTEGER:
		return types.Integer, nil
	case spec.PrimitiveType_RID:
		return types.RID, nil
	case spec.PrimitiveType_SAFELONG:
		return types.SafeLong, nil
	case spec.PrimitiveType_STRING:
		return types.String, nil
	case spec.PrimitiveType_UUID:
		return types.UUID, nil
	default:
		typ, _ := p.primitiveType.MarshalText()
		return nil, errors.New("Unsupported primitive type " + string(typ))
	}
}

func (p *primitiveVisitor) CollectionInitializationIfNeeded(_ types.PkgInfo) (*expression.CallExpression, error) {
	switch p.primitiveType.Value() {
	case spec.PrimitiveType_BINARY:
		return expression.NewCallExpression(expression.MakeBuiltIn, expression.ByteSliceType, expression.IntVal(0)), nil
	default:
		return nil, nil
	}
}

func (p *primitiveVisitor) IsSpecificType(typeCheck TypeCheck) bool {
	switch typeCheck {
	case IsString:
		return p.primitiveType.Value() == spec.PrimitiveType_STRING
	case IsBinary:
		return p.primitiveType.Value() == spec.PrimitiveType_BINARY
	case IsBoolean:
		return p.primitiveType.Value() == spec.PrimitiveType_BOOLEAN
	case IsText:
		switch p.primitiveType.Value() {
		case spec.PrimitiveType_BEARERTOKEN,
			spec.PrimitiveType_DATETIME,
			spec.PrimitiveType_RID,
			spec.PrimitiveType_STRING,
			spec.PrimitiveType_UUID:
			return true
		}
	}
	return false
}
