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

const (
	aliasReceiverName   = "a"
	aliasValueFieldName = "Value"
)

var (
	aliasOptionalValueSelector = expression.NewSelector(expression.VariableVal(aliasReceiverName), aliasValueFieldName)
)

func astForAlias(aliasDefinition spec.AliasDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	aliasTypeProvider, err := visitors.NewConjureTypeProvider(aliasDefinition.Alias)
	if err != nil {
		return nil, err
	}
	aliasTyper, err := aliasTypeProvider.ParseType(info)
	if err != nil {
		return nil, errors.Wrapf(err, "alias type %s specifies unrecognized type", aliasDefinition.TypeName.Name)
	}
	info.AddImports(aliasTyper.ImportPaths()...)
	aliasGoType := aliasTyper.GoType(info)

	isOptional := aliasTypeProvider.IsSpecificType(visitors.IsOptional)
	isString := aliasTypeProvider.IsSpecificType(visitors.IsString)
	isText := aliasTypeProvider.IsSpecificType(visitors.IsText)

	var decls []astgen.ASTDecl
	if isOptional {
		decls = append(decls, &decl.Struct{
			Name:    aliasDefinition.TypeName.Name,
			Comment: transforms.Documentation(aliasDefinition.Docs),
			StructType: *expression.NewStructType(&expression.StructField{
				Name: aliasValueFieldName,
				Type: expression.Type(aliasGoType),
			}),
		})
	} else {
		decls = append(decls, &decl.Alias{
			Name:    aliasDefinition.TypeName.Name,
			Type:    expression.Type(aliasGoType),
			Comment: transforms.Documentation(aliasDefinition.Docs),
		})
	}

	// Attach encoding methods
	switch {
	case isOptional:
		// Optionals have special method ASTs.
		valueInit, err := aliasTypeProvider.CollectionInitializationIfNeeded(info)
		if err != nil {
			return nil, err
		}
		if valueInit == nil {
			valueInit = &expression.CallExpression{
				Function: expression.VariableVal("new"),
				Args:     []astgen.ASTExpr{expression.Type(strings.TrimPrefix(aliasGoType, "*"))},
			}
		}
		switch {
		case isString:
			decls = append(decls, astForOptionalStringAliasTextMarshal(aliasDefinition))
			decls = append(decls, astForOptionalStringAliasTextUnmarshal(aliasDefinition))
		case isText:
			decls = append(decls, astForOptionalAliasTextMarshal(aliasDefinition))
			decls = append(decls, astForOptionalAliasTextUnmarshal(aliasDefinition, valueInit))
		default:
			decls = append(decls, astForOptionalAliasJSONMarshal(aliasDefinition, info))
			decls = append(decls, astForOptionalAliasJSONUnmarshal(aliasDefinition, valueInit, info))
			decls = append(decls, astForOptionalAliasYAMLMarshal(aliasDefinition))
			decls = append(decls, astForOptionalAliasYAMLUnmarshal(aliasDefinition, valueInit))
		}

	case len(aliasTyper.ImportPaths()) == 0:
		// Plain builtins do not need encoding methods; do nothing.
	case isText:
		// If we have gotten here, we have a non-go-builtin text type that implements MarshalText/UnmarshalText.
		decls = append(decls, astForAliasTextMarshal(aliasDefinition, aliasGoType))
		decls = append(decls, astForAliasTextUnmarshal(aliasDefinition, aliasGoType))
	default:
		// By default, we delegate json/yaml encoding to the aliased type.
		decls = append(decls, astForAliasJSONMarshal(aliasDefinition, aliasGoType, info))
		decls = append(decls, astForAliasJSONUnmarshal(aliasDefinition, aliasGoType, info))
		decls = append(decls, astForAliasYAMLMarshal(aliasDefinition, aliasGoType))
		decls = append(decls, astForAliasYAMLUnmarshal(aliasDefinition, aliasGoType))
	}
	return decls, nil
}

// astForAliasTextMarshal creates the MarshalText method that delegates to the aliased type.
//
//    func (a DateAlias) MarshalText() ([]byte, error) {
//        return datetime.DateTime(a).MarshalText()
//    }
func astForAliasTextMarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	return newMarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name, statement.NewReturn(
		&expression.CallExpression{
			Function: expression.NewSelector(
				&expression.CallExpression{
					Function: expression.Type(aliasGoType),
					Args: []astgen.ASTExpr{
						expression.VariableVal(aliasReceiverName),
					},
				},
				"MarshalText"),
		},
	))
}

// astForOptionalAliasTextMarshal creates the MarshalText method that delegates to the aliased type.
//
//    func (a OptionalDateAlias) MarshalText() ([]byte, error) {
//        if a.Value == nil {
//            return nil, nil
//        }
//        return a.Value.MarshalText()
//    }
func astForOptionalAliasTextMarshal(aliasDefinition spec.AliasDefinition) astgen.ASTDecl {
	return newMarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		stmtIfNilAliasValueReturnNilNil,
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.NewSelector(aliasOptionalValueSelector, "MarshalText"),
			},
		),
	)
}

// astForOptionalStringAliasTextMarshal creates the MarshalText method that delegates to the aliased type.
//
//    func (a OptionalStringAlias) MarshalText() ([]byte, error) {
//        if a.Value == nil {
//            return nil, nil
//        }
//        return []byte(*a.Value), nil
//    }
func astForOptionalStringAliasTextMarshal(aliasDefinition spec.AliasDefinition) astgen.ASTDecl {
	return newMarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		stmtIfNilAliasValueReturnNilNil,
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.Type("[]byte"),
				Args:     []astgen.ASTExpr{expression.NewStar(aliasOptionalValueSelector)},
			},
			expression.Nil,
		),
	)
}

// astForAliasTextUnmarshal creates the UnmarshalText method that delegates to the aliased type.
//
//    func (a *DateAlias) UnmarshalText(data []byte) error {
//        var rawDateAlias datetime.DateTime
//        if err := rawDateAlias.UnmarshalText(data); err != nil {
//            return err
//        }
//        *d = DateAlias(rawDateAlias)
//        return nil
//    }
func astForAliasTextUnmarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
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

// astForOptionalAliasTextUnmarshal creates the UnmarshalText method that delegates to the aliased type.
//
//    func (a *OptionalDateAlias) UnmarshalText(data []byte) error {
//        if a.Value == nil {
//            a.Value = new(datetime.DateTime)
//        }
//        return a.Value.UnmarshalText(data)
//    }
func astForOptionalAliasTextUnmarshal(aliasDefinition spec.AliasDefinition, aliasValueInit astgen.ASTExpr) astgen.ASTDecl {
	return newUnmarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		astForAliasInitializeOptional(aliasValueInit),
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.NewSelector(aliasOptionalValueSelector, "UnmarshalText"),
				Args:     []astgen.ASTExpr{expression.VariableVal(dataVarName)},
			},
		),
	)
}

// astForOptionalStringAliasTextUnmarshal creates the UnmarshalText method that delegates to the aliased type.
//
//    func (a *OptionalStringAlias) UnmarshalText(data []byte) error {
//        rawOptionalStringAlias := string(data)
//        a.Value = &rawOptionalStringAlias
//        return nil
//    }
func astForOptionalStringAliasTextUnmarshal(aliasDefinition spec.AliasDefinition) astgen.ASTDecl {
	rawVar := expression.VariableVal(fmt.Sprint("raw", aliasDefinition.TypeName.Name))
	return newUnmarshalTextMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		statement.NewAssignment(
			rawVar,
			token.DEFINE,
			&expression.CallExpression{
				Function: expression.StringType,
				Args:     []astgen.ASTExpr{expression.VariableVal(dataVarName)},
			},
		),
		statement.NewAssignment(
			aliasOptionalValueSelector,
			token.ASSIGN,
			expression.NewUnary(token.AND, rawVar),
		),
		statement.NewReturn(expression.Nil),
	)
}

// astForAliasJSONMarshal creates the MarshalJSON method that delegates to the aliased type.
//
//    func (a ObjectAlias) MarshalJSON() ([]byte, error) {
//        return safejson.Marshal(Object(a))
//    }
func astForAliasJSONMarshal(aliasDefinition spec.AliasDefinition, aliasGoType string, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.SafeJSONMarshal.ImportPaths()...)
	return newMarshalJSONMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.Type(types.SafeJSONMarshal.GoType(info)),
				Args: []astgen.ASTExpr{
					&expression.CallExpression{
						Function: expression.Type(aliasGoType),
						Args: []astgen.ASTExpr{
							expression.VariableVal(aliasReceiverName),
						},
					},
				},
			},
		),
	)
}

// astForOptionalAliasJSONMarshal creates the MarshalJSON method that delegates to the aliased type.
//
//    func (a OptionalObjectAlias) MarshalJSON() ([]byte, error) {
//        if a.Value == nil {
//            return nil, nil
//        }
//        return safejson.Marshal(a.Value)
//    }
func astForOptionalAliasJSONMarshal(aliasDefinition spec.AliasDefinition, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.SafeJSONMarshal.ImportPaths()...)
	return newMarshalJSONMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		stmtIfNilAliasValueReturnNilNil,
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.Type(types.SafeJSONMarshal.GoType(info)),
				Args: []astgen.ASTExpr{
					aliasOptionalValueSelector,
				},
			},
		),
	)
}

// astForAliasJSONUnmarshal creates the UnmarshalJSON method that delegates to the aliased type.
//
//    func (a *ObjectAlias) UnmarshalJSON(data []byte) error {
//        var rawObjectAlias Object
//        if err := safejson.Unmarshal(data, &rawObjectAlias); err != nil {
//            return err
//        }
//        *d = ObjectAlias(rawObjectAlias)
//        return nil
//    }
func astForAliasJSONUnmarshal(aliasDefinition spec.AliasDefinition, aliasGoType string, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.SafeJSONUnmarshal.ImportPaths()...)
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)
	return newUnmarshalJSONMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		// var rawObjectAlias Object
		statement.NewDecl(decl.NewVar(rawVarName, expression.Type(aliasGoType))),
		// safejson.Unmarshal(data, &rawObjectAlias)
		ifErrNotNilReturnErrStatement("err",
			statement.NewAssignment(
				expression.VariableVal("err"),
				token.DEFINE,
				&expression.CallExpression{
					Function: expression.Type(types.SafeJSONUnmarshal.GoType(info)),
					Args: []astgen.ASTExpr{
						expression.VariableVal(dataVarName),
						expression.NewUnary(token.AND, expression.VariableVal(rawVarName)),
					},
				},
			),
		),
		// *a = ObjectAlias(rawObjectAlias)
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

// astForAliasJSONUnmarshal creates the UnmarshalJSON method that delegates to the aliased type.
//
//    func (a *OptionalObjectAlias) UnmarshalJSON(data []byte) error {
//        if a.Value == nil {
//            a.Value = new(Object)
//        }
//        return safejson.Unmarshal(data, a.Value)
//    }
func astForOptionalAliasJSONUnmarshal(aliasDefinition spec.AliasDefinition, aliasValueInit astgen.ASTExpr, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.SafeJSONUnmarshal.ImportPaths()...)
	return newUnmarshalJSONMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		astForAliasInitializeOptional(aliasValueInit),
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.Type(types.SafeJSONUnmarshal.GoType(info)),
				Args: []astgen.ASTExpr{
					expression.VariableVal(dataVarName),
					aliasOptionalValueSelector,
				},
			},
		),
	)
}

// astForAliasYAMLMarshal creates the MarshalYAML method that delegates to the aliased type.
//
//    func (a ObjectAlias) MarshalYAML() (interface{}, error) {
//        return Object(a)
//    }
func astForAliasYAMLMarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	return newMarshalYAMLMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		statement.NewReturn(&expression.CallExpression{
			Function: expression.Type(aliasGoType),
			Args:     []astgen.ASTExpr{expression.VariableVal(aliasReceiverName)},
		},
			expression.Nil),
	)
}

// astForOptionalAliasYAMLMarshal creates the MarshalYAML method that delegates to the aliased type.
//
//    func (a OptionalObjectAlias) MarshalYAML() (interface{}, error) {
//        return a.Value, nil
//    }
func astForOptionalAliasYAMLMarshal(aliasDefinition spec.AliasDefinition) astgen.ASTDecl {
	return newMarshalYAMLMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		statement.NewReturn(aliasOptionalValueSelector, expression.Nil),
	)
}

// astForAliasYAMLUnmarshal creates the UnmarshalYAML method that delegates to the aliased type.
//
//    func (a *ObjectAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
//        var rawObjectAlias Object
//        if err := unmarshal(rawObjectAlias); err != nil {
//            return err
//        }
//        *d = ObjectAlias(rawObjectAlias)
//        return nil
//    }
func astForAliasYAMLUnmarshal(aliasDefinition spec.AliasDefinition, aliasGoType string) astgen.ASTDecl {
	rawVarName := fmt.Sprint("raw", aliasDefinition.TypeName.Name)

	return newUnmarshalYAMLMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		// var rawObjectAlias Object
		statement.NewDecl(decl.NewVar(rawVarName, expression.Type(aliasGoType))),
		// unmarshal(rawObjectAlias)
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
		// *a = ObjectAlias(rawObjectAlias)
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

// astForAliasYAMLUnmarshal creates the UnmarshalYAML method that delegates to the aliased type.
//
//    func (a *OptionalObjectAlias) UnmarshalYAML(unmarshal func(interface{}) error) error {
//        if a.Value == nil {
//            a.Value = new(Object)
//        }
//        return unmarshal(a.Value)
//    }
func astForOptionalAliasYAMLUnmarshal(aliasDefinition spec.AliasDefinition, aliasValueInit astgen.ASTExpr) astgen.ASTDecl {
	return newUnmarshalYAMLMethod(aliasReceiverName, aliasDefinition.TypeName.Name,
		astForAliasInitializeOptional(aliasValueInit),
		statement.NewReturn(&expression.CallExpression{
			Function: expression.VariableVal("unmarshal"),
			Args:     []astgen.ASTExpr{aliasOptionalValueSelector},
		}),
	)
}

func astForAliasInitializeOptional(valueInit astgen.ASTExpr) astgen.ASTStmt {
	// if a.Value == nil { a.Value = new(Object) }
	return &statement.If{
		Cond: expression.NewBinary(aliasOptionalValueSelector, token.EQL, expression.Nil),
		Body: []astgen.ASTStmt{
			statement.NewAssignment(aliasOptionalValueSelector, token.ASSIGN, valueInit),
		},
	}
}

// if a.Value == nil { return nil, nil }
var stmtIfNilAliasValueReturnNilNil = &statement.If{
	Cond: expression.NewBinary(aliasOptionalValueSelector, token.EQL, expression.Nil),
	Body: []astgen.ASTStmt{statement.NewReturn(expression.Nil, expression.Nil)},
}
