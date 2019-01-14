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
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
)

const (
	dataVarName = "data"
)

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

func newMarshalYAMLMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
	return &decl.Method{
		ReceiverName: receiverName,
		ReceiverType: expression.Type(receiverType),
		Function: decl.Function{
			Name: "MarshalYAML",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{expression.EmptyInterfaceType, expression.ErrorType},
			},
			Body: body,
		},
	}
}

func newUnmarshalYAMLMethod(receiverName, receiverType string, body ...astgen.ASTStmt) *decl.Method {
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
			Body: body,
		},
	}
}
