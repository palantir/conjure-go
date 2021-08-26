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
	appendMarshalBufferJSONValue(methodBody, selector, aliasType, false)
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func EnumMethodBodyAppendJSON(methodBody *jen.Group, receiverName string) {
	methodBody.Add(appendMarshalBufferQuotedString(jen.String().Call(jen.Id(receiverName).Dot("val"))))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

type JSONStructField struct {
	Spec     *types.Field
	Selector func() *jen.Statement
}

func StructMethodBodyAppendJSON(methodBody *jen.Group, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	for i, field := range fields {
		methodBody.BlockFunc(func(g *jen.Group) {
			g.Add(appendMarshalBufferVariadic(jen.Lit(safejson.QuoteString(field.Spec.Name) + ":")))
			appendMarshalBufferJSONValue(g, field.Selector, field.Spec.Type, false)

			if i < len(fields)-1 {
				g.Add(appendMarshalBufferLiteralRune(','))
			}
		})
	}
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func UnionMethodBodyAppendJSON(methodBody *jen.Group, typeFieldSelctor func() *jen.Statement, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	methodBody.Switch(typeFieldSelctor()).BlockFunc(func(g *jen.Group) {
		g.Default().Block(
			appendMarshalBufferVariadic(jen.Lit(`"type":`)),
			appendMarshalBufferQuotedString(typeFieldSelctor()),
		)
		for _, fieldDef := range fields {
			g.Case(jen.Lit(fieldDef.Spec.Name)).BlockFunc(func(g *jen.Group) {
				g.Add(appendMarshalBufferVariadic(jen.Lit(`"type":` + safejson.QuoteString(fieldDef.Spec.Name))))
				g.If(fieldDef.Selector().Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
					g.Add(appendMarshalBufferLiteralRune(','))
					g.Add(appendMarshalBufferLiteralString(fieldDef.Spec.Name))
					g.Add(appendMarshalBufferLiteralRune(':'))
					g.Id("unionVal").Op(":=").Op("*").Add(fieldDef.Selector())
					appendMarshalBufferJSONValue(g, jen.Id("unionVal").Clone, fieldDef.Spec.Type, false)
				})
			})
		}
	})
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func AnonFuncBodyAppendJSON(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	appendMarshalBufferJSONValue(funcBody, selector, valueType, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
}

func appendMarshalBufferJSONValue(g *jen.Group, selector func() *jen.Statement, valueType types.Type, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.String:
		g.Add(appendMarshalBufferQuotedString(selector()))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID:
		g.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
	case types.Any:
		g.If(
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
			g.Add(appendMarshalBufferQuotedString(jen.String().Call(selector())))
		} else {
			g.Add(appendMarshalBufferLiteralRune('"'))
			g.If(jen.Len(selector()).Op(">").Lit(0)).Block(
				jen.Id("b64out").Op(":=").Make(
					jen.Op("[]").Byte(),
					jen.Lit(0),
					snip.Base64EncodedLen().Call(jen.Len(selector())),
				),
				snip.Base64Encode().Call(jen.Id("b64out"), selector()),
				appendMarshalBufferVariadic(jen.Id("b64out")),
			)
			g.Add(appendMarshalBufferLiteralRune('"'))
		}
	case types.Boolean:
		if isMapKey {
			g.If(selector()).Block(
				g.Add(appendMarshalBufferVariadic(jen.Lit(`"true"`))),
			).Else().Block(
				g.Add(appendMarshalBufferVariadic(jen.Lit(`"false"`))),
			)
		} else {
			g.If(selector()).Block(
				g.Add(appendMarshalBufferVariadic(jen.Lit("true"))),
			).Else().Block(
				g.Add(appendMarshalBufferVariadic(jen.Lit("false"))),
			)
		}
	case types.Double:
		g.Switch().Block(
			jen.Default().BlockFunc(func(g *jen.Group) {
				if isMapKey {
					g.Add(appendMarshalBufferLiteralRune('"'))
				}
				g.Id(outName).Op("=").Add(snip.StrconvAppendFloat()).Call(jen.Id(outName), selector(), jen.Lit(-1), jen.Lit(10), jen.Lit(64))
				if isMapKey {
					g.Add(appendMarshalBufferLiteralRune('"'))
				}
			}),
			jen.Case(snip.MathIsNaN().Call(selector())).Block(appendMarshalBufferLiteralString("NaN")),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(1))).Block(appendMarshalBufferLiteralString("Infinity")),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(-1))).Block(appendMarshalBufferLiteralString("-Infinity")),
		)
	case types.Integer, types.Safelong:
		if isMapKey {
			g.Add(appendMarshalBufferLiteralRune('"'))
		}
		g.Id(outName).Op("=").Add(snip.StrconvAppendInt()).Call(jen.Id(outName), jen.Int64().Call(selector()), jen.Lit(10))
		if isMapKey {
			g.Add(appendMarshalBufferLiteralRune('"'))
		}
	case *types.Optional:
		g.If(selector().Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
			g.Id("optVal").Op(":=").Op("*").Add(selector())
			appendMarshalBufferJSONValue(g, jen.Id("optVal").Clone, typ.Item, isMapKey)
		}).Else().Block(
			appendMarshalBufferLiteralNull(),
		)
	case *types.List:
		g.Add(appendMarshalBufferLiteralRune('['))
		g.For(jen.Id("i").Op(":=").Range().Add(selector())).BlockFunc(func(g *jen.Group) {
			appendMarshalBufferJSONValue(g, selector().Index(jen.Id("i")).Clone, typ.Item, false)
			g.If(jen.Id("i").Op("<").Len(selector()).Op("-").Lit(1)).Block(
				appendMarshalBufferLiteralRune(','),
			)
		})
		g.Add(appendMarshalBufferLiteralRune(']'))
	case *types.Map:
		g.Add(appendMarshalBufferLiteralRune('{'))
		g.Block(
			jen.Var().Id("i").Int(),
			jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Add(selector())).BlockFunc(func(g *jen.Group) {
				g.BlockFunc(func(g *jen.Group) {
					appendMarshalBufferJSONValue(g, jen.Id("k").Clone, typ.Key, true)
				})
				g.Add(appendMarshalBufferLiteralRune(':'))
				g.BlockFunc(func(g *jen.Group) {
					appendMarshalBufferJSONValue(g, jen.Id("v").Clone, typ.Val, false)
				})
				g.Id("i").Op("++")
				g.If(jen.Id("i").Op("<").Len(selector())).Block(
					appendMarshalBufferLiteralRune(','),
				)
			}),
		)
		g.Add(appendMarshalBufferLiteralRune('}'))
	case *types.AliasType, *types.EnumType, *types.ObjectType, *types.UnionType:
		g.Var().Err().Error()
		g.List(jen.Id(outName), jen.Err()).Op("=").Add(selector()).Dot("AppendJSON").Call(jen.Id(outName))
		g.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		)
	case *types.External:
		g.If(
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
