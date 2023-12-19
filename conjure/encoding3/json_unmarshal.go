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
	nameUnmarshalJSONResult   = "UnmarshalJSONResult"
	nameDisallowUnknownFields = "disallowUnknownFields"
	nameUnknownFields         = "unknownFields"
	nameMissingFields         = "missingFields"
)

func AnonFuncBodyUnmarshalJSON(methodBody *jen.Group, selector func() *jen.Statement, receiverType types.Type, ctxSelector func() *jen.Statement, strict bool) {
	if ctxSelector == nil {
		ctxSelector = snip.ContextTODO().Call().Clone
	}
	methodBody.Id("ctx").Op(":=").Add(ctxSelector())
	methodBody.Id("value").Op(":=").Add(djDot("Parse")).Call(jen.Id("data"))
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
	stmts = append(stmts, newMethodUnmarshalJSON(receiverName, receiverTypeName, includeStrict))
	if includeStrict {
		stmts = append(stmts, newMethodUnmarshalJSONStrict(receiverName, receiverTypeName))
	}
	stmts = append(stmts, newMethodUnmarshalJSONString(receiverName, receiverTypeName, includeStrict))
	if includeStrict {
		stmts = append(stmts, newMethodUnmarshalJSONStringStrict(receiverName, receiverTypeName))
	}
	stmts = append(stmts, newMethodUnmarshalJSONResult(receiverName, receiverTypeName, receiverType, includeStrict))
	return stmts
}

func newMethodUnmarshalJSON(receiverName string, receiverTypeName string, includeStrict bool) *jen.Statement {
	return snip.MethodUnmarshalJSON(receiverName, receiverTypeName).Block(
		jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(djDot("Parse")).Call(jen.Id("data")),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Return(jen.Id(receiverName).Dot(nameUnmarshalJSONResult).CallFunc(func(args *jen.Group) {
			args.Id("value")
			if includeStrict {
				args.False()
			}
		})),
	)
}

func newMethodUnmarshalJSONStrict(receiverName string, receiverTypeName string) *jen.Statement {
	return snip.MethodUnmarshalJSONStrict(receiverName, receiverTypeName).Block(
		jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(djDot("Parse")).Call(jen.Id("data")),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Return(jen.Id(receiverName).Dot(nameUnmarshalJSONResult).Call(jen.Id("value"), jen.True())),
	)
}

func newMethodUnmarshalJSONString(receiverName string, receiverTypeName string, includeStrict bool) *jen.Statement {
	return snip.MethodUnmarshalJSONString(receiverName, receiverTypeName).
		Block(
			jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(djDot("Parse")).Call(jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			jen.Return(jen.Id(receiverName).Dot(nameUnmarshalJSONResult).CallFunc(func(args *jen.Group) {
				args.Id("value")
				if includeStrict {
					args.False()
				}
			})),
		)
}

func newMethodUnmarshalJSONStringStrict(receiverName string, receiverTypeName string) *jen.Statement {
	return snip.MethodUnmarshalJSONStringStrict(receiverName, receiverTypeName).
		Block(
			jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(djDot("Parse")).Call(jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			jen.Return(jen.Id(receiverName).Dot(nameUnmarshalJSONResult).Call(jen.Id("value"), jen.True())),
		)
}

func newMethodUnmarshalJSONResult(receiverName string, receiverTypeName string, receiverType types.Type, includeStrict bool) *jen.Statement {
	return jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
		Id(nameUnmarshalJSONResult).
		ParamsFunc(func(params *jen.Group) {
			params.Add(jen.Id("value").Add(djDot("Result")))
			if includeStrict {
				params.Id(nameDisallowUnknownFields).Bool()
			}
		}).
		Params(jen.Error()).
		BlockFunc(func(methodBody *jen.Group) {
			switch typ := receiverType.(type) {
			case *types.AliasType:
				// TODO: Only do collections
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
				// TODO: this is not actually necessary for enums, UnmarshalText is enough
				methodBody.List(jen.Id("enumVal"), jen.Err()).Op(":=").Id("value").Dot("String").Call()
				methodBody.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(snip.FmtErrorf().Call(jen.Lit("type "+receiverTypeName+": %w"), jen.Err())),
				)
				methodBody.Op("*").Id(receiverName).Op("=").Add(typ.Code()).Call(jen.Id("enumVal"))
				methodBody.Return(jen.Nil())
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
		})
}

func unmarshalJSONStructFields(methodBody *jen.Group, receiverName string, receiverType string, fields []jsonStructField, isUnion bool) {
	var fieldResults []unmarshalJSONStructFieldResult
	hasRequiredFields := false
	hasCollections := false
	if isUnion {
		hasRequiredFields = true
		// add type field first
		field := jsonStructField{
			Key:      "type",
			Type:     types.String{},
			Selector: jen.Id(receiverName).Dot("typ").Clone,
		}
		typeFieldDecls := unmarshalJSONStructField(receiverName, receiverType, field, "fieldValue", false)
		typeFieldDecls.Init(methodBody)
		fieldResults = append(fieldResults, typeFieldDecls)
	}
	for _, field := range fields {
		result := unmarshalJSONStructField(receiverName, receiverType, field, "fieldValue", isUnion)
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
	methodBody.Var().Id(nameUnknownFields).Index().String()
	iterName := tmpVarName("iter", 0)
	idxName := tmpVarName("idx", 0)
	methodBody.List(jen.Id(iterName), jen.Id(idxName), jen.Err()).Op(":=").Id("value").Dot("ObjectIterator").Call(jen.Lit(0))
	methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
	methodBody.For(jen.Id(iterName).Dot("HasNext").Call(jen.Id("value"), jen.Id(idxName))).
		BlockFunc(func(forBody *jen.Group) {
			keyVar := tmpVarName("fieldKey", 0)
			valueVar := tmpVarName("fieldValue", 0)
			forBody.Var().List(jen.Id(keyVar), jen.Id(valueVar)).Add(djDot("Result"))
			forBody.List(jen.Id(keyVar), jen.Id(valueVar), jen.Id(idxName), jen.Err()).Op("=").Id(iterName).Dot("Next").Call(jen.Id("value"), jen.Id(idxName))
			forBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
			if len(fieldResults) > 0 {
				forBody.Switch(jen.Id(keyVar).Dot("Str")).BlockFunc(func(cases *jen.Group) {
					for _, result := range fieldResults {
						if result.Unmarshal != nil {
							result.Unmarshal(cases)
						}
					}
					cases.Default().Block(
						jen.If(jen.Id(nameDisallowUnknownFields)).Block(
							jen.Id(nameUnknownFields).Op("=").Append(jen.Id(nameUnknownFields), jen.Id(keyVar).Dot("Str")),
						),
					)
				})
			} else {
				forBody.Id("_").Op("=").Id(valueVar)
				forBody.If(jen.Id(nameDisallowUnknownFields)).Block(
					jen.Id(nameUnknownFields).Op("=").Append(jen.Id(nameUnknownFields), jen.Id(keyVar).Dot("Str")),
				)
			}

		})
	if hasRequiredFields {
		methodBody.Var().Id(nameMissingFields).Index().String()
		for _, result := range fieldResults {
			if result.Validate != nil {
				result.Validate(methodBody)
			}
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
		methodBody.If(jen.Len(jen.Id(nameMissingFields)).Op(">").Lit(0)).Block(
			jen.Return(snip.WerrorConvert().Call(djDot("UnmarshalMissingFieldsError").Values(
				jen.Id("Index").Op(":").Id("value").Dot("Index"),
				jen.Id("Type").Op(":").Lit(receiverType),
				jen.Id("Fields").Op(":").Id(nameMissingFields),
			))))
	} else if hasCollections {
		for _, result := range fieldResults {
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
	}
	methodBody.If(jen.Id(nameDisallowUnknownFields).Op("&&").Len(jen.Id(nameUnknownFields)).Op(">").Lit(0)).Block(
		jen.Return(snip.WerrorConvert().Call(djDot("UnmarshalUnknownFieldsError").Values(
			jen.Id("Index").Op(":").Id("value").Dot("Index"),
			jen.Id("Type").Op(":").Lit(receiverType),
			jen.Id("Fields").Op(":").Id(nameUnknownFields),
		))))
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
	valueVar string,
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
				jen.Id(nameMissingFields).Op("=").Append(jen.Id(nameMissingFields), jen.Lit(field.Key)),
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
				jen.Return(djDot("UnmarshalDuplicateFieldError").Values(
					jen.Id("Index").Op(":").Id("fieldKey").Dot("Index"),
					jen.Id("Type").Op(":").Lit(receiverType),
					jen.Id("Field").Op(":").Lit(field.Key),
				)))
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
				valueVar,
				jen.Return(snip.WerrorConvert().Call(djDot("UnmarshalFieldError").Values(
					jen.Id("Index").Op(":").Id("fieldValue").Dot("Index"),
					jen.Id("Type").Op(":").Lit(receiverType),
					jen.Id("Field").Op(":").Lit(field.Key),
					jen.Id("Err").Op(":").Id("err"),
				))).Clone,
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
		methodBody.List(selector(), jen.Err()).Op("=").Id(valueVar).Dot("Value").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case types.Bearertoken:
		tokenVal := tmpVarName("tokenVal", nestDepth)
		methodBody.List(jen.Id(tokenVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		methodBody.Add(selector()).Op("=").Add(snip.BearerTokenToken()).Call(jen.Id(tokenVal))

	case types.Binary:
		binaryVal := tmpVarName("binaryVal", nestDepth)
		methodBody.List(jen.Id(binaryVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		if isMapKey {
			methodBody.Add(selector()).Op("=").Add(snip.BinaryBinary()).Call(jen.Id(binaryVal))
		} else {
			methodBody.List(selector(), jen.Err()).Op("=").Add(snip.BinaryBinary()).Call(jen.Id(binaryVal)).Dot("Bytes").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		}

	case types.Boolean:
		if isMapKey {
			boolString := tmpVarName("boolString", nestDepth)
			methodBody.List(jen.Id(boolString), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			boolVal := tmpVarName("boolVal", nestDepth)
			methodBody.List(jen.Id(boolVal), jen.Err()).Op(":=").Add(snip.StrconvParseBool()).Call(jen.Id(boolString))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			methodBody.Add(selector()).Op("=").Add(snip.BooleanBoolean()).Call(jen.Id(boolVal))
		} else {
			methodBody.List(selector(), jen.Err()).Op("=").Id(valueVar).Dot("Bool").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		}
	case types.DateTime:
		timeVal := tmpVarName("timeVal", nestDepth)
		methodBody.List(jen.Id(timeVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.DateTimeParseDateTime()).Call(jen.Id(timeVal))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case types.Double:
		if isMapKey {
			floatVal := tmpVarName("floatVal", nestDepth)
			methodBody.List(jen.Id(floatVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			methodBody.List(selector(), jen.Err()).Op("=").Add(snip.StrconvParseFloat()).Call(jen.Id(floatVal))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		} else {
			methodBody.List(selector(), jen.Err()).Op("=").Id(valueVar).Dot("Float").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		}

	case types.Integer:
		intVal := tmpVarName("intVal", nestDepth)
		if isMapKey {
			methodBody.List(jen.Id(intVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			methodBody.List(selector(), jen.Err()).Op("=").Add(snip.StrconvAtoi()).Call(jen.Id(intVal))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		} else {
			methodBody.List(jen.Id(intVal), jen.Err()).Op(":=").Id(valueVar).Dot("Int").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			methodBody.Add(selector()).Op("=").Int().Call(jen.Id(intVal))
		}
	case types.RID:
		ridVal := tmpVarName("ridVal", nestDepth)
		methodBody.List(jen.Id(ridVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.RIDParseRID()).Call(jen.Id(ridVal))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case types.Safelong:
		longVal := tmpVarName("longVal", nestDepth)
		if isMapKey {
			methodBody.List(jen.Id(longVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			methodBody.List(selector(), jen.Err()).Op("=").Add(snip.SafeLongParseSafeLong()).Call(jen.Id(longVal))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		} else {
			methodBody.List(jen.Id(longVal), jen.Err()).Op(":=").Id(valueVar).Dot("Int").Call()
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
			methodBody.List(selector(), jen.Err()).Op("=").Add(snip.SafeLongNewSafeLong()).Call(jen.Id(longVal))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		}
	case types.String:
		methodBody.List(selector(), jen.Err()).Op("=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case types.UUID:
		uuidVal := tmpVarName("uuidVal", nestDepth)
		methodBody.List(jen.Id(uuidVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		methodBody.List(selector(), jen.Err()).Op("=").Add(snip.UUIDParseUUID()).Call(jen.Id(uuidVal))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())

	case *types.Optional:
		methodBody.If(jen.Id(valueVar).Dot("Type").Op("!=").Add(djDot("Null"))).
			BlockFunc(func(ifBody *jen.Group) {
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
		methodBody.If(selector().Op("==").Nil()).Block(
			selector().Op("=").Add(typ.Make()),
		)
		iterName := tmpVarName("iter", nestDepth)
		idxName := tmpVarName("idx", nestDepth)
		methodBody.List(jen.Id(iterName), jen.Id(idxName), jen.Err()).Op(":=").
			Id(valueVar).Dot("ArrayIterator").Call(jen.Lit(0))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		methodBody.For(jen.Id(iterName).Dot("HasNext").Call(jen.Id(valueVar), jen.Id(idxName))).
			BlockFunc(func(forBody *jen.Group) {
				nestDepth := nestDepth + 1
				resultVar := tmpVarName("arrayValue", nestDepth)
				forBody.Var().Id(resultVar).Add(djDot("Result"))
				forBody.List(jen.Id(resultVar), jen.Id(idxName), jen.Err()).Op("=").
					Id(iterName).Dot("Next").Call(jen.Id(valueVar), jen.Id(idxName))
				forBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
				listElement := tmpVarName("listElement", nestDepth)
				forBody.Var().Id(listElement).Add(typ.Item.Code())
				unmarshalJSONValue(
					forBody,
					jen.Id(listElement).Clone,
					typ.Item,
					resultVar,
					returnErrStmt,
					fieldDescriptor+" list element",
					false,
					nestDepth+1,
					strict)
				forBody.Add(selector()).Op("=").Append(selector(), jen.Id(listElement))
			})

	case *types.Map:
		methodBody.If(selector().Op("==").Nil()).Block(
			selector().Op("=").Add(typ.Make()),
		)
		iterName := tmpVarName("iter", nestDepth)
		idxName := tmpVarName("idx", nestDepth)
		methodBody.List(jen.Id(iterName), jen.Id(idxName), jen.Err()).Op(":=").
			Id(valueVar).Dot("ObjectIterator").Call(jen.Lit(0))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		methodBody.For(jen.Id(iterName).Dot("HasNext").Call(jen.Id(valueVar), jen.Id(idxName))).
			BlockFunc(func(forBody *jen.Group) {
				nestDepth := nestDepth + 1
				keyVar := tmpVarName("mapKey", nestDepth)
				resultVar := tmpVarName("mapValue", nestDepth)
				forBody.Var().List(jen.Id(keyVar), jen.Id(resultVar)).Add(djDot("Result"))
				forBody.List(jen.Id(keyVar), jen.Id(resultVar), jen.Id(idxName), jen.Err()).Op("=").
					Id(iterName).Dot("Next").Call(jen.Id(valueVar), jen.Id(idxName))
				forBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
				mapKeyVal := tmpVarName("mapKeyVal", nestDepth)
				switch typ.Key.(type) {
				case types.Binary:
					// Use binary.Binary for map keys since []byte is invalid in go maps.
					forBody.Var().Id(mapKeyVal).Add(snip.BinaryBinary())
				case types.Boolean:
					forBody.Var().Id(mapKeyVal).Add(snip.BooleanBoolean())
				default:
					forBody.Var().Id(mapKeyVal).Add(typ.Key.Code())
				}
				forBody.BlockFunc(func(keyBlock *jen.Group) {
					unmarshalJSONValue(
						keyBlock,
						jen.Id(mapKeyVal).Clone,
						typ.Key,
						keyVar,
						returnErrStmt,
						fieldDescriptor+" map key",
						true,
						nestDepth+1,
						strict)
				})
				forBody.If(
					jen.List(jen.Id("_"), jen.Id("exists").Op(":=").Add(selector()).Index(jen.Id(mapKeyVal))),
					jen.Id("exists"),
				).Block(
					jen.Return(snip.WerrorConvert().Call(djDot("UnmarshalDuplicateMapKeyError").Values(
						jen.Id("Type").Op(":").Lit(fieldDescriptor),
					))),
				)
				mapVal := tmpVarName("mapVal", nestDepth)
				forBody.Var().Id(mapVal).Add(typ.Val.Code())
				forBody.BlockFunc(func(valBlock *jen.Group) {
					unmarshalJSONValue(
						valBlock,
						jen.Id(mapVal).Clone,
						typ.Val,
						resultVar,
						returnErrStmt,
						fieldDescriptor+" map value",
						false,
						nestDepth+1,
						strict)
				})
				forBody.Add(selector()).Index(jen.Id(mapKeyVal)).Op("=").Id(mapVal)
			})

	case *types.EnumType:
		enumVal := tmpVarName("enumVal", nestDepth)
		methodBody.List(jen.Id(enumVal), jen.Err()).Op(":=").Id(valueVar).Dot("String").Call()
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(snip.FmtErrorf().Call(jen.Lit("field "+fieldDescriptor+": %w"), jen.Err())),
		)
		methodBody.Add(selector()).Op("=").Add(typ.Constructor()).Call(typ.ValueType().Call(jen.Id(enumVal)))
	case *types.AliasType:
		if typ.IsOptional() {
			unmarshalJSONValue(
				methodBody,
				selector().Dot("Value").Clone,
				typ.Item,
				valueVar,
				returnErrStmt,
				fieldDescriptor,
				isMapKey,
				nestDepth+1,
				strict)
		} else {
			aliasVal := tmpVarName("aliasVal", nestDepth)
			methodBody.Var().Id(aliasVal).Add(typ.Item.Code())
			unmarshalJSONValue(
				methodBody,
				jen.Id(aliasVal).Clone,
				typ.Item,
				valueVar,
				returnErrStmt,
				fieldDescriptor,
				isMapKey,
				nestDepth+1,
				strict)
			methodBody.Add(selector()).Op("=").Add(typ.Code()).Call(jen.Id(aliasVal))
		}
	case *types.ObjectType, *types.UnionType:
		methodBody.If(
			jen.Err().Op(":=").Add(selector()).Dot(nameUnmarshalJSONResult).CallFunc(func(args *jen.Group) {
				args.Id(valueVar)
				if valueType.ContainsStrictFields() {
					if strict != nil {
						args.Lit(*strict)
					} else {
						args.Id(nameDisallowUnknownFields)
					}
				}
			}),
			jen.Err().Op("!=").Nil(),
		).Block(
			returnErrStmt(),
		)

	case *types.External:
		methodBody.Err().Op("=").Add(snip.SafeJSONUnmarshal()).Call(jen.Op("&").Add(selector()), jen.Id(valueVar).Dot("Raw"))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
	}
}
