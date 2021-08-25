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

//go:build ignore
// +build ignore

package old

import (
	"go/token"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
)

func reflectJSONMethods(receiverName string, def spec.TypeDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	addImports(info)
	var decls []astgen.ASTDecl
	if err := def.AcceptFuncs(
		func(def spec.AliasDefinition) error {
			panic("implement me")
		},
		func(def spec.EnumDefinition) error {
			panic("implement me")
		},
		func(def spec.ObjectDefinition) error {
			panic("implement me")
		},
		func(def spec.UnionDefinition) error {
			panic("implement me")
		},
		def.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return decls, nil
}

func reflectAliasTypeMarshalMethods(
	receiverName string,
	receiverType string,
	aliasType spec.Type,
	info types.PkgInfo,
) ([]astgen.ASTDecl, error) {
	panic("implement me")
}

func reflectStructFieldsMarshalMethods(
	receiverName string,
	receiverType string,
	fields []JSONField,
	info types.PkgInfo,
) ([]astgen.ASTDecl, error) {
	var body []astgen.ASTStmt
	for _, field := range fields {
		conjureTypeProvider, err := visitors.NewConjureTypeProvider(field.Type)
		if err != nil {
			return nil, err
		}

		collectionExpression, err := conjureTypeProvider.CollectionInitializationIfNeeded(info)
		if err != nil {
			return nil, err
		}

		if collectionExpression != nil {
			body = append(body, &statement.If{
				Cond: expression.NewBinary(
					expression.NewSelector(expression.VariableVal(receiverName), transforms.ExportedFieldName(field.FieldSelector)),
					token.EQL,
					expression.Nil,
				),
				Body: []astgen.ASTStmt{
					statement.NewAssignment(
						expression.NewSelector(expression.VariableVal(receiverName), transforms.ExportedFieldName(field.FieldSelector)),
						token.ASSIGN,
						collectionExpression,
					),
				},
			})
		}
	}
	aliasTypeName := receiverType + "Alias"
	body = append(body, statement.NewDecl(
		&decl.Alias{
			Name: aliasTypeName,
			Type: expression.Type(receiverType),
		},
	))

	info.AddImports(types.SafeJSONMarshal.ImportPaths()...)
	body = append(body, statement.NewReturn(
		&expression.CallExpression{
			Function: expression.Type(types.SafeJSONMarshal.GoType(info)),
			Args: []astgen.ASTExpr{
				&expression.CallExpression{
					Function: expression.VariableVal(aliasTypeName),
					Args: []astgen.ASTExpr{
						expression.VariableVal(receiverName),
					},
				},
			},
		},
	))

	return []astgen.ASTDecl{newMarshalJSONMethod(receiverName, receiverType, body...)}, nil
}

func newMarshalJSONMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType),
		Function: decl.Function{
			Name: "MarshalJSON",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.ByteSliceType, expression.ErrorType},
			},
			Body: body,
		},
	}
}

func newUnmarshalJSONMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType).Pointer(),
		Function: decl.Function{
			Name: "UnmarshalJSON",
			FuncType: expression.FuncType{
				Params:      expression.FuncParams{expression.NewFuncParam("data", expression.ByteSliceType)},
				ReturnTypes: []expression.Type{expression.ErrorType},
			},
			Body: body,
		},
	}
}
