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

package jsonencoding

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"github.com/pkg/errors"
)

// publicUnmarshalJSONMethods Creates four methods that delegate to unmarshalGJSON(value gjson.Value, strict bool) which must exist on the type.
//
//	func (o *BinaryMap) UnmarshalJSON(data []byte) error {
//		if !gjson.ValidBytes(data) { return errors.NewInvalidArgument() }
//		return o.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
//	}
//
//	func (o *BinaryMap) UnmarshalJSONStrict(data []byte) error {
//		if !gjson.ValidBytes(data) { return errors.NewInvalidArgument() }
//		return o.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
//	}
func publicUnmarshalJSONMethods(receiverName string, receiverType string) []astgen.ASTDecl {
	ctxDecl := statement.NewAssignment(expression.VariableVal("ctx"), token.DEFINE, expression.NewCallFunction("context", "TODO"))
	return []astgen.ASTDecl{
		&decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(receiverType).Pointer(),
			Function: decl.Function{
				Name: "UnmarshalJSON",
				FuncType: expression.FuncType{
					Params:      expression.FuncParams{expression.NewFuncParam("data", expression.ByteSliceType)},
					ReturnTypes: []expression.Type{expression.ErrorType},
				},
				Body: []astgen.ASTStmt{
					ctxDecl,
					&statement.If{
						Cond: expression.NewUnary(token.NOT, expression.NewCallFunction("gjson", "ValidBytes", expression.VariableVal("data"))),
						Body: []astgen.ASTStmt{
							statement.NewReturn(werrorNew("invalid json")),
						},
					},
					// return o.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
					statement.NewReturn(expression.NewCallFunction(receiverName, "unmarshalGJSON",
						expression.VariableVal("ctx"),
						expression.NewCallFunction("gjson", "ParseBytes", expression.VariableVal("data")),
						expression.VariableVal("false"))),
				},
			},
		},
		&decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(receiverType).Pointer(),
			Function: decl.Function{
				Name: "UnmarshalJSONStrict",
				FuncType: expression.FuncType{
					Params:      expression.FuncParams{expression.NewFuncParam("data", expression.ByteSliceType)},
					ReturnTypes: []expression.Type{expression.ErrorType},
				},
				Body: []astgen.ASTStmt{
					ctxDecl,
					&statement.If{
						Cond: expression.NewUnary(token.NOT, expression.NewCallFunction("gjson", "ValidBytes", expression.VariableVal("data"))),
						Body: []astgen.ASTStmt{
							statement.NewReturn(werrorNew("invalid json")),
						},
					},
					// return o.unmarshalGJSON(ctx, gjson.ParseBytes(data), true)
					statement.NewReturn(expression.NewCallFunction(receiverName, "unmarshalGJSON",
						expression.VariableVal("ctx"),
						expression.NewCallFunction("gjson", "ParseBytes", expression.VariableVal("data")),
						expression.VariableVal("true"))),
				},
			},
		},
	}
}

func unmarshalGJSONMethod(receiverName string, receiverType string, body []astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType).Pointer(),
		Function: decl.Function{
			Name: "unmarshalGJSON",
			FuncType: expression.FuncType{
				Params: expression.FuncParams{
					expression.NewFuncParam("ctx", expression.Type("context.Context")),
					expression.NewFuncParam("value", expression.Type("gjson.Result")),
					expression.NewFuncParam("strict", expression.BoolType),
				},
				ReturnTypes: []expression.Type{expression.ErrorType},
			},
			Body: body,
		},
	}
}

type literalUnmarshalJSONStmtVisitor struct {
	receiverVar  expression.VariableVal
	receiverType expression.Type
	info         types.PkgInfo

	out []astgen.ASTStmt
}

func (v literalUnmarshalJSONStmtVisitor) VisitAlias(def spec.AliasDefinition) error {
	typeProvider, err := visitors.NewConjureTypeProvider(def.Alias)
	if err != nil {
		return err
	}
	typer, err := typeProvider.ParseType(v.info)
	if err != nil {
		return err
	}
	aliasGoType := expression.Type(typer.GoType(v.info))
	isOptional := typeProvider.IsSpecificType(visitors.IsOptional)

	v.out = append(v.out, statement.NewDecl(decl.NewVar("err", expression.ErrorType)))

	var selector astgen.ASTExpr
	if isOptional {
		selector = expression.NewSelector(v.receiverVar, "Value")
	} else {
		selector = expression.VariableVal("objectValue")
		v.out = append(v.out, statement.NewDecl(decl.NewVar("objectValue", aliasGoType)))

		if collectionExpression, err := typeProvider.CollectionInitializationIfNeeded(v.info); err != nil {
			return err
		} else if collectionExpression != nil {
			v.out = append(v.out, statement.NewAssignment(selector, token.ASSIGN, collectionExpression))
		}
	}

	//objectVar := expression.VariableVal("objectValue")
	visitor := &gjsonUnmarshalValueVisitor{
		info:            v.info,
		selector:        selector,
		valueVar:        "value",
		returnErrStmt:   statement.NewReturn(expression.VariableVal("err")),
		fieldDescriptor: fmt.Sprintf("type %s ", def.TypeName.Name),
	}
	if err := def.Alias.Accept(visitor); err != nil {
		return err
	}
	if visitor.typeCheck != nil {
		v.out = append(v.out, visitor.typeCheck)
	}
	v.out = append(v.out, visitor.stmts...)

	//	*a = RidAlias(objectValue)
	if !isOptional {
		v.out = append(v.out, statement.NewAssignment(
			expression.NewUnary(token.MUL, v.receiverVar),
			token.ASSIGN,
			expression.NewCallExpression(v.receiverType, selector),
		))
	}

	v.out = append(v.out, statement.NewReturn(expression.VariableVal("err")))
	return nil
}

func (v literalUnmarshalJSONStmtVisitor) VisitEnum(def spec.EnumDefinition) error {
	panic("implement me")
}

func (v literalUnmarshalJSONStmtVisitor) VisitObject(def spec.ObjectDefinition) error {
	panic("implement me")
}

func (v literalUnmarshalJSONStmtVisitor) VisitUnion(def spec.UnionDefinition) error {
	panic("implement me")
}

func (v literalUnmarshalJSONStmtVisitor) VisitUnknown(typeName string) error {
	return errors.Errorf("unknown type definition %s", typeName)
}

func visitStructFieldsUnmarshalGJSONMethodBody(receiverName, receiverType string, fields []JSONField, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt

	// if !value.IsObject() { return werror.WrapWithContextParams(ctx, errors.NewInvalidArgument(), "expected json object for TypeName") }
	body = append(body, &statement.If{
		Cond: expression.NewUnary(token.NOT, expression.NewCallFunction("value", "IsObject")),
		Body: []astgen.ASTStmt{
			statement.NewReturn(werrorNew(fmt.Sprintf("type %s expected json type Object", receiverType))),
		},
	})
	var fieldInits []astgen.ASTStmt
	var fieldCases []statement.CaseClause
	var fieldValidates []astgen.ASTStmt
	for _, field := range fields {
		selector := expression.NewSelector(expression.VariableVal(receiverName), field.FieldSelector)
		stmts, err := visitStructFieldsUnmarshalJSONMethodStmts(receiverType, selector, field, info)
		if err != nil {
			return nil, err
		}
		fieldInits = append(fieldInits, stmts.Init...)
		fieldCases = append(fieldCases, statement.CaseClause{
			Exprs: []astgen.ASTExpr{expression.StringVal(field.JSONKey)},
			Body:  stmts.UnmarshalGJSON,
		})
		fieldValidates = append(fieldValidates, stmts.ValidateReqdField...)
	}
	unrecognizedFieldsVar := expression.VariableVal("unrecognizedFields")
	fieldCases = append(fieldCases, statement.CaseClause{
		Exprs: nil, // default case
		Body: []astgen.ASTStmt{&statement.If{
			Cond: expression.VariableVal("strict"),
			Body: []astgen.ASTStmt{
				statement.NewAssignment(unrecognizedFieldsVar, token.ASSIGN,
					expression.NewCallExpression(expression.AppendBuiltIn, unrecognizedFieldsVar, expression.VariableVal("key.Str"))),
			},
		}},
	})

	body = append(body, fieldInits...)

	// var unrecognizedFields []string
	body = append(body, statement.NewDecl(decl.NewVar("unrecognizedFields", "[]string")))
	// var err error
	body = append(body, statement.NewDecl(decl.NewVar("err", expression.ErrorType)))

	// value.ForEach(func(key, value gjson.Result) bool { switch key.Str { ... } return err == nil }
	body = append(body, statement.NewExpression(expression.NewCallFunction("value", "ForEach",
		expression.NewFuncLit(
			expression.FuncType{
				Params: expression.FuncParams{{
					Names: []string{"key", "value"},
					Type:  "gjson.Result",
				}},
				ReturnTypes: []expression.Type{expression.BoolType},
			},
			&statement.Switch{
				Expression: expression.NewSelector(expression.VariableVal("key"), "Str"),
				Cases:      fieldCases,
			},
			statement.NewReturn(expression.NewBinary(expression.VariableVal("err"), token.EQL, expression.Nil)),
		),
	)))
	// if err != nil { return err }
	body = append(body, &statement.If{
		Cond: expression.NewBinary(expression.VariableVal("err"), token.NEQ, expression.Nil),
		Body: []astgen.ASTStmt{statement.NewReturn(expression.VariableVal("err"))},
	})
	if len(fieldValidates) > 0 {
		missingFieldsVar := expression.VariableVal("missingFields")
		body = append(body, statement.NewDecl(decl.NewVar("missingFields", "[]string")))
		body = append(body, fieldValidates...)
		body = append(body, &statement.If{
			Cond: expression.NewBinary(
				expression.NewCallExpression(expression.LenBuiltIn, missingFieldsVar),
				token.GTR,
				expression.IntVal(0),
			),
			Body: []astgen.ASTStmt{statement.NewReturn(
				werrorNew(fmt.Sprintf("type %s missing required json fields", receiverType),
					expression.NewCallFunction("werror", "SafeParam", expression.StringVal("missingFields"), missingFieldsVar),
				),
			)},
		})
	}
	body = append(body, &statement.If{
		Cond: expression.NewBinary(
			expression.VariableVal("strict"),
			token.LAND,
			expression.NewBinary(
				expression.NewCallExpression(expression.LenBuiltIn, unrecognizedFieldsVar),
				token.GTR,
				expression.IntVal(0),
			),
		),
		Body: []astgen.ASTStmt{statement.NewReturn(
			werrorNew(fmt.Sprintf("type %s encountered unrecognized json fields", receiverType),
				// unrecognized user input must stay unsafe
				expression.NewCallFunction("werror", "UnsafeParam", expression.StringVal("unrecognizedFields"), unrecognizedFieldsVar),
			),
		)},
	})
	body = append(body, statement.NewReturn(expression.Nil))

	return body, nil
}

type structFieldsUnmarshalJSONMethodStmts struct {
	Init              []astgen.ASTStmt
	UnmarshalGJSON    []astgen.ASTStmt
	ValidateReqdField []astgen.ASTStmt
}

func visitStructFieldsUnmarshalJSONMethodStmts(receiverType string, selector astgen.ASTExpr, field JSONField, info types.PkgInfo) (structFieldsUnmarshalJSONMethodStmts, error) {
	result := structFieldsUnmarshalJSONMethodStmts{}

	typeProvider, err := visitors.NewConjureTypeProvider(field.Type)
	if err != nil {
		return result, err
	}
	typer, err := typeProvider.ParseType(info)
	if err != nil {
		return result, err
	}
	info.AddImports(typer.ImportPaths()...)

	collectionExpression, err := typeProvider.CollectionInitializationIfNeeded(info)
	if err != nil {
		return result, err
	}
	// If a field is not a collection or optional, it is required.
	requiredField := collectionExpression == nil && !typeProvider.IsSpecificType(visitors.IsOptional) // TODO(bmoylan) This does not handle aliases of optionals
	seenVar := "seen" + field.FieldSelector

	if requiredField {
		// Declare a 'var seenFieldName bool' which we will set to true inside the case statement.
		result.Init = append(result.Init, statement.NewDecl(decl.NewVar(seenVar, expression.BoolType)))
	}
	if collectionExpression != nil {
		result.Init = append(result.Init, statement.NewAssignment(selector, token.ASSIGN, collectionExpression))
	}

	visitor := &gjsonUnmarshalValueVisitor{
		info:            info,
		selector:        selector,
		valueVar:        "value",
		returnErrStmt:   statement.NewReturn(expression.VariableVal("false")),
		fieldDescriptor: fmt.Sprintf("field %s[%q] ", receiverType, field.JSONKey),
	}
	if err := field.Type.Accept(visitor); err != nil {
		return result, err
	}
	if visitor.typeCheck != nil {
		result.UnmarshalGJSON = append(result.UnmarshalGJSON, visitor.typeCheck)
	}
	result.UnmarshalGJSON = append(result.UnmarshalGJSON, visitor.stmts...)

	if requiredField {
		result.UnmarshalGJSON = append(
			[]astgen.ASTStmt{statement.NewAssignment(expression.VariableVal(seenVar), token.ASSIGN, expression.VariableVal("true"))},
			result.UnmarshalGJSON...)

		result.ValidateReqdField = append(result.ValidateReqdField, &statement.If{
			Cond: expression.NewUnary(token.NOT, expression.VariableVal(seenVar)),
			Body: []astgen.ASTStmt{
				statement.NewAssignment(expression.VariableVal("missingFields"), token.ASSIGN,
					expression.NewCallExpression(expression.AppendBuiltIn, expression.VariableVal("missingFields"), expression.StringVal(field.JSONKey))),
			},
		})
	}

	return result, nil
}

type gjsonUnmarshalValueVisitor struct {
	// in
	info            types.PkgInfo
	selector        astgen.ASTExpr
	valueVar        string
	selectorToken   token.Token
	isMapKey        bool
	nestDepth       int
	returnErrStmt   astgen.ASTStmt
	fieldDescriptor string // used in error messages

	// out
	typeCheck astgen.ASTStmt
	stmts     []astgen.ASTStmt
}

func (v *gjsonUnmarshalValueVisitor) VisitPrimitive(t spec.PrimitiveType) error {
	switch t.Value() {
	case spec.PrimitiveType_ANY:
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "JSON", "String", "Number", "True", "False")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewCallFunction(v.valueVar, "Value"),
		})
	case spec.PrimitiveType_STRING:
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
		})
	case spec.PrimitiveType_INTEGER:
		var rhs astgen.ASTExpr
		if v.isMapKey {
			v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
			rhs = expression.NewCallExpression(expression.IntType, expression.NewCallFunction(v.valueVar, "Int"))
		} else {
			v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "Number")
			rhs = expression.NewCallExpression(expression.IntType, expression.NewCallFunction(v.valueVar, "Int"))
		}
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: rhs,
		})
	case spec.PrimitiveType_DOUBLE:
		v.info.AddImports("math")
		assignDouble := func(rhs astgen.ASTExpr) astgen.ASTStmt {
			return &statement.Assignment{
				LHS: []astgen.ASTExpr{v.selector},
				Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
				RHS: rhs,
			}
		}
		v.stmts = append(v.stmts, &statement.Switch{
			Expression: expression.NewSelector(expression.VariableVal(v.valueVar), "Type"),
			Cases: []statement.CaseClause{
				{
					Exprs: []astgen.ASTExpr{expression.NewSelector(expression.VariableVal("gjson"), "Number")},
					Body:  []astgen.ASTStmt{assignDouble(expression.NewSelector(expression.VariableVal(v.valueVar), "Num"))},
				},
				{
					Exprs: []astgen.ASTExpr{expression.NewSelector(expression.VariableVal("gjson"), "String")},
					Body: []astgen.ASTStmt{&statement.Switch{
						Expression: expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
						Cases: []statement.CaseClause{
							{
								Exprs: []astgen.ASTExpr{expression.StringVal("NaN")},
								Body:  []astgen.ASTStmt{assignDouble(expression.NewCallFunction("math", "NaN"))},
							},
							{
								Exprs: []astgen.ASTExpr{expression.StringVal("Infinity")},
								Body:  []astgen.ASTStmt{assignDouble(expression.NewCallFunction("math", "Inf", expression.IntVal(1)))},
							},
							{
								Exprs: []astgen.ASTExpr{expression.StringVal("-Infinity")},
								Body:  []astgen.ASTStmt{assignDouble(expression.NewCallFunction("math", "Inf", expression.IntVal(-1)))},
							},
							{
								Exprs: nil, // default case
								Body: []astgen.ASTStmt{
									statement.NewAssignment(
										expression.VariableVal("err"),
										token.ASSIGN,
										werrorNew(v.fieldDescriptor+"got invalid json value for double"),
									),
								},
							},
						},
					}},
				},
				{
					Exprs: nil, // default case
					Body: []astgen.ASTStmt{
						statement.NewAssignment(
							expression.VariableVal("err"),
							token.ASSIGN,
							werrorNew(v.fieldDescriptor+"got invalid json type for double"),
						),
					},
				},
			},
		})

		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "Number")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewSelector(expression.VariableVal(v.valueVar), "Num"),
		})
	case spec.PrimitiveType_BOOLEAN:
		var rhs astgen.ASTExpr
		if v.isMapKey {
			v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
			rhs = expression.NewCallFunction("boolean", "Boolean", expression.NewCallFunction(v.valueVar, "Bool"))
		} else {
			v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "False", "True")
			rhs = expression.NewCallFunction(v.valueVar, "Bool")
		}
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: rhs,
		})
	case spec.PrimitiveType_BINARY:
		v.info.AddImports(types.BinaryPkg.ImportPaths()...)
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")

		binaryObj := expression.NewCallExpression(
			expression.Type(types.BinaryPkg.GoType(v.info)),
			expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
		)
		if v.isMapKey {
			// v = binary.Binary(value.Str)
			v.stmts = append(v.stmts, &statement.Assignment{
				LHS: []astgen.ASTExpr{v.selector},
				Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
				RHS: binaryObj,
			})
		} else {
			// v, err = binary.Binary(value.Str).Bytes()
			v.stmts = append(v.stmts, &statement.Assignment{
				LHS: []astgen.ASTExpr{v.selector, expression.VariableVal("err")},
				Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
				RHS: expression.NewCallExpression(expression.NewSelector(binaryObj, "Bytes")),
			})
		}
	case spec.PrimitiveType_BEARERTOKEN:
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewCallFunction("bearertoken", "Token",
				expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
			),
		})
	case spec.PrimitiveType_DATETIME:
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector, expression.VariableVal("err")},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewCallFunction("datetime", "ParseDateTime",
				expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
			),
		})
	case spec.PrimitiveType_RID:
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector, expression.VariableVal("err")},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewCallFunction("rid", "ParseRID",
				expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
			),
		})
	case spec.PrimitiveType_SAFELONG:
		if v.isMapKey {
			v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
			v.stmts = append(v.stmts,
				&statement.Assignment{
					LHS: []astgen.ASTExpr{v.selector, expression.VariableVal("err")},
					Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
					RHS: expression.NewCallFunction("safelong", "ParseSafeLong",
						expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
					),
				},
			)
		} else {
			v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "Number")
			v.stmts = append(v.stmts, &statement.Assignment{
				LHS: []astgen.ASTExpr{v.selector, expression.VariableVal("err")},
				Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
				RHS: expression.NewCallFunction("safelong", "NewSafeLong",
					expression.NewCallFunction(v.valueVar, "Int"),
				),
			})
		}
	case spec.PrimitiveType_UUID:
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, &statement.Assignment{
			LHS: []astgen.ASTExpr{v.selector, expression.VariableVal("err")},
			Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
			RHS: expression.NewCallFunction("uuid", "ParseUUID",
				expression.NewSelector(expression.VariableVal(v.valueVar), "Str"),
			),
		})
	case spec.PrimitiveType_UNKNOWN:
		return errors.New("Unsupported primitive type " + t.String())
	default:
		return errors.New("Unsupported primitive type " + t.String())
	}
	return nil
}

func (v *gjsonUnmarshalValueVisitor) VisitOptional(t spec.OptionalType) error {
	valVar := tmpVarName("optionalValue", v.nestDepth)
	innerVisitor := &gjsonUnmarshalValueVisitor{
		info:            v.info,
		selector:        expression.VariableVal(valVar),
		valueVar:        v.valueVar,
		nestDepth:       v.nestDepth + 1,
		returnErrStmt:   v.returnErrStmt,
		fieldDescriptor: v.fieldDescriptor,
	}
	if err := t.ItemType.Accept(innerVisitor); err != nil {
		return err
	}

	nullCheck := &statement.If{
		Cond: expression.NewBinary(
			expression.NewSelector(expression.VariableVal(v.valueVar), "Type"),
			token.NEQ,
			expression.NewSelector(expression.VariableVal("gjson"), "Null")),
	}

	if innerVisitor.typeCheck != nil {
		nullCheck.Body = append(nullCheck.Body, innerVisitor.typeCheck)
	}
	valTyper, err := visitors.NewConjureTypeProviderTyper(t.ItemType, v.info)
	if err != nil {
		return err
	}
	nullCheck.Body = append(nullCheck.Body, statement.NewDecl(decl.NewVar(valVar, expression.Type(valTyper.GoType(v.info)))))
	nullCheck.Body = append(nullCheck.Body, innerVisitor.stmts...)
	nullCheck.Body = append(nullCheck.Body, &statement.Assignment{
		LHS: []astgen.ASTExpr{v.selector},
		Tok: tokenOrDefault(v.selectorToken, token.ASSIGN),
		RHS: expression.NewUnary(token.AND, expression.VariableVal(valVar)),
	})

	v.stmts = append(v.stmts, nullCheck)
	return nil
}

func (v *gjsonUnmarshalValueVisitor) VisitList(t spec.ListType) error {
	valVar := tmpVarName("listElement", v.nestDepth)
	innerVisitor := &gjsonUnmarshalValueVisitor{
		info:            v.info,
		selector:        expression.VariableVal(valVar),
		valueVar:        "value",
		nestDepth:       v.nestDepth + 1,
		returnErrStmt:   statement.NewReturn(expression.VariableVal("false")),
		fieldDescriptor: v.fieldDescriptor + "list element ",
	}
	if err := t.ItemType.Accept(innerVisitor); err != nil {
		return err
	}

	var innerStmts []astgen.ASTStmt
	if innerVisitor.typeCheck != nil {
		innerStmts = append(innerStmts, innerVisitor.typeCheck)
	}
	valTyper, err := visitors.NewConjureTypeProviderTyper(t.ItemType, v.info)
	if err != nil {
		return err
	}
	innerStmts = append(innerStmts, statement.NewDecl(decl.NewVar(valVar, expression.Type(valTyper.GoType(v.info)))))
	innerStmts = append(innerStmts, innerVisitor.stmts...)
	// x.List = append(x.List, v)
	innerStmts = append(innerStmts, &statement.Assignment{
		LHS: []astgen.ASTExpr{v.selector},
		Tok: token.ASSIGN,
		RHS: expression.NewCallExpression(expression.AppendBuiltIn, v.selector, expression.VariableVal(valVar)),
	})
	innerStmts = append(innerStmts, statement.NewReturn(expression.NewBinary(expression.VariableVal("err"), token.EQL, expression.Nil)))

	v.typeCheck = &statement.If{
		Cond: expression.NewUnary(token.NOT, expression.NewCallFunction(v.valueVar, "IsArray")),
		Body: []astgen.ASTStmt{
			statement.NewAssignment(
				expression.VariableVal("err"),
				token.ASSIGN,
				werrorNew(v.fieldDescriptor+"expected json type Array")),
			v.returnErrStmt,
		},
	}

	// value.ForEach(func(_, value gjson.Result) bool { innerStmts...; return err == nil }
	v.stmts = append(v.stmts, statement.NewExpression(expression.NewCallFunction(v.valueVar, "ForEach",
		expression.NewFuncLit(
			expression.FuncType{
				Params: expression.FuncParams{{
					Names: []string{"_", "value"},
					Type:  "gjson.Result",
				}},
				ReturnTypes: []expression.Type{expression.BoolType},
			},
			innerStmts...,
		),
	)))
	return nil
}

func (v *gjsonUnmarshalValueVisitor) VisitSet(t spec.SetType) error {
	return v.VisitList(spec.ListType{ItemType: t.ItemType})
}

func (v *gjsonUnmarshalValueVisitor) VisitMap(t spec.MapType) error {
	mapTypeProvider, err := visitors.NewConjureTypeProvider(spec.NewTypeFromMap(t))
	if err != nil {
		return err
	}
	keyTypeProvider, err := visitors.NewConjureTypeProvider(t.KeyType)
	if err != nil {
		return err
	}
	keyTyper, err := keyTypeProvider.ParseType(v.info)
	if err != nil {
		return err
	}
	// Use binary.Binary for map keys since []byte is invalid in go maps.
	if keyTypeProvider.IsSpecificType(visitors.IsBinary) {
		keyTyper = types.BinaryPkg
	}
	// Use boolean.Boolean for map keys since conjure boolean keys are serialized as strings
	if keyTypeProvider.IsSpecificType(visitors.IsBoolean) {
		keyTyper = types.BooleanPkg
	}

	var innerStmts []astgen.ASTStmt

	keyVar := expression.VariableVal(tmpVarName("mapKey", v.nestDepth))
	valVar := expression.VariableVal(tmpVarName("mapVal", v.nestDepth))

	keyVisitor := &gjsonUnmarshalValueVisitor{
		info:            v.info,
		selector:        keyVar,
		valueVar:        "key",
		isMapKey:        true,
		nestDepth:       v.nestDepth + 1,
		returnErrStmt:   statement.NewReturn(expression.VariableVal("false")),
		fieldDescriptor: v.fieldDescriptor + "map key ",
	}
	if err := t.KeyType.Accept(keyVisitor); err != nil {
		return err
	}
	//valDecl, err := declVar(string(valVar), t.ValueType, v.info)
	//if err != nil {
	//	return err
	//}
	valVisitor := &gjsonUnmarshalValueVisitor{
		info:            v.info,
		selector:        valVar,
		valueVar:        "value",
		nestDepth:       v.nestDepth + 1,
		returnErrStmt:   statement.NewReturn(expression.VariableVal("false")),
		fieldDescriptor: v.fieldDescriptor + "map value ",
	}
	if err := t.ValueType.Accept(valVisitor); err != nil {
		return err
	}

	if keyVisitor.typeCheck != nil {
		innerStmts = append(innerStmts, keyVisitor.typeCheck)
	}
	if valVisitor.typeCheck != nil {
		innerStmts = append(innerStmts, valVisitor.typeCheck)
	}
	keyDecl := statement.NewDecl(decl.NewVar(string(keyVar), expression.Type(keyTyper.GoType(v.info))))
	innerStmts = append(innerStmts, keyDecl)
	innerStmts = append(innerStmts, keyVisitor.stmts...)
	valTyper, err := visitors.NewConjureTypeProviderTyper(t.ValueType, v.info)
	if err != nil {
		return err
	}
	innerStmts = append(innerStmts, statement.NewDecl(decl.NewVar(string(valVar), expression.Type(valTyper.GoType(v.info)))))
	innerStmts = append(innerStmts, valVisitor.stmts...)

	v.typeCheck = &statement.If{
		Cond: expression.NewUnary(token.NOT, expression.NewCallFunction(v.valueVar, "IsObject")),
		Body: []astgen.ASTStmt{
			statement.NewAssignment(expression.VariableVal("err"), token.ASSIGN, werrorNew(v.fieldDescriptor+"expected json type Object")),
			v.returnErrStmt,
		},
	}
	collectionInit, err := mapTypeProvider.CollectionInitializationIfNeeded(v.info)
	if err != nil {
		return err
	}
	variableInit := statement.NewAssignment(v.selector, tokenOrDefault(v.selectorToken, token.ASSIGN), collectionInit)

	if v.selectorToken == token.DEFINE {
		// v1 := make(map[k]v, 0)
		v.stmts = append(v.stmts, variableInit)
	} else {
		// if r.Field == nil { r.Field = make(map[k]v) }
		v.stmts = append(v.stmts, &statement.If{
			Cond: expression.NewBinary(v.selector, token.EQL, expression.Nil),
			Body: []astgen.ASTStmt{variableInit},
		})
	}

	v.stmts = append(v.stmts,
		// value.ForEach(func(key, value gjson.Result) bool { innerStmts... ; return err == nil }
		statement.NewExpression(expression.NewCallFunction(v.valueVar, "ForEach",
			expression.NewFuncLit(
				expression.FuncType{
					Params: expression.FuncParams{{
						Names: []string{"key", "value"},
						Type:  "gjson.Result",
					}},
					ReturnTypes: []expression.Type{expression.BoolType},
				},
				append(innerStmts,
					statement.NewAssignment(expression.NewIndex(v.selector, keyVar), token.ASSIGN, valVar),
					statement.NewReturn(expression.NewBinary(expression.VariableVal("err"), token.EQL, expression.Nil)),
				)...,
			),
		)),
	)
	return nil
}

func (v *gjsonUnmarshalValueVisitor) VisitExternal(_ spec.ExternalReference) error {
	v.stmts = append(v.stmts, &statement.Assignment{
		LHS: []astgen.ASTExpr{expression.VariableVal("err")},
		Tok: token.ASSIGN,
		RHS: expression.NewCallFunction("safejson", "Unmarshal",
			expression.NewUnary(token.AND, v.selector),
			expression.NewSelector(expression.VariableVal(v.valueVar), "Raw")),
	})
	return nil
}

func (v *gjsonUnmarshalValueVisitor) VisitReference(t spec.TypeName) error {
	typ, ok := v.info.CustomTypes().Get(visitors.TypeNameToTyperName(t))
	if !ok {
		return errors.Errorf("reference type not found %s", t.Name)
	}
	defVisitor := gjsonUnmarshalValueReferenceDefVisitor{
		info:            v.info,
		selector:        v.selector,
		valueVar:        v.valueVar,
		typer:           typ,
		selectorToken:   v.selectorToken,
		nestDepth:       v.nestDepth,
		returnErrStmt:   v.returnErrStmt,
		fieldDescriptor: v.fieldDescriptor,
	}
	if err := typ.Def.Accept(&defVisitor); err != nil {
		return err
	}

	v.typeCheck = defVisitor.typeCheck
	v.stmts = append(v.stmts, defVisitor.stmts...)
	return nil
}

type gjsonUnmarshalValueReferenceDefVisitor struct {
	// in
	info            types.PkgInfo
	selector        astgen.ASTExpr
	valueVar        string
	typer           types.Typer
	selectorToken   token.Token
	nestDepth       int
	returnErrStmt   astgen.ASTStmt
	fieldDescriptor string // used in error messages

	// out
	typeCheck astgen.ASTStmt
	stmts     []astgen.ASTStmt
}

func (v *gjsonUnmarshalValueReferenceDefVisitor) VisitAlias(def spec.AliasDefinition) error {
	aliasTypeProvider, err := visitors.NewConjureTypeProvider(def.Alias)
	if err != nil {
		return err
	}
	if aliasTypeProvider.IsSpecificType(visitors.IsString) {
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, statement.NewAssignment(v.selector, tokenOrDefault(v.selectorToken, token.ASSIGN),
			expression.NewCallExpression(expression.Type(v.typer.GoType(v.info)), expression.NewSelector(expression.VariableVal(v.valueVar), "Str"))))
	} else if aliasTypeProvider.IsSpecificType(visitors.IsText) {
		v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
		v.stmts = append(v.stmts, unmarshalTextValue(v.selector, v.valueVar))
	} else {
		v.stmts = append(v.stmts, v.unmarshalJSONStringValue()...)
	}
	return nil
}

func (v *gjsonUnmarshalValueReferenceDefVisitor) VisitEnum(_ spec.EnumDefinition) error {
	v.typeCheck = gjsonTypeCheck(v.fieldDescriptor, v.returnErrStmt, v.valueVar, "String")
	v.stmts = append(v.stmts, unmarshalTextValue(v.selector, v.valueVar))
	return nil
}

func (v *gjsonUnmarshalValueReferenceDefVisitor) VisitObject(def spec.ObjectDefinition) error {
	if len(def.Fields) > 0 {
		v.stmts = append(v.stmts, v.unmarshalJSONStringValue()...)
	}
	return nil
}

func (v *gjsonUnmarshalValueReferenceDefVisitor) VisitUnion(_ spec.UnionDefinition) error {
	v.stmts = append(v.stmts, v.unmarshalJSONStringValue()...)
	return nil
}

func (v *gjsonUnmarshalValueReferenceDefVisitor) VisitUnknown(typeName string) error {
	return errors.Errorf("unknown type %q", typeName)
}

func (v *gjsonUnmarshalValueVisitor) VisitUnknown(typeName string) error {
	return errors.Errorf("unknown type %q", typeName)
}

func unmarshalTextValue(selector astgen.ASTExpr, valueVar string) astgen.ASTStmt {
	return &statement.Assignment{
		LHS: []astgen.ASTExpr{expression.VariableVal("err")},
		Tok: token.ASSIGN,
		RHS: expression.NewCallExpression(expression.NewSelector(selector, "UnmarshalText"),
			expression.NewCallExpression(expression.Type("[]byte"), expression.NewSelector(expression.VariableVal(valueVar), "Str"))),
	}
}

func (v *gjsonUnmarshalValueReferenceDefVisitor) unmarshalJSONStringValue() []astgen.ASTStmt {
	return []astgen.ASTStmt{
		&statement.If{
			Cond: expression.VariableVal("strict"),
			Body: []astgen.ASTStmt{&statement.Assignment{
				LHS: []astgen.ASTExpr{expression.VariableVal("err")},
				Tok: token.ASSIGN,
				RHS: expression.NewCallExpression(expression.NewSelector(v.selector, "UnmarshalJSONStringStrict"),
					expression.NewSelector(expression.VariableVal(v.valueVar), "Raw")),
			}},
			Else: &statement.Assignment{
				LHS: []astgen.ASTExpr{expression.VariableVal("err")},
				Tok: token.ASSIGN,
				RHS: expression.NewCallExpression(expression.NewSelector(v.selector, "UnmarshalJSONString"),
					expression.NewSelector(expression.VariableVal(v.valueVar), "Raw")),
			},
		},
		&statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal("err")},
			Tok: token.ASSIGN,
			RHS: expression.NewCallFunction("werror", "WrapWithContextParams",
				expression.VariableVal("ctx"),
				expression.VariableVal("err"),
				expression.StringVal(strings.TrimSpace(v.fieldDescriptor)),
			),
		},
	}

}

// if value.Type != gjson.Number { err = werror.Wrap(errors.NewInvalidArgument(), "field \"my-bool\" expected json type True/False"); returnErrStmt }
func gjsonTypeCheck(fieldDescriptor string, returnErrStmt astgen.ASTStmt, valueVar string, typeNames ...string) astgen.ASTStmt {
	var cond astgen.ASTExpr
	for _, typeName := range typeNames {
		test := expression.NewBinary(
			expression.NewSelector(expression.VariableVal(valueVar), "Type"),
			token.NEQ,
			expression.NewSelector(expression.VariableVal("gjson"), typeName),
		)
		if cond == nil {
			cond = test
		} else {
			cond = expression.NewBinary(cond, token.LAND, test)
		}
	}
	errExpr := werrorNew(fmt.Sprintf("%sexpected json type %s", fieldDescriptor, strings.Join(typeNames, "/")))
	return &statement.If{
		Cond: cond,
		Body: []astgen.ASTStmt{
			statement.NewAssignment(expression.VariableVal("err"), token.ASSIGN, errExpr), // TODO add more type context to errors
			returnErrStmt,
		},
	}
}

// werror.ErrorWithContextParams(ctx, "message", werrorParams...)
func werrorNew(message string, werrorParams ...astgen.ASTExpr) astgen.ASTExpr {
	return expression.NewCallFunction("werror", "ErrorWithContextParams",
		append([]astgen.ASTExpr{
			expression.VariableVal("ctx"),
			expression.StringVal(message),
		}, werrorParams...)...)
}

func tokenOrDefault(t, d token.Token) token.Token {
	if t == 0 {
		return d
	}
	return t
}

func tmpVarName(base string, depth int) string {
	if depth == 0 {
		return base
	}
	return fmt.Sprintf("%s%d", base, depth)
}
