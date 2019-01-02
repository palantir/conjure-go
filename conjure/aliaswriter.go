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

func astForAlias(ctx types.TypeContext, aliasDefinition spec.AliasDefinition) ([]astgen.ASTDecl, error) {
	aliasTypeProvider, err := visitors.NewConjureTypeProvider(aliasDefinition.Alias)
	if err != nil {
		return nil, err
	}
	aliasTyper, err := aliasTypeProvider.ParseType(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "alias type %s specifies unrecognized type", aliasDefinition.TypeName.Name)
	}
	ctx.AddImports(aliasTyper.ImportPaths()...)
	aliasGoType := aliasTyper.GoType(ctx)

	decls := []astgen.ASTDecl{&decl.Alias{
		Name:    aliasDefinition.TypeName.Name,
		Type:    expression.Type(aliasGoType),
		Comment: transforms.Documentation(aliasDefinition.Docs),
	}}

	// Attach encoding methods
	switch {
	case len(aliasTyper.ImportPaths()) == 0:
		// We are aliasing a builtin, this does not require encoding methods.
	case aliasTypeProvider.IsSpecificType(visitors.IsOptional):
	// TODO(bmoylan) Implement encoding for aliased optionals.
	// Change optional aliases to struct types instead of aliasing a pointer because pointer types can not have methods.
	// For now, just do nothing.
	case aliasTypeProvider.IsSpecificType(visitors.IsBinary):
	// TODO(bmoylan) Remove this case when https://github.com/palantir/conjure-go/pull/17 (binary.Binary type) merges.
	// For now, just do nothing.
	case aliasTypeProvider.IsSpecificType(visitors.IsText):
		// If we have gotten here, we have a non-go-builtin text type that implements MarshalText/UnmarshalText.
		decls = append(decls, astForAliasTextMarshal(ctx, aliasDefinition, aliasGoType))
		decls = append(decls, astForAliasTextUnmarshal(ctx, aliasDefinition, aliasGoType))
	default:
		decls = append(decls, astForAliasJSONMarshal(ctx, aliasDefinition, aliasGoType))
		decls = append(decls, astForAliasJSONUnmarshal(ctx, aliasDefinition, aliasGoType))
		decls = append(decls, astForAliasYAMLMarshal(aliasDefinition, aliasGoType))
		decls = append(decls, astForAliasYAMLUnmarshal(aliasDefinition, aliasGoType))
	}
	return decls, nil
}

// astForAliasTextMarshal creates the MarshalText method that delegates to the aliased type.
//
//    func (a DateAlias) MarshalText() ([]byte, error) {
//	      return datetime.DateTime(a).MarshalText()
//    }
func astForAliasTextMarshal(ctx types.TypeContext, aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	return newMarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name, statement.NewReturn(&expression.CallExpression{
		Function: expression.NewSelector(&expression.CallExpression{
			Function: expression.Type(aliasGoType),
			Args:     []astgen.ASTExpr{expression.VariableVal(aliasReceiverName)},
		}, "MarshalText"),
	}))
}

// astForAliasTextUnmarshal creates the UnmarshalText method that delegates to the aliased type.
//
//    func (a *DateAlias) UnmarshalText(data []byte) error {
//        var rawDateAlias datetime.DateTime
//	      if err := rawDateAlias.UnmarshalText(data); err != nil {
//            return err
//	      }
//	      *d = DateAlias(rawDateAlias)
//	      return nil
//    }
func astForAliasTextUnmarshal(ctx types.TypeContext, aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)
	return newUnmarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		// var rawAliasType AliasType
		statement.NewDecl(decl.NewVar(rawVarName, expression.Type(aliasGoType))),
		// rawAliasType.UnmarshalText(data)
		ifErrNotNilReturnErrStatement("err",
			statement.NewAssignment(
				expression.VariableVal("err"),
				token.DEFINE,
				expression.NewCallFunction(
					rawVarName,
					"UnmarshalText",
					expression.VariableVal(dataVarName),
				),
			),
		),
		// *a = Type(rawAliasType)
		statement.NewAssignment(
			expression.NewStar(expression.VariableVal(aliasReceiverName)),
			token.ASSIGN,
			&expression.CallExpression{
				Function: expression.Type(aliasDefinition.TypeName.Name),
				Args: []astgen.ASTExpr{
					expression.VariableVal(rawVarName),
				},
			},
		),
		// return nil
		statement.NewReturn(expression.Nil),
	)
}

func astForAliasJSONMarshal(ctx types.TypeContext, aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	ctx.AddImports(types.CodecJSON.ImportPaths()...)
	return newMarshalJSONMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		statement.NewReturn(
			expression.NewCallFunction(
				types.CodecJSON.GoType(ctx),
				"Marshal",
				&expression.CallExpression{
					Function: expression.Type(aliasGoType),
					Args: []astgen.ASTExpr{
						expression.VariableVal(aliasReceiverName),
					},
				},
			),
		),
	)
}

func astForAliasJSONUnmarshal(ctx types.TypeContext, aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	ctx.AddImports(types.CodecJSON.ImportPaths()...)
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)
	return newUnmarshalJSONMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		// var rawAliasType AliasType
		statement.NewDecl(decl.NewVar(rawVarName, expression.Type(aliasGoType))),
		// codecs.JSON.Unmarshal(data, &rawAliasType)
		ifErrNotNilReturnErrStatement("err",
			statement.NewAssignment(
				expression.VariableVal("err"),
				token.DEFINE,
				expression.NewCallFunction(
					types.CodecJSON.GoType(ctx),
					"Unmarshal",
					expression.VariableVal(dataVarName),
					expression.NewUnary(token.AND, expression.VariableVal(rawVarName)),
				),
			),
		),
		// *a = Type(rawAliasType)
		statement.NewAssignment(
			expression.NewStar(expression.VariableVal(aliasReceiverName)),
			token.ASSIGN,
			&expression.CallExpression{
				Function: expression.Type(aliasDefinition.TypeName.Name),
				Args: []astgen.ASTExpr{
					expression.VariableVal(rawVarName),
				},
			},
		),
		// return nil
		statement.NewReturn(expression.Nil),
	)
}

func astForAliasYAMLMarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	return newMarshalYAMLMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		statement.NewReturn(&expression.CallExpression{
			Function: expression.Type(aliasGoType),
			Args:     []astgen.ASTExpr{expression.VariableVal(aliasReceiverName)},
		},
			expression.Nil),
	)
}

func astForAliasYAMLUnmarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)

	return newUnmarshalYAMLMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
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
	)
}
