package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/pkg/safejson"
)

const (
	outName = "out"
)

func MarshalJSONValue(methodBody *jen.Group, valueType types.Type, pointerSelector jen.Code) {

}

type JSONStructField struct {
	Spec     *types.Field
	Selector jen.Code
}

func MarshalJSONStruct(methodBody *jen.Group, fields []JSONStructField) {
	methodBody.Add(appendMarshalBufferLiteralRune('{'))
	for i, field := range fields {
		methodBody.Add(appendMarshalBufferLiteralString(field.Spec.Name))
		methodBody.Add(appendMarshalBufferLiteralRune(':'))
		//body = append(body, appendMarshalBufferQuotedString(expression.StringVal(field.JSONKey)))
		//body = append(body, appendMarshalBuffer(expression.VariableVal(`':'`)))
		visitor := &jsonMarshalValueVisitor{
			info:     info,
			selector: expression.NewSelector(expression.VariableVal(receiverName), field.FieldSelector),
		}
		if err := field.Type.Accept(visitor); err != nil {
			return nil, err
		}
		body = append(body, visitor.stmts...)

		if i < len(fields)-1 {
			methodBody.Add(appendMarshalBufferLiteralRune(','))
		}
	}
	methodBody.Add(appendMarshalBufferLiteralRune('}'))
	methodBody.Return(jen.Id(outName))
}

func appendMarshalBufferValue(g *jen.Group, selector func() *jen.Statement, valueType types.Type, isMapKey bool) {
	switch typ := valueType.(type) {
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
	case types.Bearertoken:
		g.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
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
				g.Add(appendMarshalBufferQuotedString(jen.Lit("true"))),
			).Else().Block(
				g.Add(appendMarshalBufferQuotedString(jen.Lit("false"))),
			)
		} else {
			g.If(selector()).Block(
				g.Add(appendMarshalBufferVariadic(jen.Lit("true"))),
			).Else().Block(
				g.Add(appendMarshalBufferVariadic(jen.Lit("false"))),
			)
		}
	case types.DateTime:
		g.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
	case types.Double:
	case types.Integer:
	case types.RID:
		g.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
	case types.Safelong:
	case types.String:
	case types.UUID:
		g.Add(appendMarshalBufferQuotedString(selector().Dot("String").Call()))
	case *types.Optional:
	case *types.List:
	case *types.Map:
	case *types.AliasType:
	case *types.EnumType:
	case *types.ObjectType:
	case *types.UnionType:
	case *types.External:
	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

func appendMarshalBufferLiteralNull() *jen.Statement {
	return appendMarshalBufferVariadic(jen.Lit("null"))
}

func appendMarshalBufferLiteralRune(r rune) *jen.Statement {
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.Lit(fmt.Sprintf("'%c'", r)))
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
