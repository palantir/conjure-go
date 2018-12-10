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
	"github.com/palantir/goastwriter/spec"
	"github.com/palantir/goastwriter/statement"

	"github.com/palantir/conjure-go/conjure/transforms"
)

const (
	enumReceiverName = "e"
	unknownEnumValue = "UNKNOWN"
)

type Enum struct {
	Name    string
	Values  []Value
	Comment string
}

func (e *Enum) ASTDeclers() ([]astgen.ASTDecl, StringSet) {
	var returnVal []astgen.ASTDecl

	returnVal = append(returnVal, &decl.Alias{
		Comment: e.Comment,
		Name:    e.Name,
		Type:    expression.StringType,
	})

	toCamelCase := varcaser.Caser{From: varcaser.ScreamingSnakeCase, To: varcaser.UpperCamelCase}.String

	var vals []*spec.Value
	for _, currVal := range e.Values {
		vals = append(vals, &spec.Value{
			Comment: currVal.Docs,
			Names:   []string{e.Name + toCamelCase(currVal.Value)},
			Type:    expression.Type(e.Name),
			Values:  []astgen.ASTExpr{expression.StringVal(currVal.Value)},
		})
	}

	vals = append(vals, &spec.Value{
		Names:  []string{e.Name + toCamelCase(unknownEnumValue)},
		Type:   expression.Type(e.Name),
		Values: []astgen.ASTExpr{expression.StringVal(unknownEnumValue)},
	})

	unmarshalDecl, imports := e.unmarshalJSONAST()

	return append(
		returnVal,
		&decl.Const{Values: vals},
		unmarshalDecl,
	), imports
}

func (e *Enum) unmarshalJSONAST() (astgen.ASTDecl, StringSet) {
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
					expression.VariableVal(e.Name+toCamelCase(unknownEnumValue)),
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
				expression.VariableVal(e.Name+toCamelCase(currVal.Value)),
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
		ReceiverType: expression.Type(transforms.Export(e.Name)).Pointer(),
	}, imports

}
