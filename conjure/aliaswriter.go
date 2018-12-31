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

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/transforms"
	"github.com/palantir/conjure-go/conjure/types"
	"github.com/palantir/conjure-go/conjure/visitors"
)

const (
	aliasReceiverName = "a"
)

func astForAlias(ctx types.TypeContext, aliasDefinition spec.AliasDefinition) ([]astgen.ASTDecl, StringSet, error) {
	conjureTypeProvider, err := visitors.NewConjureTypeProvider(aliasDefinition.Alias)
	if err != nil {
		return nil, nil, err
	}
	aliasTyper, err := conjureTypeProvider.ParseType(ctx)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "alias type %s specifies unrecognized type", aliasDefinition.TypeName.Name)
	}
	aliasGoType := aliasTyper.GoType(ctx)

	imports := NewStringSet(aliasTyper.ImportPaths()...)

	var decls []astgen.ASTDecl

	decls = append(decls, &decl.Alias{
		Name:    aliasDefinition.TypeName.Name,
		Type:    expression.Type(aliasTyper.GoType(ctx)),
		Comment: transforms.Documentation(aliasDefinition.Docs),
	})

	if aliasTyper.ImportPaths() != nil {
		// We are aliasing a non-builtin. Create marshal/unmarshal functions which delegate to the aliased type.
		for _, f := range []aliasSerdeFunc{
			astForAliasJSONMarshal,
			astForAliasJSONUnmarshal,
			astForAliasYAMLMarshal,
			astForAliasYAMLUnmarshal,
		} {
			serdeDecl, currImports, err := f(aliasDefinition, aliasGoType)
			if err != nil {
				return nil, nil, err
			}
			decls = append(decls, serdeDecl)
			imports.AddAll(currImports)
		}
	}

	return decls, imports, nil
}

type aliasSerdeFunc func(aliasDefinition spec.AliasDefinition, aliasGoType string) (astgen.ASTDecl, StringSet, error)

func astForAliasJSONMarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) (astgen.ASTDecl, StringSet, error) {
	return &decl.Method{
		Function: decl.Function{
			Name: "MarshalJSON",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					expression.Type("[]byte"),
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCallFunction(
						"json",
						"Marshal",
						&expression.CallExpression{
							Function: expression.VariableVal(aliasGoType),
							Args: []astgen.ASTExpr{
								expression.VariableVal(aliasReceiverName),
							},
						},
					),
				),
			},
		},
		ReceiverName: aliasReceiverName,
		ReceiverType: expression.Type(aliasDefinition.TypeName.Name),
	}, NewStringSet("encoding/json"), nil
}

func astForAliasJSONUnmarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) (astgen.ASTDecl, StringSet, error) {
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)
	return &decl.Method{
		Function: decl.Function{
			Name: "UnmarshalJSON",
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam("data", expression.Type("[]byte")),
				},
				ReturnTypes: []expression.Type{
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{
				// var rawAlias alias.RawType
				statement.NewDecl(decl.NewVar(rawVarName, expression.Type(aliasGoType))),
				ifErrNotNilReturnErrStatement("err",
					statement.NewAssignment(
						expression.VariableVal("err"),
						token.DEFINE,
						expression.NewCallFunction(
							"json",
							"Unmarshal",
							expression.VariableVal("data"), expression.NewUnary(token.AND, expression.VariableVal(rawVarName)),
						),
					),
				),
				statement.NewAssignment(
					expression.NewStar(expression.VariableVal(aliasReceiverName)),
					token.ASSIGN,
					&expression.CallExpression{
						Function: expression.VariableVal(aliasDefinition.TypeName.Name),
						Args: []astgen.ASTExpr{
							expression.VariableVal(rawVarName),
						},
					},
				),
				statement.NewReturn(expression.Nil),
			},
		},
		ReceiverName: aliasReceiverName,
		ReceiverType: expression.Type(aliasDefinition.TypeName.Name).Pointer(),
	}, NewStringSet("encoding/json"), nil
}

func astForAliasYAMLMarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) (astgen.ASTDecl, StringSet, error) {
	return &decl.Method{
		Function: decl.Function{
			Name: "MarshalYAML",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					expression.Type("interface{}"),
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{statement.NewReturn(&expression.CallExpression{
				Function: expression.VariableVal(aliasGoType),
				Args:     []astgen.ASTExpr{expression.VariableVal(aliasReceiverName)},
			}, expression.Nil)},
		},
		ReceiverName: aliasReceiverName,
		ReceiverType: expression.Type(aliasDefinition.TypeName.Name),
	}, nil, nil
}

func astForAliasYAMLUnmarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) (astgen.ASTDecl, StringSet, error) {
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)
	return &decl.Method{
		Function: decl.Function{
			Name: "UnmarshalYAML",
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam("unmarshal", expression.Type("func(interface{}) error")),
				},
				ReturnTypes: []expression.Type{
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{
				// var rawAlias alias.RawType
				statement.NewDecl(decl.NewVar(rawVarName, expression.Type(aliasGoType))),
				// unmarshal(rawAlias)
				ifErrNotNilReturnErrStatement("err",
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
				),
				// *a = AliasType(rawAlias)
				statement.NewAssignment(
					expression.NewStar(expression.VariableVal(aliasReceiverName)),
					token.ASSIGN,
					&expression.CallExpression{
						Function: expression.VariableVal(aliasDefinition.TypeName.Name),
						Args: []astgen.ASTExpr{
							expression.VariableVal(rawVarName),
						},
					},
				),
				statement.NewReturn(expression.Nil),
			},
		},
		ReceiverName: aliasReceiverName,
		ReceiverType: expression.Type(aliasDefinition.TypeName.Name).Pointer(),
	}, nil, nil
}
