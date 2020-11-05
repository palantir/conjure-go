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

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"

	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	dataVarName = "data"
)

func newStringMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType),
		Function: decl.Function{
			Name: "String",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.StringType},
			},
			Body: body,
		},
	}
}

func newMarshalTextMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType),
		Function: decl.Function{
			Name: "MarshalText",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.ByteSliceType, expression.ErrorType},
			},
			Body: body,
		},
	}
}

func newUnmarshalTextMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType).Pointer(),
		Function: decl.Function{
			Name: "UnmarshalText",
			FuncType: expression.FuncType{
				Params:      expression.FuncParams{expression.NewFuncParam(dataVarName, expression.ByteSliceType)},
				ReturnTypes: []expression.Type{expression.ErrorType},
			},
			Body: body,
		},
	}
}

func newMarshalJSONMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType),
		Function: decl.Function{
			Name: "MarshalJSON",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.ByteSliceType, expression.ErrorType},
			},
			Body: body,
		},
	}
}

func newUnmarshalJSONMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType).Pointer(),
		Function: decl.Function{
			Name: "UnmarshalJSON",
			FuncType: expression.FuncType{
				Params:      expression.FuncParams{expression.NewFuncParam(dataVarName, expression.ByteSliceType)},
				ReturnTypes: []expression.Type{expression.ErrorType},
			},
			Body: body,
		},
	}
}

/*
func (o Foo) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(o)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
*/
func newMarshalYAMLMethod(receiverName, receiverType string, info types.PkgInfo) *decl.Method {
	info.AddImports("github.com/palantir/pkg/safejson", "github.com/palantir/pkg/safeyaml")
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType),
		Function: decl.Function{
			Name: "MarshalYAML",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.EmptyInterfaceType, expression.ErrorType},
			},
			Body: []astgen.ASTStmt{
				&statement.Assignment{
					LHS: []astgen.ASTExpr{expression.VariableVal("jsonBytes"), expression.VariableVal("err")},
					Tok: token.DEFINE,
					RHS: expression.NewCallFunction("safejson", "Marshal", expression.VariableVal(receiverName)),
				},
				ifErrNotNilReturnHelper(true, "nil", "err", nil),
				statement.NewReturn(expression.NewCallFunction("safeyaml", "JSONtoYAMLMapSlice", expression.VariableVal("jsonBytes"))),
			},
		},
	}
}

/*
func (o *Foo) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&o)
}
*/
func newUnmarshalYAMLMethod(receiverName, receiverType string, info types.PkgInfo) *decl.Method {
	info.AddImports("github.com/palantir/pkg/safejson", "github.com/palantir/pkg/safeyaml")
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType).Pointer(),
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
				&statement.Assignment{
					LHS: []astgen.ASTExpr{expression.VariableVal("jsonBytes"), expression.VariableVal("err")},
					Tok: token.DEFINE,
					RHS: expression.NewCallFunction("safeyaml", "UnmarshalerToJSONBytes", expression.VariableVal("unmarshal")),
				},
				ifErrNotNilReturnErrStatement("err", nil),
				statement.NewReturn(
					expression.NewCallFunction("safejson", "Unmarshal",
						expression.VariableVal("jsonBytes"),
						expression.NewStar(expression.NewUnary(token.AND, expression.VariableVal(receiverName))))),
			},
		},
	}
}
