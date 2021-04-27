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

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	astspec "github.com/palantir/goastwriter/spec"
	"github.com/palantir/goastwriter/statement"
)

const (
	enumReceiverName    = "e"
	enumUpperVarName    = "v"
	enumUnknownValue    = "UNKNOWN"
	enumStructFieldName = "val"
)

func astForEnum(enumDefinition spec.EnumDefinition, info types.PkgInfo) []astgen.ASTDecl {
	enumName := enumDefinition.TypeName.Name

	valueType := expression.Type(enumName + "_Value")
	valueTypeDef := &decl.Alias{
		Comment: transforms.Documentation(enumDefinition.Docs),
		Name:    string(valueType),
		Type:    expression.StringType,
	}

	structType := expression.Type(enumName)
	structTypeDef := &decl.Struct{
		Comment: transforms.Documentation(enumDefinition.Docs),
		Name:    string(structType),
		StructType: expression.StructType{
			Fields: []*expression.StructField{
				{
					Name: enumStructFieldName,
					Type: valueType,
				},
			},
		},
	}

	var vals []*astspec.Value
	for _, currVal := range enumDefinition.Values {
		vals = append(vals, &astspec.Value{
			Comment: transforms.Documentation(currVal.Docs),
			Names:   []string{enumName + "_" + currVal.Value},
			Type:    valueType,
			Values:  []astgen.ASTExpr{expression.StringVal(currVal.Value)},
		})
	}
	vals = append(vals, &astspec.Value{
		Names:  []string{enumName + "_" + enumUnknownValue},
		Type:   valueType,
		Values: []astgen.ASTExpr{expression.StringVal(enumUnknownValue)},
	})
	valsDecl := &decl.Const{Values: vals}

	enumConstructorDecl := &decl.Function{
		Name: "New_" + structTypeDef.Name,
		FuncType: expression.FuncType{
			Params: []*expression.FuncParam{
				{
					Names: []string{"value"},
					Type:  valueType,
				},
			},
			ReturnTypes: []expression.Type{structType},
		},
		Body: []astgen.ASTStmt{
			statement.NewReturn(&expression.CompositeLit{
				Type: structType,
				Elements: []astgen.ASTExpr{
					expression.NewKeyValue(enumStructFieldName, expression.VariableVal("value")),
				},
			}),
		},
	}

	valueStructMethod := &decl.Method{
		ReceiverName: enumReceiverName,
		ReceiverType: structType,
		Function: decl.Function{
			Name: "Value",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{valueType},
			},
			Body: []astgen.ASTStmt{
				&statement.If{
					Cond: expression.NewCallFunction(enumReceiverName, "IsUnknown"),
					Body: []astgen.ASTStmt{
						statement.NewReturn(expression.VariableVal(enumName + "_" + enumUnknownValue)),
					},
				},
				statement.NewReturn(expression.NewSelector(expression.VariableVal(enumReceiverName), enumStructFieldName)),
			},
		},
	}

	return []astgen.ASTDecl{
		structTypeDef,
		valueTypeDef,
		valsDecl,
		enumValuesAST(enumDefinition),
		enumConstructorDecl,
		enumIsUnknownAST(enumDefinition),
		valueStructMethod,
		enumStringAST(string(structType)),
		enumMarshalTextAST(string(structType)),
		enumUnmarshalTextAST(enumDefinition, info),
	}
}

func enumStringAST(enumName string) astgen.ASTDecl {
	return newStringMethod(enumReceiverName, enumName, statement.NewReturn(
		expression.NewCallExpression(expression.StringType,
			expression.NewSelector(expression.VariableVal(enumReceiverName), enumStructFieldName))))
}

func enumMarshalTextAST(enumName string) astgen.ASTDecl {
	return newMarshalTextMethod(enumReceiverName, enumName, statement.NewReturn(
		expression.NewCallExpression(expression.ByteSliceType,
			expression.NewSelector(expression.VariableVal(enumReceiverName), enumStructFieldName)),
		expression.Nil))
}

func enumUnmarshalTextAST(e spec.EnumDefinition, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports("strings")
	info.AddImports("github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors")
	info.AddImports("github.com/palantir/witchcraft-go-error")
	info.AddImports("github.com/palantir/witchcraft-go-params")

	enumName := e.TypeName.Name

	switchStmt := &statement.Switch{
		Init:       statement.NewAssignment(expression.VariableVal(enumUpperVarName), token.DEFINE, expression.NewCallFunction("strings", "ToUpper", expression.NewCallExpression(expression.StringType, expression.VariableVal(dataVarName)))),
		Expression: expression.VariableVal(enumUpperVarName),
		Cases: []statement.CaseClause{
			// default case
			{
				Body: []astgen.ASTStmt{
					statement.NewAssignment(
						expression.NewUnary(token.MUL, expression.VariableVal(enumReceiverName)),
						token.ASSIGN,
						expression.NewCallExpression(
							expression.Type("New_"+enumName),
							expression.NewCallExpression(expression.Type(enumName+"_Value"), expression.VariableVal(enumUpperVarName)),
						),
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
				expression.NewCallExpression(
					expression.Type("New_"+enumName),
					expression.VariableVal(enumName+"_"+currVal.Value),
				),
			),
		))
	}
	return newUnmarshalTextMethod(enumReceiverName, enumName, switchStmt, statement.NewReturn(expression.Nil))
}

func enumIsUnknownAST(e spec.EnumDefinition) astgen.ASTDecl {
	return &decl.Method{
		ReceiverName: enumReceiverName,
		ReceiverType: expression.Type(e.TypeName.Name),
		Function: decl.Function{
			Comment: fmt.Sprintf("IsUnknown returns false for all known variants of %s and true otherwise.", e.TypeName.Name),
			Name:    "IsUnknown",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.BoolType},
			},
			Body: []astgen.ASTStmt{
				&statement.Switch{
					Expression: expression.NewSelector(expression.VariableVal(enumReceiverName), enumStructFieldName),
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
	typeName := e.TypeName.Name
	funcName := typeName + "_Values"
	sliceType := expression.Type("[]" + typeName + "_Value")
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
	values := make([]astgen.ASTExpr, 0, len(e.Values))
	for _, currVal := range e.Values {
		values = append(values, expression.VariableVal(e.TypeName.Name+"_"+currVal.Value))
	}
	return values
}
