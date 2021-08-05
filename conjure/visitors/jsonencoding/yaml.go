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

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
)

func unmarshalYAML(receiverName, receiverType string) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType).Pointer(),
		Function: decl.Function{
			Name:    "UnmarshalYAML",
			Comment: "UnmarshalYAML implements yaml.Unmarshaler. It converts the YAML to JSON, then runs UnmarshalJSON.",
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam("unmarshal", expression.Type("func(interface{}) error")),
				},
				ReturnTypes: []expression.Type{
					expression.ErrorType,
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewAssignment(expression.VariableVal("ctx"), token.DEFINE, expression.NewCallFunction("context", "TODO")),
				&statement.Assignment{
					LHS: []astgen.ASTExpr{expression.VariableVal("data"), expression.VariableVal("err")},
					Tok: token.DEFINE,
					RHS: expression.NewCallFunction("safeyaml", "UnmarshalerToJSONBytes", expression.VariableVal("unmarshal")),
				},
				&statement.If{
					Cond: expression.NewBinary(expression.VariableVal("err"), token.NEQ, expression.Nil),
					Body: []astgen.ASTStmt{statement.NewReturn(
						expression.NewCallFunction("werror", "WrapWithContextParams",
							expression.VariableVal("ctx"),
							expression.VariableVal("err"),
							expression.StringVal(fmt.Sprintf("type %s failed to convert YAML to JSON", receiverType)),
						),
					)},
				},
				// return o.unmarshalGJSON(ctx, gjson.ParseBytes(data), false)
				statement.NewReturn(expression.NewCallFunction(receiverName, "unmarshalGJSON",
					expression.VariableVal("ctx"),
					expression.NewCallFunction("gjson", "ParseBytes", expression.VariableVal("data")),
					expression.VariableVal("false"))),
			},
		},
	}
}
