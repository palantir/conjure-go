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

package conjure

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/transforms"
	"github.com/palantir/conjure-go/conjure/types"
	"github.com/palantir/conjure-go/conjure/visitors"
)

func astForObject(objectDefinition spec.ObjectDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	containsCollection := false
	var structFields []*expression.StructField

	for _, fieldDefinition := range objectDefinition.Fields {
		conjureTypeProvider, err := visitors.NewConjureTypeProvider(fieldDefinition.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create type provider for field %s for object %s",
				fieldDefinition.FieldName,
				objectDefinition.TypeName.Name,
			)
		}
		typer, err := conjureTypeProvider.ParseType(info)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse type field %s for object %s",
				fieldDefinition.FieldName,
				objectDefinition.TypeName.Name,
			)
		}
		info.AddImports(typer.ImportPaths()...)

		collectionExpression, err := conjureTypeProvider.CollectionInitializationIfNeeded(info)
		if err != nil {
			return nil, err
		}
		if collectionExpression != nil {
			// if there is a map or slice field, the struct contains a collection
			containsCollection = true
		}
		fieldName := string(fieldDefinition.FieldName)
		tags := []string{
			fmt.Sprintf("json:%q", fieldName),
			fmt.Sprintf(`yaml:"%s,omitempty"`, fieldName),
		}

		comment := transforms.Documentation(fieldDefinition.Docs)
		if comment != "" {
			// backtick characters ("`") are really painful to deal with in struct tags
			// (which are themselves defined within backtick literals), so replace with
			// double quotes instead.
			tags = append(tags, fmt.Sprintf("conjure-docs:%q", strings.Replace(comment, "`", `"`, -1)))
		}
		structFields = append(structFields, &expression.StructField{
			Name:    transforms.ExportedFieldName(fieldName),
			Type:    expression.Type(typer.GoType(info)),
			Tag:     strings.Join(tags, " "),
			Comment: comment,
		})
	}

	comment := transforms.Documentation(objectDefinition.Docs)
	decls := []astgen.ASTDecl{
		decl.NewStruct(objectDefinition.TypeName.Name, structFields, comment),
	}
	if containsCollection {
		for _, f := range []serdeFunc{
			astForStructJSONMarshal,
			astForStructJSONUnmarshal,
			astForStructYAMLMarshal,
			astForStructYAMLUnmarshal,
		} {
			serdeDecl, err := f(objectDefinition, info)
			if err != nil {
				return nil, err
			}
			decls = append(decls, serdeDecl)
		}
	}
	return decls, nil
}

const (
	objReceiverName = "o"
)

type serdeFunc func(objectDefinition spec.ObjectDefinition, info types.PkgInfo) (astgen.ASTDecl, error)

func astForStructJSONMarshal(objectDefinition spec.ObjectDefinition, info types.PkgInfo) (astgen.ASTDecl, error) {
	var body []astgen.ASTStmt
	marshalInit, err := structMarshalInitDecls(objectDefinition, objReceiverName, info)
	if err != nil {
		return nil, err
	}
	body = append(body, marshalInit...)

	aliasTypeName := objectDefinition.TypeName.Name + "Alias"
	body = append(body, statement.NewDecl(
		&decl.Alias{
			Name: aliasTypeName,
			Type: expression.Type(objectDefinition.TypeName.Name),
		},
	))

	info.AddImports(types.CodecJSON.ImportPaths()...)
	body = append(body, statement.NewReturn(
		expression.NewCallFunction(
			types.CodecJSON.GoType(info),
			"Marshal",
			&expression.CallExpression{
				Function: expression.VariableVal(aliasTypeName),
				Args: []astgen.ASTExpr{
					expression.VariableVal(objReceiverName),
				},
			},
		),
	))

	return newMarshalJSONMethod(objReceiverName, objectDefinition.TypeName.Name, body...), nil
}

func astForStructJSONUnmarshal(objectDefinition spec.ObjectDefinition, info types.PkgInfo) (astgen.ASTDecl, error) {
	var body []astgen.ASTStmt
	aliasTypeName := objectDefinition.TypeName.Name + "Alias"
	body = append(body, statement.NewDecl(
		&decl.Alias{
			Name: aliasTypeName,
			Type: expression.Type(objectDefinition.TypeName.Name),
		},
	))

	rawVarName := fmt.Sprint("raw", objectDefinition.TypeName.Name)
	body = append(body, statement.NewDecl(
		decl.NewVar(rawVarName, expression.Type(aliasTypeName)),
	))

	info.AddImports(types.CodecJSON.ImportPaths()...)
	body = append(body, ifErrNotNilReturnErrStatement("err",
		statement.NewAssignment(
			expression.VariableVal("err"),
			token.DEFINE,
			expression.NewCallFunction(
				types.CodecJSON.GoType(info),
				"Unmarshal",
				expression.VariableVal(dataVarName), expression.NewUnary(token.AND, expression.VariableVal(rawVarName)),
			),
		),
	))

	marshalInit, err := structMarshalInitDecls(objectDefinition, rawVarName, info)
	if err != nil {
		return nil, err
	}
	body = append(body, marshalInit...)

	body = append(body, statement.NewAssignment(
		expression.NewStar(expression.VariableVal(objReceiverName)),
		token.ASSIGN,
		&expression.CallExpression{
			Function: expression.VariableVal(objectDefinition.TypeName.Name),
			Args: []astgen.ASTExpr{
				expression.VariableVal(rawVarName),
			},
		},
	))

	body = append(body, statement.NewReturn(expression.Nil))

	return newUnmarshalJSONMethod(objReceiverName, objectDefinition.TypeName.Name, body...), nil
}

func astForStructYAMLMarshal(objectDefinition spec.ObjectDefinition, info types.PkgInfo) (astgen.ASTDecl, error) {
	var body []astgen.ASTStmt
	marshalInit, err := structMarshalInitDecls(objectDefinition, objReceiverName, info)
	if err != nil {
		return nil, err
	}
	body = append(body, marshalInit...)

	aliasTypeName := objectDefinition.TypeName.Name + "Alias"
	body = append(body, statement.NewDecl(
		&decl.Alias{
			Name: aliasTypeName,
			Type: expression.Type(objectDefinition.TypeName.Name),
		},
	))

	body = append(body, statement.NewReturn(
		&expression.CallExpression{
			Function: expression.VariableVal(aliasTypeName),
			Args: []astgen.ASTExpr{
				expression.VariableVal(objReceiverName),
			},
		},
		expression.Nil,
	))

	return newMarshalYAMLMethod(objReceiverName, objectDefinition.TypeName.Name, body...), nil
}

func astForStructYAMLUnmarshal(objectDefinition spec.ObjectDefinition, info types.PkgInfo) (astgen.ASTDecl, error) {
	var body []astgen.ASTStmt
	aliasTypeName := objectDefinition.TypeName.Name + "Alias"
	body = append(body, statement.NewDecl(
		&decl.Alias{
			Name: aliasTypeName,
			Type: expression.Type(objectDefinition.TypeName.Name),
		},
	))

	rawVarName := fmt.Sprint("raw", objectDefinition.TypeName.Name)
	body = append(body, statement.NewDecl(
		decl.NewVar(rawVarName, expression.Type(aliasTypeName)),
	))

	body = append(body, ifErrNotNilReturnErrStatement("err",
		statement.NewAssignment(
			expression.VariableVal("err"),
			token.DEFINE,
			&expression.CallExpression{
				Function: expression.Type("unmarshal"),
				Args: []astgen.ASTExpr{
					expression.NewUnary(token.AND, expression.VariableVal(rawVarName)),
				},
			},
		),
	))

	marshalInit, err := structMarshalInitDecls(objectDefinition, rawVarName, info)
	if err != nil {
		return nil, err
	}
	body = append(body, marshalInit...)

	body = append(body, statement.NewAssignment(
		expression.NewStar(expression.VariableVal(objReceiverName)),
		token.ASSIGN,
		&expression.CallExpression{
			Function: expression.VariableVal(objectDefinition.TypeName.Name),
			Args: []astgen.ASTExpr{
				expression.VariableVal(rawVarName),
			},
		},
	))

	body = append(body, statement.NewReturn(expression.Nil))

	return newUnmarshalYAMLMethod(objReceiverName, objectDefinition.TypeName.Name, body...), nil
}

func structMarshalInitDecls(objectDefinition spec.ObjectDefinition, variableVal string, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var decls []astgen.ASTStmt
	for _, fieldDefinition := range objectDefinition.Fields {
		conjureTypeProvider, err := visitors.NewConjureTypeProvider(fieldDefinition.Type)
		if err != nil {
			return nil, err
		}

		collectionExpression, err := conjureTypeProvider.CollectionInitializationIfNeeded(info)
		if err != nil {
			return nil, err
		}

		if collectionExpression != nil {
			currFieldName := string(fieldDefinition.FieldName)
			decls = append(decls, &statement.If{
				Cond: expression.NewBinary(
					expression.NewSelector(expression.VariableVal(variableVal), transforms.ExportedFieldName(currFieldName)),
					token.EQL,
					expression.Nil,
				),
				Body: []astgen.ASTStmt{
					statement.NewAssignment(
						expression.NewSelector(expression.VariableVal(variableVal), transforms.ExportedFieldName(currFieldName)),
						token.ASSIGN,
						collectionExpression,
					),
				},
			})
		}
	}
	return decls, nil
}
