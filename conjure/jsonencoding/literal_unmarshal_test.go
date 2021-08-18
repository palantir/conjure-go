// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package jsonencoding

import (
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/goastwriter"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructFieldJSONMethods(t *testing.T) {

	stmts, err := visitStructFieldsUnmarshalGJSONMethodBody("x", "Type", []JSONField{
		{
			JSONKey: "fieldAny",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_ANY)),
		},
		{
			JSONKey: "fieldString",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
		},
		{
			JSONKey: "fieldInt",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER)),
		},
		{
			JSONKey: "fieldDatetime",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME)),
		},
		{
			JSONKey: "fieldSafelong",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_SAFELONG)),
		},
		{
			JSONKey: "fieldUUID",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_UUID)),
		},
		{
			JSONKey: "fieldBinary",
			Type:    spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_BINARY)),
		},
		{
			JSONKey: "fieldOptionalString",
			Type:    spec.NewTypeFromOptional(spec.OptionalType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING))}),
		},
		{
			JSONKey: "fieldListString",
			Type:    spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING))}),
		},
		{
			JSONKey: "fieldListInteger",
			Type:    spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_INTEGER))}),
		},
		{
			JSONKey: "fieldListDatetime",
			Type:    spec.NewTypeFromList(spec.ListType{ItemType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME))}),
		},
		{
			JSONKey: "fieldMapStringString",
			Type: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
				ValueType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
			}),
		},
		{
			JSONKey: "fieldMapDatetimeSafelong",
			Type: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_DATETIME)),
				ValueType: spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_SAFELONG)),
			}),
		},
	}, types.NewPkgInfo("main", types.NewCustomConjureTypes()))
	require.NoError(t, err)
	out, err := goastwriter.Write("main", &decl.Method{
		Function: decl.Function{
			Name: "UnmarshalJSON",
			FuncType: expression.FuncType{
				Params:      expression.FuncParams{expression.NewFuncParam("data", expression.ByteSliceType)},
				ReturnTypes: []expression.Type{expression.ErrorType},
			},
			Body: stmts,
		},
		ReceiverName: "x",
		ReceiverType: "*Foo",
	})
	assert.NoError(t, err)
	t.Log(string(out))
}
