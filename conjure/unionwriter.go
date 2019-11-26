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

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/v4/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v4/conjure/transforms"
	"github.com/palantir/conjure-go/v4/conjure/types"
	"github.com/palantir/conjure-go/v4/conjure/visitors"
)

const (
	unionReceiverName = "u"
)

func astForUnion(unionDefinition spec.UnionDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	if err := addImportPathsFromFields(unionDefinition.Union, info); err != nil {
		return nil, err
	}

	unionTypeName := unionDefinition.TypeName.Name
	fieldNameToGoType := make(map[string]string)

	for _, fieldDefinition := range unionDefinition.Union {
		typer, err := fieldDefinitionToTyper(fieldDefinition, info)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to process object %s", unionTypeName)
		}
		goType := typer.GoType(info)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine type for type %q in union type %q", goType, unionTypeName)
		}
		fieldNameToGoType[string(fieldDefinition.FieldName)] = goType
	}

	components := []astgen.ASTDecl{
		unionStructAST(unionTypeName, unionDefinition, fieldNameToGoType),
		unionStructDeserializerAST(unionTypeName, unionDefinition, fieldNameToGoType),
		unionStructDeserializerToStructAST(unionTypeName, unionDefinition),
		toSerializerFuncAST(unionTypeName, unionDefinition, fieldNameToGoType),
		unionMarshalJSONAST(unionTypeName, info),
		unionUnmarshalJSONAST(unionTypeName, info),
		newMarshalYAMLMethod(unionReceiverName, transforms.Export(unionTypeName), info),
		newUnmarshalYAMLMethod(unionReceiverName, transforms.Export(unionTypeName), info),
		acceptMethodAST(unionTypeName, unionDefinition, fieldNameToGoType, info),
		unionTypeVisitorInterfaceAST(unionTypeName, unionDefinition, fieldNameToGoType),
	}
	components = append(components, newFunctionASTs(unionTypeName, unionDefinition, fieldNameToGoType)...)
	return components, nil
}

func fieldDefinitionToTyper(fieldDefinition spec.FieldDefinition, info types.PkgInfo) (types.Typer, error) {
	newConjureTypeProvider, err := visitors.NewConjureTypeProvider(fieldDefinition.Type)
	name := string(fieldDefinition.FieldName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process object %s", name)
	}
	typer, err := newConjureTypeProvider.ParseType(info)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process object %s", name)
	}
	return typer, nil
}

func unionStructAST(unionTypeName string, unionDefinition spec.UnionDefinition, fieldNameToGoType map[string]string) astgen.ASTDecl {
	structFields := []*expression.StructField{
		{
			Name: "typ",
			Type: expression.StringType,
		},
	}
	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		structFields = append(structFields, &expression.StructField{
			Name: transforms.PrivateFieldName(fieldName),
			Type: expression.Type(fieldNameToGoType[fieldName]).Pointer(),
		})
	}
	return decl.NewStruct(
		unionTypeName,
		structFields,
		transforms.Documentation(unionDefinition.Docs),
	)
}

func unionStructDeserializerAST(unionTypeName string, unionDefinition spec.UnionDefinition, fieldNameToGoType map[string]string) astgen.ASTDecl {
	structFields := []*expression.StructField{
		{
			Name: "Type",
			Type: expression.StringType,
			Tag:  fmt.Sprintf(`json:%q`, "type"),
		},
	}
	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		structFields = append(structFields, &expression.StructField{
			Name: transforms.Export(fieldName),
			Type: expression.Type(fieldNameToGoType[fieldName]).Pointer(),
			Tag:  fmt.Sprintf(`json:%q`, fieldName),
		})
	}
	return decl.NewStruct(
		deserializerStructName(unionTypeName),
		structFields,
		"",
	)
}

func unionStructDeserializerToStructAST(unionTypeName string, unionDefinition spec.UnionDefinition) astgen.ASTDecl {
	keyVals := []astgen.ASTExpr{
		expression.NewKeyValue("typ", expression.NewSelector(expression.VariableVal(unionReceiverName), "Type")),
	}
	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		keyVals = append(keyVals,
			expression.NewKeyValue(transforms.PrivateFieldName(fieldName), expression.NewSelector(expression.VariableVal(unionReceiverName), transforms.Export(fieldName))),
		)
	}
	return &decl.Method{
		Function: decl.Function{
			Name: "toStruct",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					expression.Type(transforms.ExportedFieldName(unionTypeName)),
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.Type(transforms.Export(unionTypeName)),
						keyVals...,
					),
				),
			},
		},
		ReceiverName: unionReceiverName,
		ReceiverType: expression.Type(deserializerStructName(unionTypeName)).Pointer(),
	}
}

func toSerializerFuncAST(unionTypeName string, unionDefinition spec.UnionDefinition, fieldNameToGoType map[string]string) astgen.ASTDecl {
	// start with default case
	cases := []statement.CaseClause{
		{
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.Nil,
					expression.NewCallFunction("fmt", "Errorf", expression.StringVal("unknown type %s"), expression.NewSelector(expression.VariableVal(unionReceiverName), "typ")),
				),
			},
		},
	}

	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		fieldNameVarName := transforms.PrivateFieldName(fieldName)

		var caseStmtBody []astgen.ASTStmt
		// TODO(nmiyake): handle case where type is an alias that resolves to an optional
		isOptional, _ := visitors.IsSpecificConjureType(fieldDefinition.Type, visitors.IsOptional)
		if isOptional {
			caseStmtBody = []astgen.ASTStmt{
				statement.NewDecl(
					decl.NewVar(fieldNameVarName, expression.Type(fieldNameToGoType[fieldName])),
				),
				&statement.If{
					Cond: expression.NewBinary(
						expression.NewSelector(expression.VariableVal(unionReceiverName), fieldNameVarName),
						token.NEQ,
						expression.Nil,
					),
					Body: []astgen.ASTStmt{
						statement.NewAssignment(
							expression.VariableVal(fieldNameVarName),
							token.ASSIGN,
							expression.NewUnary(token.MUL, expression.NewSelector(expression.VariableVal(unionReceiverName), fieldNameVarName)),
						),
					},
				},
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.NewStructType(
							&expression.StructField{
								Name: "Type",
								Type: expression.StringType,
								Tag:  `json:"type"`,
							},
							&expression.StructField{
								Name: transforms.Export(fieldName),
								Type: expression.Type(fieldNameToGoType[fieldName]),
								Tag:  fmt.Sprintf(`json:%q`, fieldName),
							},
						),
						expression.NewKeyValue("Type", expression.StringVal(fieldName)),
						expression.NewKeyValue(transforms.Export(fieldName), expression.VariableVal(fieldNameVarName)),
					),
					expression.Nil,
				),
			}
		} else {
			caseStmtBody = []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.NewStructType(
							&expression.StructField{
								Name: "Type",
								Type: expression.StringType,
								Tag:  `json:"type"`,
							},
							&expression.StructField{
								Name: transforms.Export(fieldName),
								Type: expression.Type(fieldNameToGoType[fieldName]),
								Tag:  fmt.Sprintf(`json:%q`, fieldName),
							},
						),
						expression.NewKeyValue("Type", expression.StringVal(fieldName)),
						expression.NewKeyValue(transforms.Export(fieldName), expression.NewUnary(token.MUL, expression.NewSelector(expression.VariableVal(unionReceiverName), fieldNameVarName))),
					),
					expression.Nil,
				),
			}
		}

		cases = append(cases, *statement.NewCaseClause(
			expression.StringVal(fieldName),
			caseStmtBody...,
		))
	}

	return &decl.Method{
		Function: decl.Function{
			Name: "toSerializer",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					expression.Type("interface{}"),
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{
				&statement.Switch{
					Expression: expression.NewSelector(expression.VariableVal(unionReceiverName), "typ"),
					Cases:      cases,
				},
			},
		},
		ReceiverName: unionReceiverName,
		ReceiverType: expression.Type(transforms.Export(unionTypeName)).Pointer(),
	}
}

func deserializerStructName(unionTypeName string) string {
	return transforms.Private(transforms.ExportedFieldName(unionTypeName) + "Deserializer")
}

func unionMarshalJSONAST(unionTypeName string, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.SafeJSONMarshal.ImportPaths()...)
	return newMarshalJSONMethod(unionReceiverName, transforms.Export(unionTypeName),
		&statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.VariableVal("ser"),
				expression.VariableVal("err"),
			},
			Tok: token.DEFINE,
			RHS: expression.NewCallFunction(
				unionReceiverName,
				"toSerializer",
			),
		},
		ifErrNotNilReturnHelper(true, "nil", "err", nil),
		statement.NewReturn(&expression.CallExpression{
			Function: expression.Type(types.SafeJSONMarshal.GoType(info)),
			Args: []astgen.ASTExpr{
				expression.VariableVal("ser"),
			},
		}),
	)
}

func unionUnmarshalJSONAST(unionTypeName string, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.SafeJSONUnmarshal.ImportPaths()...)
	return newUnmarshalJSONMethod(unionReceiverName, transforms.Export(unionTypeName),
		statement.NewDecl(decl.NewVar("deser", expression.Type(deserializerStructName(unionTypeName)))),
		ifErrNotNilReturnErrStatement("err", statement.NewAssignment(
			expression.VariableVal("err"),
			token.DEFINE,
			&expression.CallExpression{
				Function: expression.Type(types.SafeJSONUnmarshal.GoType(info)),
				Args: []astgen.ASTExpr{
					expression.VariableVal(dataVarName),
					expression.NewUnary(token.AND, expression.VariableVal("deser")),
				},
			},
		)),
		statement.NewAssignment(
			expression.NewUnary(token.MUL, expression.VariableVal(unionReceiverName)),
			token.ASSIGN,
			expression.NewCallFunction("deser", "toStruct"),
		),
		statement.NewReturn(expression.Nil),
	)
}

func unionTypeVisitorInterfaceAST(unionTypeName string, unionDefinition spec.UnionDefinition, fieldNameToGoType map[string]string) astgen.ASTDecl {
	var funcs []*expression.InterfaceFunctionDecl

	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		goType := fieldNameToGoType[fieldName]
		funcs = append(funcs, &expression.InterfaceFunctionDecl{
			Name: "Visit" + transforms.ExportedFieldName(fieldName),
			Params: []*expression.FuncParam{
				expression.NewFuncParam("v", expression.Type(goType)),
			},
			ReturnTypes: []expression.Type{
				expression.ErrorType,
			},
		})
	}

	funcs = append(funcs, &expression.InterfaceFunctionDecl{
		Name: "VisitUnknown",
		Params: []*expression.FuncParam{
			expression.NewFuncParam("typeName", expression.StringType),
		},
		ReturnTypes: []expression.Type{
			expression.ErrorType,
		},
	})

	return &decl.Interface{
		Name:          visitorInterfaceName(unionTypeName),
		InterfaceType: *expression.NewInterfaceType(funcs...),
	}
}

func visitorInterfaceName(unionTypeName string) string {
	return transforms.Export(unionTypeName) + "Visitor"
}

func acceptMethodAST(unionTypeName string, unionDefinition spec.UnionDefinition, fieldNameToGoType map[string]string, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports("fmt")
	// start with default case
	cases := []statement.CaseClause{
		{
			Body: []astgen.ASTStmt{
				&statement.If{
					Cond: expression.NewBinary(
						expression.NewSelector(expression.VariableVal(unionReceiverName), "typ"),
						token.EQL,
						expression.StringVal(""),
					),
					Body: []astgen.ASTStmt{
						statement.NewReturn(expression.NewCallFunction("fmt", "Errorf", expression.StringVal("invalid value in union type"))),
					},
				},
				statement.NewReturn(expression.NewCallFunction("v", "VisitUnknown", expression.NewSelector(expression.VariableVal(unionReceiverName), "typ"))),
			},
		},
	}

	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		fieldNameVarName := transforms.PrivateFieldName(fieldName)

		var caseStmtBody []astgen.ASTStmt
		// TODO(nmiyake): handle case where type is an alias that resolves to an optional
		isOptional, _ := visitors.IsSpecificConjureType(fieldDefinition.Type, visitors.IsOptional)
		if isOptional {
			// if the type is an optional and is nil, the value should not be dereferenced
			caseStmtBody = []astgen.ASTStmt{
				statement.NewDecl(
					decl.NewVar(fieldNameVarName, expression.Type(fieldNameToGoType[fieldName])),
				),
				&statement.If{
					Cond: expression.NewBinary(
						expression.NewSelector(expression.VariableVal(unionReceiverName), fieldNameVarName),
						token.NEQ,
						expression.Nil,
					),
					Body: []astgen.ASTStmt{
						statement.NewAssignment(
							expression.VariableVal(fieldNameVarName),
							token.ASSIGN,
							expression.NewUnary(token.MUL, expression.NewSelector(expression.VariableVal(unionReceiverName), fieldNameVarName)),
						),
					},
				},
				statement.NewReturn(expression.NewCallFunction(
					"v",
					"Visit"+transforms.ExportedFieldName(fieldName),
					expression.VariableVal(fieldNameVarName))),
			}
		} else {
			// return dereferenced value directly
			caseStmtBody = []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCallFunction(
						"v",
						"Visit"+transforms.ExportedFieldName(fieldName),
						expression.NewUnary(token.MUL, expression.NewSelector(expression.VariableVal(unionReceiverName), fieldNameVarName)))),
			}
		}
		cases = append(cases, *statement.NewCaseClause(expression.StringVal(fieldName), caseStmtBody...))
	}

	return &decl.Method{
		Function: decl.Function{
			Name: "Accept",
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam("v", expression.Type(visitorInterfaceName(unionTypeName))),
				},
				ReturnTypes: []expression.Type{
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{
				&statement.Switch{
					Expression: expression.NewSelector(expression.VariableVal(unionReceiverName), "typ"),
					Cases:      cases,
				},
			},
		},
		ReceiverName: unionReceiverName,
		ReceiverType: expression.Type(transforms.Export(unionTypeName)).Pointer(),
	}
}

func newFunctionASTs(unionTypeName string, unionDefinition spec.UnionDefinition, fieldNameToGoType map[string]string) []astgen.ASTDecl {
	var decls []astgen.ASTDecl
	for _, fieldDefinition := range unionDefinition.Union {
		fieldName := string(fieldDefinition.FieldName)
		goType := fieldNameToGoType[fieldName]
		decls = append(decls, &decl.Function{
			Name: fmt.Sprintf("New%sFrom%s", transforms.ExportedFieldName(unionTypeName), transforms.ExportedFieldName(fieldName)),
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam("v", expression.Type(goType)),
				},
				ReturnTypes: []expression.Type{
					expression.Type(transforms.Export(unionTypeName)),
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.Type(transforms.Export(unionTypeName)),
						expression.NewKeyValue("typ", expression.StringVal(fieldName)),
						expression.NewKeyValue(transforms.PrivateFieldName(fieldName), expression.NewUnary(token.AND, expression.VariableVal("v"))),
					),
				),
			},
		})
	}
	return decls
}
