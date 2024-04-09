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

func UnmarshalJSONMethods(receiverName string, receiverTypeName string, receiverType types.Type, withAuxiliary bool) []*jen.Statement {
	includeStrict := receiverType.ContainsStrictFields()
	var stmts []*jen.Statement
	stmts = append(stmts, newMethodUnmarshalJSON(receiverName, receiverTypeName, includeStrict))
	if includeStrict {
		stmts = append(stmts, newMethodUnmarshalJSONStrict(receiverName, receiverTypeName))
	}
	if withAuxiliary {
		stmts = append(stmts, newMethodUnmarshalJSONString(receiverName, receiverTypeName, includeStrict))
	}
	if includeStrict && withAuxiliary {
		stmts = append(stmts, newMethodUnmarshalJSONStringStrict(receiverName, receiverTypeName))
	}
	stmts = append(stmts, newMethodUnmarshalJSONResult(receiverName, receiverTypeName, receiverType, includeStrict))

	if withAuxiliary {
		stmts = append(stmts,
			snip.MethodUnmarshalYAMLSig(receiverName, receiverTypeName).Block(
				jen.Return(snip.DJUnmarshalYAML().Call(
					jen.Id(receiverName),
					jen.Id("unmarshal"),
				)),
			),
		)
	}
	return stmts
}

func newMethodUnmarshalJSON(receiverName string, receiverTypeName string, includeStrict bool) *jen.Statement {
	return snip.MethodUnmarshalJSON(receiverName, receiverTypeName).Block(
		jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(snip.DJParse()).Call(jen.Id(dataName)),
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
		jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(snip.DJParse()).Call(jen.Id(dataName)),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Return(jen.Id(receiverName).Dot(nameUnmarshalJSONResult).Call(jen.Id("value"), jen.True())),
	)
}

func newMethodUnmarshalJSONString(receiverName string, receiverTypeName string, includeStrict bool) *jen.Statement {
	return snip.MethodUnmarshalJSONString(receiverName, receiverTypeName).
		Block(
			jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(snip.DJParse()).Call(jen.Id(dataName)),
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
			jen.List(jen.Id("value"), jen.Err()).Op(":=").Add(snip.DJParse()).Call(jen.Id(dataName)),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			jen.Return(jen.Id(receiverName).Dot(nameUnmarshalJSONResult).Call(jen.Id("value"), jen.True())),
		)
}

func newMethodUnmarshalJSONResult(receiverName string, receiverTypeName string, receiverType types.Type, includeStrict bool) *jen.Statement {
	return jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
		Id(nameUnmarshalJSONResult).
		ParamsFunc(func(params *jen.Group) {
			params.Add(jen.Id("value").Add(snip.DJResult()))
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
				unmarshalJSONValue(
					methodBody,
					jen.Id(rawVarName).Clone,
					typ.Item,
					"value",
					jen.Return(snip.DJNewUnmarshalFieldError().Call(
						jen.Id("value"),
						jen.Lit("type "+typ.Name),
						jen.Err(),
					)).Clone,
					typ.Name,
					false,
					0,
					&includeStrict,
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
	nestDepth := 0
	keyVar := tmpVarName("fieldKey", nestDepth)
	valueVar := tmpVarName("fieldValue", nestDepth)

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
		typeFieldDecls := unmarshalJSONStructField(receiverName, receiverType, field, keyVar, valueVar, false)
		typeFieldDecls.Init(methodBody)
		fieldResults = append(fieldResults, typeFieldDecls)
	}
	for _, field := range fields {
		result := unmarshalJSONStructField(receiverName, receiverType, field, keyVar, valueVar, isUnion)
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

	idxName := tmpVarName("idx", nestDepth)
	methodBody.Var().Id(idxName).Int()
	methodBody.For().BlockFunc(func(forBody *jen.Group) {
		okVar := tmpVarName("ok", nestDepth)
		errVar := tmpVarName("err", nestDepth)
		forBody.Var().List(jen.Id(keyVar), jen.Id(valueVar)).Add(snip.DJResult())
		forBody.Var().Id(okVar).Bool()
		forBody.Var().Id(errVar).Error()
		forBody.List(jen.Id(keyVar), jen.Id(valueVar), jen.Id(idxName), jen.Id(okVar), jen.Id(errVar)).Op("=").
			Id("value").Dot("NextObjectEntry").Call(jen.Id(idxName))
		forBody.If(jen.Id(errVar).Op("!=").Nil()).Block(jen.Return(jen.Id(errVar)))
		forBody.If(jen.Op("!").Id(okVar)).Block(jen.Break())
		keyString := tmpVarName("keyString", nestDepth)
		forBody.List(jen.Id(keyString), jen.Err()).Op(":=").Id(keyVar).Dot("String").Call()
		forBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
		if len(fieldResults) > 0 {
			forBody.Switch(jen.Id(keyString)).BlockFunc(func(cases *jen.Group) {
				for _, result := range fieldResults {
					if result.Unmarshal != nil {
						result.Unmarshal(cases)
					}
				}
				cases.Default().Block(
					jen.If(jen.Id(nameDisallowUnknownFields)).Block(
						jen.Id(nameUnknownFields).Op("=").Append(jen.Id(nameUnknownFields), jen.Id(keyString)),
					),
				)
			})
		} else {
			forBody.Id("_").Op("=").Id(valueVar)
			forBody.If(jen.Id(nameDisallowUnknownFields)).Block(
				jen.Id(nameUnknownFields).Op("=").Append(jen.Id(nameUnknownFields), jen.Id(keyString)),
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
			jen.Return(snip.DJNewUnmarshalMissingFieldsError().Call(
				jen.Id("value").Dot("Index").Call(),
				jen.Lit(receiverType),
				jen.Id(nameMissingFields),
			)),
		)
	} else if hasCollections {
		for _, result := range fieldResults {
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
	}
	methodBody.If(jen.Id(nameDisallowUnknownFields).Op("&&").Len(jen.Id(nameUnknownFields)).Op(">").Lit(0)).Block(
		jen.Return(snip.DJNewUnmarshalUnknownFieldsError().Call(
			jen.Id("value").Dot("Index").Call(),
			jen.Lit(receiverType),
			jen.Id(nameUnknownFields),
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
	keyVar string,
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
			fieldDescriptor := fmt.Sprintf("field %s[%q]", receiverType, field.Key)
			caseBody.If(jen.Id(seenVar)).Block(
				jen.Return(snip.DJNewUnmarshalDuplicateFieldError().Call(
					jen.Id(keyVar).Dot("Index").Call(),
					jen.Lit(fieldDescriptor),
				)),
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
				valueVar,
				jen.Return(snip.DJNewUnmarshalFieldError().Call(
					jen.Id(valueVar).Dot("Index").Call(),
					jen.Lit(fieldDescriptor),
					jen.Err(),
				)).Clone,
				fieldDescriptor,
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
		errVal := tmpVarName("err", nestDepth)
		methodBody.Var().Id(errVal).Error()
		methodBody.List(selector(), jen.Id(errVal)).Op("=").Id(valueVar).Dot("Value").Call()
		methodBody.If(jen.Id(errVal).Op("!=").Nil()).Block(returnErrStmt())

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
			errVal := tmpVarName("err", nestDepth)
			methodBody.Var().Id(errVal).Error()
			methodBody.List(selector(), jen.Id(errVal)).Op("=").Add(snip.BinaryBinary()).Call(jen.Id(binaryVal)).Dot("Bytes").Call()
			methodBody.If(jen.Id(errVal).Op("!=").Nil()).Block(returnErrStmt())
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
			errVal := tmpVarName("err", nestDepth)
			methodBody.Var().Id(errVal).Error()
			methodBody.List(selector(), jen.Id(errVal)).Op("=").Id(valueVar).Dot("Bool").Call()
			methodBody.If(jen.Id(errVal).Op("!=").Nil()).Block(returnErrStmt())
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
			methodBody.List(selector(), jen.Err()).Op("=").Add(snip.StrconvParseFloat()).Call(jen.Id(floatVal), jen.Lit(64))
			methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
		} else {
			errVal := tmpVarName("err", nestDepth)
			methodBody.Var().Id(errVal).Error()
			methodBody.List(selector(), jen.Id(errVal)).Op("=").Id(valueVar).Dot("Float").Call()
			methodBody.If(jen.Id(errVal).Op("!=").Nil()).Block(returnErrStmt())
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
		methodBody.If(jen.Op("!").Id(valueVar).Dot("IsNull").Call()).
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
		idxName := tmpVarName("idx", nestDepth)
		methodBody.Var().Id(idxName).Int()
		methodBody.For().BlockFunc(func(forBody *jen.Group) {
			nestDepth := nestDepth + 1
			resultVar := tmpVarName("arrayValue", nestDepth)
			okVar := tmpVarName("ok", nestDepth)
			errVar := tmpVarName("err", nestDepth)
			listElement := tmpVarName("listElement", nestDepth)
			forBody.Var().Id(resultVar).Add(snip.DJResult())
			forBody.Var().Id(okVar).Bool()
			forBody.Var().Id(errVar).Error()
			forBody.List(jen.Id(resultVar), jen.Id(idxName), jen.Id(okVar), jen.Id(errVar)).Op("=").
				Id(valueVar).Dot("NextArrayEntry").Call(jen.Id(idxName))
			forBody.If(jen.Id(errVar).Op("!=").Nil()).Block(returnErrStmt())
			forBody.If(jen.Op("!").Id(okVar)).Block(jen.Break())
			forBody.Var().Id(listElement).Add(typ.Item.Code())
			unmarshalJSONValue(
				forBody,
				jen.Id(listElement).Clone,
				typ.Item,
				resultVar,
				jen.Return(snip.DJNewUnmarshalFieldError().Call(
					jen.Id(resultVar).Dot("Index").Call(),
					jen.Lit(fieldDescriptor+" list element"),
					jen.Err(),
				)).Clone,
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
		idxName := tmpVarName("idx", nestDepth)
		methodBody.Var().Id(idxName).Int()
		methodBody.For().BlockFunc(func(forBody *jen.Group) {
			nestDepth := nestDepth + 1
			keyVar := tmpVarName("mapKey", nestDepth)
			resultVar := tmpVarName("mapValue", nestDepth)
			okVar := tmpVarName("ok", nestDepth)
			errVar := tmpVarName("err", nestDepth)
			mapKeyVal := tmpVarName("mapKeyVal", nestDepth)
			forBody.Var().List(jen.Id(keyVar), jen.Id(resultVar)).Add(snip.DJResult())
			forBody.Var().Id(okVar).Bool()
			forBody.Var().Id(errVar).Error()
			forBody.List(jen.Id(keyVar), jen.Id(resultVar), jen.Id(idxName), jen.Id(okVar), jen.Id(errVar)).Op("=").
				Id(valueVar).Dot("NextObjectEntry").Call(jen.Id(idxName))
			forBody.If(jen.Id(errVar).Op("!=").Nil()).Block(returnErrStmt())
			forBody.If(jen.Op("!").Id(okVar)).Block(jen.Break())
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
					jen.Return(snip.DJNewUnmarshalFieldError().Call(
						jen.Id(keyVar).Dot("Index").Call(),
						jen.Lit(fieldDescriptor+" map key"),
						jen.Err(),
					)).Clone,
					fieldDescriptor+" map key",
					true,
					nestDepth+1,
					strict)
			})
			forBody.If(
				jen.List(jen.Id("_"), jen.Id("exists").Op(":=").Add(selector()).Index(jen.Id(mapKeyVal))),
				jen.Id("exists"),
			).Block(
				jen.Return(snip.DJNewUnmarshalDuplicateMapKeyError().Call(
					jen.Id(keyVar).Dot("Index").Call(),
					jen.Lit(fieldDescriptor),
				)),
			)
			mapVal := tmpVarName("mapVal", nestDepth)
			forBody.Var().Id(mapVal).Add(typ.Val.Code())
			forBody.BlockFunc(func(valBlock *jen.Group) {
				unmarshalJSONValue(
					valBlock,
					jen.Id(mapVal).Clone,
					typ.Val,
					resultVar,
					jen.Return(snip.DJNewUnmarshalFieldError().Call(
						jen.Id(resultVar).Dot("Index").Call(),
						jen.Lit(fieldDescriptor+" map value"),
						jen.Err(),
					)).Clone,
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
		methodBody.Err().Op("=").Add(snip.SafeJSONUnmarshal()).Call(jen.Index().Byte().Call(jen.Id(valueVar).Dot("Raw")), jen.Op("&").Add(selector()))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(returnErrStmt())
	}
}
