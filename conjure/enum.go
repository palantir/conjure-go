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
	"go/ast"
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

	unmarshalDecl := enumUnmarshalTextAST(enumDefinition, info)

	return []astgen.ASTDecl{typeDef, valsDecl, unmarshalDecl}
}

func astForEnumPattern(info types.PkgInfo) astgen.ASTDecl {
	info.AddImports("regexp")
	matchString := expression.NewCallFunction("regexp", "MustCompile", expression.StringVal(enumValuePattern))
	return &varDecl{
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

// Goastwriter does not allow values in var declarations, so implement it here.
// This should be contributed back to goastwriter.
type varDecl struct {
	Name  string
	Type  expression.Type
	Value astgen.ASTExpr
}

func (v *varDecl) ASTDecl() ast.Decl {
	valueSpec := &ast.ValueSpec{
		Names: []*ast.Ident{ast.NewIdent(v.Name)},
	}
	if v.Type != "" {
		valueSpec.Type = v.Type.ToIdent()
	}
	if v.Value != nil {
		valueSpec.Values = []ast.Expr{v.Value.ASTExpr()}
	}
	return &ast.GenDecl{
		Tok:   token.VAR,
		Specs: []ast.Spec{valueSpec},
	}
}
