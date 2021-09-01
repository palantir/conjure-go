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

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	errorReceiverName    = "e"
	errorInstanceIDField = "errorInstanceID"
	causeField           = "cause"
	stackField           = "stack"
	errorInstanceIDParam = "errorInstanceId"
	errorNameParam       = "errorName"
)

func writeErrorType(file *jen.Group, def *types.ErrorDefinition, cfg OutputConfiguration) {
	astErrorInternalStructType(file, def, cfg)
	astErrorConstructorFuncs(file, def)
	astErrorExportedStructType(file, def)
	astIsErrorTypeFunc(file, def)
	astErrorErrorMethod(file, def)
	astErrorCauseMethod(file, def)
	astErrorStackTraceMethod(file, def)
	astErrorMessageMethod(file, def)
	astErrorFormatMethod(file, def)
	astErrorCodeMethod(file, def)
	astErrorNameMethod(file, def)
	astErrorInstanceIDMethod(file, def)
	astErrorParametersMethod(file, def)
	astErrorHelperSafeParamsMethod(file, def)
	astErrorSafeParamsMethod(file, def)
	astErrorHelperUnsafeParamsMethod(file, def)
	astErrorUnsafeParamsMethod(file, def)
	astErrorMarshalJSON(file, def)
	astErrorUnmarshalJSON(file, def)
}

// Create private *myInternal object containing known params.
func astErrorInternalStructType(file *jen.Group, def *types.ErrorDefinition, cfg OutputConfiguration) {
	allArgs := append(append([]*types.Field{}, def.SafeArgs...), def.UnsafeArgs...)
	// Use object generator to create a struct implementing JSON encoding for the error.
	writeObjectType(file, &types.ObjectType{Name: transforms.Private(def.Name), Fields: allArgs}, cfg)
}

// Declare New and Wrap constructors
func astErrorConstructorFuncs(file *jen.Group, def *types.ErrorDefinition) {
	allArgs := append(append([]*types.Field{}, def.SafeArgs...), def.UnsafeArgs...)
	// Declare New and Wrap constructors
	constructorParams := func(params *jen.Group, includeErr bool) {
		if includeErr {
			params.Err().Error()
		}
		for _, fieldDef := range allArgs {
			params.Id(argNameTransform(fieldDef.Name)).Add(fieldDef.Type.Code())
		}
	}
	paramToFieldAssignments := func(assignments *jen.Group) {
		for _, fieldDef := range allArgs {
			assignments.Id(transforms.Export(fieldDef.Name)).Op(":").Id(argNameTransform(fieldDef.Name))
		}
	}
	newStructLiteralValues := func(values *jen.Group, includeCause bool) {
		values.Id(errorInstanceIDField).Op(":").Add(snip.UUIDNewUUID()).Call()
		values.Id(stackField).Op(":").Add(snip.WerrorNewStackTrace()).Call()
		if includeCause {
			values.Id(causeField).Op(":").Err()
		}
		values.Id(transforms.Private(def.Name)).Op(":").Id(transforms.Private(def.Name)).ValuesFunc(paramToFieldAssignments)
	}

	file.Commentf("New%s returns new instance of %s error.", def.Name, def.Name)
	file.Func().
		Id("New" + def.Name).
		ParamsFunc(func(params *jen.Group) {
			constructorParams(params, false)
		}).
		Params(jen.Op("*").Id(def.Name)).
		Block(jen.Return(jen.Op("&").Id(def.Name).ValuesFunc(func(values *jen.Group) {
			newStructLiteralValues(values, false)
		})))

	file.Commentf("WrapWith%s returns new instance of %s error wrapping an existing error.", def.Name, def.Name)
	file.Func().
		Id("WrapWith" + def.Name).
		ParamsFunc(func(params *jen.Group) {
			constructorParams(params, true)
		}).
		Params(jen.Op("*").Id(def.Name)).
		Block(jen.Return(jen.Op("&").Id(def.Name).ValuesFunc(func(values *jen.Group) {
			newStructLiteralValues(values, true)
		})))

}

func astErrorExportedStructType(file *jen.Group, def *types.ErrorDefinition) {
	file.Commentf("%s is an error type.", def.Name)
	file.Add(def.Docs.CommentLine()).Type().Id(def.Name).Struct(
		jen.Id(errorInstanceIDField).Add(types.UUID{}.Code()),
		jen.Id(transforms.Private(def.Name)),
		jen.Id(causeField).Error(),
		jen.Id(stackField).Add(snip.WerrorStackTrace()),
	)
}

// astErrorErrorMethod generates Code function for an error, for example:
//
//  func (e *MyNotFound) Error() string {
//  	return fmt.Sprintf("NOT_FOUND MyNamespace:MyNotFound (%s)", e.errorInstanceID)
//  }
func astErrorErrorMethod(file *jen.Group, def *types.ErrorDefinition) {
	errFmt := fmt.Sprintf("%s %s:%s (%%s)", def.ErrorCode, def.ErrorNamespace, def.Name)
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Error").
		Params().
		Params(jen.String()).
		Block(
			jen.Return(snip.FmtSprintf().Call(jen.Lit(errFmt), jen.Id(errorReceiverName).Dot(errorInstanceIDField))),
		)
}

// astErrorCodeMethod generates Code function for an error, for example:
//
//  // Code returns an enum describing error category.
//  func (e *MyNotFound) Code() errors.ErrorCode {
//  	return errors.ErrorCodeNotFound
//  }
func astErrorCodeMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("Code returns an enum describing error category.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Code").
		Params().
		Params(snip.CGRErrorsErrorCode()).
		Block(
			jen.Return(selectorForErrorCode(def.ErrorCode)),
		)
}

// astErrorCauseMethod generates Cause function for an error, for example:
//
//  func (e *MyNotFound) Cause() error {
//  	return e.cause
//  }
func astErrorCauseMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("Cause returns the underlying cause of the error, or nil if none.")
	file.Comment("Note that cause is not serialized and sent over the wire.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Cause").
		Params().
		Params(jen.Error()).
		Block(
			jen.Return(jen.Id(errorReceiverName).Dot(causeField)),
		)
}

// astErrorStackTraceMethod generates StackTrace function for an error, for example:
//
//  func (e *MyNotFound) StackTrace() werror.StackTrace {
//  	return e.stack
//  }
func astErrorStackTraceMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("StackTrace returns the StackTrace for the error, or nil if none.")
	file.Comment("Note that stack traces are not serialized and sent over the wire.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("StackTrace").
		Params().
		Params(snip.WerrorStackTrace()).
		Block(
			jen.Return(jen.Id(errorReceiverName).Dot(stackField)),
		)
}

// astErrorMessageMethod generates Message function for an error, for example:
//
//  func (e *MyNotFound) Message() string {
//  	return "NOT_FOUND MyNamespace:MyNotFound"
//  }
func astErrorMessageMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("Message returns the message body for the error.")
	message := fmt.Sprintf("%s %s:%s", def.ErrorCode, def.ErrorNamespace, def.Name)
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Message").
		Params().
		Params(jen.String()).
		Block(
			jen.Return(jen.Lit(message)),
		)
}

// astErrorFormatMethod generates Format function for an error, for example:
//
//  func (e *MyNotFound) Format(state fmt.State, verb rune) {
//  	werror.Format(e, state, verb)
//  }
func astErrorFormatMethod(file *jen.Group, def *types.ErrorDefinition) {
	const (
		stateParam = "state"
		verbParam  = "verb"
		safeVar    = "safeParams"
	)
	file.Comment("Format implements fmt.Formatter, a requirement of werror.Werror.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Format").
		Params(jen.Id(stateParam).Qual("fmt", "State"), jen.Id(verbParam).Rune()).
		Params().
		Block(
			snip.WerrorFormat().Call(
				jen.Id(errorReceiverName),
				jen.Id(errorReceiverName).Dot(safeVar).Call(),
				jen.Id(stateParam),
				jen.Id(verbParam),
			),
		)
}

// astErrorNameMethod generates Name function for an error, for example:
//
//  func (e *MyNotFound) Name() string {
//  	return "MyNamespace:MyNotFound"
//  }
func astErrorNameMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("Name returns an error name identifying error type.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Name").
		Params().
		Params(jen.String()).
		Block(
			jen.Return(jen.Lit(fmt.Sprintf("%s:%s", def.ErrorNamespace, def.Name))),
		)
}

// astErrorInstanceIDMethod generates InstanceID function for an error, for example:
//
//  func (e *MyNotFound) InstanceID() errors.ErrorInstanceID {
//  	return e.errorInstanceID
//  }
func astErrorInstanceIDMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("InstanceID returns unique identifier of this particular error instance.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("InstanceID").
		Params().
		Params(types.UUID{}.Code()).
		Block(
			jen.Return(jen.Id(errorReceiverName).Dot(errorInstanceIDField)),
		)
}

// astErrorParametersMethod generates Parameters function for an error, for example:
//
//  func (e *MyNotFound) Parameters() map[string]interface{} {
//  	return map[string]interface{}{"safeArgA": e.safeArgA, "unsafeArgA": e.unsafeArgA}
//  }
func astErrorParametersMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("Parameters returns a set of named parameters detailing this particular error instance.")
	allArgs := append(append([]*types.Field{}, def.SafeArgs...), def.UnsafeArgs...)
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("Parameters").
		Params().
		Params(jen.Map(jen.String()).Interface()).
		Block(
			jen.Return(jen.Map(jen.String()).Interface().ValuesFunc(func(values *jen.Group) {
				for _, fieldDef := range allArgs {
					values.Lit(fieldDef.Name).Op(":").Id(errorReceiverName).Dot(transforms.Export(fieldDef.Name))
				}
			})),
		)
}

// astErrorSafeParamsMethod generates SafeParams function for an error, for example:
//
//  func (e *MyNotFound) SafeParams() map[string]interface{} {
//  	safeParams, _ := werror.ParamsFromError(e.cause)
//      for k, v := range e.safeParams() {
//          if _, exists := safeParams[k]; !exists {
//              safeParams[k] = v
//          }
//      }
//      return safeParams
//  }
func astErrorSafeParamsMethod(file *jen.Group, def *types.ErrorDefinition) {
	const (
		k, v, exists, safeParams = "k", "v", "exists", "safeParams"
	)
	file.Comment("SafeParams returns a set of named safe parameters detailing this particular error instance and")
	file.Comment("any underlying causes.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("SafeParams").
		Params().
		Params(jen.Map(jen.String()).Interface()).
		Block(
			jen.List(jen.Id(safeParams), jen.Id("_")).Op(":=").
				Add(snip.WerrorParamsFromError()).Call(jen.Id(errorReceiverName).Dot(causeField)),
			jen.For(jen.List(jen.Id(k), jen.Id(v)).Op(":=").
				Range().Id(errorReceiverName).Dot(safeParams).Call()).
				Block(
					jen.If(
						jen.List(jen.Id("_"), jen.Id(exists)).Op(":=").Id(safeParams).Index(jen.Id(k)),
						jen.Op("!").Id(exists),
					).Block(
						jen.Id(safeParams).Index(jen.Id(k)).Op("=").Id(v),
					),
				),
			jen.Return(jen.Id(safeParams)),
		)
}

// astErrorHelperSafeParamsMethod generates safeParams function for an error, for example:
//
//  func (e *MyNotFound) safeParams() map[string]interface{} {
//  	return map[string]interface{}{"safeArgA": e.SafeArgA, "safeArgB": e.SafeArgB}
//  }
func astErrorHelperSafeParamsMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("safeParams returns a set of named safe parameters detailing this particular error instance.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("safeParams").
		Params().
		Params(jen.Map(jen.String()).Interface()).
		Block(
			jen.Return(jen.Map(jen.String()).Interface().ValuesFunc(func(values *jen.Group) {
				for _, safeArg := range def.SafeArgs {
					values.Lit(safeArg.Name).Op(":").Id(errorReceiverName).Dot(transforms.Export(safeArg.Name))
				}
				values.Lit(errorInstanceIDParam).Op(":").Id(errorReceiverName).Dot(errorInstanceIDField)
				values.Lit(errorNameParam).Op(":").Id(errorReceiverName).Dot("Name").Call()
			})),
		)
}

// astErrorUnsafeParamsMethod generates SafeParams function for an error, for example:
//
//  func (e *MyNotFound) UnsafeParams() map[string]interface{} {
//  	_, unsafeParams := werror.ParamsFromError(e.cause)
//      for k, v := range e.unsafeParams() {
//          if _, exists := unsafeParams[k]; !exists {
//              unsafeParams[k] = v
//          }
//      }
//      return unsafeParams
//  }
func astErrorUnsafeParamsMethod(file *jen.Group, def *types.ErrorDefinition) {
	const (
		k, v, exists, unsafeParams = "k", "v", "exists", "unsafeParams"
	)
	file.Comment("UnsafeParams returns a set of named unsafe parameters detailing this particular error instance and")
	file.Comment("any underlying causes.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("UnsafeParams").
		Params().
		Params(jen.Map(jen.String()).Interface()).
		Block(
			jen.List(jen.Id("_"), jen.Id(unsafeParams)).Op(":=").
				Add(snip.WerrorParamsFromError()).Call(jen.Id(errorReceiverName).Dot(causeField)),
			jen.For(jen.List(jen.Id(k), jen.Id(v)).Op(":=").
				Range().Id(errorReceiverName).Dot(unsafeParams).Call()).
				Block(
					jen.If(
						jen.List(jen.Id("_"), jen.Id(exists)).Op(":=").Id(unsafeParams).Index(jen.Id(k)),
						jen.Op("!").Id(exists),
					).Block(
						jen.Id(unsafeParams).Index(jen.Id(k)).Op("=").Id(v),
					),
				),
			jen.Return(jen.Id(unsafeParams)),
		)
}

// astErrorHelperUnsafeParamsMethod generates unsafeParams function for an error, for example:
//
//  func (e *MyNotFound) unsafeParams() map[string]interface{} {
//  	return map[string]interface{}{"unsafeArgA": e.UnsafeArgA, "unsafeArgB": e.UnsafeArgB}
//  }
func astErrorHelperUnsafeParamsMethod(file *jen.Group, def *types.ErrorDefinition) {
	file.Comment("unsafeParams returns a set of named unsafe parameters detailing this particular error instance.")
	file.Func().
		Params(jen.Id(errorReceiverName).Op("*").Id(def.Name)).
		Id("unsafeParams").
		Params().
		Params(jen.Map(jen.String()).Interface()).
		Block(
			jen.Return(jen.Map(jen.String()).Interface().ValuesFunc(func(values *jen.Group) {
				for _, unsafeArg := range def.UnsafeArgs {
					values.Lit(unsafeArg.Name).Op(":").Id(errorReceiverName).Dot(transforms.Export(unsafeArg.Name))
				}
			})),
		)
}

// astErrorMarshalJSON generates MarshalJSON function for an error, for example:
//
//  func (e *MyNotFound) MarshalJSON() ([]byte, error) {
//    parameters, err := safejson.Marshal(e.myNotFound)
//    if err != nil {
//      return nil, err
//    }
//    return safejson.Marshal(errors.SerializableError{
//      ErrorCode: errors.NotFound,
//      ErrorName: "MyNamespace:MyNotFound",
//      ErrorInstanceID: e.errorInstanceID,
//      Parameters: json.RawMessage(parameters),
//    })
//  }
func astErrorMarshalJSON(file *jen.Group, def *types.ErrorDefinition) {
	file.Add(snip.MethodMarshalJSON(errorReceiverName, def.Name)).Block(
		jen.List(jen.Id("parameters"), jen.Err()).Op(":=").
			Add(snip.SafeJSONMarshal()).Call(jen.Id(errorReceiverName).Dot(transforms.Private(def.Name))),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
		jen.Return(snip.SafeJSONMarshal().Call(snip.CGRErrorsSerializableError().Values(
			jen.Id("ErrorCode").Op(":").Add(selectorForErrorCode(def.ErrorCode)),
			jen.Id("ErrorName").Op(":").Lit(fmt.Sprintf("%s:%s", def.ErrorNamespace, def.Name)),
			jen.Id("ErrorInstanceID").Op(":").Id(errorReceiverName).Dot(errorInstanceIDField),
			jen.Id("Parameters").Op(":").Qual("encoding/json", "RawMessage").Call(jen.Id("parameters")),
		))),
	)
}

// astErrorUnmarshalJSON generates UnmarshalJSON function for an error, for example:
//
//  func (e *MyNotFound) UnmarshalJSON(data []byte) error {
//    var serializableError errors.SerializableError
//    if err := safejson.Unmarshal(data, &serializableError); err != nil {
//      return err
//    }
//    var parameters myNotFound
//    if err := safejson.Unmarshal([]byte(serializableError.Parameters), &parameters); err != nil {
//      return err
//    }
//    e.errorInstanceID = serializableError.ErrorInstanceID
//    e.myNotFound = parameters
//    return nil
//  }
func astErrorUnmarshalJSON(file *jen.Group, def *types.ErrorDefinition) {
	file.Add(snip.MethodUnmarshalJSON(errorReceiverName, def.Name)).Block(
		jen.Var().Id("serializableError").Add(snip.CGRErrorsSerializableError()),
		jen.If(
			jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal()).Call(jen.Id(dataVarName), jen.Op("&").Id("serializableError")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		jen.Var().Id("parameters").Id(transforms.Private(def.Name)),
		jen.If(
			jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal()).Call(
				jen.Id("[]byte").Call(jen.Id("serializableError").Dot("Parameters")),
				jen.Op("&").Id("parameters")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		jen.Id(errorReceiverName).Dot(errorInstanceIDField).Op("=").Id("serializableError").Dot("ErrorInstanceID"),
		jen.Id(errorReceiverName).Dot(transforms.Private(def.Name)).Op("=").Id("parameters"),
		jen.Return(jen.Nil()),
	)
}

// astErrorInitFunc generates init func that registers each error type in the conjure-go-runtime
// error type registry, for example:
//
// func init() {
//     errors.RegisterErrorType("MyNamespace:MyInternal", reflect.TypeOf(MyInternal{}))
//     errors.RegisterErrorType("MyNamespace:MyNotFound", reflect.TypeOf(MyNotFound{}))
// }
func astErrorInitFunc(file *jen.Group, defs []*types.ErrorDefinition) {
	file.Func().Id("init").Params().BlockFunc(func(funcBody *jen.Group) {
		for _, def := range defs {
			funcBody.Add(snip.CGRErrorsRegisterErrorType()).Call(
				jen.Lit(fmt.Sprintf("%s:%s", def.ErrorNamespace, def.Name)),
				snip.ReflectTypeOf().Call(jen.Id(def.Name).Values()),
			)
		}
	})
}

// astIsErrorTypeFunc generates a helper func that checks whether an error is of the error type:
//
// func IsMyNotFound(err error) bool {
//	   if err == nil {
//		   return false
//	   }
//	   _, ok := errors.GetConjureError(err).(*MyNotFound)
//	   return ok
// }
func astIsErrorTypeFunc(file *jen.Group, def *types.ErrorDefinition) {
	name := "Is" + def.Name
	file.Commentf("%s returns true if err is an instance of %s.", name, def.Name)
	file.Func().Id(name).Params(jen.Err().Error()).Params(jen.Bool()).Block(
		jen.If(jen.Err().Op("==").Nil()).Block(jen.Return(jen.False())),
		jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Add(snip.CGRErrorsGetConjureError()).Call(jen.Err()).Assert(jen.Op("*").Id(def.Name)),
		jen.Return(jen.Id("ok")),
	)
}

func selectorForErrorCode(errorCode spec.ErrorCode) *jen.Statement {
	switch errorCode.Value() {
	case spec.ErrorCode_PERMISSION_DENIED:
		return snip.CGRErrorsPermissionDenied()
	case spec.ErrorCode_INVALID_ARGUMENT:
		return snip.CGRErrorsInvalidArgument()
	case spec.ErrorCode_NOT_FOUND:
		return snip.CGRErrorsNotFound()
	case spec.ErrorCode_CONFLICT:
		return snip.CGRErrorsConflict()
	case spec.ErrorCode_REQUEST_ENTITY_TOO_LARGE:
		return snip.CGRErrorsRequestEntityTooLarge()
	case spec.ErrorCode_FAILED_PRECONDITION:
		return snip.CGRErrorsFailedPrecondition()
	case spec.ErrorCode_INTERNAL:
		return snip.CGRErrorsInternal()
	case spec.ErrorCode_TIMEOUT:
		return snip.CGRErrorsTimeout()
	case spec.ErrorCode_CUSTOM_CLIENT:
		return snip.CGRErrorsCustomClient()
	case spec.ErrorCode_CUSTOM_SERVER:
		return snip.CGRErrorsCustomServer()
	default:
		panic(fmt.Sprintf(`unknown error code string %q`, errorCode))
	}
}
