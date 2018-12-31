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
	"go/token"

	"github.com/danverbraganza/varcaser/varcaser"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	astspec "github.com/palantir/goastwriter/spec"
	"github.com/palantir/goastwriter/statement"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/transforms"
)

const (
	enumReceiverName = "e"
	unknownEnumValue = "UNKNOWN"
)

func astForEnum(enumDefinition spec.EnumDefinition) ([]astgen.ASTDecl, StringSet) {
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
	vals = append(vals, &astspec.Value{
		Names:  []string{enumName + toCamelCase(unknownEnumValue)},
		Type:   expression.Type(enumName),
		Values: []astgen.ASTExpr{expression.StringVal(unknownEnumValue)},
	})
	valsDecl := &decl.Const{Values: vals}

	unmarshalDecl, imports := enumUnmarshalJSONAST(enumDefinition)

	return []astgen.ASTDecl{typeDef, valsDecl, unmarshalDecl}, imports
}

func enumUnmarshalJSONAST(e spec.EnumDefinition) (astgen.ASTDecl, StringSet) {
	const (
		stringVar = "s"
		errVar    = "err"
		dataVar   = "data"
	)

	imports := NewStringSet()

	var stmts []astgen.ASTStmt

	toCamelCase := varcaser.Caser{From: varcaser.ScreamingSnakeCase, To: varcaser.UpperCamelCase}.String

	stmts = append(stmts,
		statement.NewDecl(decl.NewVar(stringVar, expression.StringType)),
		ifErrNotNilReturnErrStatement(errVar, statement.NewAssignment(expression.VariableVal(errVar), token.DEFINE, expression.NewCallFunction(
			"json",
			"Unmarshal",
			expression.VariableVal(dataVar),
			expression.NewUnary(token.AND, expression.VariableVal(stringVar)),
		))),
	)
	imports["encoding/json"] = struct{}{}

	// start with default case
	cases := []statement.CaseClause{
		// default case
		{
			Body: []astgen.ASTStmt{
				statement.NewAssignment(
					expression.NewUnary(token.MUL, expression.VariableVal(enumReceiverName)),
					token.ASSIGN,
					expression.VariableVal(e.TypeName.Name+toCamelCase(unknownEnumValue)),
				),
			},
		},
	}

	for _, currVal := range e.Values {
		cases = append(cases, *statement.NewCaseClause(
			expression.StringVal(currVal.Value),
			statement.NewAssignment(
				expression.NewUnary(token.MUL, expression.VariableVal(enumReceiverName)),
				token.ASSIGN,
				expression.VariableVal(e.TypeName.Name+toCamelCase(currVal.Value)),
			),
		))
	}

	stmts = append(stmts, &statement.Switch{
		Expression: expression.NewCallFunction("strings", "ToUpper", expression.VariableVal(stringVar)),
		Cases:      cases,
	})
	imports["strings"] = struct{}{}

	stmts = append(stmts, statement.NewReturn(expression.Nil))

	return &decl.Method{
		Function: decl.Function{
			Name: "UnmarshalJSON",
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam(dataVar, expression.Type("[]byte")),
				},
				ReturnTypes: []expression.Type{
					expression.ErrorType,
				},
			},
			Body: stmts,
		},
		ReceiverName: enumReceiverName,
		ReceiverType: expression.Type(transforms.Export(e.TypeName.Name)).Pointer(),
	}, imports

}
