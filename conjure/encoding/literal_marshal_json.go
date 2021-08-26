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

func AliasMethodBodyAppendJSON(methodBody *jen.Group, alias *types.AliasType, selector func() *jen.Statement) {
	appendMarshalBufferValue(methodBody, selector, alias.Item, false)
	methodBody.Return(jen.Id(outName), jen.Nil())
}

type JSONStructField struct {
	Spec     *types.Field
	Selector func() *jen.Statement
}

func StructMethodBodyAppendJSON(methodBody *jen.Group, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	for i, field := range fields {
		methodBody.Add(appendMarshalBufferLiteralString(field.Spec.Name))
		methodBody.Add(appendMarshalBufferLiteralRune(':'))

		appendMarshalBufferValue(methodBody, field.Selector, field.Spec.Type, false)

		if i < len(fields)-1 {
			methodBody.Add(appendMarshalBufferLiteralRune(','))
		}
	}
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func appendMarshalBufferValue(g *jen.Group, selector func() *jen.Statement, valueType types.Type, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.String:
		g.Add(appendMarshalBufferQuotedString(selector()))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID, *types.EnumType:
		g.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
	case types.Any:
		g.If(
			selector().Op("==").Nil(),
		).Block(
			appendMarshalBufferLiteralNull(),
		).Else().If(
			jen.List(jen.Id("appender"), jen.Id("ok")).Op(":=").Add(selector()).Assert(jen.Interface(
				jen.Id("AppendJSON").Params(jen.Op("[]").Byte()).Params(jen.Op("[]").Byte(), jen.Error()),
			)),
			jen.Id("ok"),
		).Block(
			jen.Var().Err().Error(),
			jen.List(jen.Id(outName), jen.Err()).Op("=").Id("appender").Dot("AppendJSON").Call(jen.Id(outName)),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
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
				g.Id(outName).Op("=").Add(snip.StrconvAppendFloat()).Call(jen.Id(outName), jen.Int64().Call(selector()), jen.Lit(-1), jen.Lit(10), jen.Lit(64))
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
			appendMarshalBufferValue(g, jen.Parens(jen.Op("*").Add(selector())).Clone, typ.Item, isMapKey)
		}).Else().Block(
			appendMarshalBufferLiteralNull(),
		)
	case *types.List:
		g.Add(appendMarshalBufferLiteralRune('['))
		g.BlockFunc(func(g *jen.Group) {
			g.For(jen.Id("i").Op(":=").Range().Add(selector())).BlockFunc(func(g *jen.Group) {
				appendMarshalBufferValue(g, selector().Index(jen.Id("i")).Clone, typ.Item, false)
				g.If(jen.Id("i").Op("<").Len(selector()).Op("-").Lit(1)).Block(
					appendMarshalBufferLiteralRune(','),
				)
			})
		})
		g.Add(appendMarshalBufferLiteralRune(']'))
	case *types.Map:
		g.Add(appendMarshalBufferLiteralRune('{'))
		g.BlockFunc(func(g *jen.Group) {
			g.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Add(selector())).BlockFunc(func(g *jen.Group) {
				appendMarshalBufferValue(g, selector().Index(jen.Id("k")).Clone, typ.Key, true)
				g.Add(appendMarshalBufferLiteralRune(':'))
				appendMarshalBufferValue(g, selector().Index(jen.Id("k")).Clone, typ.Val, false)
				g.If(jen.Id("i").Op("<").Len(selector()).Op("-").Lit(1)).Block(
					appendMarshalBufferLiteralRune(','),
				)
			})
		})
		g.Add(appendMarshalBufferLiteralRune('}'))
	case *types.AliasType, *types.ObjectType, *types.UnionType:
		g.If(
			jen.List(jen.Id("tmpOut"), jen.Err()).Op(":=").Add(selector()).Dot("AppendJSON").Call(jen.Id(outName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		).Else().Block(
			jen.Id(outName).Op("=").Id("tmpOut"),
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
