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

	"github.com/danverbraganza/varcaser/varcaser"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	astspec "github.com/palantir/goastwriter/spec"
	"github.com/palantir/goastwriter/statement"

	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v5/conjure/transforms"
	"github.com/palantir/conjure-go/v5/conjure/types"
)

const (
	enumReceiverName   = "e"
	enumUpperVarName   = "v"
	enumPatternVarName = "enumValuePattern"
	enumValuePattern   = "^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"
)

func astForEnum(enumDefinition spec.EnumDefinition, info types.PkgInfo) []astgen.ASTDecl {
	enumName := enumDefinition.TypeName.Name

	typeDef := &decl.Alias{
		Comment: transforms.Documentation(enumDefinition.Docs),
		Name:    enumName,
		Type:    expression.StringType,
	}

	toCamelCase := varcaser.Caser{From: varcaser.ScreamingSnakeCase, To: varcaser.UpperCamelCase}.String

	var vals []*astspec.Value
	for _, currVal := range enumDefinition.Values {
		vals = append(vals, &astspec.Value{
			Comment: transforms.Documentation(currVal.Docs),
			Names:   []string{enumName + toCamelCase(currVal.Value)},
			Type:    expression.Type(enumName),
			Values:  []astgen.ASTExpr{expression.StringVal(currVal.Value)},
		})
	}
	valsDecl := &decl.Const{Values: vals}

	valuesFuncDecl := enumValuesAST(enumDefinition)
	isUnknownDecl := enumIsUnknownAST(enumDefinition)
	unmarshalDecl := enumUnmarshalTextAST(enumDefinition, info)

	return []astgen.ASTDecl{typeDef, valsDecl, valuesFuncDecl, isUnknownDecl, unmarshalDecl}
}

func astForEnumPattern(info types.PkgInfo) astgen.ASTDecl {
	info.AddImports("regexp")
	matchString := expression.NewCallFunction("regexp", "MustCompile", expression.StringVal(enumValuePattern))
	return &decl.Var{
		Name:  enumPatternVarName,
		Value: matchString,
	}
}

func enumUnmarshalTextAST(e spec.EnumDefinition, info types.PkgInfo) astgen.ASTDecl {
	mapStringInterface := expression.Type(types.NewMapType(types.String, types.Any).GoType(info))
	toCamelCase := varcaser.Caser{From: varcaser.ScreamingSnakeCase, To: varcaser.UpperCamelCase}.String

	info.AddImports("strings")
	info.AddImports("github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors")
	info.AddImports("github.com/palantir/witchcraft-go-error")
	info.AddImports("github.com/palantir/witchcraft-go-params")

	switchStmt := &statement.Switch{
		Init:       statement.NewAssignment(expression.VariableVal(enumUpperVarName), token.DEFINE, expression.NewCallFunction("strings", "ToUpper", expression.NewCallExpression(expression.StringType, expression.VariableVal(dataVarName)))),
		Expression: expression.VariableVal(enumUpperVarName),
		Cases: []statement.CaseClause{
			// default case
			{
				Body: []astgen.ASTStmt{
					&statement.If{
						Cond: expression.NewUnary(token.NOT,
							expression.NewCallExpression(
								expression.NewSelector(expression.VariableVal(enumPatternVarName), "MatchString"),
								expression.VariableVal(enumUpperVarName),
							),
						),
						Body: []astgen.ASTStmt{
							statement.NewReturn(
								expression.NewCallFunction("werror", "Convert",
									expression.NewCallFunction("errors", "NewInvalidArgument",
										expression.NewCallFunction("wparams", "NewSafeAndUnsafeParamStorer",
											expression.NewCompositeLit(mapStringInterface,
												expression.NewKeyValue(`"enumType"`, expression.StringVal(e.TypeName.Name)),
												expression.NewKeyValue(`"message"`, expression.StringVal("enum value must match pattern "+enumValuePattern)),
											),
											expression.NewCompositeLit(mapStringInterface,
												expression.NewKeyValue(`"enumValue"`, expression.NewCallExpression(expression.StringType, expression.VariableVal(dataVarName))),
											),
										),
									),
								),
							),
						},
					},
					statement.NewAssignment(
						expression.NewUnary(token.MUL, expression.VariableVal(enumReceiverName)),
						token.ASSIGN,
						expression.NewCallExpression(expression.Type(e.TypeName.Name), expression.VariableVal(enumUpperVarName)),
					),
				},
			},
		},
	}
	for _, currVal := range e.Values {
		switchStmt.Cases = append(switchStmt.Cases, *statement.NewCaseClause(
			expression.StringVal(currVal.Value),
			statement.NewAssignment(
				expression.NewUnary(token.MUL, expression.VariableVal(enumReceiverName)),
				token.ASSIGN,
				expression.VariableVal(e.TypeName.Name+toCamelCase(currVal.Value)),
			),
		))
	}
	return newUnmarshalTextMethod(enumReceiverName, transforms.Export(e.TypeName.Name), switchStmt, statement.NewReturn(expression.Nil))
}

func enumIsUnknownAST(e spec.EnumDefinition) astgen.ASTDecl {
	typeName := transforms.Export(e.TypeName.Name)
	return &decl.Method{
		ReceiverName: enumReceiverName,
		ReceiverType: expression.Type(transforms.Export(e.TypeName.Name)),
		Function: decl.Function{
			Comment: fmt.Sprintf("IsUnknown returns false for all known variants of %s and true otherwise.", typeName),
			Name:    "IsUnknown",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.BoolType},
			},
			Body: []astgen.ASTStmt{
				&statement.Switch{
					Expression: expression.VariableVal(enumReceiverName),
					Cases: []statement.CaseClause{
						{
							Exprs: enumValuesToExprs(e),
							Body:  []astgen.ASTStmt{statement.NewReturn(expression.VariableVal("false"))},
						},
					},
				},
				statement.NewReturn(expression.VariableVal("true")),
			},
		},
	}
}

func enumValuesAST(e spec.EnumDefinition) astgen.ASTDecl {
	typeName := transforms.Export(e.TypeName.Name)
	funcName := typeName + "_Values"
	sliceType := expression.Type("[]" + transforms.Export(e.TypeName.Name))
	return &decl.Function{
		Comment: fmt.Sprintf("%s returns all known variants of %s.", funcName, typeName),
		Name:    funcName,
		FuncType: expression.FuncType{
			ReturnTypes: []expression.Type{sliceType},
		},
		Body: []astgen.ASTStmt{
			statement.NewReturn(expression.NewCompositeLit(sliceType, enumValuesToExprs(e)...)),
		},
	}
}

func enumValuesToExprs(e spec.EnumDefinition) []astgen.ASTExpr {
	toCamelCase := varcaser.Caser{From: varcaser.ScreamingSnakeCase, To: varcaser.UpperCamelCase}.String
	values := make([]astgen.ASTExpr, 0, len(e.Values))
	for _, currVal := range e.Values {
		values = append(values, expression.VariableVal(e.TypeName.Name+toCamelCase(currVal.Value)))
	}
	return values
}
