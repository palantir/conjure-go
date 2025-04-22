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
				marshalJSONValue(methodBody, jen.Id(receiverName).Clone, typ, 0, false)
				methodBody.Return(jen.Nil())
			case *types.AliasType:
				marshalJSONValue(methodBody, aliasTypeItemSelector(typ, jen.Id(receiverName)), typ.Item, 0, false)
				methodBody.Return(jen.Nil())
			case *types.EnumType:
				marshalJSONValue(methodBody, jen.Id(receiverName).Clone, typ, 0, false)
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
				methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
				for fieldIdx := range fields {
					marshalJSONStructField(methodBody, fields, fieldIdx)
				}
				methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
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
				methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
				methodBody.Switch(jen.Id(receiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
					for _, field := range fields {
						cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
							caseBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit("type"))))
							caseBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
							caseBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
								ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
								ifBody.Id("unionVal").Op(":=").Op("*").Add(field.Selector())
								marshalJSONValue(ifBody, jen.Id("unionVal").Clone, field.Type, 0, false)
							})
						})
					}
					cases.Default().Block(
						mustWriteToken(snip.JSONV2String().Call(jen.Lit("type"))),
						mustWriteToken(snip.JSONV2String().Call(jen.Id(receiverName).Dot("typ"))),
					)
				})
				methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
				methodBody.Return(jen.Nil())
			}
		})
}

func marshalJSONStructField(methodBody *jen.Group, fields []jsonStructField, fieldIdx int) {
	field := fields[fieldIdx]
	if field.Type.IsOptional() {
		switch typ := field.Type.(type) {
		case *types.Optional:
			methodBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
				ifBody.Id("optVal").Op(":=").Op("*").Add(field.Selector())
				marshalJSONValue(ifBody, jen.Id("optVal").Clone, typ.Item, 0, false)
			})
		case *types.AliasType:
			methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
				marshalJSONValue(ifBody, field.Selector, typ, 0, false)
			})
		default:
			panic(fmt.Sprintf("unexpected optional type %T", field.Type))
		}
	} else {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			fieldBlock.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
			marshalJSONValue(fieldBlock, field.Selector, field.Type, 0, false)
		})
	}
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

const (
	cjImport = "github.com/palantir/conjure-go/v6/cj"
)

var (
	snipCJTypeArshaler      = jen.Qual(cjImport, "TypeArshaler").Clone
	snipCJTypeEncoder       = jen.Qual(cjImport, "TypeEncoder").Clone
	snipCJTypeDecoder       = jen.Qual(cjImport, "TypeDecoder").Clone
	snipCJTypeAny           = jen.Qual(cjImport, "TypeAny").Clone
	snipCJTypeBinary        = jen.Qual(cjImport, "TypeBinary").Clone
	snipCJTypeBoolean       = jen.Qual(cjImport, "TypeBoolean").Clone
	snipCJTypeBooleanMapKey = jen.Qual(cjImport, "TypeBooleanMapKey").Clone
	snipCJTypeDateTime      = jen.Qual(cjImport, "TypeDateTime").Clone
	snipCJTypeFloat         = jen.Qual(cjImport, "TypeFloat").Clone
	snipCJTypeFloatMapKey   = jen.Qual(cjImport, "TypeFloatMapKey").Clone
	snipCJTypeInt           = jen.Qual(cjImport, "TypeInt").Clone
	snipCJTypeIntMapKey     = jen.Qual(cjImport, "TypeIntMapKey").Clone
	snipCJTypeUint          = jen.Qual(cjImport, "TypeUint").Clone
	snipCJTypeUintMapKey    = jen.Qual(cjImport, "TypeUintMapKey").Clone
	snipCJTypeRID           = jen.Qual(cjImport, "TypeRID").Clone
	snipCJTypeString        = jen.Qual(cjImport, "TypeString").Clone
	snipCJTypeUUID          = jen.Qual(cjImport, "TypeUUID").Clone
	snipCJTypeOptional      = jen.Qual(cjImport, "TypeOptional").Clone
	snipCJTypeList          = jen.Qual(cjImport, "TypeList").Clone
	snipCJTypeSortedMap     = jen.Qual(cjImport, "TypeSortedMap").Clone
	snipCJTypeOrderedMap    = jen.Qual(cjImport, "TypeOrderedMap").Clone
	snipCJTypeTextArshaler  = jen.Qual(cjImport, "TypeTextArshaler").Clone
)

//type TypeArshaler[T any] interface {
//type TypeEncoder[T any] interface {
//type TypeDecoder[T any] interface {
//type TypeAny[T any] struct{}
//type TypeBinary[T ~[]byte] struct{}
//type TypeBoolean[T ~bool] struct{}
//type TypeBooleanMapKey[T ~bool] struct{}
//type TypeDateTime[T time.Time | datetime.DateTime] struct{}
//type TypeFloat[T ~float64] struct{}
//type TypeFloatMapKey[T ~float64] struct{}
//type TypeInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}
//type TypeIntMapKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64] struct{}
//type TypeUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}
//type TypeUintMapKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64] struct{}
//type TypeRID[T rid.ResourceIdentifier] struct{}
//type TypeString[T ~string] struct{}
//type TypeUUID[T ~[16]byte] struct{}
//type TypeOptional[T any, ITEM TypeArshaler[T]] struct{}
//type TypeList[T any, ITEM TypeArshaler[T]] struct{}
//type TypeSortedMap[K comparable, V any, KEY MapKeyArshaler[K], VAL TypeArshaler[V]] struct{}
//type TypeOrderedMap[K cmp.Ordered, V any, KEY TypeArshaler[K], VAL TypeArshaler[V]] struct{}

func marshalJSONValue2(methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	methodBody.Parens(getTypeArshaler(valueType, valueType.Code(), isMapKey).Values()).Dot("MarshalJSONTo").Call(selector(), jen.Id(encName))
}

func getTypeArshaler(valueType types.Type, declType *jen.Statement, isMapKey bool) *jen.Statement {
	switch typ := valueType.(type) {
	case types.Any:
		return snipCJTypeAny().Types(declType)
	case types.String, types.Bearertoken:
		return snipCJTypeString().Types(declType)
	case types.DateTime:
		return snipCJTypeDateTime().Types(declType)
	case types.RID:
		return snipCJTypeRID().Types(declType)
	case types.UUID:
		return snipCJTypeUUID().Types(declType)
	case types.Boolean:
		if isMapKey {
			return snipCJTypeBooleanMapKey().Types(declType)
		}
		return snipCJTypeBoolean().Types(declType)
	case types.Double:
		if isMapKey {
			return snipCJTypeFloatMapKey().Types(declType)
		}
		return snipCJTypeFloat().Types(declType)
	case types.Integer, types.Safelong:
		if isMapKey {
			return snipCJTypeIntMapKey().Types(declType)
		}
		return snipCJTypeInt().Types(declType)
	case types.Binary:
		return snipCJTypeBinary().Types(declType)
	case *types.Optional:
		return snipCJTypeOptional().Types(declType, getTypeArshaler(typ.Item, typ.Item.Code(), isMapKey))
	case *types.List:
		return snipCJTypeList().Types(declType, getTypeArshaler(typ.Item, typ.Item.Code(), false))
	case *types.Set:
		return snipCJTypeList().Types(declType, getTypeArshaler(typ.Item, typ.Item.Code(), false))
	case *types.Map:
		if typ.Key.IsOrdered() {
			return snipCJTypeOrderedMap().Types(typ.Key.Code(), typ.Val.Code(), getTypeArshaler(typ.Key, typ.Key.Code(), true), getTypeArshaler(typ.Val, typ.Val.Code(), false))
		}
		return snipCJTypeSortedMap().Types(typ.Key.Code(), typ.Val.Code(), getTypeArshaler(typ.Key, typ.Key.Code(), true), getTypeArshaler(typ.Val, typ.Val.Code(), false))
	case *types.AliasType:
		// simple alias types don't have encoding methods, so access them directly
		return getTypeArshaler(typ.Item, declType, isMapKey)
	case *types.EnumType:
		return snipCJTypeTextArshaler().Types(declType)
	case *types.ObjectType:
		return snipCJTypeArshaler().Types(declType)
	case *types.UnionType:
		return snipCJTypeArshaler().Types(declType)
	case *types.External:
		if typ.ExternalHasGoType() {
			return snipCJTypeAny().Types(declType)
		}
		return getTypeArshaler(typ.Fallback, declType, isMapKey)
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

func marshalJSONValue(methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.String:
		methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector())))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID, *types.EnumType:
		methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector().Dot("String").Call())))
	case types.Binary:
		if isMapKey {
			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector().Dot("String").Call())))
		} else {
			methodBody.If(jen.Len(selector()).Op(">").Lit(0)).Block(
				jen.Id("b64out").Op(":=").Id(encName).Dot("UnusedBuffer").Call(),
				jen.Id("b64out").Op("=").Append(jen.Id("b64out"), jen.LitRune('"')),
				jen.Id("b64out").Op("=").Add(snip.Base64StdEncoding()).Dot("AppendEncode").Call(jen.Id("b64out"), selector()),
				jen.Id("b64out").Op("=").Append(jen.Id("b64out"), jen.LitRune('"')),
				mustWriteValue(jen.Id("b64out")),
			).Else().Block(
				mustWriteToken(snip.JSONV2String().Call(jen.Lit(""))),
			)
		}
	case types.Boolean:
		if isMapKey {
			methodBody.If(selector()).Block(
				mustWriteToken(snip.JSONV2String().Call(jen.Lit("true"))),
			).Else().Block(
				mustWriteToken(snip.JSONV2String().Call(jen.Lit("false"))),
			)
		} else {
			methodBody.If(selector()).Block(
				mustWriteToken(snip.JSONV2True()),
			).Else().Block(
				mustWriteToken(snip.JSONV2False()),
			)
		}
	case types.Double:
		if isMapKey {
			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(snip.StrconvFormatFloat().Call(
				selector(), jen.LitRune('g'), jen.Lit(-1), jen.Lit(64),
			))))
		} else {
			methodBody.Add(mustWriteToken(snip.JSONV2Float().Call(selector())))
		}
	case types.Integer, types.Safelong:
		if isMapKey {
			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(snip.StrconvFormatInt().Call(
				jen.Int64().Call(selector()), jen.Lit(10),
			))))
		} else {
			methodBody.Add(mustWriteToken(snip.JSONV2Int().Call(jen.Int64().Call(selector()))))
		}
	case *types.Optional:
		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
			marshalJSONValue(ifBody, jen.Id("optVal").Clone, typ.Item, nestDepth+1, isMapKey)
		}).Else().Block(
			mustWriteToken(snip.JSONV2Null()),
		)
	case *types.List:
		methodBody.Add(mustWriteToken(snip.JSONV2BeginArray()))
		i := tmpVarName("i", nestDepth)
		methodBody.For(jen.List(jen.Op("_"), jen.Id(i)).
			Op(":=").
			Range().
			Add(selector())).
			BlockFunc(func(rangeBody *jen.Group) {
				marshalJSONValue(rangeBody, jen.Id(i).Clone, typ.Item, nestDepth+1, false)
			})
		methodBody.Add(mustWriteToken(snip.JSONV2EndArray()))
	case *types.Set:
		methodBody.Add(mustWriteToken(snip.JSONV2BeginArray()))
		i := tmpVarName("i", nestDepth)
		methodBody.For(jen.List(jen.Op("_"), jen.Id(i)).
			Op(":=").
			Range().
			Add(selector())).
			BlockFunc(func(rangeBody *jen.Group) {
				marshalJSONValue(rangeBody, jen.Id(i).Clone, typ.Item, nestDepth+1, false)
			})
		methodBody.Add(mustWriteToken(snip.JSONV2EndArray()))
	case *types.Map:
		methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
		kVar := tmpVarName("k", nestDepth)
		methodBody.For(jen.List(jen.Op("_"), jen.Id(kVar))).
			Op(":=").
			Range().
			Add(sortedMapKeySequence(typ.Key, snip.MapsKeys().Call(selector()).Clone, nestDepth)).
			BlockFunc(func(rangeBody *jen.Group) {
				rangeBody.BlockFunc(func(keyBlock *jen.Group) {
					marshalJSONValue(keyBlock, jen.Id(kVar).Clone, typ.Key, nestDepth+1, true)
				})
				rangeBody.BlockFunc(func(valueBlock *jen.Group) {
					marshalJSONValue(valueBlock, selector().Index(jen.Id(kVar)).Clone, typ.Val, nestDepth+1, false)
				})
			})
		methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
	case *types.AliasType:
		if typ.IsSimpleAliasType() || typ.IsBinary() {
			// simple alias types don't have encoding methods, so access them directly
			marshalJSONValue(methodBody, aliasTypeItemSelector(typ, selector()), typ.Item, nestDepth+1, isMapKey)
		} else {
			methodBody.If(
				jen.Err().Op(":=").Add(selector().Dot("MarshalJSONTo").Call(jen.Id(encName))),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)
		}
	case *types.ObjectType, *types.UnionType:
		methodBody.If(
			jen.Err().Op(":=").Add(selector().Dot("MarshalJSONTo").Call(jen.Id(encName))),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		)
	case types.Any:
		methodBody.If(
			jen.Err().Op(":=").Add(snip.JSONV2MarshalEncode()).Call(jen.Id(encName), selector()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		)
	case *types.External:
		if typ.ExternalHasGoType() {
			methodBody.If(
				jen.Err().Op(":=").Add(snip.JSONV2MarshalEncode()).Call(jen.Id(encName), selector()),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)
		} else {
			marshalJSONValue(methodBody, selector, typ.Fallback, nestDepth, isMapKey)
		}
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

func sortedMapKeySequence(typ types.Type, selector func() *jen.Statement, nestedDepth int) *jen.Statement {
	if typ.IsOrdered() {
		return snip.SlicesSorted().Call(selector())
	}
	switch keyType := typ.(type) {
	case *types.AliasType:
		return sortedMapKeySequence(keyType.Item, aliasTypeItemSelector(keyType, selector()), nestedDepth+1)
	case *types.EnumType:
		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
		return snip.SlicesSortedFunc().Call(
			selector(),
			jen.Func().
				Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(typ.Code()))).
				Int().
				Block(
					jen.Return(snip.CmpCompare().Call(
						jen.Id(aVar).Dot("String").Call(),
						jen.Id(bVar).Dot("String").Call(),
					)),
				))
	case types.Boolean:
		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
		return snip.SlicesSortedFunc().Call(selector(), jen.Func().
			Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(snip.BooleanBoolean()))).
			Int().
			BlockFunc(func(cmp *jen.Group) {
				cmp.If(jen.Id(aVar).Op("!=").Id(bVar)).Block(
					jen.If(jen.Id(aVar)).Block(jen.Return(jen.Lit(1))),
					jen.Return(jen.Lit(-1)),
				)
				cmp.Return(jen.Lit(0))
			}))
	case types.DateTime:
		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
		return snip.SlicesSortedFunc().Call(selector(), jen.Func().
			Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(snip.DateTimeDateTime()))).
			Int().
			Block(
				jen.Return(snip.CmpCompare().Call(
					snip.TimeTime().Call(jen.Id(aVar)).Dot("UnixNano").Call(),
					snip.TimeTime().Call(jen.Id(bVar)).Dot("UnixNano").Call(),
				)),
			),
		)
	case types.RID:
		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
		return snip.SlicesSortedFunc().Call(
			selector(),
			jen.Func().
				Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(snip.RIDResourceIdentifier()))).
				Int().
				BlockFunc(func(cmp *jen.Group) {
					for _, field := range []string{"Service", "Type", "Instance", "Locator"} {
						cmp.If(
							jen.Id("c").Op(":=").Add(snip.StringsCompare()).Call(jen.Id(aVar).Dot(field), jen.Id(bVar).Dot(field)),
							jen.Id("c").Op("!=").Lit(0),
						).Block(
							jen.Return(jen.Id("c")),
						)
					}
					cmp.Return(jen.Lit(0))
				}))
	case types.UUID:
		aVar, bVar := tmpVarName("a", nestedDepth), tmpVarName("b", nestedDepth)
		return snip.SlicesSortedFunc().Call(
			selector(),
			jen.Func().
				Params(jen.List(jen.Id(aVar), jen.Id(bVar).Add(typ.Code()))).
				Int().
				Block(
					jen.Return(snip.CmpCompare().Call(
						jen.Id(aVar).Dot("String").Call(),
						jen.Id(bVar).Dot("String").Call(),
					)),
				),
		)
	default:
		panic(fmt.Errorf("cannot sort map key of type %T", typ))
	}
}

func aliasTypeItemSelector(typ *types.AliasType, selector *jen.Statement) func() *jen.Statement {
	if typ.IsOptional() {
		return selector.Dot("Value").Clone
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
