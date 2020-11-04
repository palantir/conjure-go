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

package visitors

import (
	"errors"
	"fmt"
	"go/token"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"

	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v5/conjure/types"
	"github.com/palantir/conjure-go/v5/conjure/werrorexpressions"
)

// StatementsForHTTPParam returns the AST statements converting an HTTP parameter (path/query/string) to the proper
// type. Supports generating statements for primitives and containers (set/list) of primitives.
func StatementsForHTTPParam(argName spec.ArgumentName, argType spec.Type, stringExpr astgen.ASTExpr, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	v := httpParamVisitor{
		argName:    argName,
		stringExpr: stringExpr,
		info:       info,
	}
	if err := argType.Accept(&v); err != nil {
		return nil, err
	}
	return v.result, nil
}

type httpParamVisitor struct {
	// argName is the argument we are unmarshaling. The result statements will conclude with a variable of the correct type with this arg name.
	argName spec.ArgumentName
	// stringExpr is an expression which results in a string, e.g. req.URL.Query().Get("myParam")
	stringExpr astgen.ASTExpr
	// result will be populated by the visitor
	result []astgen.ASTStmt

	// imports will be added by each visitor to the PkgInfo
	info types.PkgInfo
}

func (v *httpParamVisitor) VisitPrimitive(t spec.PrimitiveType) error {
	var typer types.Typer
	var returnsErr bool
	args := []astgen.ASTExpr{v.stringExpr}
	switch t {
	case spec.PrimitiveTypeString:
		typer = types.String
	case spec.PrimitiveTypeInteger:
		typer = types.ParseInt
		returnsErr = true
	case spec.PrimitiveTypeDouble:
		typer = types.ParseFloat
		returnsErr = true
		args = append(args, expression.IntVal(64))
	case spec.PrimitiveTypeBoolean:
		typer = types.ParseBool
		returnsErr = true
	case spec.PrimitiveTypeBearertoken:
		typer = types.Bearertoken
	case spec.PrimitiveTypeDatetime:
		typer = types.ParseDateTime
		returnsErr = true
	case spec.PrimitiveTypeRid:
		typer = types.ParseRID
		returnsErr = true
	case spec.PrimitiveTypeSafelong:
		typer = types.ParseSafeLong
		returnsErr = true
	case spec.PrimitiveTypeUuid:
		typer = types.ParseUUID
		returnsErr = true
	case spec.PrimitiveTypeAny:
		typer = types.Any
	case spec.PrimitiveTypeBinary:
		typer = types.BinaryType
	default:
		return errors.New("Unsupported primitive type " + string(t))
	}

	var rhs astgen.ASTExpr
	switch typer {
	case types.String, types.Any:
		rhs = v.stringExpr
	default:
		rhs = &expression.CallExpression{
			Function: expression.VariableVal(typer.GoType(v.info)),
			Args:     args,
		}
	}

	if !returnsErr {
		v.result = append(v.result, &statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal(v.argName)},
			Tok: token.DEFINE,
			RHS: rhs,
		})
	} else {
		errVar := expression.VariableVal("err")
		v.result = append(v.result, &statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal(v.argName), errVar},
			Tok: token.DEFINE,
			RHS: rhs,
		}, &statement.If{
			Cond: expression.NewBinary(errVar, token.NEQ, expression.Nil),
			Body: []astgen.ASTStmt{statement.NewReturn(errVar)}, // TODO should be a 400 Bad Request
		})
	}

	v.info.AddImports(typer.ImportPaths()...)
	return nil
}

func (v *httpParamVisitor) VisitOptional(t spec.OptionalType) error {
	typer, err := newOptionalVisitor(t).ParseType(v.info)
	if err != nil {
		return err
	}

	v.info.AddImports(typer.ImportPaths()...)
	v.result = append(v.result, &statement.Decl{
		Decl: &decl.Var{
			Name: string(v.argName),
			Type: expression.Type(typer.GoType(v.info)),
		},
	})

	ifStmt := &statement.If{
		Init: &statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.VariableVal(v.argName + "Str"),
			},
			Tok: token.DEFINE,
			RHS: v.stringExpr,
		},
		Cond: &expression.Binary{
			LHS: expression.VariableVal(v.argName + "Str"),
			Op:  token.NEQ,
			RHS: expression.StringVal(""),
		},
	}

	internalVisitor := httpParamVisitor{
		argName:    v.argName + "Internal",
		stringExpr: expression.VariableVal(v.argName + "Str"),
		info:       v.info,
	}
	if err := t.ItemType.Accept(&internalVisitor); err != nil {
		return err
	}

	ifStmt.Body = append(ifStmt.Body, internalVisitor.result...)
	// assign arg variable to address of internal result
	// argName = &argNameInternal
	ifStmt.Body = append(ifStmt.Body, &statement.Assignment{
		LHS: []astgen.ASTExpr{expression.VariableVal(v.argName)},
		Tok: token.ASSIGN,
		RHS: expression.NewUnary(token.AND, expression.VariableVal(internalVisitor.argName)),
	})
	v.result = append(v.result, ifStmt)
	return nil
}

func (v *httpParamVisitor) VisitExternal(t spec.ExternalReference) error {
	// If the fallback is something we recognize, use that and cast to the external type
	fallbackVisitor := httpParamVisitor{
		argName:    v.argName + "Internal",
		stringExpr: v.stringExpr,
		info:       v.info,
	}
	if err := t.Fallback.Accept(&fallbackVisitor); err != nil {
		return err
	}

	typer, err := newExternalVisitor(t).ParseType(v.info)
	if err != nil {
		return err
	}
	v.info.AddImports(typer.ImportPaths()...)
	v.result = append(v.result, fallbackVisitor.result...)
	v.result = append(v.result, &statement.Assignment{
		LHS: []astgen.ASTExpr{expression.VariableVal(v.argName)},
		Tok: token.DEFINE,
		RHS: &expression.CallExpression{
			Function: expression.VariableVal(typer.GoType(v.info)),
			Args:     []astgen.ASTExpr{expression.VariableVal(fallbackVisitor.argName)},
		},
	})
	return nil
}

func (v *httpParamVisitor) VisitReference(t spec.TypeName) error {
	typer, err := newReferenceVisitor(t).ParseType(v.info)
	if err != nil {
		return err
	}
	v.info.AddImports(typer.ImportPaths()...)
	v.result = append(v.result, &statement.Decl{
		Decl: &decl.Var{
			Name: string(v.argName),
			Type: expression.Type(typer.GoType(v.info)),
		},
	})
	v.info.AddImports("strconv")
	v.info.AddImports(types.SafeJSONUnmarshal.ImportPaths()...)
	v.result = append(v.result, &statement.Assignment{
		LHS: []astgen.ASTExpr{expression.VariableVal(v.argName + "Quote")},
		Tok: token.DEFINE,
		RHS: expression.NewCallFunction("strconv", "Quote", v.stringExpr),
	})
	v.result = append(v.result, &statement.If{
		Init: &statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal("err")},
			Tok: token.DEFINE,
			RHS: expression.NewCallExpression(expression.Type(types.SafeJSONUnmarshal.GoType(v.info)),
				&expression.CallExpression{
					Function: expression.VariableVal("[]byte"),
					Args:     []astgen.ASTExpr{expression.VariableVal(v.argName + "Quote")},
				},
				expression.NewUnary(token.AND, expression.VariableVal(v.argName)),
			),
		},
		Cond: &expression.Binary{
			LHS: expression.VariableVal("err"),
			Op:  token.NEQ,
			RHS: expression.Nil,
		},
		Body: []astgen.ASTStmt{statement.NewReturn(werrorexpressions.CreateWrapWErrorExpression(
			expression.VariableVal("err"), //TODO(bmoylan): This should be a conjure 400 error, right now it will return 500
			"failed to unmarshal argument",
			map[string]string{
				"argName": string(v.argName),
				"argType": typer.GoType(v.info),
			})),
		},
	})

	return nil
}

func (v *httpParamVisitor) VisitList(t spec.ListType) error {
	return v.visitCollectionType(t.ItemType)
}

func (v *httpParamVisitor) VisitSet(t spec.SetType) error {
	return v.visitCollectionType(t.ItemType)
}

func (v *httpParamVisitor) VisitMap(t spec.MapType) error {
	return errors.New("can not assign string expression to map type")
}

func (v *httpParamVisitor) VisitUnknown(typeName string) error {
	return fmt.Errorf("can not create httpParamVisitor for unknown type %s", typeName)
}

func (v *httpParamVisitor) visitCollectionType(itemType spec.Type) error {
	provider, err := NewConjureTypeProvider(itemType)
	if err != nil {
		return err
	}

	if provider.IsSpecificType(IsString) {
		// get query arguments
		v.result = append(v.result, &statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.VariableVal(v.argName),
			},
			Tok: token.DEFINE,
			RHS: v.stringExpr,
		})
		return nil
	}

	parsedType, err := provider.ParseType(v.info)
	if err != nil {
		return err
	}
	goTypeOutput := parsedType.GoType(v.info)
	v.result = append(v.result, statement.NewDecl(decl.NewVar(string(v.argName), expression.Type("[]"+goTypeOutput))))

	astStmts, err := StatementsForHTTPParam("convertedVal", itemType, expression.VariableVal("v"), v.info)
	if err != nil {
		return err
	}

	v.result = append(v.result, &statement.Range{
		Key:   expression.VariableVal("_"),
		Value: expression.VariableVal("v"),
		Tok:   token.DEFINE,
		Expr:  v.stringExpr,
		Body: append(astStmts,
			statement.NewAssignment(
				expression.VariableVal(v.argName),
				token.ASSIGN,
				expression.NewCallExpression(
					expression.AppendBuiltIn,
					expression.VariableVal(v.argName),
					expression.VariableVal("convertedVal"),
				),
			),
		),
	})
	return nil
}
