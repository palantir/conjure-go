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
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/dj"
)

func MarshalJSONMethods(receiverName string, receiverTypeName string, receiverType types.Type, withAuxiliary bool) []*jen.Statement {
	stmts := []*jen.Statement{
		snip.MethodMarshalJSON(receiverName, receiverTypeName).Block(
			jen.Id(outName).Op(":=").Make(jen.Index().Byte(), jen.Lit(0)), // jen.Id("size")),
			jen.If(
				jen.List(jen.Id("_"), jen.Err()).Op(":=").Id(receiverName).Dot("WriteJSON").Call(
					snip.DJNewAppender().Call(jen.Op("&").Id(outName)),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Id(outName), snip.DJValid().Call(jen.Id(outName))),
		),
		//snip.MethodAppendJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
		//	methodBody.List(jen.Id("_"), jen.Err()).
		//		Op(":=").
		//		Id(receiverName).
		//		Dot("WriteJSON").
		//		Call(djDot("NewAppender").Call(jen.Op("&").Id(outName)))
		//	methodBody.Return(jen.Id(outName), jen.Err())
		//}),
		//snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
		//	methodBody.Return(jen.Id(receiverName).Dot("WriteJSON").Call(snip.IODiscard()))
		//}),
	}
	switch v := receiverType.(type) {
	case *types.AliasType:
		var selector *jen.Statement
		if v.IsOptional() {
			selector = jen.Id(receiverName).Dot("Value")
		} else {
			selector = v.Item.Code().Call(jen.Id(receiverName))
		}
		stmts = append(stmts,
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				methodBody.Var().Id(outName).Int()
				n := 0
				marshalJSONValue(&n, methodBody, selector.Clone, v.Item, 0, false)
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
		)
	case *types.EnumType:
		stmts = append(stmts,
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				methodBody.Var().Id(outName).Int()
				n := 0
				methodBody.Add(djMarshalFunc(&n, snip.DJWriteString(), jen.Id(receiverName).Dot("String").Call()))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
		)
	case *types.ObjectType:
		var fields []jsonStructField
		for _, field := range v.Fields {
			fields = append(fields, jsonStructField{
				Key:      field.Name,
				Type:     field.Type,
				Selector: jen.Id(receiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
			})
		}
		stmts = append(stmts,
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				methodBody.Var().Id(outName).Int()
				n := 0
				methodBody.Add(djMarshalFunc(&n, snip.DJWriteOpenObject()))
				for i := range fields {
					marshalJSONStructField(&n, methodBody, fields, i)
				}
				methodBody.Add(djMarshalFunc(&n, snip.DJWriteCloseObject()))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
		)
	case *types.UnionType:
		var fields []jsonStructField
		for _, field := range v.Fields {
			fields = append(fields, jsonStructField{
				Key:      field.Name,
				Type:     field.Type,
				Selector: jen.Id(receiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone,
			})
		}
		stmts = append(stmts,
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				methodBody.Var().Id(outName).Int()
				n := 0
				methodBody.Add(djMarshalFunc(&n, snip.DJWriteOpenObject()))
				methodBody.Switch(jen.Id(receiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
					for _, field := range fields {
						cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
							caseBody.Add(djMarshalFunc(&n, snip.DJWriteLiteral(), jen.Lit(`"type":`+quoteJSONString(field.Key))))
							caseBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
								ifBody.Add(djMarshalFunc(&n, snip.DJWriteLiteral(), jen.Lit(","+quoteJSONString(field.Key)+":")))
								ifBody.Id("unionVal").Op(":=").Op("*").Add(field.Selector())
								marshalJSONValue(&n, ifBody, jen.Id("unionVal").Clone, field.Type, 0, false)
							})
						})
					}
					cases.Default().Block(
						djMarshalFunc(&n, snip.DJWriteLiteral(), jen.Lit(`"type":`)),
						djMarshalFunc(&n, snip.DJWriteString(), jen.Call(jen.Id(receiverName).Dot("typ"))),
					)
				})
				methodBody.Add(djMarshalFunc(&n, snip.DJWriteCloseObject()))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
		)
	default:
		panic(receiverType)
	}
	if withAuxiliary {
		stmts = append(stmts,
			snip.MethodMarshalYAMLSig(receiverName, receiverTypeName).Block(
				jen.Return(snip.DJMarshalYAML().Call(
					jen.Id(receiverName),
				)),
			),
		)
	}
	return stmts
}

type jsonStructField struct {
	Key      string
	Type     types.Type
	Selector func() *jen.Statement
}

func marshalJSONStructField(n *int, methodBody *jen.Group, fields []jsonStructField, fieldIdx int) {
	field := fields[fieldIdx]
	if field.Type.IsOptional() {
		switch typ := field.Type.(type) {
		case *types.Optional:
			methodBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(n, ifBody, fields, fieldIdx)
				ifBody.Add(djMarshalFunc(n, snip.DJWriteLiteral(), jen.Lit(quoteJSONString(field.Key)+":")))
				ifBody.Id("optVal").Op(":=").Op("*").Add(field.Selector())
				marshalJSONValue(n, ifBody, jen.Id("optVal").Clone, typ.Item, 0, false)
			})
		case *types.AliasType:
			methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(n, ifBody, fields, fieldIdx)
				ifBody.Add(djMarshalFunc(n, snip.DJWriteLiteral(), jen.Lit(quoteJSONString(field.Key)+":")))
				marshalJSONValue(n, ifBody, field.Selector, typ, 0, false)
			})
		default:
			panic(fmt.Sprintf("unexpected optional type %T", field.Type))
		}
	} else {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			appendCommaIfNotFirstField(n, fieldBlock, fields, fieldIdx)
			fieldBlock.Add(djMarshalFunc(n, snip.DJWriteLiteral(), jen.Lit(quoteJSONString(field.Key)+":")))
			marshalJSONValue(n, fieldBlock, field.Selector, field.Type, 0, false)
		})
	}
}

func marshalJSONValue(n *int, methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.Any:
		methodBody.If(
			selector().Op("==").Nil(),
		).Block(
			djMarshalFunc(n, snip.DJWriteNull()),
		).Else().Block(
			djMarshalFunc(n, snip.DJWriteObject(), selector()),
		)
	case types.String:
		methodBody.Add(djMarshalFunc(n, snip.DJWriteString(), selector()))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID, *types.EnumType:
		methodBody.Add(djMarshalFunc(n, snip.DJWriteString(), selector().Dot("String").Call()))
	case types.Binary:
		if isMapKey {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteString(), jen.String().Call(selector())))
		} else {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteBase64(), selector()))
		}
	case types.Boolean:
		if isMapKey {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteBoolString(), jen.Bool().Call(selector())))
		} else {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteBool(), jen.Bool().Call(selector())))
		}
	case types.Double:
		if isMapKey {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteFloatString(), selector()))
		} else {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteFloat(), selector()))
		}
	case types.Integer, types.Safelong:
		if isMapKey {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteIntString(), jen.Int64().Call(selector())))
		} else {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteInt(), jen.Int64().Call(selector())))
		}
	case *types.Optional:
		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
			marshalJSONValue(n, ifBody, jen.Id("optVal").Clone, typ.Item, nestDepth+1, isMapKey)
		}).Else().Block(
			djMarshalFunc(n, snip.DJWriteNull()),
		)
	case *types.List:
		i := tmpVarName("i", nestDepth)
		methodBody.Add(djMarshalFunc(n, snip.DJWriteOpenArray()))
		methodBody.For(jen.Id(i).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
			marshalJSONValue(n, rangeBody, selector().Index(jen.Id(i)).Clone, typ.Item, nestDepth+1, false)
			rangeBody.If(jen.Id(i).Op("<").Len(selector()).Op("-").Lit(1)).Block(
				djMarshalFunc(n, snip.DJWriteComma()),
			)
		})
		methodBody.Add(djMarshalFunc(n, snip.DJWriteCloseArray()))
	case *types.Map:
		methodBody.Add(djMarshalFunc(n, snip.DJWriteOpenObject()))
		methodBody.BlockFunc(func(mapBlock *jen.Group) {
			keyType := typ.Key.Code
			if typ.Key.IsBinary() {
				keyType = snip.BinaryBinary
			}
			if typ.Key.IsBoolean() {
				keyType = snip.BooleanBoolean
			}
			nestDepth := nestDepth + 1
			keyIdxVar := tmpVarName("k", nestDepth)
			mapKeys := tmpVarName("mapKeys", nestDepth)
			mapIdx := tmpVarName("i", nestDepth)
			// sort map keys
			var keyIdxVal func() *jen.Statement

			if typ.Key.IsOrdered() {
				// directly sortable
				mapBlock.Id(mapKeys).Op(":=").Make(jen.Index().Add(keyType()), jen.Lit(0), jen.Len(selector()))
				mapBlock.For(jen.Id(keyIdxVar).Op(":=").Range().Add(selector())).Block(
					jen.Id(mapKeys).Op("=").Append(jen.Id(mapKeys), jen.Id(keyIdxVar)),
				)
				mapBlock.Add(snip.SlicesSort().Call(jen.Id(mapKeys)))
				keyIdxVal = jen.Id(keyIdxVar).Clone
			} else if typ.Key.IsText() || typ.Key.IsBoolean() {
				// need to stringify before sorting
				mapKeysByString := tmpVarName("mapKeysByString", nestDepth)
				mapBlock.Id(mapKeysByString).Op(":=").Make(jen.Map(jen.String()).Add(keyType()), jen.Len(selector()))
				mapBlock.Id(mapKeys).Op(":=").Make(jen.Index().String(), jen.Lit(0), jen.Len(selector()))
				mapBlock.For(jen.Id(keyIdxVar).Op(":=").Range().Add(selector())).Block(
					jen.List(jen.Id("text"), jen.Err()).Op(":=").Id(keyIdxVar).Dot("MarshalText").Call(),
					jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Lit(0), jen.Err())),
					jen.Id("s").Op(":=").String().Call(jen.Id("text")),
					jen.Id(mapKeysByString).Index(jen.Id("s")).Op("=").Id(keyIdxVar),
					jen.Id(mapKeys).Op("=").Append(jen.Id(mapKeys), jen.Id("s")),
				)
				mapBlock.Add(snip.SlicesSort().Call(jen.Id(mapKeys)))
				keyIdxVal = jen.Id(mapKeysByString).Index(jen.Id(keyIdxVar)).Clone
			} else {
				panic("map key type is not ordered or text: " + typ.Key.String())
			}

			mapBlock.For(jen.List(jen.Id(mapIdx), jen.Id(keyIdxVar)).Op(":=").Range().Id(mapKeys)).
				BlockFunc(func(rangeBody *jen.Group) {
					rangeBody.If(jen.Id(mapIdx).Op(">").Lit(0)).Block(
						djMarshalFunc(n, snip.DJWriteComma()),
					)
					rangeBody.BlockFunc(func(keyBlock *jen.Group) {
						marshalJSONValue(n, keyBlock, keyIdxVal, typ.Key, nestDepth+1, true)
					})
					rangeBody.Add(djMarshalFunc(n, snip.DJWriteColon()))
					rangeBody.BlockFunc(func(valueBlock *jen.Group) {
						marshalJSONValue(n, valueBlock, selector().Index(keyIdxVal()).Clone, typ.Val, nestDepth+1, false)
					})
				})
		})
		methodBody.Add(djMarshalFunc(n, snip.DJWriteCloseObject()))
	case *types.AliasType:
		if typ.IsOptional() {
			marshalJSONValue(n, methodBody, selector().Dot("Value").Clone, typ.Item, nestDepth, isMapKey)
		} else {
			marshalJSONValue(n, methodBody, typ.Item.Code().Call(selector()).Clone, typ.Item, nestDepth, isMapKey)
		}
	case *types.ObjectType, *types.UnionType:
		nArg := "n" + strconv.Itoa(*n)
		*n++
		methodBody.List(jen.Id(nArg), jen.Err()).Op(":=").Add(selector().Dot("WriteJSON").Call(jen.Id(wName)))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Lit(0), jen.Err()))
		methodBody.Id(outName).Op("+=").Id(nArg)
	case *types.External:
		if typ.ExternalHasGoType() {
			methodBody.Add(djMarshalFunc(n, snip.DJWriteObject(), selector()))
		} else {
			marshalJSONValue(n, methodBody, selector, typ.Fallback, nestDepth, isMapKey)
		}

	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

// appendCommaIfNotFirstField adds the comma separator between struct fields.
// It should come before writing the field's key string.
// If fieldIdx == 0, this is a noop.
func appendCommaIfNotFirstField(n *int, g *jen.Group, fields []jsonStructField, fieldIdx int) {
	if fieldIdx == 0 {
		return
	}
	// If our field is preceded only by optional fields, we check whether any are non-nil
	// and only write the comma if we are not opening the object with '{'
	precedingRequired := false
	var anyPrecedingNotNil *jen.Statement
	for i := 0; i < fieldIdx; i++ {
		if !fields[i].Type.IsOptional() {
			precedingRequired = true
			break
		}
		// if optional field is an alias, need to add .Value
		selector := fields[i].Selector()
		if _, isAlias := fields[i].Type.(*types.AliasType); isAlias {
			selector = fields[i].Selector().Dot("Value")
		}
		// Chain them together with ||
		if anyPrecedingNotNil == nil {
			anyPrecedingNotNil = selector.Op("!=").Nil()
		} else {
			anyPrecedingNotNil.Op("||").Add(selector).Op("!=").Nil()
		}
	}
	if precedingRequired {
		g.Add(djMarshalFunc(n, snip.DJWriteComma()))
	} else {
		g.If(anyPrecedingNotNil).Block(djMarshalFunc(n, snip.DJWriteComma()))
	}
}

func djMarshalFunc(n *int, djFunc *jen.Statement, args ...jen.Code) *jen.Statement {
	nArg := "n" + strconv.Itoa(*n)
	*n++

	return jen.List(jen.Id(nArg), jen.Err()).Op(":=").Add(djFunc).Call(
		append([]jen.Code{jen.Id(wName)}, args...)...).
		Line().
		If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Lit(0), jen.Err())).
		Line().
		Id(outName).Op("+=").Id(nArg)
}

func quoteJSONString(s string) string {
	return dj.QuoteString(s)
}
