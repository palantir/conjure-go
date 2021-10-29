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

package types

import (
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypes(t *testing.T) {
	for _, test := range []struct {
		// in
		Name string
		Spec spec.Type
		// out
		Type Type
		// assert
		Code         string
		Make         string
		IsString     bool
		IsText       bool
		IsBinary     bool
		IsBoolean    bool
		IsOptional   bool
		IsCollection bool
		IsList       bool
	}{
		// simple primitives
		{
			Name: "any",
			Spec: newPrimitive(spec.PrimitiveType_ANY),
			Type: Any{},
			Code: "interface{}",
		},
		{
			Name:   "bearertoken",
			Spec:   newPrimitive(spec.PrimitiveType_BEARERTOKEN),
			Type:   Bearertoken{},
			Code:   "bearertoken.Token",
			IsText: true,
		},
		{
			Name:     "binary",
			Spec:     newPrimitive(spec.PrimitiveType_BINARY),
			Type:     Binary{},
			Code:     "[]byte",
			IsText:   true,
			IsBinary: true,
		},
		{
			Name:      "boolean",
			Spec:      newPrimitive(spec.PrimitiveType_BOOLEAN),
			Type:      Boolean{},
			Code:      "bool",
			IsBoolean: true,
		},
		{
			Name:   "datetime",
			Spec:   newPrimitive(spec.PrimitiveType_DATETIME),
			Type:   DateTime{},
			Code:   "datetime.DateTime",
			IsText: true,
		},
		{
			Name: "double",
			Spec: newPrimitive(spec.PrimitiveType_DOUBLE),
			Type: Double{},
			Code: "float64",
		},
		{
			Name: "integer",
			Spec: newPrimitive(spec.PrimitiveType_INTEGER),
			Type: Integer{},
			Code: "int",
		},
		{
			Name:   "rid",
			Spec:   newPrimitive(spec.PrimitiveType_RID),
			Type:   RID{},
			Code:   "rid.ResourceIdentifier",
			IsText: true,
		},
		{
			Name: "safelong",
			Spec: newPrimitive(spec.PrimitiveType_SAFELONG),
			Type: Safelong{},
			Code: "safelong.SafeLong",
		},
		{
			Name:     "string",
			Spec:     newPrimitive(spec.PrimitiveType_STRING),
			Type:     String{},
			Code:     "string",
			IsString: true,
			IsText:   true,
		},
		{
			Name:   "uuid",
			Spec:   newPrimitive(spec.PrimitiveType_UUID),
			Type:   UUID{},
			Code:   "uuid.UUID",
			IsText: true,
		},
		// optional primitives
		{
			Name:       "optional<any>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_ANY)}),
			Type:       &Optional{Item: Any{}},
			Code:       "*interface{}",
			IsOptional: true,
		},
		{
			Name:       "optional<bearertoken>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_BEARERTOKEN)}),
			Type:       &Optional{Item: Bearertoken{}},
			Code:       "*bearertoken.Token",
			IsText:     true,
			IsOptional: true,
		},
		{
			Name:       "optional<binary>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_BINARY)}),
			Type:       &Optional{Item: Binary{}},
			Code:       "*[]byte",
			IsText:     true,
			IsBinary:   true,
			IsOptional: true,
		},
		{
			Name:       "optional<boolean>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_BOOLEAN)}),
			Type:       &Optional{Item: Boolean{}},
			Code:       "*bool",
			IsBoolean:  true,
			IsOptional: true,
		},
		{
			Name:       "optional<datetime>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_DATETIME)}),
			Type:       &Optional{Item: DateTime{}},
			Code:       "*datetime.DateTime",
			IsText:     true,
			IsOptional: true,
		},
		{
			Name:       "optional<double>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_DOUBLE)}),
			Type:       &Optional{Item: Double{}},
			Code:       "*float64",
			IsOptional: true,
		},
		{
			Name:       "optional<integer>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_INTEGER)}),
			Type:       &Optional{Item: Integer{}},
			Code:       "*int",
			IsOptional: true,
		},
		{
			Name:       "optional<rid>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_RID)}),
			Type:       &Optional{Item: RID{}},
			Code:       "*rid.ResourceIdentifier",
			IsText:     true,
			IsOptional: true,
		},
		{
			Name:       "optional<safelong>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_SAFELONG)}),
			Type:       &Optional{Item: Safelong{}},
			Code:       "*safelong.SafeLong",
			IsOptional: true,
		},
		{
			Name:       "optional<string>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_STRING)}),
			Type:       &Optional{Item: String{}},
			Code:       "*string",
			IsString:   true,
			IsText:     true,
			IsOptional: true,
		},
		{
			Name:       "optional<uuid>",
			Spec:       spec.NewTypeFromOptional(spec.OptionalType{ItemType: newPrimitive(spec.PrimitiveType_UUID)}),
			Type:       &Optional{Item: UUID{}},
			Code:       "*uuid.UUID",
			IsText:     true,
			IsOptional: true,
		},
		// lists of primitives
		{
			Name:         "list<any>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_ANY)}),
			Type:         &List{Item: Any{}},
			Code:         "[]interface{}",
			Make:         "make([]interface{}, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<bearertoken>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_BEARERTOKEN)}),
			Type:         &List{Item: Bearertoken{}},
			Code:         "[]bearertoken.Token",
			Make:         "make([]bearertoken.Token, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<binary>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_BINARY)}),
			Type:         &List{Item: Binary{}},
			Code:         "[][]byte",
			Make:         "make([][]byte, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<boolean>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_BOOLEAN)}),
			Type:         &List{Item: Boolean{}},
			Code:         "[]bool",
			Make:         "make([]bool, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<datetime>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_DATETIME)}),
			Type:         &List{Item: DateTime{}},
			Code:         "[]datetime.DateTime",
			Make:         "make([]datetime.DateTime, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<double>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_DOUBLE)}),
			Type:         &List{Item: Double{}},
			Code:         "[]float64",
			Make:         "make([]float64, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<integer>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_INTEGER)}),
			Type:         &List{Item: Integer{}},
			Code:         "[]int",
			Make:         "make([]int, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<rid>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_RID)}),
			Type:         &List{Item: RID{}},
			Code:         "[]rid.ResourceIdentifier",
			Make:         "make([]rid.ResourceIdentifier, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<safelong>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_SAFELONG)}),
			Type:         &List{Item: Safelong{}},
			Code:         "[]safelong.SafeLong",
			Make:         "make([]safelong.SafeLong, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<string>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_STRING)}),
			Type:         &List{Item: String{}},
			Code:         "[]string",
			Make:         "make([]string, 0)",
			IsList:       true,
			IsCollection: true,
		},
		{
			Name:         "list<uuid>",
			Spec:         spec.NewTypeFromList(spec.ListType{ItemType: newPrimitive(spec.PrimitiveType_UUID)}),
			Type:         &List{Item: UUID{}},
			Code:         "[]uuid.UUID",
			Make:         "make([]uuid.UUID, 0)",
			IsList:       true,
			IsCollection: true,
		},
		// maps of primitives
		{
			Name: "map<any, any>",
			Spec: spec.NewTypeFromMap(spec.MapType{
				KeyType:   newPrimitive(spec.PrimitiveType_ANY),
				ValueType: newPrimitive(spec.PrimitiveType_ANY),
			}),
			Type:         &Map{Key: Any{}, Val: Any{}},
			Code:         "map[interface{}]interface{}",
			Make:         "make(map[interface{}]interface{}, 0)",
			IsCollection: true,
		},
		{
			Name: "map<binary, binary>",
			Spec: spec.NewTypeFromMap(spec.MapType{
				KeyType:   newPrimitive(spec.PrimitiveType_BINARY),
				ValueType: newPrimitive(spec.PrimitiveType_BINARY),
			}),
			Type:         &Map{Key: Binary{}, Val: Binary{}},
			Code:         "map[binary.Binary][]byte",
			Make:         "make(map[binary.Binary][]byte, 0)",
			IsCollection: true,
		},
		{
			Name: "map<bool, bool>",
			Spec: spec.NewTypeFromMap(spec.MapType{
				KeyType:   newPrimitive(spec.PrimitiveType_BOOLEAN),
				ValueType: newPrimitive(spec.PrimitiveType_BOOLEAN),
			}),
			Type:         &Map{Key: Boolean{}, Val: Boolean{}},
			Code:         "map[boolean.Boolean]bool",
			Make:         "make(map[boolean.Boolean]bool, 0)",
			IsCollection: true,
		},
		{
			Name: "map<string, string>",
			Spec: spec.NewTypeFromMap(spec.MapType{
				KeyType:   newPrimitive(spec.PrimitiveType_STRING),
				ValueType: newPrimitive(spec.PrimitiveType_STRING),
			}),
			Type:         &Map{Key: String{}, Val: String{}},
			Code:         "map[string]string",
			Make:         "make(map[string]string, 0)",
			IsCollection: true,
		},
		// aliases of primitives
		// aliases of optionals
		// aliases of collections
	} {
		t.Run(test.Name, func(t *testing.T) {
			typ := (&namedTypes{}).GetBySpec(test.Spec)
			require.Equal(t, test.Type, typ)

			stmt := jen.Var().Id("_").Add(typ.Code())
			expect := "var _ " + test.Code
			require.Equal(t, expect, stmt.GoString(), "Code() output")

			mk := typ.Make()
			if mk == nil {
				assert.Equal(t, test.Make, "", "Make() should not be nil")
			} else {
				assert.Equal(t, test.Make, mk.GoString(), "Make() should not be nil")
			}

			assert.Equal(t, test.IsString, typ.IsString(), "IsString() output")
			assert.Equal(t, test.IsText, typ.IsText(), "IsText() output")
			assert.Equal(t, test.IsBinary, typ.IsBinary(), "IsBinary() output")
			assert.Equal(t, test.IsBoolean, typ.IsBoolean(), "IsBoolean() output")
			assert.Equal(t, test.IsOptional, typ.IsOptional(), "IsOptional() output")
			assert.Equal(t, test.IsCollection, typ.IsCollection(), "IsCollection() output")
			assert.Equal(t, test.IsList, typ.IsList(), "IsList() output")
		})
	}
}
