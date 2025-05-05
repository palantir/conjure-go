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

func aliasTypeItemSelector(typ *types.AliasType, selector *jen.Statement) func() *jen.Statement {
	if typ.IsOptional() {
		return selector.Dot("Value").Clone
	}
	if typ.IsCollection() {
		return selector.Clone
	}
	return typ.Item.Code().Call(selector).Clone
}
