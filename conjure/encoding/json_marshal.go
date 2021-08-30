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
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/pkg/safejson"
)

const (
	outName = "out"
)

func MarshalJSONMethodBody(receiverName string) *jen.Statement {
	return jen.Return(jen.Id(receiverName).Dot("AppendJSON").Call(jen.Nil()))
}

func AliasMethodBodyAppendJSON(methodBody *jen.Group, aliasType types.Type, selector func() *jen.Statement) {
	appendMarshalBufferJSONValue(methodBody, selector, aliasType, 0, false)
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func EnumMethodBodyAppendJSON(methodBody *jen.Group, receiverName string) {
	methodBody.Add(appendMarshalBufferQuotedString(jen.String().Call(jen.Id(receiverName).Dot("val"))))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

type JSONStructField struct {
	Key      string
	Type     types.Type
	Selector func() *jen.Statement
}

func StructMethodBodyAppendJSON(methodBody *jen.Group, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	for i, field := range fields {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			fieldBlock.Add(appendMarshalBufferVariadic(jen.Lit(safejson.QuoteString(field.Key) + ":")))
			appendMarshalBufferJSONValue(fieldBlock, field.Selector, field.Type, 0, false)

			if i < len(fields)-1 {
				fieldBlock.Add(appendMarshalBufferLiteralRune(','))
			}
		})
	}
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func UnionMethodBodyAppendJSON(methodBody *jen.Group, typeFieldSelctor func() *jen.Statement, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	methodBody.Switch(typeFieldSelctor()).BlockFunc(func(cases *jen.Group) {
		cases.Default().Block(
			appendMarshalBufferVariadic(jen.Lit(`"type":`)),
			appendMarshalBufferQuotedString(typeFieldSelctor()),
		)
		for _, fieldDef := range fields {
			cases.Case(jen.Lit(fieldDef.Key)).BlockFunc(func(caseBody *jen.Group) {
				caseBody.Add(appendMarshalBufferVariadic(jen.Lit(`"type":` + safejson.QuoteString(fieldDef.Key))))
				caseBody.If(fieldDef.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
					ifBody.Add(appendMarshalBufferLiteralRune(','))
					ifBody.Add(appendMarshalBufferLiteralString(fieldDef.Key))
					ifBody.Add(appendMarshalBufferLiteralRune(':'))
					ifBody.Id("unionVal").Op(":=").Op("*").Add(fieldDef.Selector())
					appendMarshalBufferJSONValue(ifBody, jen.Id("unionVal").Clone, fieldDef.Type, 0, false)
				})
			})
		}
	})
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func AnonFuncBodyAppendJSON(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	appendMarshalBufferJSONValue(funcBody, selector, valueType, 0, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
}

func appendMarshalBufferJSONValue(methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.String:
		methodBody.Add(appendMarshalBufferQuotedString(selector()))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID:
		methodBody.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
	case types.Any:
		methodBody.If(
			selector().Op("==").Nil(),
		).Block(
			appendMarshalBufferLiteralNull(),
		).Else().If(
			jen.List(jen.Id("jsonBytes"), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		).Else().Block(
			appendMarshalBufferVariadic(jen.Id("jsonBytes")),
		)
	case types.Binary:
		if isMapKey {
			methodBody.Add(appendMarshalBufferQuotedString(jen.String().Call(selector())))
		} else {
			methodBody.Add(appendMarshalBufferLiteralRune('"'))
			methodBody.If(jen.Len(selector()).Op(">").Lit(0)).Block(
				jen.Id("b64out").Op(":=").Make(
					jen.Op("[]").Byte(),
					snip.Base64EncodedLen().Call(jen.Len(selector())),
				),
				snip.Base64Encode().Call(jen.Id("b64out"), selector()),
				appendMarshalBufferVariadic(jen.Id("b64out")),
			)
			methodBody.Add(appendMarshalBufferLiteralRune('"'))
		}
	case types.Boolean:
		if isMapKey {
			methodBody.If(selector()).Block(
				appendMarshalBufferVariadic(jen.Lit(`"true"`)),
			).Else().Block(
				appendMarshalBufferVariadic(jen.Lit(`"false"`)),
			)
		} else {
			methodBody.If(selector()).Block(
				appendMarshalBufferVariadic(jen.Lit("true")),
			).Else().Block(
				appendMarshalBufferVariadic(jen.Lit("false")),
			)
		}
	case types.Double:
		methodBody.Switch().Block(
			jen.Default().BlockFunc(func(caseBody *jen.Group) {
				if isMapKey {
					caseBody.Add(appendMarshalBufferLiteralRune('"'))
				}
				caseBody.Id(outName).Op("=").Add(snip.StrconvAppendFloat()).Call(jen.Id(outName), selector(), jen.LitRune('g'), jen.Lit(-1), jen.Lit(64))
				if isMapKey {
					caseBody.Add(appendMarshalBufferLiteralRune('"'))
				}
			}),
			jen.Case(snip.MathIsNaN().Call(selector())).Block(appendMarshalBufferLiteralString("NaN")),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(1))).Block(appendMarshalBufferLiteralString("Infinity")),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(-1))).Block(appendMarshalBufferLiteralString("-Infinity")),
		)
	case types.Integer, types.Safelong:
		if isMapKey {
			methodBody.Add(appendMarshalBufferLiteralRune('"'))
		}
		methodBody.Id(outName).Op("=").Add(snip.StrconvAppendInt()).Call(jen.Id(outName), jen.Int64().Call(selector()), jen.Lit(10))
		if isMapKey {
			methodBody.Add(appendMarshalBufferLiteralRune('"'))
		}
	case *types.Optional:
		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
			appendMarshalBufferJSONValue(ifBody, jen.Id("optVal").Clone, typ.Item, nestDepth+1, isMapKey)
		}).Else().Block(
			appendMarshalBufferLiteralNull(),
		)
	case *types.List:
		i := tmpVarName("i", nestDepth)
		methodBody.Add(appendMarshalBufferLiteralRune('['))
		methodBody.For(jen.Id(i).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
			appendMarshalBufferJSONValue(rangeBody, selector().Index(jen.Id(i)).Clone, typ.Item, nestDepth+1, false)
			rangeBody.If(jen.Id(i).Op("<").Len(selector()).Op("-").Lit(1)).Block(
				appendMarshalBufferLiteralRune(','),
			)
		})
		methodBody.Add(appendMarshalBufferLiteralRune(']'))
	case *types.Map:
		methodBody.Add(appendMarshalBufferLiteralRune('{'))
		i := tmpVarName("i", nestDepth)
		methodBody.Block(
			jen.Var().Id(i).Int(),
			jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
				rangeBody.BlockFunc(func(keyBlock *jen.Group) {
					appendMarshalBufferJSONValue(keyBlock, jen.Id("k").Clone, typ.Key, nestDepth+1, true)
				})
				rangeBody.Add(appendMarshalBufferLiteralRune(':'))
				rangeBody.BlockFunc(func(valueBlock *jen.Group) {
					appendMarshalBufferJSONValue(valueBlock, jen.Id("v").Clone, typ.Val, nestDepth+1, false)
				})
				rangeBody.Id(i).Op("++")
				rangeBody.If(jen.Id(i).Op("<").Len(selector())).Block(
					appendMarshalBufferLiteralRune(','),
				)
			}),
		)
		methodBody.Add(appendMarshalBufferLiteralRune('}'))
	case *types.AliasType, *types.EnumType, *types.ObjectType, *types.UnionType:
		methodBody.Var().Err().Error()
		methodBody.List(jen.Id(outName), jen.Err()).Op("=").Add(selector()).Dot("AppendJSON").Call(jen.Id(outName))
		methodBody.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		)
	case *types.External:
		methodBody.If(
			jen.List(jen.Id("jsonBytes"), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		).Else().Block(
			appendMarshalBufferVariadic(jen.Id("jsonBytes")),
		)
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

func appendMarshalBufferLiteralNull() *jen.Statement {
	return appendMarshalBufferVariadic(jen.Lit("null"))
}

func appendMarshalBufferLiteralRune(r rune) *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.LitRune(r))
}

func appendMarshalBufferLiteralString(s string) *jen.Statement {
	return appendMarshalBufferVariadic(jen.Lit(safejson.QuoteString(s)))
}

func appendMarshalBufferQuotedString(selector *jen.Statement) *jen.Statement {
	return jen.Id(outName).Op("=").Add(snip.SafeJSONAppendQuotedString()).Call(jen.Id(outName), selector)
}

func appendMarshalBufferVariadic(selector *jen.Statement) *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), selector.Op("..."))
}
