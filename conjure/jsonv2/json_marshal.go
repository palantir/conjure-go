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
				jen.Id("b64len").Op(":=").Add(snip.Base64StdEncoding()).Dot("EncodedLen").Call(jen.Len(selector())),
				jen.Id("b64out").Op(":=").Make(
					jen.Index().Byte(),
					jen.Id("b64len").Op("+2"),
				),
				jen.Id("b64out").Index(jen.Lit(0)).Op("=").LitRune('"'),
				snip.Base64StdEncoding().Dot("Encode").Call(jen.Id("b64out").Index(jen.Op("1:")), selector()),
				jen.Id("b64out").Index(jen.Id("b64len").Op("+1")).Op("=").LitRune('"'),
				mustWriteValue(snip.JSONV2Value().Call(jen.Id("b64out"))),
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
		if typ.IsSimpleAliasType() {
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
					jen.If(jen.Id(bVar)).Block(jen.Return(jen.Lit(1))),
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
		// type cannot be sorted
		//return selector()
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
