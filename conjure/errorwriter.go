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

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/transforms"
	"github.com/palantir/conjure-go/conjure/types"
	"github.com/palantir/conjure-go/conjure/visitors"
)

const (
	errorReceiverName    = "e"
	errorInstanceIDField = "errorInstanceID"
)

const (
	codecsPackagePath = "github.com/palantir/conjure-go-runtime/conjure-go-contract/codecs"
	errorsPackagePath = "github.com/palantir/conjure-go-runtime/conjure-go-contract/errors"
	uuidPackagePath   = "github.com/palantir/conjure-go-runtime/conjure-go-contract/uuid"
)

func astForError(errorDefinition spec.ErrorDefinition, customTypes types.CustomConjureTypes, goPkgImportPath string, importToAlias map[string]string) ([]astgen.ASTDecl, StringSet, error) {
	allArgs := make([]spec.FieldDefinition, 0, len(errorDefinition.SafeArgs)+len(errorDefinition.UnsafeArgs))
	allArgs = append(allArgs, errorDefinition.SafeArgs...)
	allArgs = append(allArgs, errorDefinition.UnsafeArgs...)
	decls, imports, err := astForObject(
		spec.ObjectDefinition{
			TypeName: spec.TypeName{
				Name:    transforms.Private(errorDefinition.ErrorName.Name),
				Package: errorDefinition.ErrorName.Package,
			},
			Fields: allArgs,
		},
		customTypes,
		goPkgImportPath,
		importToAlias,
	)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to generate object for error %q parameters",
			errorDefinition.ErrorName.Name,
		)
	}
	imports.AddAll(NewStringSet(
		errorsPackagePath,
		uuidPackagePath,
	))
	var constructorParams []*expression.FuncParam
	var paramToFieldAssignments []astgen.ASTExpr
	for _, fieldDefinition := range allArgs {
		newConjureTypeProvider, err := visitors.NewConjureTypeProvider(fieldDefinition.Type)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create type provider for argument %s for error %s",
				fieldDefinition.FieldName,
				errorDefinition.ErrorName.Name,
			)
		}
		typer, err := newConjureTypeProvider.ParseType(customTypes)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to parse type argument %s for error %s",
				fieldDefinition.FieldName,
				errorDefinition.ErrorName.Name,
			)
		}
		goType := typer.GoType(goPkgImportPath, importToAlias)
		constructorParams = append(constructorParams, &expression.FuncParam{
			Names: []string{string(fieldDefinition.FieldName)},
			Type:  expression.Type(goType),
		})
		paramToFieldAssignments = append(paramToFieldAssignments, expression.NewKeyValue(
			transforms.Export(string(fieldDefinition.FieldName)),
			expression.VariableVal(string(fieldDefinition.FieldName)),
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
		astErrorMarshalJSON,
		astErrorUnmarshalJSON,
	} {
		decl, currImports := f(errorDefinition)
		decls = append(decls, decl)
		imports.AddAll(currImports)
	}

	decls = append(decls)
	return decls, imports, nil
}

type errorDeclFunc func(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet)

// astErrorErrorMethod generates Code function for an error, for example:
//
//  func (e *MyNotFound) Error() string {
//  	return fmt.Sprintf("NOT_FOUND MyNamespace:MyNotFound (%s)", e.errorInstanceID)
//  }
func astErrorErrorMethod(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
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
	}, NewStringSet("fmt")
}

// astErrorCodeMethod generates Code function for an error, for example:
//
//  func (e *MyNotFound) Code() errors.ErrorCode {
//  	return errors.ErrorCodeNotFound
//  }
func astErrorCodeMethod(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
	return &decl.Method{
		Function: decl.Function{
			Name: "Code",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"errors.ErrorCode",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewReturn(
					selectorForErrorCode(errorDefinition.Code),
				),
			},
			Comment: "Code returns an enum describing error category.",
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}, NewStringSet(errorsPackagePath)
}

func selectorForErrorCode(errorCode spec.ErrorCode) *expression.Selector {
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
	return expression.NewSelector(
		expression.VariableVal("errors"),
		varName,
	)
}

// astErrorNameMethod generates Name function for an error, for example:
//
//  func (e *MyNotFound) Name() string {
//  	return "MyNamespace:MyNotFound"
//  }
func astErrorNameMethod(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
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
	}, nil
}

// astErrorInstanceIDMethod generates InstanceID function for an error, for example:
//
//  func (e *MyNotFound) InstanceID() errors.ErrorInstanceID {
//  	return e.errorInstanceID
//  }
func astErrorInstanceIDMethod(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
	return &decl.Method{
		Function: decl.Function{
			Name: "InstanceID",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"uuid.UUID",
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
	}, NewStringSet(uuidPackagePath)
}

// astErrorParametersMethod generates Parameters function for an error, for example:
//
//  func (e *MyNotFound) Parameters() map[string]interface{} {
//  	return map[string]interface{}{"safeArgA": e.safeArgA, "unsafeArgA": e.unsafeArgA}
//  }
func astErrorParametersMethod(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
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
	}, nil
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
func astErrorMarshalJSON(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
	return &decl.Method{
		Function: decl.Function{
			Name: "MarshalJSON",
			FuncType: expression.FuncType{
				ReturnTypes: []expression.Type{
					"[]byte",
					"error",
				},
			},
			Body: []astgen.ASTStmt{
				&statement.Assignment{
					LHS: []astgen.ASTExpr{
						expression.VariableVal("parameters"),
						expression.VariableVal("err"),
					},
					Tok: token.DEFINE,
					RHS: &expression.CallExpression{
						Function: expression.NewSelector(
							expression.VariableVal("codecs.JSON"),
							"Marshal",
						),
						Args: []astgen.ASTExpr{
							expression.NewSelector(
								expression.VariableVal(errorReceiverName),
								transforms.Private(errorDefinition.ErrorName.Name),
							),
						},
					},
				},
				&statement.If{
					Cond: &expression.Binary{
						LHS: expression.VariableVal("err"),
						Op:  token.NEQ,
						RHS: expression.Nil,
					},
					Body: []astgen.ASTStmt{
						&statement.Return{
							Values: []astgen.ASTExpr{
								expression.Nil,
								expression.VariableVal("err"),
							},
						},
					},
				},
				statement.NewReturn(
					expression.NewCallFunction("codecs.JSON", "Marshal",
						expression.NewCompositeLit(
							expression.NewSelector(
								expression.VariableVal("errors"),
								"SerializableError",
							),
							expression.NewKeyValue(
								"ErrorCode",
								selectorForErrorCode(errorDefinition.Code),
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
									Function: expression.NewSelector(
										expression.VariableVal("json"),
										"RawMessage",
									),
									Args: []astgen.ASTExpr{
										expression.VariableVal("parameters"),
									},
								},
							),
						),
					),
				),
			},
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}, NewStringSet("encoding/json", errorsPackagePath, codecsPackagePath)
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
func astErrorUnmarshalJSON(errorDefinition spec.ErrorDefinition) (astgen.ASTDecl, StringSet) {
	return &decl.Method{
		Function: decl.Function{
			Name: "UnmarshalJSON",
			FuncType: expression.FuncType{
				Params: []*expression.FuncParam{
					expression.NewFuncParam(
						"data",
						expression.Type("[]byte"),
					),
				},
				ReturnTypes: []expression.Type{
					"error",
				},
			},
			Body: []astgen.ASTStmt{
				statement.NewDecl(
					decl.NewVar(
						"serializableError",
						expression.Type("errors.SerializableError"),
					),
				),
				ifErrNotNilReturnErrStatement("err",
					statement.NewAssignment(
						expression.VariableVal("err"),
						token.DEFINE,
						expression.NewCallFunction(
							"codecs.JSON",
							"Unmarshal",
							expression.VariableVal("data"),
							expression.NewUnary(token.AND, expression.VariableVal("serializableError")),
						),
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
						expression.NewCallFunction(
							"codecs.JSON",
							"Unmarshal",
							&expression.CallExpression{
								Function: expression.Type("[]byte"),
								Args: []astgen.ASTExpr{
									expression.NewSelector(
										expression.VariableVal("serializableError"),
										"Parameters",
									),
								},
							},
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
			},
		},
		ReceiverName: errorReceiverName,
		ReceiverType: expression.Type(errorDefinition.ErrorName.Name).Pointer(),
	}, NewStringSet("encoding/json", errorsPackagePath, codecsPackagePath)
}
