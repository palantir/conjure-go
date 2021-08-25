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
	"fmt"
	"go/token"
	"strconv"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"github.com/pkg/errors"
)

var (
	marshalBufVar = expression.VariableVal("buf")
)

func publicMarshalJSONMethods(receiverName string, receiverType string, bodyMarshalJSONBuffer []astgen.ASTStmt) []astgen.ASTDecl {
	return []astgen.ASTDecl{
		&decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(receiverType),
			Function: decl.Function{
				Name: "MarshalJSON",
				FuncType: expression.FuncType{
					ReturnTypes: []expression.Type{expression.ByteSliceType, expression.ErrorType},
				},
				Body: []astgen.ASTStmt{statement.NewReturn(
					expression.NewCallFunction(receiverName, "MarshalJSONBuffer", expression.Nil),
				)},
			},
		},
		&decl.Method{
			ReceiverName: receiverName,
			ReceiverType: expression.Type(receiverType),
			Function: decl.Function{
				Name: "MarshalJSONBuffer",
				FuncType: expression.FuncType{
					Params: expression.FuncParams{&expression.FuncParam{
						Names: []string{string(marshalBufVar)},
						Type:  expression.ByteSliceType,
					}},
					ReturnTypes: []expression.Type{expression.ByteSliceType, expression.ErrorType},
				},
				Body: bodyMarshalJSONBuffer,
			},
		},
	}
}

func visitAliasMarshalGJSONMethodBody(receiverName string, aliasType spec.Type, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	typeProvider, err := visitors.NewConjureTypeProvider(aliasType)
	if err != nil {
		return nil, err
	}
	typer, err := typeProvider.ParseType(info)
	if err != nil {
		return nil, err
	}

	var valueSel astgen.ASTExpr
	if typeProvider.IsSpecificType(visitors.IsOptional) {
		valueSel = expression.NewSelector(expression.VariableVal(receiverName), "Value")
	} else {
		valueSel = expression.NewCallExpression(expression.Type(typer.GoType(info)), expression.VariableVal(receiverName))
	}

	visitor := &jsonMarshalValueVisitor{info: info, selector: valueSel}
	if err := aliasType.Accept(visitor); err != nil {
		return nil, err
	}
	body := append(visitor.stmts, statement.NewReturn(marshalBufVar, expression.Nil))
	return body, nil
}

func visitStructFieldsMarshalGJSONMethodBody(receiverName, receiverType string, fields []JSONField, info types.PkgInfo) ([]astgen.ASTStmt, error) {
	var body []astgen.ASTStmt
	body = append(body, appendMarshalBufferLiteralRune('{'))
	//if len(fields) > 0 {
	//	body = append(body, trailingElemVarDecl)
	//}
	//trackTrailingElements, err := fieldsContainOptional(fields)
	//if err != nil {
	//	return nil, err
	//}
	for i, field := range fields {
		body = append(body, appendMarshalBufferQuotedString(expression.StringVal(field.JSONKey)))
		body = append(body, appendMarshalBuffer(expression.VariableVal(`':'`)))
		visitor := &jsonMarshalValueVisitor{
			info:     info,
			selector: expression.NewSelector(expression.VariableVal(receiverName), field.FieldSelector),
		}
		if err := field.Type.Accept(visitor); err != nil {
			return nil, err
		}
		body = append(body, visitor.stmts...)
		if i < len(fields)-1 {
			body = append(body, appendMarshalBufferLiteralRune(','))
		} else {
			body = append(body, appendMarshalBufferLiteralRune('}'))
		}
	}
	body = append(body, statement.NewReturn(marshalBufVar, expression.Nil))
	return body, nil
}

func appendMarshalBuffer(exprs ...astgen.ASTExpr) astgen.ASTStmt {
	return statement.NewAssignment(
		marshalBufVar,
		token.ASSIGN,
		expression.NewCallExpression(expression.AppendBuiltIn, append([]astgen.ASTExpr{marshalBufVar}, exprs...)...),
	)
}

func appendMarshalBufferLiteralRune(r rune) astgen.ASTStmt {
	return appendMarshalBuffer(expression.VariableVal(fmt.Sprintf("'%c'", r)))
}

func appendMarshalBufferLiteralString(s string) astgen.ASTStmt {
	return appendMarshalBuffer(expression.VariableVal(strconv.Quote(s) + "..."))
}

func appendMarshalBufferQuotedString(expr astgen.ASTExpr) astgen.ASTStmt {
	return statement.NewAssignment(
		marshalBufVar,
		token.ASSIGN,
		expression.NewCallFunction("safejson", "AppendQuotedString", marshalBufVar, expr),
	)
}

var (
	trailingElemVar     = expression.VariableVal("trailingElem")
	trailingElemVarDecl = statement.NewDecl(decl.NewVar(string(trailingElemVar), expression.BoolType))
)

func trailingCommaIfStatement() *statement.If {
	return &statement.If{
		Cond: trailingElemVar,
		Body: []astgen.ASTStmt{appendMarshalBufferLiteralRune(',')},
		Else: statement.NewBlock(
			statement.NewAssignment(trailingElemVar, token.ASSIGN, expression.VariableVal("true")),
		),
	}
}

type jsonMarshalValueVisitor struct {
	// in
	info     types.PkgInfo
	selector astgen.ASTExpr
	isMapKey bool

	// out
	stmts []astgen.ASTStmt
}

func (v *jsonMarshalValueVisitor) VisitPrimitive(t spec.PrimitiveType) error {
	switch t.Value() {
	case spec.PrimitiveType_ANY:
		v.stmts = append(v.stmts, &statement.If{
			Init: &statement.Assignment{
				LHS: []astgen.ASTExpr{expression.VariableVal("jsonBytes"), expression.VariableVal("err")},
				Tok: token.DEFINE,
				RHS: expression.NewCallFunction("safejson", "Marshal", v.selector),
			},
			Cond: expression.NewBinary(expression.VariableVal("err"), token.NEQ, expression.Nil),
			Body: []astgen.ASTStmt{statement.NewReturn(expression.Nil, expression.VariableVal("err"))},
			Else: statement.NewBlock(appendMarshalBuffer(expression.VariableVal("jsonBytes..."))),
		})
	case spec.PrimitiveType_STRING:
		v.stmts = append(v.stmts, appendMarshalBufferQuotedString(v.selector))
	case spec.PrimitiveType_DATETIME:
		v.stmts = append(v.stmts, appendMarshalBufferQuotedString(expression.NewCallExpression(expression.NewSelector(v.selector, "String"))))
	case spec.PrimitiveType_INTEGER, spec.PrimitiveType_SAFELONG:
		appendInt := statement.NewAssignment(marshalBufVar, token.ASSIGN,
			expression.NewCallFunction("strconv", "AppendInt",
				marshalBufVar,
				expression.NewCallExpression(expression.Type("int64"), v.selector),
				expression.IntVal(10)))
		if v.isMapKey {
			v.stmts = append(v.stmts, appendMarshalBufferLiteralRune('"'), appendInt, appendMarshalBufferLiteralRune('"'))
		} else {
			v.stmts = append(v.stmts, appendInt)
		}
	case spec.PrimitiveType_DOUBLE:
		appendFloat := statement.NewAssignment(marshalBufVar, token.ASSIGN,
			expression.NewCallFunction("strconv", "AppendFloat",
				marshalBufVar, v.selector, expression.IntVal(-1), expression.IntVal(10), expression.IntVal(64)))
		var defaultCase []astgen.ASTStmt
		if v.isMapKey {
			defaultCase = []astgen.ASTStmt{
				appendMarshalBufferLiteralRune('"'),
				appendFloat,
				appendMarshalBufferLiteralRune('"'),
			}
		} else {
			defaultCase = []astgen.ASTStmt{appendFloat}
		}
		v.stmts = append(v.stmts, &statement.Switch{
			Expression: v.selector,
			Cases: []statement.CaseClause{
				{
					Body: defaultCase,
				},
				{
					Exprs: []astgen.ASTExpr{expression.NewCallFunction("math", "IsNaN", v.selector)},
					Body:  []astgen.ASTStmt{appendMarshalBufferLiteralString("NaN")},
				},
				{
					Exprs: []astgen.ASTExpr{expression.NewCallFunction("math", "IsInf", v.selector, expression.IntVal(1))},
					Body:  []astgen.ASTStmt{appendMarshalBufferLiteralString("Infinity")},
				},
				{
					Exprs: []astgen.ASTExpr{expression.NewCallFunction("math", "IsInf", v.selector, expression.IntVal(-1))},
					Body:  []astgen.ASTStmt{appendMarshalBufferLiteralString("-Infinity")},
				},
			},
		})
	case spec.PrimitiveType_BINARY:
		if v.isMapKey {
			v.stmts = append(v.stmts, appendMarshalBufferQuotedString(expression.NewCallExpression(expression.StringType, v.selector)))
		} else {
			v.stmts = append(v.stmts,
				appendMarshalBufferLiteralRune('"'),
				&statement.If{
					Cond: expression.NewBinary(
						expression.NewCallExpression(expression.LenBuiltIn, v.selector),
						token.GTR,
						expression.IntVal(0)),
					Body: []astgen.ASTStmt{
						statement.NewAssignment(
							expression.VariableVal("b64out"),
							token.DEFINE,
							expression.NewCallExpression(
								expression.MakeBuiltIn,
								expression.ByteSliceType,
								expression.IntVal(0),
								expression.NewCallFunction("base64.StdEncoding", "EncodedLen",
									expression.NewCallExpression(expression.LenBuiltIn, v.selector),
								),
							),
						),
						statement.NewExpression(expression.NewCallFunction(
							"base64.StdEncoding", "Encode", expression.VariableVal("b64out"), v.selector)),
						appendMarshalBuffer(expression.VariableVal("b64out...")),
					},
				},
				appendMarshalBufferLiteralRune('"'),
			)
		}
	case spec.PrimitiveType_BOOLEAN:
		if v.isMapKey {
			v.stmts = append(v.stmts, &statement.If{
				Cond: v.selector,
				Body: []astgen.ASTStmt{appendMarshalBuffer(expression.VariableVal(`"\"true\""...`))},
				Else: statement.NewBlock(appendMarshalBuffer(expression.VariableVal(`"\"false\""...`))),
			})
		} else {
			v.stmts = append(v.stmts, &statement.If{
				Cond: v.selector,
				Body: []astgen.ASTStmt{appendMarshalBuffer(expression.VariableVal(`"true"...`))},
				Else: statement.NewBlock(appendMarshalBuffer(expression.VariableVal(`"false"...`))),
			})
		}
	case spec.PrimitiveType_UUID:
		v.stmts = append(v.stmts, appendMarshalBufferQuotedString(
			expression.NewCallExpression(expression.NewSelector(v.selector, "String"))))
	case spec.PrimitiveType_RID:
		v.stmts = append(v.stmts, appendMarshalBufferQuotedString(
			expression.NewCallExpression(expression.NewSelector(v.selector, "String"))))
	case spec.PrimitiveType_BEARERTOKEN:
		v.stmts = append(v.stmts, appendMarshalBufferQuotedString(
			expression.NewCallExpression(expression.NewSelector(v.selector, "String"))))
	case spec.PrimitiveType_UNKNOWN:
		return errors.Errorf("unknown type %q", t.String())
	}
	return nil
}

func (v *jsonMarshalValueVisitor) VisitOptional(t spec.OptionalType) error {
	visitor := &jsonMarshalValueVisitor{
		info:     v.info,
		selector: expression.NewUnary(token.MUL, v.selector),
	}
	if err := t.ItemType.Accept(visitor); err != nil {
		return err
	}
	v.stmts = append(v.stmts, &statement.If{
		Cond: expression.NewBinary(v.selector, token.NEQ, expression.Nil),
		Body: visitor.stmts,
		Else: statement.NewBlock(appendMarshalBuffer(expression.VariableVal(`"null"...`))),
	})
	return nil
}

func (v *jsonMarshalValueVisitor) VisitList(t spec.ListType) error {
	visitor := &jsonMarshalValueVisitor{
		info:     v.info,
		selector: expression.NewIndex(v.selector, expression.VariableVal("i")),
	}
	if err := t.ItemType.Accept(visitor); err != nil {
		return err
	}
	v.stmts = append(v.stmts,
		appendMarshalBufferLiteralRune('['),
		statement.NewBlock(
			trailingElemVarDecl,
			&statement.Range{
				Key:  expression.VariableVal("i"),
				Tok:  token.DEFINE,
				Expr: v.selector,
				Body: append([]astgen.ASTStmt{trailingCommaIfStatement()}, visitor.stmts...),
			},
		),
		appendMarshalBufferLiteralRune(']'),
	)
	return nil
}

func (v *jsonMarshalValueVisitor) VisitSet(t spec.SetType) error {
	return v.VisitList(spec.ListType{ItemType: t.ItemType})
}

func (v *jsonMarshalValueVisitor) VisitMap(t spec.MapType) error {
	keyVisitor := &jsonMarshalValueVisitor{
		info:     v.info,
		selector: expression.VariableVal("k"),
		isMapKey: true,
	}
	if err := t.KeyType.Accept(keyVisitor); err != nil {
		return err
	}
	valVisitor := &jsonMarshalValueVisitor{
		info:     v.info,
		selector: expression.VariableVal("v"),
	}
	if err := t.ValueType.Accept(valVisitor); err != nil {
		return err
	}
	var innerStmts []astgen.ASTStmt
	innerStmts = append(innerStmts, keyVisitor.stmts...)
	innerStmts = append(innerStmts, appendMarshalBufferLiteralRune(':'))
	innerStmts = append(innerStmts, valVisitor.stmts...)
	innerStmts = append(innerStmts,
		statement.NewExpression(expression.NewBinary(expression.VariableVal("i"), token.ADD_ASSIGN, expression.IntVal(1))),
		&statement.If{
			Cond: expression.NewBinary(expression.VariableVal("i"), token.LSS,
				expression.NewBinary(
					expression.NewCallExpression(expression.LenBuiltIn, v.selector),
					token.SUB,
					expression.IntVal(1),
				),
			),
			Body: []astgen.ASTStmt{appendMarshalBufferLiteralRune(',')},
		})
	v.stmts = append(v.stmts, statement.NewBlock(
		appendMarshalBufferLiteralRune('{'),
		statement.NewAssignment(expression.VariableVal("i"), token.DEFINE, expression.IntVal(0)),
		&statement.Range{
			Key:   expression.VariableVal("k"),
			Value: expression.VariableVal("v"),
			Tok:   token.DEFINE,
			Expr:  v.selector,
			Body:  innerStmts,
		},
		appendMarshalBufferLiteralRune('}'),
	))
	return nil
}

func (v *jsonMarshalValueVisitor) VisitReference(t spec.TypeName) error {
	typ, ok := v.info.CustomTypes().Get(visitors.TypeNameToTyperName(t))
	if !ok {
		return errors.Errorf("reference type not found %s", t.Name)
	}
	visitor := &jsonMarshalValueReferenceDefVisitor{selector: v.selector}
	if err := typ.Def.Accept(visitor); err != nil {
		return err
	}
	v.stmts = append(v.stmts, visitor.stmts...)
	return nil
}

func (v *jsonMarshalValueVisitor) VisitExternal(t spec.ExternalReference) error {
	panic("implement me")
}

func (v *jsonMarshalValueVisitor) VisitUnknown(typeName string) error {
	return errors.Errorf("unknown type %q", typeName)
}

type jsonMarshalValueReferenceDefVisitor struct {
	selector astgen.ASTExpr
	stmts    []astgen.ASTStmt
}

func (v *jsonMarshalValueReferenceDefVisitor) VisitAlias(d spec.AliasDefinition) error {
	return v.visitJSONMarshaler()
}

func (v *jsonMarshalValueReferenceDefVisitor) VisitEnum(d spec.EnumDefinition) error {
	v.stmts = append(v.stmts, appendMarshalBufferQuotedString(expression.NewCallExpression(expression.NewSelector(v.selector, "String"))))
	return nil
}

func (v *jsonMarshalValueReferenceDefVisitor) VisitObject(d spec.ObjectDefinition) error {
	return v.visitJSONMarshaler()
}

func (v *jsonMarshalValueReferenceDefVisitor) VisitUnion(d spec.UnionDefinition) error {
	return v.visitJSONMarshaler()
}

func (v *jsonMarshalValueReferenceDefVisitor) visitJSONMarshaler() error {
	v.stmts = append(v.stmts, &statement.If{
		Init: &statement.Assignment{
			LHS: []astgen.ASTExpr{expression.VariableVal("out"), expression.VariableVal("err")},
			Tok: token.DEFINE,
			RHS: expression.NewCallExpression(expression.NewSelector(v.selector, "MarshalJSONBuffer"), marshalBufVar),
		},
		Cond: expression.NewBinary(expression.VariableVal("err"), token.NEQ, expression.Nil),
		Body: []astgen.ASTStmt{statement.NewReturn(expression.Nil, expression.VariableVal("err"))},
		Else: statement.NewBlock(statement.NewAssignment(marshalBufVar, token.ASSIGN, expression.VariableVal("out"))),
	})
	return nil
}

func (v *jsonMarshalValueReferenceDefVisitor) VisitUnknown(typeName string) error {
	return errors.Errorf("unknown type %q", typeName)
}

func fieldsBeginsWithOptionals(fields []JSONField) (n int, found bool) {
	n = -1
	for i, field := range fields {
		if err := field.Type.AcceptFuncs(
			field.Type.PrimitiveNoopSuccess,
			func(optionalType spec.OptionalType) error {
				n = i
				found = true
				return nil
			},
			field.Type.ListNoopSuccess,
			field.Type.SetNoopSuccess,
			field.Type.MapNoopSuccess,
			field.Type.ReferenceNoopSuccess,
			field.Type.ExternalNoopSuccess,
			field.Type.ErrorOnUnknown,
		); err != nil {
			panic(err)
		}
		if found {
			break
		}
	}
	return
}
