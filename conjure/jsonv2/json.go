// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package jsonv2

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	encName = "enc"
	decName = "dec"
)

func MarshalJSONMethod(receiverName string, receiverTypeName string) *jen.Statement {
	return snip.MethodMarshalJSON(receiverName, receiverTypeName).Block(
		jen.Return(snip.JSONV2Marshal().Call(snip.JSONV2MarshalerTo().Call(jen.Id(receiverName)))),
	)
}

func MarshalJSONToMethod(receiverName string, receiverTypeName string, receiverType types.Type) *jen.Statement {
	return jen.Func().
		Params(jen.Id(receiverName).Id(receiverTypeName)).
		Id("MarshalJSONTo").
		Params(jen.Id(encName).Op("*").Add(snip.JSONV2Encoder())).
		Error().
		BlockFunc(func(methodBody *jen.Group) {
			switch typ := receiverType.(type) {
			default:
				methodBody.Return(marshalJSONValue2(jen.Id(receiverName).Clone, typ, false))
			case *types.AliasType:
				methodBody.Return(marshalJSONValue2(aliasTypeItemSelector(typ, jen.Id(receiverName), false), typ.Item, false))
			case *types.EnumType:
				methodBody.Return(marshalJSONValue2(jen.Id(receiverName).Clone, typ, false))
			case *types.ObjectType:
				var fields []jsonStructField
				for _, field := range typ.Fields {
					fields = append(fields, jsonStructField{
						Key:      field.Name,
						Type:     field.Type,
						Selector: jen.Id(receiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
					})
				}
				methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
				for _, field := range fields {
					if field.Type.IsOptional() {
						switch typ := field.Type.(type) {
						case *types.Optional:
							methodBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
								ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
								ifBody.If(
									jen.Err().Op(":=").Add(marshalJSONValue2(jen.Op("*").Add(field.Selector()).Clone, typ.Item, false)),
									jen.Err().Op("!=").Nil(),
								).Block(
									jen.Return(jen.Err()),
								)
							})
						case *types.AliasType:
							methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
								ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
								ifBody.If(
									jen.Err().Op(":=").Add(marshalJSONValue2(field.Selector().Dot("Value").Clone, typ.Item, false)),
									jen.Err().Op("!=").Nil(),
								).Block(
									jen.Return(jen.Err()),
								)
							})
						default:
							panic(fmt.Sprintf("unexpected optional type %T", field.Type))
						}
					} else {
						methodBody.BlockFunc(func(fieldBlock *jen.Group) {
							fieldBlock.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
							fieldBlock.If(
								jen.Err().Op(":=").Add(marshalJSONValue2(field.Selector, field.Type, false)),
								jen.Err().Op("!=").Nil(),
							).Block(
								jen.Return(jen.Err()),
							)
						})
					}
				}
				methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
				methodBody.Return(jen.Nil())
			case *types.UnionType:
				methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
				methodBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit("type"))))
				methodBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Id(receiverName).Dot("typ"))))
				methodBody.Switch(jen.Id(receiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
					for _, field := range typ.Fields {
						fieldSelector := jen.Id(receiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone
						cases.Case(jen.Lit(field.Name)).BlockFunc(func(caseBody *jen.Group) {
							caseBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Name))))
							caseBody.If(fieldSelector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
								ifBody.If(
									jen.Err().Op(":=").Add(marshalJSONValue2(jen.Op("*").Add(fieldSelector()).Clone, field.Type, false)),
									jen.Err().Op("!=").Nil(),
								).Block(
									jen.Return(jen.Err()),
								)
							}).Else().Block(
								mustWriteToken(snip.JSONV2Null()),
							)
						})
					}
				})
				methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
				methodBody.Return(jen.Nil())
			}
		})
}

func mustWriteToken(token jen.Code) *jen.Statement {
	return jen.If(
		jen.Err().Op(":=").Id(encName).Dot("WriteToken").Call(token),
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(jen.Err()),
	)
}

func mustWriteValue(value jen.Code) *jen.Statement {
	return jen.If(
		jen.Err().Op(":=").Id(encName).Dot("WriteValue").Call(value),
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(jen.Err()),
	)
}

type jsonStructField struct {
	Key      string
	Type     types.Type
	Selector func() *jen.Statement
}

func marshalJSONValue2(selector func() *jen.Statement, valueType types.Type, isMapKey bool) *jen.Statement {
	return jen.Parens(getTypeArshaler(valueType, valueType.Code, isMapKey, false).Values()).Dot("MarshalJSONTo").Call(selector(), jen.Id(encName))
}

func UnmarshalJSONMethod(receiverName string, receiverTypeName string) *jen.Statement {
	return snip.MethodUnmarshalJSON(receiverName, receiverTypeName).Block(
		jen.Return(snip.JSONV2Unmarshal().Call(jen.Id("data"), snip.JSONV2UnmarshalerFrom().Call(jen.Id(receiverName)))),
	)
}

func UnmarshalJSONFromMethod(receiverName string, receiverTypeName string, receiverType types.Type) *jen.Statement {
	return jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverTypeName)).
		Id("UnmarshalJSONFrom").
		Params(jen.Id(decName).Op("*").Add(snip.JSONV2Decoder())).
		Error().
		BlockFunc(func(methodBody *jen.Group) {
			switch typ := receiverType.(type) {
			default:
				methodBody.Return(unmarshalJSONValue(jen.Id(receiverName).Clone, receiverType, false))
			case *types.AliasType:
				methodBody.Return(unmarshalJSONValue(aliasTypeItemSelector(typ, jen.Id(receiverName), true), typ.Item, false))
			case *types.EnumType:
				methodBody.Return(unmarshalJSONValue(jen.Id(receiverName).Clone, typ, false))
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
			}
		})
}

func unmarshalJSONStructFields(methodBody *jen.Group, receiverName string, receiverType string, fields []jsonStructField, isUnion bool) {
	methodBody.If(
		jen.List(jen.Id("tok"), jen.Err()).Op(":=").Id(decName).Dot("ReadToken").Call(),
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(jen.Err()),
	).Else().If(
		jen.Id("tok").Dot("Kind").Call().Op("!=").LitRune('{'),
	).Block(
		jen.Return(snip.CJNewSyntaxError().Call(jen.Id(decName), jen.Lit(receiverType+" expected opening brace"))),
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
		fieldResults = append(fieldResults, result)
	}
	for _, result := range fieldResults {
		if result.Init != nil {
			result.Init(methodBody)
		}
	}
	methodBody.List(jen.Id("strict"), jen.Op("_")).Op(":=").Add(snip.JSONV2GetOption()).Call(
		jen.Id(decName).Dot("Options").Call(),
		snip.JSONV2RejectUnknownMembers())
	methodBody.Var().Id("unknownMembers").Index().String()
	methodBody.For().BlockFunc(func(forBody *jen.Group) {
		forBody.List(jen.Id("key"), jen.Err()).Op(":=").Id(decName).Dot("ReadToken").Call()
		forBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		)
		forBody.If(
			jen.Id("kind").Op(":=").Id("key").Dot("Kind").Call(),
			jen.Id("kind").Op("!=").LitRune('"'),
		).Block(
			jen.If(jen.Id("kind").Op("==").LitRune('}')).Block(
				jen.Break().Comment("End of object"),
			),
			jen.Return(snip.CJNewSyntaxError().Call(jen.Id(decName), jen.Lit(receiverName+" expected string key or closing brace"))),
		)
		forBody.Switch(jen.Id("key").Dot("String").Call()).BlockFunc(func(cases *jen.Group) {
			for _, result := range fieldResults {
				if result.Unmarshal != nil {
					result.Unmarshal(cases)
				}
			}
			cases.Default().Block(
				jen.If(jen.Id("strict")).Block(
					jen.Id("unknownMembers").
						Op("=").
						Append(jen.Id("unknownMembers"), jen.Id("key").Dot("String").Call()),
				),
			)
		})
	})

	if hasRequiredFields {
		methodBody.Var().Id("missingFields").Index().String()
		for _, result := range fieldResults {
			if result.Validate != nil {
				result.Validate(methodBody)
			}
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
		methodBody.If(jen.Len(jen.Id("missingFields")).Op(">").Lit(0)).Block(
			jen.Return(snip.CJNewMissingRequiredFieldsError().Call(
				jen.Id(decName),
				jen.Lit(receiverType),
				jen.Id("missingFields"),
			)),
		)
	} else if hasCollections {
		for _, result := range fieldResults {
			if result.DefaultCollection != nil {
				result.DefaultCollection(methodBody)
			}
		}
	}
	methodBody.If(jen.Id("strict").Op("&&").Len(jen.Id("unknownMembers")).Op(">").Lit(0)).Block(
		jen.Return(snip.CJNewUnknownFieldsError().Call(
			jen.Id(decName),
			jen.Lit(receiverType),
			jen.Id("unknownMembers"),
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
	seenVar := "seen" + transforms.ExportedFieldName(field.Key)
	result.Init = func(methodBody *jen.Group) {
		methodBody.Var().Id(seenVar).Bool()
	}
	requiredField := !(field.Type.IsCollection() || field.Type.IsOptional())
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
				jen.Return(snip.CJNewDuplicateFieldKeyError().Call(
					jen.Id(decName),
					jen.Lit(receiverType),
					jen.Lit(field.Key),
				)),
			)
			caseBody.Id(seenVar).Op("=").True()

			selector := jen.Op("&").Add(field.Selector()).Clone
			if isUnionField {
				caseBody.Add(field.Selector()).Op("=").New(field.Type.Code())
				selector = field.Selector
			}
			caseBody.If(
				jen.Err().Op(":=").Add(unmarshalJSONValue(selector, field.Type, false)),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)
		})
	}

	return result
}

func unmarshalJSONValue(
	selector func() *jen.Statement,
	valueType types.Type,
	isMapKey bool,
) *jen.Statement {
	return jen.Parens(getTypeArshaler(valueType, valueType.Code, isMapKey, true).Values()).Dot("UnmarshalJSONFrom").Call(selector(), jen.Id(decName))
}

func aliasTypeItemSelector(typ *types.AliasType, selector *jen.Statement, isUnmarshal bool) func() *jen.Statement {
	if typ.IsOptional() {
		if isUnmarshal {
			return jen.Op("&").Add(selector).Dot("Value").Clone
		}
		return selector.Dot("Value").Clone
	}
	if isUnmarshal {
		return jen.Parens(jen.Op("*").Add(typ.Item.Code())).Call(selector).Clone
	}
	if typ.IsCollection() {
		return selector.Clone
	}
	return typ.Item.Code().Call(selector).Clone
}

func getTypeArshaler(valueType types.Type, declType func() *jen.Statement, isMapKey bool, isUnmarshal bool) *jen.Statement {
	switch typ := valueType.(type) {
	case types.Any:
		return snip.CJTypeAny().Types(declType())
	case types.String, types.Bearertoken:
		return snip.CJTypeString().Types(declType())
	case types.DateTime:
		return snip.CJTypeDateTime().Types(declType())
	case types.RID:
		return snip.CJTypeRID().Types(declType())
	case types.UUID:
		return snip.CJTypeUUID().Types(declType())
	case types.Boolean:
		if isMapKey {
			return snip.CJTypeBooleanMapKey().Types(declType())
		}
		return snip.CJTypeBoolean().Types(declType())
	case types.Double:
		if isMapKey {
			return snip.CJTypeFloatMapKey().Types(declType())
		}
		return snip.CJTypeFloat().Types(declType())
	case types.Integer:
		if isMapKey {
			return snip.CJTypeInt32MapKey().Types(declType())
		}
		return snip.CJTypeInt32().Types(declType())
	case types.Safelong:
		if isMapKey {
			return snip.CJTypeSafeLongMapKey().Types(declType())
		}
		return snip.CJTypeSafeLong().Types(declType())
	case types.Binary:
		if isMapKey {
			return snip.CJTypeBinaryMapKey().Types(declType())
		}
		return snip.CJTypeBinary().Types(declType())
	case *types.Optional:
		if isUnmarshal {
			return snip.CJTypeOptionalUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, isMapKey, isUnmarshal))
		}
		return snip.CJTypeOptionalMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, isMapKey, isUnmarshal))
	case *types.List:
		if isUnmarshal {
			return snip.CJTypeListUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
		}
		return snip.CJTypeListMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
	case *types.Set:
		if isUnmarshal {
			return snip.CJTypeListUnmarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
		}
		return snip.CJTypeListMarshaler().Types(declType(), typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code, false, isUnmarshal))
	case *types.Map:
		var keyType *jen.Statement
		if typ.Key.IsBinary() {
			keyType = snip.BinaryBinary()
		} else if typ.Key.IsBoolean() {
			keyType = snip.BooleanBoolean()
		} else {
			keyType = typ.Key.Code()
		}
		key := getTypeArshaler(typ.Key, keyType.Clone, true, isUnmarshal)
		val := getTypeArshaler(typ.Val, typ.Val.Code, false, isUnmarshal)
		typeArgs := jen.Types(declType(), keyType, typ.Val.Code(), key, val)
		switch {
		case isUnmarshal:
			return snip.CJTypeMapUnmarshaler().Add(typeArgs)
		case typ.Key.IsOrdered():
			return snip.CJTypeOrderedMapMarshaler().Add(typeArgs)
		default:
			return snip.CJTypeComparableMapMarshaler().Add(typeArgs)
		}
	case *types.External:
		if typ.ExternalHasGoType() {
			return snip.CJTypeAny().Types(declType())
		}
		return getTypeArshaler(typ.Fallback, declType, isMapKey, isUnmarshal)
	case *types.AliasType:
		if typ.Item.IsOptional() {
			if isUnmarshal {
				return snip.CJTypeStructUnmarshaler().Types(jen.Op("*").Add(declType()))
			}
			return snip.CJTypeStructMarshaler().Types(declType())
		}
		return getTypeArshaler(typ.Item, declType, isMapKey, isUnmarshal)
	case *types.EnumType:
		if isUnmarshal {
			return snip.CJTypeTextUnmarshaler().Types(jen.Op("*").Add(declType()))
		}
		return snip.CJTypeStringerMarshaler().Types(declType())
	case *types.ObjectType, *types.UnionType:
		if isUnmarshal {
			return snip.CJTypeStructUnmarshaler().Types(jen.Op("*").Add(declType()))
		}
		return snip.CJTypeStructMarshaler().Types(declType())
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}
