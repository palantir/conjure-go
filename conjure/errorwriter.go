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

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/v5/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v5/conjure/transforms"
	"github.com/palantir/conjure-go/v5/conjure/types"
	"github.com/palantir/conjure-go/v5/conjure/visitors"
)

const (
	errorReceiverName    = "e"
	errorInstanceIDField = "errorInstanceID"
	errorInstanceIDParam = "errorInstanceId"
)

const (
	errorsPackagePath  = "github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/errors"
	reflectPackagePath = "reflect"
)

func astForError(errorDefinition spec.ErrorDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	allArgs := make([]spec.FieldDefinition, 0, len(errorDefinition.SafeArgs)+len(errorDefinition.UnsafeArgs))
	allArgs = append(allArgs, errorDefinition.SafeArgs...)
	allArgs = append(allArgs, errorDefinition.UnsafeArgs...)
	decls, err := astForObject(
		spec.ObjectDefinition{
			TypeName: spec.TypeName{
				Name:    transforms.Private(errorDefinition.ErrorName.Name),
				Package: errorDefinition.ErrorName.Package,
			},
			Fields: allArgs,
		},
		info,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate object for error %q parameters",
			errorDefinition.ErrorName.Name,
		)
	}
	var constructorParams []*expression.FuncParam
	var paramToFieldAssignments []astgen.ASTExpr
	for _, fieldDefinition := range allArgs {
		newConjureTypeProvider, err := visitors.NewConjureTypeProvider(fieldDefinition.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create type provider for argument %s for error %s",
				fieldDefinition.FieldName,
				errorDefinition.ErrorName.Name,
			)
		}
		typer, err := newConjureTypeProvider.ParseType(info)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse type argument %s for error %s",
				fieldDefinition.FieldName,
				errorDefinition.ErrorName.Name,
			)
		}
		goType := typer.GoType(info)
		constructorParams = append(constructorParams, &expression.FuncParam{
			Names: []string{argNameTransform(string(fieldDefinition.FieldName))},
			Type:  expression.Type(goType),
		})
		paramToFieldAssignments = append(paramToFieldAssignments, expression.NewKeyValue(
			transforms.Export(string(fieldDefinition.FieldName)),
			expression.VariableVal(argNameTransform(string(fieldDefinition.FieldName))),
		))
	}
	decls = append(decls,
		&decl.Function{
			Name: "New" + errorDefinition.ErrorName.Name,
			FuncType: expression.FuncType{
				Params: constructorParams,
				ReturnTypes: []expression.Type{
					expression.Type(errorDefinition.ErrorName.Name).Pointer(),
				},
			},
			Comment: fmt.Sprintf("New%s returns new instance of %s error.",
				errorDefinition.ErrorName.Name,
				errorDefinition.ErrorName.Name,
			),
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewUnary(token.AND, expression.NewCompositeLit(
						expression.Type(errorDefinition.ErrorName.Name),
						expression.NewKeyValue(
							errorInstanceIDField,
							expression.NewCallFunction("uuid", "NewUUID"),
						),
						expression.NewKeyValue(
							transforms.Private(errorDefinition.ErrorName.Name),
							expression.NewCompositeLit(
								expression.Type(transforms.Private(errorDefinition.ErrorName.Name)),
								paramToFieldAssignments...,
							),
						),
					)),
				),
			},
		},
		decl.NewStruct(
			errorDefinition.ErrorName.Name,
			expression.StructFields{
				&expression.StructField{
					Name: errorInstanceIDField,
					Type: expression.Type("uuid.UUID"),
				},
				&expression.StructField{
					Type: expression.Type(transforms.Private(errorDefinition.ErrorName.Name)),
				},
			},
			errorDefinition.ErrorName.Name+" is an error type.\n\n"+transforms.Documentation(errorDefinition.Docs),
		),
	)
	for _, f := range []errorDeclFunc{
		astErrorErrorMethod,
		astErrorCodeMethod,
		astErrorNameMethod,
		astErrorInstanceIDMethod,
		astErrorParametersMethod,
		astErrorSafeParamsMethod,
		astErrorUnsafeParamsMethod,
		astErrorMarshalJSON,
		astErrorUnmarshalJSON,
	} {
		methodDecl := f(errorDefinition, info)
		decls = append(decls, methodDecl)
	}

	return decls, nil
}

type errorDeclFunc func(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl

// astErrorErrorMethod generates Code function for an error, for example:
//
//  func (e *MyNotFound) Error() string {
//  	return fmt.Sprintf("NOT_FOUND MyNamespace:MyNotFound (%s)", e.errorInstanceID)
//  }
func astErrorErrorMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports("fmt")
	return &decl.Method{
		Function: decl.Function{
			Name: "Error",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"string",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCallFunction("fmt", "Sprintf",
						expression.StringVal(
							fmt.Sprintf("%s %s:%s (%%s)", errorDefinition.Code, errorDefinition.Namespace, errorDefinition.ErrorName.Name),
						),
						expression.NewSelector(
							expression.VariableVal(errorReceiverName),
							errorInstanceIDField,
						),
					),
				),
			},
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

// astErrorCodeMethod generates Code function for an error, for example:
//
//  func (e *MyNotFound) Code() errors.ErrorCode {
//  	return errors.ErrorCodeNotFound
//  }
func astErrorCodeMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	errCode := types.NewGoType("ErrorCode", errorsPackagePath)
	info.AddImports(errCode.ImportPaths()...)
	return &decl.Method{
		Function: decl.Function{
			Name: "Code",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					expression.Type(errCode.GoType(info)),
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					selectorForErrorCode(errorDefinition.Code, info),
				),
			},
			Comment: "Code returns an enum describing error category.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

func selectorForErrorCode(errorCode spec.ErrorCode, info types.PkgInfo) astgen.ASTExpr {
	var varName string
	switch errorCode {
	case spec.ErrorCodePermissionDenied:
		varName = "PermissionDenied"
	case spec.ErrorCodeInvalidArgument:
		varName = "InvalidArgument"
	case spec.ErrorCodeNotFound:
		varName = "NotFound"
	case spec.ErrorCodeConflict:
		varName = "Conflict"
	case spec.ErrorCodeRequestEntityTooLarge:
		varName = "RequestEntityTooLarge"
	case spec.ErrorCodeFailedPrecondition:
		varName = "FailedPrecondition"
	case spec.ErrorCodeInternal:
		varName = "Internal"
	case spec.ErrorCodeTimeout:
		varName = "Timeout"
	case spec.ErrorCodeCustomClient:
		varName = "CustomClient"
	case spec.ErrorCodeCustomServer:
		varName = "CustomServer"
	default:
		panic(fmt.Sprintf(`unknown error code string %q`, errorCode))
	}
	typ := types.NewGoType(varName, errorsPackagePath)
	info.AddImports(typ.ImportPaths()...)
	return expression.VariableVal(typ.GoType(info))
}

// astErrorNameMethod generates Name function for an error, for example:
//
//  func (e *MyNotFound) Name() string {
//  	return "MyNamespace:MyNotFound"
//  }
func astErrorNameMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	return &decl.Method{
		Function: decl.Function{
			Name: "Name",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"string",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.StringVal(fmt.Sprintf("%s:%s", errorDefinition.Namespace, errorDefinition.ErrorName.Name)),
				),
			},
			Comment: "Name returns an error name identifying error type.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

// astErrorInstanceIDMethod generates InstanceID function for an error, for example:
//
//  func (e *MyNotFound) InstanceID() errors.ErrorInstanceID {
//  	return e.errorInstanceID
//  }
func astErrorInstanceIDMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(types.UUID.ImportPaths()...)
	return &decl.Method{
		Function: decl.Function{
			Name: "InstanceID",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					expression.Type(types.UUID.GoType(info)),
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewSelector(
						expression.VariableVal(errorReceiverName),
						errorInstanceIDField,
					),
				),
			},
			Comment: "InstanceID returns unique identifier of this particular error instance.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

// astErrorParametersMethod generates Parameters function for an error, for example:
//
//  func (e *MyNotFound) Parameters() map[string]interface{} {
//  	return map[string]interface{}{"safeArgA": e.safeArgA, "unsafeArgA": e.unsafeArgA}
//  }
func astErrorParametersMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	var keyValues []astgen.ASTExpr
	allArgs := make([]spec.FieldDefinition, 0, len(errorDefinition.SafeArgs)+len(errorDefinition.UnsafeArgs))
	allArgs = append(allArgs, errorDefinition.SafeArgs...)
	allArgs = append(allArgs, errorDefinition.UnsafeArgs...)
	for _, fieldDefinition := range allArgs {
		keyValues = append(keyValues, expression.NewKeyValue(
			fmt.Sprintf("%q", fieldDefinition.FieldName),
			expression.NewSelector(
				expression.VariableVal(errorReceiverName),
				transforms.Export(string(fieldDefinition.FieldName)),
			),
		))
	}
	return &decl.Method{
		Function: decl.Function{
			Name: "Parameters",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"map[string]interface{}",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.Type("map[string]interface{}"),
						keyValues...,
					),
				),
			},
			Comment: "Parameters returns a set of named parameters detailing this particular error instance.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

// astErrorSafeParamsMethod generates SafeParams function for an error, for example:
//
//  func (e *MyNotFound) SafeParams() map[string]interface{} {
//  	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB}
//  }
func astErrorSafeParamsMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	keyValues := make([]astgen.ASTExpr, 0, len(errorDefinition.SafeArgs)+1)
	for _, fieldDefinition := range errorDefinition.SafeArgs {
		keyValues = append(keyValues, expression.NewKeyValue(
			fmt.Sprintf("%q", fieldDefinition.FieldName),
			expression.NewSelector(
				expression.VariableVal(errorReceiverName),
				transforms.Export(string(fieldDefinition.FieldName)),
			),
		))
	}
	keyValues = append(keyValues, expression.NewKeyValue(
		fmt.Sprintf("%q", errorInstanceIDParam),
		expression.NewSelector(
			expression.VariableVal(errorReceiverName),
			errorInstanceIDField,
		),
	))
	return &decl.Method{
		Function: decl.Function{
			Name: "SafeParams",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"map[string]interface{}",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.Type("map[string]interface{}"),
						keyValues...,
					),
				),
			},
			Comment: "SafeParams returns a set of named safe parameters detailing this particular error instance.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

// astErrorUnsafeParamsMethod generates UnsafeParams function for an error, for example:
//
//  func (e *MyNotFound) UnsafeParams() map[string]interface{} {
//  	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA}
//  }
func astErrorUnsafeParamsMethod(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	var keyValues []astgen.ASTExpr
	for _, fieldDefinition := range errorDefinition.UnsafeArgs {
		keyValues = append(keyValues, expression.NewKeyValue(
			fmt.Sprintf("%q", fieldDefinition.FieldName),
			expression.NewSelector(
				expression.VariableVal(errorReceiverName),
				transforms.Export(string(fieldDefinition.FieldName)),
			),
		))
	}
	return &decl.Method{
		Function: decl.Function{
			Name: "UnsafeParams",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"map[string]interface{}",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					expression.NewCompositeLit(
						expression.Type("map[string]interface{}"),
						keyValues...,
					),
				),
			},
			Comment: "UnsafeParams returns a set of named unsafe parameters detailing this particular error instance.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}
}

// astErrorMarshalJSON generates MarshalJSON function for an error, for example:
//
//  func (e *MyNotFound) MarshalJSON() ([]byte, error) {
//    parameters, err := codecs.JSON.Marshal(e.myNotFound)
//    if err != nil {
//      return nil, err
//    }
//    return codecs.JSON.Marshal(errors.SerializableError{
//      ErrorCode: errors.NotFound,
//      ErrorName: "MyNamespace:MyNotFound",
//      ErrorInstanceID: e.errorInstanceID,
//      Parameters: json.RawMessage(parameters),
//    })
//  }
func astErrorMarshalJSON(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	serError := types.NewGoType("SerializableError", errorsPackagePath)
	info.AddImports(serError.ImportPaths()...)
	jsonMessage := types.NewGoType("RawMessage", "encoding/json")
	info.AddImports(jsonMessage.ImportPaths()...)
	info.AddImports(types.SafeJSONMarshal.ImportPaths()...)

	return newMarshalJSONMethod(errorReceiverName, errorDefinition.ErrorName.Name,
		&statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.VariableVal("parameters"),
				expression.VariableVal("err"),
			},
			Tok: token.DEFINE,

			RHS: expression.NewCallExpression(expression.Type(types.SafeJSONMarshal.GoType(info)),
				expression.NewSelector(
					expression.VariableVal(errorReceiverName),
					transforms.Private(errorDefinition.ErrorName.Name),
				),
			),
		},
		ifErrNotNilReturnHelper(true, "nil", "err", nil),
		statement.NewReturn(
			&expression.CallExpression{
				Function: expression.Type(types.SafeJSONMarshal.GoType(info)),
				Args: []astgen.ASTExpr{
					expression.NewCompositeLit(expression.Type(serError.GoType(info)),
						expression.NewKeyValue(
							"ErrorCode",
							selectorForErrorCode(errorDefinition.Code, info),
						),
						expression.NewKeyValue(
							"ErrorName",
							expression.StringVal(
								fmt.Sprintf("%s:%s", errorDefinition.Namespace, errorDefinition.ErrorName.Name)),
						),
						expression.NewKeyValue(
							"ErrorInstanceID",
							expression.NewSelector(
								expression.VariableVal(errorReceiverName),
								errorInstanceIDField,
							),
						),
						expression.NewKeyValue(
							"Parameters",
							&expression.CallExpression{
								Function: expression.Type(jsonMessage.GoType(info)),
								Args: []astgen.ASTExpr{
									expression.VariableVal("parameters"),
								},
							},
						),
					),
				},
			},
		),
	)
}

// astErrorUnmarshalJSON generates UnmarshalJSON function for an error, for example:
//
//  func (e *MyNotFound) UnmarshalJSON(data []byte) error {
//    var serializableError errors.SerializableError
//    if err := codecs.JSON.Unmarshal(data, &serializableError); err != nil {
//      return err
//    }
//    var parameters myNotFound
//    if err := codecs.JSON.Unmarshal([]byte(serializableError.Parameters), &parameters); err != nil {
//      return err
//    }
//    e.errorInstanceID = serializableError.ErrorInstanceID
//    e.myNotFound = parameters
//    return nil
//  }
func astErrorUnmarshalJSON(errorDefinition spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	serError := types.NewGoType("SerializableError", errorsPackagePath)
	info.AddImports(serError.ImportPaths()...)
	info.AddImports(types.SafeJSONUnmarshal.ImportPaths()...)
	return newUnmarshalJSONMethod(errorReceiverName, errorDefinition.ErrorName.Name,
		statement.NewDecl(
			decl.NewVar(
				"serializableError",
				expression.Type(serError.GoType(info)),
			),
		),
		ifErrNotNilReturnErrStatement("err",
			statement.NewAssignment(
				expression.VariableVal("err"),
				token.DEFINE,
				&expression.CallExpression{
					Function: expression.Type(types.SafeJSONUnmarshal.GoType(info)),
					Args: []astgen.ASTExpr{
						expression.VariableVal(dataVarName),
						expression.NewUnary(token.AND, expression.VariableVal("serializableError")),
					},
				},
			),
		),
		statement.NewDecl(
			decl.NewVar(
				"parameters",
				expression.Type(transforms.Private(errorDefinition.ErrorName.Name)),
			),
		),
		ifErrNotNilReturnErrStatement("err",
			statement.NewAssignment(
				expression.VariableVal("err"),
				token.DEFINE,
				expression.NewCallExpression(expression.Type(types.SafeJSONUnmarshal.GoType(info)),
					expression.NewCallExpression(expression.ByteSliceType,
						expression.NewSelector(
							expression.VariableVal("serializableError"),
							"Parameters",
						),
					),
					expression.NewUnary(token.AND, expression.VariableVal("parameters")),
				),
			),
		),
		&statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.NewSelector(
					expression.VariableVal(errorReceiverName),
					errorInstanceIDField,
				),
			},
			Tok: token.ASSIGN,
			RHS: expression.NewSelector(
				expression.VariableVal("serializableError"),
				"ErrorInstanceID",
			),
		},
		&statement.Assignment{
			LHS: []astgen.ASTExpr{
				expression.NewSelector(
					expression.VariableVal(errorReceiverName),
					transforms.Private(errorDefinition.ErrorName.Name),
				),
			},
			Tok: token.ASSIGN,
			RHS: expression.VariableVal("parameters"),
		},
		statement.NewReturn(expression.Nil),
	)
}

// astErrorInitFunc generates init func that registers each error type in the conjure-go-runtime
// error type registry, for example:
//
// func init() {
//     errors.RegisterErrorType("MyNamespace:MyNotFound", reflect.TypeOf(MyNotFound{}))
// }
func astErrorInitFunc(errorDefinitions []spec.ErrorDefinition, info types.PkgInfo) astgen.ASTDecl {
	info.AddImports(reflectPackagePath)
	stmts := make([]astgen.ASTStmt, 0, len(errorDefinitions))
	for _, def := range errorDefinitions {
		stmts = append(stmts, &statement.Expression{
			Expr: expression.NewCallFunction("errors", "RegisterErrorType",
				expression.StringVal(fmt.Sprintf("%s:%s", def.Namespace, def.ErrorName.Name)),
				expression.NewCallFunction("reflect", "TypeOf",
					expression.NewCompositeLit(expression.Type(def.ErrorName.Name)),
				))})
	}

	return &decl.Function{
		Name: "init",
		Body: stmts,
	}
}
