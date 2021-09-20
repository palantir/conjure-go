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

package encoding

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	dataName = "data"
)

func AnonFuncBodyUnmarshalJSON(methodBody *jen.Group, selector func() *jen.Statement, receiverType types.Type, ctxSelector func() *jen.Statement, strict bool) {
	if ctxSelector == nil {
		ctxSelector = snip.ContextTODO().Call().Clone
	}
	methodBody.Id("ctx").Op(":=").Add(ctxSelector())
	methodBody.Add(unmarshalJSONValidBytes(receiverType.String()))
	methodBody.Id("value").Op(":=").Add(snip.GJSONParseBytes()).Call(jen.Id("data"))
	methodBody.Var().Err().Error()
	unmarshalJSONValue(
		methodBody,
		selector,
		receiverType,
		"value",
		jen.Return(jen.Err()).Clone,
		receiverType.String(),
		false,
		0,
		&strict,
	)
	methodBody.Return(jen.Nil())
}

func UnmarshalJSONMethods(receiverName string, receiverTypeName string, receiverType types.Type) []*jen.Statement {
	includeStrict := receiverType.ContainsStrictFields()
	var stmts []*jen.Statement
	stmts = append(stmts, jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
		Id("UnmarshalJSON").
		Params(jen.Id("data").Op("[]").Byte()).
		Params(jen.Error()).
		BlockFunc(func(methodBody *jen.Group) {
			methodBody.Id("ctx").Op(":=").Add(snip.ContextTODO()).Call()
			methodBody.Add(unmarshalJSONValidBytes(receiverTypeName))
			methodBody.Return(jen.Id(receiverName).Dot("unmarshalJSONResult").CallFunc(func(args *jen.Group) {
				args.Id("ctx")
				args.Add(snip.GJSONParseBytes().Call(jen.Id("data")))
				if includeStrict {
					args.False()
				}
			}))
		}),
	)
	if includeStrict {
		stmts = append(stmts, jen.Func().
			Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
			Id("UnmarshalJSONStrict").
			Params(jen.Id("data").Op("[]").Byte()).
			Params(jen.Error()).
			BlockFunc(func(methodBody *jen.Group) {
				methodBody.Id("ctx").Op(":=").Add(snip.ContextTODO()).Call()
				methodBody.Add(unmarshalJSONValidBytes(receiverTypeName))
				methodBody.Return(jen.Id(receiverName).Dot("unmarshalJSONResult").Call(
					jen.Id("ctx"),
					snip.GJSONParseBytes().Call(jen.Id("data")),
					jen.True(),
				))
			}),
		)
	}
	stmts = append(stmts, jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
		Id("UnmarshalJSONString").
		Params(jen.Id("data").String()).
		Params(jen.Error()).
		BlockFunc(func(methodBody *jen.Group) {
			methodBody.Id("ctx").Op(":=").Add(snip.ContextTODO()).Call()
			methodBody.Add(unmarshalJSONValid(receiverTypeName))
			methodBody.Return(jen.Id(receiverName).Dot("unmarshalJSONResult").CallFunc(func(args *jen.Group) {
				args.Id("ctx")
				args.Add(snip.GJSONParse()).Call(jen.Id("data"))
				if includeStrict {
					args.False()
				}
			}))
		}),
	)
	if includeStrict {
		stmts = append(stmts, jen.Func().
			Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
			Id("UnmarshalJSONStringStrict").
			Params(jen.Id("data").String()).
			Params(jen.Error()).
			BlockFunc(func(methodBody *jen.Group) {
				methodBody.Id("ctx").Op(":=").Add(snip.ContextTODO()).Call()
				methodBody.Add(unmarshalJSONValid(receiverTypeName))
				methodBody.Return(jen.Id(receiverName).Dot("unmarshalJSONResult").Call(
					jen.Id("ctx"),
					snip.GJSONParse().Call(jen.Id("data")),
					jen.True(),
				))
			}),
		)
	}
	stmts = append(stmts, jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
		Id("unmarshalJSONResult").
		ParamsFunc(func(params *jen.Group) {
			params.Add(snip.ContextVar())
			params.Add(jen.Id("value").Add(snip.GJSONResult()))
			if includeStrict {
				params.Id("strict").Bool()
			}
		}).
		Params(jen.Error()).
		BlockFunc(func(methodBody *jen.Group) {
			switch typ := receiverType.(type) {
			case *types.AliasType:
				rawVarName := "raw" + typ.Name
				methodBody.Var().Id(rawVarName).Add(typ.Item.Code())
				methodBody.Var().Err().Error()
				unmarshalJSONValue(
					methodBody,
					jen.Id(rawVarName).Clone,
					typ.Item,
					"value",
					jen.Return(jen.Err()).Clone,
					typ.Name,
					false,
					0,
					nil,
				)
				if typ.IsOptional() {
					methodBody.Id(receiverName).Dot("Value").Op("=").Id(rawVarName)
				} else {
					methodBody.Op("*").Id(receiverName).Op("=").Id(typ.Name).Call(jen.Id(rawVarName))
				}
				methodBody.Return(jen.Nil())
			case *types.EnumType:
				methodBody.Var().Err().Error()
				methodBody.Add(unmarshalJSONTypeCheck("value", jen.Return(jen.Err()).Clone, "type "+typ.Name, "string", snip.GJSONString))
				methodBody.Return(jen.Id(receiverName).Dot("UnmarshalString").Call(jen.Id("value").Dot("Str")))
			case *types.ObjectType:
				var fields []jsonStructField
				for _, field := range typ.Fields {
					fields = append(fields, jsonStructField{
						Key:      field.Name,
						Type:     field.Type,
						Selector: jen.Id(receiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
					})
				}
				unmarshalJSONStructFields(methodBody, receiverName, receiverTypeName, fields, false)
				methodBody.Return(jen.Nil())
			case *types.UnionType:
				var fields []jsonStructField
				for _, field := range typ.Fields {
					fields = append(fields, jsonStructField{
						Key:      field.Name,
						Type:     field.Type,
						Selector: jen.Id(receiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone,
					})
				}
				unmarshalJSONStructFields(methodBody, receiverName, receiverTypeName, fields, true)
				methodBody.Return(jen.Nil())
			default:
				panic("cannot generate methods for non-named type " + receiverType.String())
			}
		}),
	)
	return stmts
}

func unmarshalJSONValid(receiverType string) *jen.Statement {
	return jen.If(jen.Op("!").Add(snip.GJSONValid()).Call(jen.Id(dataName))).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"),
			jen.Lit(fmt.Sprintf("invalid JSON for %s", receiverType)),
		)),
	)
}

func unmarshalJSONValidBytes(receiverType string) *jen.Statement {
	return jen.If(jen.Op("!").Add(snip.GJSONValidBytes()).Call(jen.Id(dataName))).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"),
			jen.Lit(fmt.Sprintf("invalid JSON for %s", receiverType)),
		)),
	)
}

func unmarshalJSONStructFields(methodBody *jen.Group, receiverName string, receiverType string, fields []jsonStructField, isUnion bool) {
	methodBody.If(jen.Op("!").Id("value").Dot("IsObject").Call()).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"),
			jen.Lit(fmt.Sprintf("type %s expected JSON object", receiverType)),
		)),
	)
	var fieldResults []unmarshalJSONStructFieldResult
	hasRequiredFields := false
	hasCollections := false
	if isUnion {
		hasRequiredFields = true
		result := unmarshalJSONStructField(receiverName, receiverType, jsonStructField{
			Key:      "type",
			Type:     types.String{},
			Selector: jen.Id(receiverName).Dot("typ").Clone,
		}, false)
		result.Init(methodBody)
		fieldResults = append(fieldResults, result)
	}
	for _, field := range fields {
		result := unmarshalJSONStructField(receiverName, receiverType, field, isUnion)
		if result.Validate != nil {
			hasRequiredFields = true
		}
		if result.DefaultCollection != nil {
			hasCollections = true
		}
		if result.Init != nil {
			result.Init(methodBody)
		}
		fieldResults = append(fieldResults, result)
	}
	methodBody.Var().Id("unrecognizedFields").Op("[]").String()
	methodBody.Var().Err().Error()
	methodBody.Id("value").Dot("ForEach").Call(
		jen.Func().
			Params(jen.Id("key"), jen.Id("value").Add(snip.GJSONResult())).
			Params(jen.Bool()).
			BlockFunc(func(rangeBody *jen.Group) {
				rangeBody.Switch(jen.Id("key").Dot("Str")).BlockFunc(func(cases *jen.Group) {
					for _, result := range fieldResults {
						if result.Unmarshal != nil {
							result.Unmarshal(cases)
						}
					}
					cases.Default().Block(
						jen.If(jen.Id("strict")).Block(
							jen.Id("unrecognizedFields").
								Op("=").
								Append(jen.Id("unrecognizedFields"), jen.Id("key").Dot("Str")),
						),
					)
				})
				rangeBody.Return(jen.Err().Op("==").Nil())
			}),
	)
	methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
	if hasRequiredFields {
		methodBody.Var().Id("missingFields").Op("[]").String()
		for _, result := range fieldResults {
			if result.Validate != nil {
				result.Validate(methodBody)
			}
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
		methodBody.If(jen.Len(jen.Id("missingFields")).Op(">").Lit(0)).Block(
			jen.Return(snip.WerrorErrorContext().Call(
				jen.Id("ctx"),
				jen.Lit(fmt.Sprintf("type %s missing required JSON fields", receiverType)),
				snip.WerrorSafeParam().Call(jen.Lit("missingFields"), jen.Id("missingFields")),
			)),
		)
	} else if hasCollections {
		for _, result := range fieldResults {
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
	}
	methodBody.If(jen.Id("strict").Op("&&").Len(jen.Id("unrecognizedFields")).Op(">").Lit(0)).Block(
		jen.Return(snip.WerrorErrorContext().Call(
			jen.Id("ctx"),
			jen.Lit(fmt.Sprintf("type %s encountered unrecognized JSON fields", receiverType)),
			// unrecognized user input must stay unsafe
			snip.WerrorUnsafeParam().Call(jen.Lit("unrecognizedFields"), jen.Id("unrecognizedFields")),
		)),
	)
}

type unmarshalJSONStructFieldResult struct {
	Init              func(*jen.Group)
	Unmarshal         func(*jen.Group)
	Validate          func(*jen.Group)
	DefaultCollection func(*jen.Group)
}

func unmarshalJSONStructField(
	receiverName string,
	receiverType string,
	field jsonStructField,
	isUnionField bool,
) (result unmarshalJSONStructFieldResult) {
	requiredField := !(field.Type.IsCollection() || field.Type.IsOptional())
	seenVar := "seen" + transforms.ExportedFieldName(field.Key)
	result.Init = func(methodBody *jen.Group) {
		methodBody.Var().Id(seenVar).Bool()
	}
	if requiredField {
		result.Validate = func(methodBody *jen.Group) {
			methodBody.IfFunc(func(conds *jen.Group) {
				if isUnionField {
					conds.Id(receiverName).Dot("typ").Op("==").Lit(field.Key).
						Op("&&").
						Op("!").Id(seenVar)
				} else {
					conds.Op("!").Id(seenVar)
				}
			}).Block(
				jen.Id("missingFields").Op("=").Append(jen.Id("missingFields"), jen.Lit(field.Key)),
			)
		}
	} else if mk := field.Type.Make(); mk != nil && !isUnionField {
		result.DefaultCollection = func(methodBody *jen.Group) {
			methodBody.If(jen.Op("!").Id(seenVar)).Block(
				field.Selector().Op("=").Add(mk),
			)
		}
	}
	result.Unmarshal = func(cases *jen.Group) {
		cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
			caseBody.If(jen.Id(seenVar)).Block(
				jen.Err().Op("=").Add(snip.WerrorErrorContext().Call(jen.Id("ctx"), jen.Lit(
					fmt.Sprintf("type %s encountered duplicate %q field", receiverType, field.Key),
				))),
				jen.Return(jen.False()),
			)
			caseBody.Id(seenVar).Op("=").True()

			selector := field.Selector
			if isUnionField {
				caseBody.Var().Id("unionVal").Add(field.Type.Code())
				selector = jen.Id("unionVal").Clone
			}
			unmarshalJSONValue(
				caseBody,
				selector,
				field.Type,
				"value",
				jen.Return(jen.False()).Clone,
				fmt.Sprintf("field %s[%q]", receiverType, field.Key),
				false,
				0,
				nil)
			if isUnionField {
				caseBody.Add(field.Selector()).Op("=").Op("&").Id("unionVal")
			}
		})
	}

	return result
}

func unmarshalJSONValue(
	methodBody *jen.Group,
	selector func() *jen.Statement,
	valueType types.Type,
	valueVar string,
	returnErrStmt func() *jen.Statement,
	fieldDescriptor string,
	isMapKey bool,
	nestDepth int,
	strict *bool, // if set, force strictness
) {
	switch typ := valueType.(type) {
	case types.Any:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "non-null value",
			snip.GJSONJSON, snip.GJSONString, snip.GJSONNumber, snip.GJSONTrue, snip.GJSONFalse))
		methodBody.Add(selector()).Op("=").Id(valueVar).Dot("Value").Call()

	case types.Bearertoken:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.BearerTokenNew()).Call(jen.Id(valueVar).Dot("Str"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case types.Binary:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
		if isMapKey {
			methodBody.Add(selector()).Op("=").Add(snip.BinaryBinary()).Call(jen.Id(valueVar).Dot("Str"))
		} else {
			methodBody.List(selector(), jen.Err()).
				Op("=").
				Add(snip.BinaryBinary()).Call(jen.Id(valueVar).Dot("Str")).Dot("Bytes").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(
				jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
				returnErrStmt())
		}

	case types.Boolean:
		if isMapKey {
			methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
			methodBody.Var().Id("boolVal").Bool()
			methodBody.List(jen.Id("boolVal"), jen.Err()).Op("=").Add(snip.StrconvParseBool()).Call(jen.Id(valueVar).Dot("Str"))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(
				jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
				returnErrStmt())
			methodBody.Add(selector()).Op("=").Add(snip.BooleanBoolean()).Call(jen.Id("boolVal"))
		} else {
			methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "boolean", snip.GJSONTrue, snip.GJSONFalse))
			methodBody.Add(selector()).Op("=").Id(valueVar).Dot("Type").Op("==").Add(snip.GJSONTrue())
		}
	case types.DateTime:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.DateTimeParseDateTime()).Call(jen.Id(valueVar).Dot("Str"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case types.Double:
		methodBody.Switch(jen.Id(valueVar).Dot("Str")).Block(
			jen.Case(jen.Lit("NaN")).Block(
				selector().Op("=").Add(snip.MathNaN()).Call(),
			),
			jen.Case(jen.Lit("Infinity")).Block(
				selector().Op("=").Add(snip.MathInf()).Call(jen.Lit(1)),
			),
			jen.Case(jen.Lit("-Infinity")).Block(
				selector().Op("=").Add(snip.MathInf()).Call(jen.Lit(-1)),
			),
			jen.Default().BlockFunc(func(defaultBody *jen.Group) {
				var valueField string
				if isMapKey {
					defaultBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
					valueField = "Str"
				} else {
					defaultBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "number", snip.GJSONNumber))
					valueField = "Raw"
				}
				defaultBody.List(selector(), jen.Err()).Op("=").Add(snip.StrconvParseFloat()).Call(jen.Id(valueVar).Dot(valueField), jen.Lit(64))
				defaultBody.If(jen.Err().Op("!=").Nil()).Block(
					jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
					returnErrStmt())
			}),
		)

	case types.Integer:
		var valueField string
		if isMapKey {
			methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
			valueField = "Str"
		} else {
			methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "number", snip.GJSONNumber))
			valueField = "Raw"
		}
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.StrconvAtoi()).Call(jen.Id(valueVar).Dot(valueField))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())

	case types.RID:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.RIDParseRID()).Call(jen.Id(valueVar).Dot("Str"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())

	case types.Safelong:
		var valueField string
		if isMapKey {
			methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
			valueField = "Str"
		} else {
			methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "number", snip.GJSONNumber))
			valueField = "Raw"
		}
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.SafeLongParseSafeLong()).Call(jen.Id(valueVar).Dot(valueField))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())

	case types.String:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
		methodBody.Add(selector()).Op("=").Id(valueVar).Dot("Str")

	case types.UUID:
		methodBody.Add(unmarshalJSONTypeCheck(valueVar, returnErrStmt, fieldDescriptor, "string", snip.GJSONString))
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.UUIDParseUUID()).Call(jen.Id(valueVar).Dot("Str"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())

	case *types.Optional:
		methodBody.If(jen.Id(valueVar).Dot("Type").Op("!=").Add(snip.GJSONNull())).BlockFunc(func(ifBody *jen.Group) {
			optVal := tmpVarName("optVal", nestDepth)
			ifBody.Var().Id(optVal).Add(typ.Item.Code())
			unmarshalJSONValue(
				ifBody,
				jen.Id(optVal).Clone,
				typ.Item,
				valueVar,
				returnErrStmt,
				fieldDescriptor,
				isMapKey,
				nestDepth+1,
				strict)
			ifBody.Add(selector()).Op("=").Op("&").Id(optVal)
		})
	case *types.List:
		methodBody.If(jen.Op("!").Id(valueVar).Dot("IsArray").Call()).Block(
			jen.Err().Op("=").Add(snip.WerrorErrorContext()).Call(
				jen.Id("ctx"),
				jen.Lit(fmt.Sprintf("%s expected JSON array", fieldDescriptor)),
			),
			returnErrStmt(),
		)
		listElement := tmpVarName("listElement", nestDepth)
		methodBody.Id(valueVar).Dot("ForEach").Call(
			jen.Func().
				Params(jen.Id("_"), jen.Id("value").Add(snip.GJSONResult())).
				Params(jen.Bool()).
				BlockFunc(func(rangeBody *jen.Group) {
					rangeBody.Var().Id(listElement).Add(typ.Item.Code())
					unmarshalJSONValue(
						rangeBody,
						jen.Id(listElement).Clone,
						typ.Item,
						"value",
						jen.Return(jen.False()).Clone,
						fieldDescriptor+" list element",
						false,
						nestDepth+1,
						strict)
					rangeBody.Add(selector()).Op("=").Append(selector(), jen.Id(listElement))
					rangeBody.Return(jen.Err().Op("==").Nil())
				}),
		)
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
	case *types.Map:
		methodBody.If(jen.Op("!").Id(valueVar).Dot("IsObject").Call()).Block(
			jen.Err().Op("=").Add(snip.WerrorErrorContext()).Call(
				jen.Id("ctx"),
				jen.Lit(fmt.Sprintf("%s expected JSON object", fieldDescriptor)),
			),
			returnErrStmt(),
		)
		methodBody.If(selector().Op("==").Nil()).Block(
			selector().Op("=").Add(typ.Make()),
		)
		mapKey := tmpVarName("mapKey", nestDepth)
		mapVal := tmpVarName("mapVal", nestDepth)
		methodBody.Id(valueVar).Dot("ForEach").Call(
			jen.Func().
				Params(jen.Id("key"), jen.Id("value").Add(snip.GJSONResult())).
				Params(jen.Bool()).
				BlockFunc(func(rangeBody *jen.Group) {
					returnErrStmt := jen.Return(jen.False()).Clone
					switch typ.Key.(type) {
					case types.Binary:
						// Use binary.Binary for map keys since []byte is invalid in go maps.
						rangeBody.Var().Id(mapKey).Add(snip.BinaryBinary())
					case types.Boolean:
						rangeBody.Var().Id(mapKey).Add(snip.BooleanBoolean())
					default:
						rangeBody.Var().Id(mapKey).Add(typ.Key.Code())
					}
					rangeBody.BlockFunc(func(keyBlock *jen.Group) {
						unmarshalJSONValue(
							keyBlock,
							jen.Id(mapKey).Clone,
							typ.Key,
							"key",
							returnErrStmt,
							fieldDescriptor+" map key",
							true,
							nestDepth+1,
							strict)
					})
					rangeBody.If(
						jen.List(jen.Id("_"), jen.Id("exists").Op(":=").Add(selector()).Index(jen.Id(mapKey))),
						jen.Id("exists"),
					).Block(
						jen.Err().Op("=").Add(snip.WerrorErrorContext().Call(
							jen.Id("ctx"),
							jen.Lit(fmt.Sprintf("%s encountered duplicate map key", fieldDescriptor)),
						)),
						returnErrStmt(),
					)
					rangeBody.Var().Id(mapVal).Add(typ.Val.Code())
					rangeBody.BlockFunc(func(valBlock *jen.Group) {
						unmarshalJSONValue(
							valBlock,
							jen.Id(mapVal).Clone,
							typ.Val,
							"value",
							returnErrStmt,
							fieldDescriptor+" map value",
							false,
							nestDepth+1,
							strict)
					})

					rangeBody.Add(selector()).Index(jen.Id(mapKey)).Op("=").Id(mapVal)
					rangeBody.Return(jen.Err().Op("==").Nil())
				}),
		)
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
	case *types.AliasType, *types.EnumType, *types.ObjectType, *types.UnionType:
		unmarshalStrict := jen.If(
			jen.Err().Op("=").Add(selector()).Dot("UnmarshalJSONStringStrict").Call(jen.Id(valueVar).Dot("Raw")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())
		unmarshalNotStrict := jen.If(
			jen.Err().Op("=").Add(selector()).Dot("UnmarshalJSONString").Call(jen.Id(valueVar).Dot("Raw")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())

		if valueType.ContainsStrictFields() {
			if strict == nil {
				methodBody.If(jen.Id("strict")).Block(unmarshalStrict).Else().Block(unmarshalNotStrict)
			} else if *strict {
				methodBody.Add(unmarshalStrict)
			} else {
				methodBody.Add(unmarshalNotStrict)
			}
		} else {
			methodBody.Add(unmarshalNotStrict)
		}
	case *types.External:
		methodBody.Err().Op("=").Add(snip.SafeJSONUnmarshal()).Call(jen.Op("&").Add(selector()), jen.Id(valueVar).Dot("Raw"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Add(snip.WerrorWrapContext()).Call(jen.Id("ctx"), jen.Err(), jen.Lit(fieldDescriptor)),
			returnErrStmt())
	}
}

func unmarshalJSONTypeCheck(
	valueVar string,
	returnErrStmt func() *jen.Statement,
	fieldDescriptor string,
	typeDescriptor string,
	typeStmts ...func() *jen.Statement,
) *jen.Statement {
	return jen.IfFunc(func(conds *jen.Group) {
		cond := jen.Empty()
		for i, typeStmt := range typeStmts {
			if i > 0 {
				cond = cond.Op("&&")
			}
			cond = cond.Id(valueVar).Dot("Type").Op("!=").Add(typeStmt())
		}
		conds.Add(cond)
	}).Block(
		jen.Err().Op("=").Add(snip.WerrorErrorContext()).Call(
			jen.Id("ctx"),
			jen.Lit(fmt.Sprintf("%s expected JSON %s", fieldDescriptor, typeDescriptor)),
		),
		returnErrStmt(),
	)
}

func tmpVarName(base string, depth int) string {
	if depth == 0 {
		return base
	}
	return fmt.Sprintf("%s%d", base, depth)
}
