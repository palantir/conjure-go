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
	"github.com/palantir/pkg/safejson"
)

const (
	outName = "out"
)

func AnonFuncBodyAppendJSON(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	marshalJSONValue(false, funcBody, selector, valueType, 0, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
}

func AnonFuncBodyJSONSize(funcBody *jen.Group, selector func() *jen.Statement, valueType types.Type) {
	funcBody.Var().Id(outName).Int()
	marshalJSONValue(true, funcBody, selector, valueType, 0, false)
	funcBody.Return(jen.Id(outName), jen.Nil())
}

func MarshalJSONMethods(receiverName string, receiverTypeName string, receiverType types.Type) []*jen.Statement {
	stmts := []*jen.Statement{snip.MethodMarshalJSON(receiverName, receiverTypeName).Block(
		jen.Return(jen.Id(receiverName).Dot("AppendJSON").Call(
			jen.Make(jen.Op("[]").Byte(), jen.Lit(0), jen.Id(receiverName).Dot("JSONSize").Call()),
		)),
	)}
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
				marshalJSONAlias(false, methodBody, v, selector.Clone)
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				marshalJSONAlias(true, methodBody, v, selector.Clone)
			}),
		)
	case *types.EnumType:
		return append(stmts,
			snip.MethodAppendJSON(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				marshalJSONEnum(false, methodBody, receiverName)
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				marshalJSONEnum(true, methodBody, receiverName)
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
				marshalJSONStruct(false, methodBody, fields)
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				marshalJSONStruct(true, methodBody, fields)
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
				marshalJSONUnion(false, methodBody, jen.Id(receiverName).Dot("typ").Clone, fields)
			}),
			snip.MethodJSONSize(receiverName, receiverTypeName).BlockFunc(func(methodBody *jen.Group) {
				marshalJSONUnion(true, methodBody, jen.Id(receiverName).Dot("typ").Clone, fields)
			}),
		)
	default:
		panic(receiverType)
	}
}

func marshalJSONAlias(isJSONSize bool, methodBody *jen.Group, aliasType types.Type, selector func() *jen.Statement) {
	if isJSONSize {
		methodBody.Var().Id(outName).Int()
	}
	marshalJSONValue(isJSONSize, methodBody, selector, aliasType, 0, false)
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func marshalJSONEnum(isJSONSize bool, methodBody *jen.Group, receiverName string) {
	if isJSONSize {
		methodBody.Var().Id(outName).Int()
	}
	methodBody.Add(appendMarshalBufferQuotedString(isJSONSize, jen.String().Call(jen.Id(receiverName).Dot("val"))))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

type jsonStructField struct {
	Key      string
	Type     types.Type
	Selector func() *jen.Statement
}

func marshalJSONStruct(isJSONSize bool, methodBody *jen.Group, fields []jsonStructField) {
	if isJSONSize {
		methodBody.Var().Id(outName).Int()
	}
	methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '{'))
	for i := range fields {
		marshalJSONStructField(isJSONSize, methodBody, fields, i)
	}
	methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func marshalJSONUnion(isJSONSize bool, methodBody *jen.Group, typeFieldSelctor func() *jen.Statement, fields []jsonStructField) {
	if isJSONSize {
		methodBody.Var().Id(outName).Int()
	}
	methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '{'))
	methodBody.Switch(typeFieldSelctor()).BlockFunc(func(cases *jen.Group) {
		for _, fieldDef := range fields {
			cases.Case(jen.Lit(fieldDef.Key)).BlockFunc(func(caseBody *jen.Group) {
				caseBody.Add(appendMarshalBufferVariadic(isJSONSize, jen.Lit(`"type":`+safejson.QuoteString(fieldDef.Key))))
				caseBody.If(fieldDef.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
					ifBody.Add(appendMarshalBufferLiteralRune(isJSONSize, ','))
					ifBody.Add(appendMarshalBufferLiteralString(isJSONSize, fieldDef.Key))
					ifBody.Add(appendMarshalBufferLiteralRune(isJSONSize, ':'))
					ifBody.Id("unionVal").Op(":=").Op("*").Add(fieldDef.Selector())
					marshalJSONValue(isJSONSize, ifBody, jen.Id("unionVal").Clone, fieldDef.Type, 0, false)
				})
			})
		}
		cases.Default().Block(
			appendMarshalBufferVariadic(isJSONSize, jen.Lit(`"type":`)),
			appendMarshalBufferQuotedString(isJSONSize, typeFieldSelctor()),
		)
	})
	methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '}'))
	methodBody.Return(jen.Id(outName), jen.Nil())
}

func marshalJSONStructField(isJSONSize bool, methodBody *jen.Group, fields []jsonStructField, fieldIdx int) {
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
			g.Add(appendMarshalBufferLiteralRune(isJSONSize, ','))
		} else {
			// Creates `if out[len(out)-1] != '{' {	out = append(out, ',') }`
			// No need to check for empty because the '{' byte is definitely written by now.
			g.If(
				jen.Id(outName).Index(jen.Len(jen.Id(outName)).Op("-").Lit(1)).Op("!=").LitRune('{'),
			).Block(
				appendMarshalBufferLiteralRune(isJSONSize, ','),
			)
		}
	}

	field := fields[fieldIdx]
	if field.Type.IsOptional() {
		switch typ := field.Type.(type) {
		case *types.Optional:
			methodBody.If(field.Selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(ifBody)
				ifBody.Add(appendMarshalBufferVariadic(isJSONSize, jen.Lit(safejson.QuoteString(field.Key)+":")))
				ifBody.Id("optVal").Op(":=").Op("*").Add(field.Selector())
				marshalJSONValue(isJSONSize, ifBody, jen.Id("optVal").Clone, typ.Item, 0, false)
			})
		case *types.AliasType:
			methodBody.If(field.Selector().Dot("Value").Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
				appendCommaIfNotFirstField(ifBody)
				ifBody.Add(appendMarshalBufferVariadic(isJSONSize, jen.Lit(safejson.QuoteString(field.Key)+":")))
				marshalJSONValue(isJSONSize, ifBody, field.Selector, field.Type, 0, false)
			})
		default:
			panic(fmt.Sprintf("unexpected optional type %T", field.Type))
		}
	} else {
		methodBody.BlockFunc(func(fieldBlock *jen.Group) {
			appendCommaIfNotFirstField(fieldBlock)
			fieldBlock.Add(appendMarshalBufferVariadic(isJSONSize, jen.Lit(safejson.QuoteString(field.Key)+":")))
			marshalJSONValue(isJSONSize, fieldBlock, field.Selector, field.Type, 0, false)
		})
	}
}

func marshalJSONValue(isJSONSize bool, methodBody *jen.Group, selector func() *jen.Statement, valueType types.Type, nestDepth int, isMapKey bool) {
	switch typ := valueType.(type) {
	case types.String:
		methodBody.Add(appendMarshalBufferQuotedString(isJSONSize, selector()))
	case types.Bearertoken, types.DateTime, types.RID, types.UUID:
		methodBody.Add(appendMarshalBufferQuotedString(isJSONSize, selector().Dot("String").Call()))
	case types.Any:
		methodBody.If(
			selector().Op("==").Nil(),
		).Block(
			appendMarshalBufferLiteralNull(isJSONSize),
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
			appendMarshalBufferVariadic(isJSONSize, jen.Id(dataName)),
		).Else().If(
			jen.List(jen.Id(dataName), jen.Err()).Op(":=").Add(snip.SafeJSONMarshal()).Call(selector()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		).Else().Block(
			appendMarshalBufferVariadic(isJSONSize, jen.Id(dataName)),
		)
	case types.Binary:
		if isMapKey {
			methodBody.Add(appendMarshalBufferQuotedString(isJSONSize, jen.String().Call(selector())))
		} else {
			methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '"'))
			methodBody.If(jen.Len(selector()).Op(">").Lit(0)).Block(
				jen.Id("b64out").Op(":=").Make(
					jen.Op("[]").Byte(),
					snip.Base64EncodedLen().Call(jen.Len(selector())),
				),
				snip.Base64Encode().Call(jen.Id("b64out"), selector()),
				appendMarshalBufferVariadic(isJSONSize, jen.Id("b64out")),
			)
			methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '"'))
		}
	case types.Boolean:
		if isMapKey {
			methodBody.If(selector()).Block(
				appendMarshalBufferVariadic(isJSONSize, jen.Lit(`"true"`)),
			).Else().Block(
				appendMarshalBufferVariadic(isJSONSize, jen.Lit(`"false"`)),
			)
		} else {
			methodBody.If(selector()).Block(
				appendMarshalBufferVariadic(isJSONSize, jen.Lit("true")),
			).Else().Block(
				appendMarshalBufferVariadic(isJSONSize, jen.Lit("false")),
			)
		}
	case types.Double:
		methodBody.Switch().Block(
			jen.Default().BlockFunc(func(caseBody *jen.Group) {
				if isMapKey {
					caseBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '"'))
				}
				caseBody.Id(outName).Op("=").Add(snip.StrconvAppendFloat()).Call(jen.Id(outName), selector(), jen.LitRune('g'), jen.Lit(-1), jen.Lit(64))
				if isMapKey {
					caseBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '"'))
				}
			}),
			jen.Case(snip.MathIsNaN().Call(selector())).Block(
				appendMarshalBufferLiteralString(isJSONSize, "NaN"),
			),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(1))).Block(
				appendMarshalBufferLiteralString(isJSONSize, "Infinity"),
			),
			jen.Case(snip.MathIsInf().Call(selector(), jen.Lit(-1))).Block(
				appendMarshalBufferLiteralString(isJSONSize, "-Infinity"),
			),
		)
	case types.Integer, types.Safelong:
		if isMapKey {
			methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '"'))
		}
		methodBody.Id(outName).Op("=").Add(snip.StrconvAppendInt()).Call(jen.Id(outName), jen.Int64().Call(selector()), jen.Lit(10))
		if isMapKey {
			methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '"'))
		}
	case *types.Optional:
		methodBody.If(selector().Op("!=").Nil()).BlockFunc(func(ifBody *jen.Group) {
			ifBody.Id("optVal").Op(":=").Op("*").Add(selector())
			marshalJSONValue(isJSONSize, ifBody, jen.Id("optVal").Clone, typ.Item, nestDepth+1, isMapKey)
		}).Else().Block(
			appendMarshalBufferLiteralNull(isJSONSize),
		)
	case *types.List:
		i := tmpVarName("i", nestDepth)
		methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '['))
		methodBody.For(jen.Id(i).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
			marshalJSONValue(isJSONSize, rangeBody, selector().Index(jen.Id(i)).Clone, typ.Item, nestDepth+1, false)
			rangeBody.If(jen.Id(i).Op("<").Len(selector()).Op("-").Lit(1)).Block(
				appendMarshalBufferLiteralRune(isJSONSize, ','),
			)
		})
		methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, ']'))
	case *types.Map:
		methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '{'))
		mapIdx := tmpVarName("mapIdx", nestDepth)
		methodBody.Block(
			jen.Var().Id(mapIdx).Int(),
			jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Add(selector())).BlockFunc(func(rangeBody *jen.Group) {
				rangeBody.BlockFunc(func(keyBlock *jen.Group) {
					marshalJSONValue(isJSONSize, keyBlock, jen.Id("k").Clone, typ.Key, nestDepth+1, true)
				})
				rangeBody.Add(appendMarshalBufferLiteralRune(isJSONSize, ':'))
				rangeBody.BlockFunc(func(valueBlock *jen.Group) {
					marshalJSONValue(isJSONSize, valueBlock, jen.Id("v").Clone, typ.Val, nestDepth+1, false)
				})
				rangeBody.Id(mapIdx).Op("++")
				rangeBody.If(jen.Id(mapIdx).Op("<").Len(selector())).Block(
					appendMarshalBufferLiteralRune(isJSONSize, ','),
				)
			}),
		)
		methodBody.Add(appendMarshalBufferLiteralRune(isJSONSize, '}'))
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
				appendMarshalBufferVariadic(isJSONSize, jen.Id(dataName)),
			)
		} else {
			marshalJSONValue(isJSONSize, methodBody, selector, typ.Fallback, nestDepth, isMapKey)
		}

	default:
		panic(fmt.Sprintf("unknown type %T", typ))
	}
}

func appendMarshalBufferLiteralNull(isJSONSize bool) *jen.Statement {
	if isJSONSize {
		return jen.Id(outName).Op("+=").Len(jen.Lit("null"))
	}
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.Lit("null").Op("..."))
}

func appendMarshalBufferLiteralRune(isJSONSize bool, r rune) *jen.Statement {
	if isJSONSize {
		return jen.Id(outName).Op("++")
	}
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.LitRune(r))
}

func appendMarshalBufferLiteralString(isJSONSize bool, s string) *jen.Statement {
	if isJSONSize {
		return jen.Id(outName).Op("+=").Len(jen.Lit(safejson.QuoteString(s)))
	}
	return jen.Id(outName).Op("=").Append(jen.Id(outName), jen.Lit(safejson.QuoteString(s)).Op("..."))
}

func appendMarshalBufferQuotedString(isJSONSize bool, selector *jen.Statement) *jen.Statement {
	if isJSONSize {
		return jen.Id(outName).Op("+=").Len(snip.SafeJSONQuotedStringLength().Call(selector))
	}
	return jen.Id(outName).Op("=").Add(snip.SafeJSONAppendQuotedString()).Call(jen.Id(outName), selector)
}

func appendMarshalBufferVariadic(isJSONSize bool, selector *jen.Statement) *jen.Statement {
	if isJSONSize {
		return jen.Id(outName).Op("+=").Len(selector)
	}
	return jen.Id(outName).Op("=").Append(jen.Id(outName), selector.Op("..."))
}
