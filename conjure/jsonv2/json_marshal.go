package jsonv2

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	_ "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	encName = "enc"
	decName = "dec"
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
			switch v := receiverType.(type) {
			case *types.AliasType:
				marshalJSONValue(methodBody, jen.Id(receiverName), v, 0, false)
				methodBody.Return(jen.Nil())
			case *types.EnumType:
				marshalJSONValue(methodBody, jen.Id(receiverName), v, 0, false)
				methodBody.Return(jen.Nil())
			case *types.ObjectType:
				var fields []jsonStructField
				for _, field := range v.Fields {
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
				for _, field := range v.Fields {
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
								marshalJSONValue(ifBody, jen.Id("unionVal"), field.Type, 0, false)
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
			default:
				panic(receiverType)
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
				marshalJSONValue(ifBody, jen.Id("optVal"), typ.Item, 0, false)
			})
		case *types.AliasType:
			methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				ifBody.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
				marshalJSONValue(ifBody, field.Selector(), typ, 0, false)
			})
		default:
			panic(fmt.Sprintf("unexpected optional type %T", field.Type))
		}
	} else {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			fieldBlock.Add(mustWriteToken(snip.JSONV2String().Call(jen.Lit(field.Key))))
			marshalJSONValue(fieldBlock, field.Selector(), field.Type, 0, false)
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
		jen.Err().Op(":=").Id(encName).Dot("WriteValue").Call(snip.JSONV2String().Call(value)),
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

func marshalJSONValue(methodBody *jen.Group, selectorV *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	selector := selectorV.Clone
	switch typ := valueType.(type) {
	case types.String:
		methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector())))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID, *types.EnumType:
		methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector().Dot("String").Call())))
	case types.Binary:
		if isMapKey {
			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(selector())))
		} else {
			methodBody.If(jen.Len(selector()).Op(">").Lit(0)).Block(
				jen.Id("b64out").Op(":=").Make(
					jen.Index().Byte(),
					snip.Base64StdEncoding().Dot("EncodedLen").Call(jen.Len(selector())).Op("+2"),
				),
				jen.Id("b64out").Index(jen.Lit(0)).Op("=").Lit('"'),
				snip.Base64StdEncoding().Dot("Encode").Call(jen.Id("b64out").Index(jen.Op("1:")), selector()),
				jen.Id("b64out").Index(jen.Len(jen.Id("b64out")).Op("-1")).Op("=").Lit('"'),
				mustWriteValue(snip.JSONV2Value().Call(jen.Id("b64out"))),
			).Else().Block(
				mustWriteToken(snip.JSONV2String().Call(jen.Lit(""))),
			)
		}
	case types.Boolean:
		if isMapKey {
			methodBody.If(selector()).Block(
				mustWriteValue(snip.JSONV2String().Call(jen.Lit(`"true"`))),
			).Else().Block(
				mustWriteValue(snip.JSONV2String().Call(jen.Lit(`"false"`))),
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
			methodBody.Add(mustWriteToken(snip.JSONV2String().Call(snip.StrconvItoa().Call(selector()))))
		} else {
			methodBody.Add(mustWriteToken(snip.JSONV2Int().Call(jen.Int64().Call(selector()))))
		}
	case *types.Optional:
		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
			marshalJSONValue(ifBody, jen.Id("optVal"), typ.Item, nestDepth+1, isMapKey)
		}).Else().Block(
			mustWriteToken(snip.JSONV2Null()),
		)
	case *types.List:
		methodBody.Add(mustWriteToken(snip.JSONV2BeginArray()))
		i := tmpVarName("i", nestDepth)
		methodBody.For(jen.Id(i).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
			marshalJSONValue(rangeBody, selector().Index(jen.Id(i)), typ.Item, nestDepth+1, false)
		})
		methodBody.Add(mustWriteToken(snip.JSONV2EndArray()))
	case *types.Set:
		methodBody.Add(mustWriteToken(snip.JSONV2BeginArray()))
		i := tmpVarName("i", nestDepth)
		methodBody.For(jen.Id(i).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
			marshalJSONValue(rangeBody, selector().Index(jen.Id(i)), typ.Item, nestDepth+1, false)
		})
		methodBody.Add(mustWriteToken(snip.JSONV2EndArray()))
	case *types.Map:
		methodBody.Add(mustWriteToken(snip.JSONV2BeginObject()))
		methodBody.Block(
			// TODO: sort map keys
			jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
				rangeBody.BlockFunc(func(keyBlock *jen.Group) {
					marshalJSONValue(keyBlock, jen.Id("k"), typ.Key, nestDepth+1, true)
				})
				rangeBody.BlockFunc(func(valueBlock *jen.Group) {
					marshalJSONValue(valueBlock, jen.Id("v"), typ.Val, nestDepth+1, false)
				})
			}),
		)
		methodBody.Add(mustWriteToken(snip.JSONV2EndObject()))
	case *types.AliasType:
		switch {
		case typ.IsString(), typ.IsText(), typ.IsBinary(), typ.IsBoolean():
			var itemSelector *jen.Statement
			if typ.IsOptional() {
				itemSelector = selector().Dot("Value")
			} else {
				itemSelector = typ.Item.Code().Call(selector())
			}
			// Simple types can be converted directly
			marshalJSONValue(methodBody, itemSelector, typ.Item, nestDepth, isMapKey)
		default:
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
			jen.Err().Op(":=").Add(snip.JSONV2MarshalEncode()).Call(jen.Id(encName), selector(), jen.Id(encName).Dot("Options").Call()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		)
	case *types.External:
		if typ.ExternalHasGoType() {
			methodBody.If(
				jen.Err().Op(":=").Add(snip.JSONV2MarshalEncode()).Call(jen.Id(encName), selector(), jen.Id(encName).Dot("Options").Call()),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)
		} else {
			marshalJSONValue(methodBody, selector(), typ.Fallback, nestDepth, isMapKey)
		}

	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

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
