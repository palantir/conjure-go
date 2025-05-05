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
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	encName = "enc"
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
				methodBody.Return(marshalJSONValue2(aliasTypeItemSelector(typ, jen.Id(receiverName)), typ.Item, false))
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
	return jen.Parens(getTypeArshaler(valueType, valueType.Code(), isMapKey, false).Values()).Dot("MarshalJSONTo").Call(selector(), jen.Id(encName))
}

func getTypeArshaler(valueType types.Type, declType *jen.Statement, isMapKey bool, isUnmarshal bool) *jen.Statement {
	switch typ := valueType.(type) {
	case types.Any:
		return snip.CJTypeAny().Types(declType)
	case types.String, types.Bearertoken:
		return snip.CJTypeString().Types(declType)
	case types.DateTime:
		return snip.CJTypeDateTime().Types(declType)
	case types.RID:
		return snip.CJTypeRID().Types(declType)
	case types.UUID:
		return snip.CJTypeUUID().Types(declType)
	case types.Boolean:
		if isMapKey {
			return snip.CJTypeBooleanMapKey().Types(declType)
		}
		return snip.CJTypeBoolean().Types(declType)
	case types.Double:
		if isMapKey {
			return snip.CJTypeFloatMapKey().Types(declType)
		}
		return snip.CJTypeFloat().Types(declType)
	case types.Integer, types.Safelong:
		if isMapKey {
			return snip.CJTypeIntMapKey().Types(declType)
		}
		return snip.CJTypeInt().Types(declType)
	case types.Binary:
		if isMapKey {
			return snip.CJTypeString().Types(snip.BinaryBinary())
		}
		return snip.CJTypeBinary().Types(declType)
	case *types.Optional:
		if isUnmarshal {
			return snip.CJTypeOptionalUnmarshaler().Types(typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), isMapKey, isUnmarshal))
		}
		return snip.CJTypeOptionalMarshaler().Types(typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), isMapKey, isUnmarshal))
	case *types.List:
		if isUnmarshal {
			return snip.CJTypeListUnmarshaler().Types(typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
		}
		return snip.CJTypeListMarshaler().Types(typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
	case *types.Set:
		if isUnmarshal {
			return snip.CJTypeListUnmarshaler().Types(typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
		}
		return snip.CJTypeListMarshaler().Types(typ.Item.Code(), getTypeArshaler(typ.Item, typ.Item.Code(), false, isUnmarshal))
	case *types.Map:
		var keyType *jen.Statement
		if typ.Key.IsBinary() {
			keyType = snip.BinaryBinary()
		} else if typ.Key.IsBoolean() {
			keyType = snip.BooleanBoolean()
		} else {
			keyType = typ.Key.Code()
		}
		key := getTypeArshaler(typ.Key, keyType, true, isUnmarshal)
		val := getTypeArshaler(typ.Val, typ.Val.Code(), false, isUnmarshal)
		if isUnmarshal {
			return snip.CJTypeMapUnmarshaler().Types(keyType, typ.Val.Code(), key, val)
		}
		if typ.Key.IsOrdered() {
			return snip.CJTypeOrderedMapMarshaler().Types(keyType, typ.Val.Code(), key, val)
		}
		return snip.CJTypeComparableMapMarshaler().Types(keyType, typ.Val.Code(), key, val)
	case *types.External:
		if typ.ExternalHasGoType() {
			return snip.CJTypeAny().Types(declType)
		}
		return getTypeArshaler(typ.Fallback, declType, isMapKey, isUnmarshal)
	case *types.AliasType:
		if typ.IsOptional() {
			if isUnmarshal {
				return snip.CJTypeStructUnmarshaler().Types(declType)
			}
			return snip.CJTypeStructMarshaler().Types(declType)
		}
		return getTypeArshaler(typ.Item, declType, isMapKey, isUnmarshal)
	case *types.EnumType:
		if isUnmarshal {
			return snip.CJTypeTextUnmarshaler().Types(jen.Op("*").Add(declType))
		}
		return snip.CJTypeTextMarshaler().Types(declType)
	case *types.ObjectType, *types.UnionType:
		if isUnmarshal {
			return snip.CJTypeStructUnmarshaler().Types(jen.Op("*").Add(declType))
		}
		return snip.CJTypeStructMarshaler().Types(declType)
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

//func marshalJSONValue(methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
//	switch typ := valueType.(type) {
//	case types.String:
//		methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector())))
//	case types.Bearertoken, types.DateTime, types.RID, types.UUID, *types.EnumType:
//		methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector().Dot("String").Call())))
//	case types.Binary:
//		if isMapKey {
//			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector().Dot("String").Call())))
//		} else {
//			methodBody.If(jen.Len(selector()).Op(">").Lit(0)).Block(
//				jen.Id("b64out").Op(":=").Id(encName).Dot("UnusedBuffer").Call(),
//				jen.Id("b64out").Op("=").Append(jen.Id("b64out"), jen.LitRune('"')),
//				jen.Id("b64out").Op("=").Add(snip.Base64StdEncoding()).Dot("AppendEncode").Call(jen.Id("b64out"), selector()),
//				jen.Id("b64out").Op("=").Append(jen.Id("b64out"), jen.LitRune('"')),
//				mustWriteValue(jen.Id("b64out")),
//			).Else().Block(
//				mustWriteToken(snip.JSONV2String().Call(jen.Lit(""))),
//			)
//		}
//	case types.Boolean:
//		if isMapKey {
//			methodBody.If(selector()).Block(
//				mustWriteToken(snip.JSONV2String().Call(jen.Lit("true"))),
//			).Else().Block(
//				mustWriteToken(snip.JSONV2String().Call(jen.Lit("false"))),
//			)
//		} else {
//			methodBody.If(selector()).Block(
//				mustWriteToken(snip.JSONV2True()),
//			).Else().Block(
//				mustWriteToken(snip.JSONV2False()),
//			)
//		}
//	case types.Double:
//		if isMapKey {
//			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(snip.StrconvFormatFloat().Call(
//				selector(), jen.LitRune('g'), jen.Lit(-1), jen.Lit(64),
//			))))
//		} else {
//			methodBody.Add(mustWriteToken(snip.JSONV2Float().Call(selector())))
//		}
//	case types.Integer, types.Safelong:
//		if isMapKey {
//			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(snip.StrconvFormatInt().Call(
//				jen.Int64().Call(selector()), jen.Lit(10),
//			))))
//		} else {
//			methodBody.Add(mustWriteToken(snip.JSONV2Int().Call(jen.Int64().Call(selector()))))
//		}
//	case *types.Optional:
//		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
//			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
//			marshalJSONValue(ifBody, jen.Id("optVal").Clone, typ.Item, nestDepth+1, isMapKey)
//		}).Else().Block(
//			mustWriteToken(snip.JSONV2Null()),
//		)
//	case *types.List:
//		methodBody.Add(mustWriteToken(snip.JSONV2BeginArray()))
//		i := tmpVarName("i", nestDepth)
//		methodBody.For(jen.List(jen.Op("_"), jen.Id(i)).
//			Op(":=").
//			Range().
//			Add(selector())).
//			BlockFunc(func(rangeBody *jen.Group) {
//				marshalJSONValue(rangeBody, jen.Id(i).Clone, typ.Item, nestDepth+1, false)
//			})
//		methodBody.Add(mustWriteToken(snip.JSONV2EndArray()))
//	case *types.Set:
//		methodBody.Add(mustWriteToken(snip.JSONV2BeginArray()))
//		i := tmpVarName("i", nestDepth)
//		methodBody.For(jen.List(jen.Op("_"), jen.Id(i)).
//			Op(":=").
//			Range().
//			Add(selector())).
//			BlockFunc(func(rangeBody *jen.Group) {
//				marshalJSONValue(rangeBody, jen.Id(i).Clone, typ.Item, nestDepth+1, false)
//			})
//		methodBody.Add(mustWriteToken(snip.JSONV2EndArray()))
//	case *types.Map:
//		methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
//		kVar := tmpVarName("k", nestDepth)
//		methodBody.For(jen.List(jen.Op("_"), jen.Id(kVar))).
//			Op(":=").
//			Range().
//			Add(sortedMapKeySequence(typ.Key, snip.MapsKeys().Call(selector()).Clone, nestDepth)).
//			BlockFunc(func(rangeBody *jen.Group) {
//				rangeBody.BlockFunc(func(keyBlock *jen.Group) {
//					marshalJSONValue(keyBlock, jen.Id(kVar).Clone, typ.Key, nestDepth+1, true)
//				})
//				rangeBody.BlockFunc(func(valueBlock *jen.Group) {
//					marshalJSONValue(valueBlock, selector().Index(jen.Id(kVar)).Clone, typ.Val, nestDepth+1, false)
//				})
//			})
//		methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
//	case *types.AliasType:
//		if typ.IsSimpleAliasType() || typ.IsBinary() {
//			// simple alias types don't have encoding methods, so access them directly
//			marshalJSONValue(methodBody, aliasTypeItemSelector(typ, selector()), typ.Item, nestDepth+1, isMapKey)
//		} else {
//			methodBody.If(
//				jen.Err().Op(":=").Add(selector().Dot("MarshalJSONTo").Call(jen.Id(encName))),
//				jen.Err().Op("!=").Nil(),
//			).Block(
//				jen.Return(jen.Err()),
//			)
//		}
//	case *types.ObjectType, *types.UnionType:
//		methodBody.If(
//			jen.Err().Op(":=").Add(selector().Dot("MarshalJSONTo").Call(jen.Id(encName))),
//			jen.Err().Op("!=").Nil(),
//		).Block(
//			jen.Return(jen.Err()),
//		)
//	case types.Any:
//		methodBody.If(
//			jen.Err().Op(":=").Add(snip.JSONV2MarshalEncode()).Call(jen.Id(encName), selector()),
//			jen.Err().Op("!=").Nil(),
//		).Block(
//			jen.Return(jen.Err()),
//		)
//	case *types.External:
//		if typ.ExternalHasGoType() {
//			methodBody.If(
//				jen.Err().Op(":=").Add(snip.JSONV2MarshalEncode()).Call(jen.Id(encName), selector()),
//				jen.Err().Op("!=").Nil(),
//			).Block(
//				jen.Return(jen.Err()),
//			)
//		} else {
//			marshalJSONValue(methodBody, selector, typ.Fallback, nestDepth, isMapKey)
//		}
//	default:
//		panic(fmt.Sprintf("unknown type %T", typ))
//	}
//}
//
//func sortedMapKeySequence(typ types.Type, selector func() *jen.Statement, nestedDepth int) *jen.Statement {
//	if typ.IsOrdered() {
//		return snip.SlicesSorted().Call(selector())
//	}
//	switch keyType := typ.(type) {
//	case *types.AliasType:
//		return sortedMapKeySequence(keyType.Item, aliasTypeItemSelector(keyType, selector()), nestedDepth+1)
//	case *types.EnumType:
//		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
//		return snip.SlicesSortedFunc().Call(
//			selector(),
//			jen.Func().
//				Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(typ.Code()))).
//				Int().
//				Block(
//					jen.Return(snip.CmpCompare().Call(
//						jen.Id(aVar).Dot("String").Call(),
//						jen.Id(bVar).Dot("String").Call(),
//					)),
//				))
//	case types.Boolean:
//		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
//		return snip.SlicesSortedFunc().Call(selector(), jen.Func().
//			Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(snip.BooleanBoolean()))).
//			Int().
//			BlockFunc(func(cmp *jen.Group) {
//				cmp.If(jen.Id(aVar).Op("!=").Id(bVar)).Block(
//					jen.If(jen.Id(aVar)).Block(jen.Return(jen.Lit(1))),
//					jen.Return(jen.Lit(-1)),
//				)
//				cmp.Return(jen.Lit(0))
//			}))
//	case types.DateTime:
//		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
//		return snip.SlicesSortedFunc().Call(selector(), jen.Func().
//			Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(snip.DateTimeDateTime()))).
//			Int().
//			Block(
//				jen.Return(snip.CmpCompare().Call(
//					snip.TimeTime().Call(jen.Id(aVar)).Dot("UnixNano").Call(),
//					snip.TimeTime().Call(jen.Id(bVar)).Dot("UnixNano").Call(),
//				)),
//			),
//		)
//	case types.RID:
//		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
//		return snip.SlicesSortedFunc().Call(
//			selector(),
//			jen.Func().
//				Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(snip.RIDResourceIdentifier()))).
//				Int().
//				BlockFunc(func(cmp *jen.Group) {
//					for _, field := range []string{"Service", "Type", "Instance", "Locator"} {
//						cmp.If(
//							jen.Id("c").Op(":=").Add(snip.StringsCompare()).Call(jen.Id(aVar).Dot(field), jen.Id(bVar).Dot(field)),
//							jen.Id("c").Op("!=").Lit(0),
//						).Block(
//							jen.Return(jen.Id("c")),
//						)
//					}
//					cmp.Return(jen.Lit(0))
//				}))
//	case types.UUID:
//		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
//		return snip.SlicesSortedFunc().Call(
//			selector(),
//			jen.Func().
//				Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(typ.Code()))).
//				Int().
//				Block(
//					jen.Return(snip.CmpCompare().Call(
//						jen.Id(aVar).Dot("String").Call(),
//						jen.Id(bVar).Dot("String").Call(),
//					)),
//				),
//		)
//	default:
//		panic(fmt.Errorf("cannot sort map key of type %T", typ))
//	}
//}

func aliasTypeItemSelector(typ *types.AliasType, selector *jen.Statement) func() *jen.Statement {
	if typ.IsOptional() {
		return selector.Dot("Value").Clone
	}
	if typ.IsCollection() {
		return selector.Clone
	}
	return typ.Item.Code().Call(selector).Clone
}

//	func(a, b rid.ResourceIdentifier) int {
//			return cmp.Compare(a.String(), b.String())
//		}
func quoteJSONString(s string) string {
	q, err := jsontext.AppendQuote(nil, s)
	if err != nil {
		panic(fmt.Sprintf("error quoting JSON string: %s", err))
	}
	return string(q)
}

func tmpVarName(base string, depth int) string {
	if depth == 0 {
		return base
	}
	return fmt.Sprintf("%s%d", base, depth)
}
