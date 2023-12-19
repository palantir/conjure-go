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
	"github.com/tidwall/gjson"
)

func AnonFuncBodyAppendJSON(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	marshalJSONValue(marshalContext{isAppendJSON: true}, funcBody, selector, valueType, 0, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
}

func AnonFuncBodyJSONSize(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	funcBody.Var().Id(outName).Int()
	marshalJSONValue(marshalContext{isJSONSize: true}, funcBody, selector, valueType, 0, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
}

func MarshalJSONMethods(receiverName string, receiverTypeName string, receiverType types.Type) []*jen.Statement {
	stmts := []*jen.Statement{
		snip.MethodMarshalJSON(receiverName, receiverTypeName).Block(
			jen.List(jen.Id("size"), jen.Err()).Op(":=").Id(receiverName).Dot("JSONSize").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Return(jen.Id(receiverName).Dot("AppendJSON").Call(
				jen.Make(jen.Index().Byte(), jen.Lit(0), jen.Id("size")),
			)),
		),
	}
	switch v := receiverType.(type) {
	case *types.AliasType:
		var selector *jen.Statement
		if v.IsOptional() {
			selector = jen.Id(receiverName).Dot("Value")
		} else {
			selector = v.Item.Code().Call(jen.Id(receiverName))
		}
		return append(stmts,
			snip.MethodAppendJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isAppendJSON: true}
				marshalJSONValue(ctx, methodBody, selector.Clone, v.Item, 0, false)
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isWriteJSON: true}
				marshalJSONValue(ctx, methodBody, selector.Clone, v.Item, 0, false)
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isJSONSize: true}
				marshalJSONValue(ctx, methodBody, selector.Clone, v.Item, 0, false)
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
		)
	case *types.EnumType:
		return append(stmts,
			snip.MethodAppendJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isAppendJSON: true}
				methodBody.Add(ctx.quotedString(jen.String().Call(jen.Id(receiverName).Dot("val"))))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isWriteJSON: true}
				methodBody.Add(ctx.quotedString(jen.String().Call(jen.Id(receiverName).Dot("val"))))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isJSONSize: true}
				methodBody.Add(ctx.quotedString(jen.String().Call(jen.Id(receiverName).Dot("val"))))
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
		return append(stmts,
			snip.MethodAppendJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isAppendJSON: true}
				methodBody.Add(ctx.literalRune('{'))
				for i := range fields {
					marshalJSONStructField(ctx, methodBody, fields, i)
				}
				methodBody.Add(ctx.literalRune('}'))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isWriteJSON: true}
				methodBody.Add(ctx.literalRune('{'))
				for i := range fields {
					marshalJSONStructField(ctx, methodBody, fields, i)
				}
				methodBody.Add(ctx.literalRune('}'))
				methodBody.Return(jen.Id(outName), jen.Nil())
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isJSONSize: true}
				methodBody.Add(ctx.literalRune('{'))
				for i := range fields {
					marshalJSONStructField(ctx, methodBody, fields, i)
				}
				methodBody.Add(ctx.literalRune('}'))
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
		return append(stmts,
			snip.MethodAppendJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isAppendJSON: true}
				marshalJSONUnion(ctx, methodBody, jen.Id(receiverName).Dot("typ").Clone, fields)
			}),
			snip.MethodWriteJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isWriteJSON: true}
				marshalJSONUnion(ctx, methodBody, jen.Id(receiverName).Dot("typ").Clone, fields)
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				ctx := marshalContext{isJSONSize: true}
				marshalJSONUnion(ctx, methodBody, jen.Id(receiverName).Dot("typ").Clone, fields)
			}),
		)
	default:
		panic(receiverType)
	}
}

type jsonStructField struct {
	Key      string
	Type     types.Type
	Selector func() *jen.Statement
}

func marshalJSONUnion(ctx marshalContext, methodBody *jen.Group, typeFieldSelctor func() *jen.Statement, fields []jsonStructField) {
	methodBody.Add(ctx.literalRune('{'))
	methodBody.Switch(typeFieldSelctor()).BlockFunc(func(cases *jen.Group) {
		for _, field := range fields {
			cases.Case(jen.Lit(field.Key)).BlockFunc(func(caseBody *jen.Group) {
				caseBody.Add(ctx.literalString(`"type":` + quoteJSONString(field.Key)))
				caseBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
					ifBody.Add(ctx.literalString("," + quoteJSONString(field.Key) + ":"))
					ifBody.Id("unionVal").Op(":=").Op("*").Add(field.Selector())
					marshalJSONValue(ctx, ifBody, jen.Id("unionVal").Clone, field.Type, 0, false)
				})
			})
		}
		cases.Default().Block(
			ctx.literalString(`"type":`),
			ctx.quotedString(typeFieldSelctor()),
		)
	})
	methodBody.Add(ctx.literalRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func marshalJSONStructField(ctx marshalContext, methodBody *jen.Group, fields []jsonStructField, fieldIdx int) {
	field := fields[fieldIdx]
	if field.Type.IsOptional() {
		switch typ := field.Type.(type) {
		case *types.Optional:
			methodBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(ctx, ifBody, fields, fieldIdx)
				ifBody.Add(ctx.literalString(quoteJSONString(field.Key) + ":"))
				ifBody.Id("optVal").Op(":=").Op("*").Add(field.Selector())
				marshalJSONValue(ctx, ifBody, jen.Id("optVal").Clone, typ.Item, 0, false)
			})
		case *types.AliasType:
			methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(ctx, ifBody, fields, fieldIdx)
				ifBody.Add(ctx.literalString(quoteJSONString(field.Key) + ":"))
				marshalJSONValue(ctx, ifBody, field.Selector, typ, 0, false)
			})
		default:
			panic(fmt.Sprintf("unexpected optional type %T", field.Type))
		}
	} else {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			appendCommaIfNotFirstField(ctx, fieldBlock, fields, fieldIdx)
			fieldBlock.Add(ctx.literalString(quoteJSONString(field.Key) + ":"))
			marshalJSONValue(ctx, fieldBlock, field.Selector, field.Type, 0, false)
		})
	}
}

func marshalJSONValue(ctx marshalContext, methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.String:
		methodBody.Add(ctx.quotedString(selector()))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID:
		methodBody.Add(ctx.quotedString(selector().Dot("String").Call()))
	case types.Any:
		methodBody.If(
			selector().Op("==").Nil(),
		).Block(
			ctx.literalString("null"),
		).Else().If(ctx.checkInterface(selector()), jen.Id("ok")).Block(
			ctx.callInterface(),
		).Else().If(
			jen.List(jen.Id("marshaler"), jen.Id("ok")).Op(":=").Add(selector()).Assert(snip.JSONMarshaler()),
			jen.Id("ok"),
		).Block(
			jen.List(jen.Id(dataName), jen.Err()).Op(":=").Id("marshaler").Dot("MarshalJSON").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				ctx.returnErr(jen.Err()),
			),
			ctx.variadicSlice(jen.Id(dataName)),
		).Else().If(
			jen.List(jen.Id(dataName), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
			jen.Err().Op("!=").Nil(),
		).Block(
			ctx.returnErr(jen.Err()),
		).Else().Block(
			ctx.variadicSlice(jen.Id(dataName)),
		)
	case types.Binary:
		if isMapKey {
			methodBody.Add(ctx.quotedString(jen.String().Call(selector())))
		} else {
			methodBody.Add(ctx.literalRune('"'))
			methodBody.If(jen.Len(selector()).Op(">").Lit(0)).Block(
				jen.Id("b64out").Op(":=").Make(
					jen.Index().Byte(),
					snip.Base64StdEncoding().Dot("EncodedLen").Call(jen.Len(selector())),
				),
				snip.Base64StdEncoding().Dot("Encode").Call(jen.Id("b64out"), selector()),
				ctx.variadicSlice(jen.Id("b64out")),
			)
			methodBody.Add(ctx.literalRune('"'))
		}
	case types.Boolean:
		if isMapKey {
			methodBody.If(selector()).Block(
				ctx.literalString(`"true"`),
			).Else().Block(
				ctx.literalString(`"false"`),
			)
		} else {
			methodBody.If(selector()).Block(
				ctx.literalString("true"),
			).Else().Block(
				ctx.literalString("false"),
			)
		}
	case types.Double:
		methodBody.Switch().Block(
			jen.Default().BlockFunc(func(caseBody *jen.Group) {
				if isMapKey {
					caseBody.Add(ctx.literalRune('"'))
				}
				caseBody.Add(ctx.float(selector()))
				if isMapKey {
					caseBody.Add(ctx.literalRune('"'))
				}
			}),
			jen.Case(snip.MathIsNaN().Call(selector())).Block(
				ctx.literalString(`"NaN"`),
			),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(1))).Block(
				ctx.literalString(`"Infinity"`),
			),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(-1))).Block(
				ctx.literalString(`"-Infinity"`),
			),
		)
	case types.Integer, types.Safelong:
		if isMapKey {
			methodBody.Add(ctx.literalRune('"'))
		}
		methodBody.Add(ctx.integer(selector()))
		if isMapKey {
			methodBody.Add(ctx.literalRune('"'))
		}
	case *types.Optional:
		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
			marshalJSONValue(ctx, ifBody, jen.Id("optVal").Clone, typ.Item, nestDepth+1, isMapKey)
		}).Else().Block(
			ctx.literalString("null"),
		)
	case *types.List:
		i := tmpVarName("i", nestDepth)
		methodBody.Add(ctx.literalRune('['))
		methodBody.For(jen.Id(i).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
			marshalJSONValue(ctx, rangeBody, selector().Index(jen.Id(i)).Clone, typ.Item, nestDepth+1, false)
			rangeBody.If(jen.Id(i).Op("<").Len(selector()).Op("-").Lit(1)).Block(
				ctx.literalRune(','),
			)
		})
		methodBody.Add(ctx.literalRune(']'))
	case *types.Map:
		methodBody.Add(ctx.literalRune('{'))
		mapIdx := tmpVarName("mapIdx", nestDepth)
		methodBody.Block(
			jen.Var().Id(mapIdx).Int(),
			jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
				rangeBody.BlockFunc(func(keyBlock *jen.Group) {
					marshalJSONValue(ctx, keyBlock, jen.Id("k").Clone, typ.Key, nestDepth+1, true)
				})
				rangeBody.Add(ctx.literalRune(':'))
				rangeBody.BlockFunc(func(valueBlock *jen.Group) {
					marshalJSONValue(ctx, valueBlock, jen.Id("v").Clone, typ.Val, nestDepth+1, false)
				})
				rangeBody.Id(mapIdx).Op("++")
				rangeBody.If(jen.Id(mapIdx).Op("<").Len(selector())).Block(
					ctx.literalRune(','),
				)
			}),
		)
		methodBody.Add(ctx.literalRune('}'))
	case *types.EnumType:
		methodBody.Add(ctx.quotedString(jen.String().Call(selector().Dot("String").Call())))
	case *types.AliasType, *types.ObjectType, *types.UnionType:
		ctx.callDelegate(selector())
	case *types.External:
		if typ.ExternalHasGoType() {
			methodBody.If(
				jen.List(jen.Id(dataName), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Nil(), jen.Err()),
			).Else().Block(
				ctx.variadicSlice(jen.Id(dataName)),
			)
		} else {
			marshalJSONValue(ctx, methodBody, selector, typ.Fallback, nestDepth, isMapKey)
		}

	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

// appendCommaIfNotFirstField adds the comma separator between struct fields.
// It should come before writing the field's key string.
// If fieldIdx == 0, this is a noop.
func appendCommaIfNotFirstField(ctx marshalContext, g *jen.Group, fields []jsonStructField, fieldIdx int) {
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
		g.Add(ctx.literalRune(','))
	} else {
		g.If(anyPrecedingNotNil).Block(ctx.literalRune(','))
	}
}

type marshalContext struct {
	isAppendJSON bool
	isWriteJSON  bool
	isJSONSize   bool
}

func (ctx marshalContext) returnErr(err *jen.Statement) *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.Return(jen.Nil(), err)
	case ctx.isJSONSize:
		return jen.Return(jen.Lit(0), err)
	case ctx.isWriteJSON:
		return jen.Return(jen.Lit(0), err)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) literalRune(r byte) *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.LitRune(rune(r)))
	case ctx.isJSONSize:
		return jen.Id(outName).Op("++").Commentf("'%c'", r)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("bw"), jen.Id("ok")).Op(":=").Id(wName).Assert(jen.Qual("io", "ByteWriter")),
			jen.Id("ok"),
		).Block(
			jen.If(
				jen.Err().Op(":=").Id("bw").Dot("WriteByte").Call(jen.LitRune(rune(r))),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Lit(0), jen.Err()),
			).Else().Block(
				jen.Id("out").Op("++"),
			),
		).Else().Block(
			jen.If(
				jen.List(jen.Op("_"), jen.Err()).Op(":=").Id(wName).Dot("Write").Call(jen.Index().Byte().Values(jen.LitRune(rune(r)))),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Lit(0), jen.Err()),
			).Else().Block(
				jen.Id("out").Op("++"),
			),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) literalString(s string) *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.Lit(s).Op("..."))
	case ctx.isJSONSize:
		return jen.Id(outName).Op("+=").Lit(len(s)).Comment(s)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("n"), jen.Err()).Op(":=").Qual("io", "WriteString").Call(jen.Id(wName), jen.Lit(s)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id("out").Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) quotedString(selector *jen.Statement) *jen.Statement {
	str := snip.GJSONAppendJSONString().Call(jen.Nil(), selector)
	switch {
	case ctx.isAppendJSON:
		return jen.Id(outName).Op("=").Add(str)
	case ctx.isJSONSize:
		return jen.Id(outName).Op("+=").Len(str)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("n"), jen.Err()).Op(":=").Id(wName).Dot("Write").Call(str),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id("out").Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) variadicSlice(selector *jen.Statement) *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.Id(outName).Op("=").Append(jen.Id(outName), selector.Op("..."))
	case ctx.isJSONSize:
		return jen.Id(outName).Op("+=").Len(selector)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("n"), jen.Err()).Op(":=").Id(wName).Dot("Write").Call(selector),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id("out").Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) float(selector *jen.Statement) *jen.Statement {
	floatBytes := snip.StrconvAppendFloat().Call(jen.Id(outName), selector, jen.LitRune('g'), jen.Lit(-1), jen.Lit(64))
	switch {
	case ctx.isAppendJSON:
		return jen.Id(outName).Op("=").Add(floatBytes)
	case ctx.isJSONSize:
		return jen.Id(outName).Op("+=").Len(floatBytes)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("n"), jen.Err()).Op(":=").Id(wName).Dot("Write").Call(floatBytes),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Nil(), jen.Err())).Else().Block(
			jen.Id("out").Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) integer(selector *jen.Statement) *jen.Statement {
	intBytes := snip.StrconvAppendInt().Call(jen.Nil(), jen.Int64().Call(selector), jen.Lit(10))
	switch {
	case ctx.isAppendJSON:
		return jen.Id(outName).Op("=").Add(intBytes)
	case ctx.isJSONSize:
		return jen.Id(outName).Op("+=").Len(intBytes)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("n"), jen.Err()).Op(":=").Id(wName).Dot("Write").Call(intBytes),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		).Else().Block(
			jen.Id("out").Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) checkInterface(selector *jen.Statement) *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.List(jen.Id("appender"), jen.Id("ok")).Op(":=").Add(selector).Assert(
			jen.Interface(
				jen.Id("AppendJSON").
					Params(jen.Index().Byte()).
					Params(jen.Index().Byte(), jen.Error()),
			),
		)
	case ctx.isJSONSize:
		return jen.List(jen.Id("sizer"), jen.Id("ok")).Op(":=").Add(selector).Assert(
			jen.Interface(
				jen.Id("JSONSize").
					Params().
					Params(jen.Int(), jen.Error()),
			),
		)
	case ctx.isWriteJSON:
		return jen.List(jen.Id("writer"), jen.Id("ok")).Op(":=").Add(selector).Assert(
			jen.Interface(
				jen.Id("WriteJSON").
					Params(snip.IOWriter()).
					Params(jen.Int(), jen.Error()),
			),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) callInterface() *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.Var().Err().Error().Line().
			If(
				jen.List(jen.Id(outName), jen.Err()).Op("=").Id("appender").Dot("AppendJSON").Call(jen.Id(outName)),
				jen.Err().Op("!=").Nil(),
			).Block(
			jen.Return(jen.Nil(), jen.Err()),
		)
	case ctx.isJSONSize:
		return jen.If(
			jen.Id("size"), jen.Err().Op(":=").Id("sizer").Dot("JSONSize").Call(),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id(outName).Op("+=").Id("size"),
		)
	case ctx.isWriteJSON:
		return jen.If(
			jen.Id("n"), jen.Err().Op(":=").Id("writer").Dot("WriteJSON").Call(jen.Id(wName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id(outName).Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func (ctx marshalContext) callDelegate(selector *jen.Statement) *jen.Statement {
	switch {
	case ctx.isAppendJSON:
		return jen.Var().Err().Error().
			Line().If(
			jen.List(jen.Id(outName), jen.Err()).Op("=").Add(selector).Dot("AppendJSON").Call(jen.Id(outName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		)
	case ctx.isJSONSize:
		return jen.If(
			jen.List(jen.Id("size"), jen.Err()).Op(":=").Add(selector).Dot("JSONSize").Call(),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id(outName).Op("+=").Id("size"),
		)
	case ctx.isWriteJSON:
		return jen.If(
			jen.List(jen.Id("n"), jen.Err()).Op(":=").Id(wName).Dot("WriteJSON").Call(selector),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Lit(0), jen.Err()),
		).Else().Block(
			jen.Id(outName).Op("+=").Id("n"),
		)
	default:
		panic("bad context")
	}
}

func quoteJSONString(s string) string {
	return string(gjson.AppendJSONString(nil, s))
}
