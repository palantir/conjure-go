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

func AnonFuncBodyAppendJSON(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	appendMarshalBufferJSONValue(funcBody, selector, valueType, 0, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
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
	for i := range fields {
		appendMarshalBufferJSONStructField(methodBody, fields, i)
	}
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func UnionMethodBodyAppendJSON(methodBody *jen.Group, typeFieldSelctor func() *jen.Statement, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	methodBody.Switch(typeFieldSelctor()).BlockFunc(func(cases *jen.Group) {
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
		cases.Default().Block(
			appendMarshalBufferVariadic(jen.Lit(`"type":`)),
			appendMarshalBufferQuotedString(typeFieldSelctor()),
		)
	})
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func appendMarshalBufferJSONStructField(methodBody *jen.Group, fields []JSONStructField, fieldIdx int) {
	// appendCommaIfNotFirstField adds the comma separator between fields.
	// It should come before writing the field's key string.
	// If fieldIdx == 0, this is a noop.
	appendCommaIfNotFirstField := func(g *jen.Group) {
		if fieldIdx == 0 {
			return
		}
		// If our field is preceded only by optional fields, we check the last byte of the out buffer
		// and only write the comma if we are not opening the object with '{'
		precedingRequired := false
		for i := 0; i < fieldIdx; i++ {
			if !fields[i].Type.IsOptional() {
				precedingRequired = true
			}
		}
		if precedingRequired {
			g.Add(appendMarshalBufferLiteralRune(','))
		} else {
			// Creates `if out[len(out)-1] != '{' {	out = append(out, ',') }`
			// No need to check for empty because we have definitely written the '{' byte by now.
			g.If(
				jen.Id(outName).Index(jen.Len(jen.Id(outName)).Op("-").Lit(1)).Op("!=").LitRune('{'),
			).Block(
				appendMarshalBufferLiteralRune(','),
			)
		}
	}

	field := fields[fieldIdx]
	if field.Type.IsOptional() {
		switch typ := field.Type.(type) {
		case *types.Optional:
			methodBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(ifBody)
				ifBody.Add(appendMarshalBufferVariadic(jen.Lit(safejson.QuoteString(field.Key) + ":")))
				ifBody.Id("optVal").Op(":=").Op("*").Add(field.Selector())
				appendMarshalBufferJSONValue(ifBody, jen.Id("optVal").Clone, typ.Item, 0, false)
			})
		case *types.AliasType:
			methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(ifBody)
				ifBody.Add(appendMarshalBufferVariadic(jen.Lit(safejson.QuoteString(field.Key) + ":")))
				appendMarshalBufferJSONValue(ifBody, field.Selector, field.Type, 0, false)
			})
		default:
			panic(fmt.Sprintf("unexpected optional type %T", field.Type))
		}
	} else {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			appendCommaIfNotFirstField(fieldBlock)
			fieldBlock.Add(appendMarshalBufferVariadic(jen.Lit(safejson.QuoteString(field.Key) + ":")))
			appendMarshalBufferJSONValue(fieldBlock, field.Selector, field.Type, 0, false)
		})
	}
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
			jen.List(jen.Id("appender"), jen.Id("ok")).Op(":=").Add(selector()).Assert(
				jen.Interface(
					jen.Id("AppendJSON").
						Params(jen.Op("[]").Byte()).
						Params(jen.Op("[]").Byte(), jen.Error()),
				),
			),
			jen.Id("ok"),
		).Block(
			jen.Var().Err().Error(),
			jen.List(jen.Id(outName), jen.Err()).Op("=").Id("appender").Dot("AppendJSON").Call(jen.Id(outName)),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
		).Else().If(
			jen.List(jen.Id("marshaler"), jen.Id("ok")).Op(":=").Add(selector()).Assert(snip.JSONMarshaler()),
			jen.Id("ok"),
		).Block(
			jen.List(jen.Id(dataName), jen.Err()).Op(":=").Id("marshaler").Dot("MarshalJSON").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			appendMarshalBufferVariadic(jen.Id(dataName)),
		).Else().If(
			jen.List(jen.Id(dataName), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		).Else().Block(
			appendMarshalBufferVariadic(jen.Id(dataName)),
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
		if typ.ExternalHasGoType() {
			methodBody.If(
				jen.List(jen.Id(dataName), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Nil(), jen.Err()),
			).Else().Block(
				appendMarshalBufferVariadic(jen.Id(dataName)),
			)
		} else {
			appendMarshalBufferJSONValue(methodBody, selector, typ.Fallback, nestDepth, isMapKey)
		}

	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

func appendMarshalBufferLiteralNull() *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.Lit("null").Op("..."))
}

func appendMarshalBufferLiteralRune(r rune) *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.LitRune(r))
}

func appendMarshalBufferLiteralString(s string) *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.Lit(safejson.QuoteString(s)).Op("..."))
}

func appendMarshalBufferQuotedString(selector *jen.Statement) *jen.Statement {
	return jen.Id(outName).Op("=").Add(snip.SafeJSONAppendQuotedString()).Call(jen.Id(outName), selector)
}

func appendMarshalBufferVariadic(selector *jen.Statement) *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), selector.Op("..."))
}
