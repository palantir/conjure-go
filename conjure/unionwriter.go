// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package conjure

import (
	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	unionReceiverName = "u"
	withContextSuffix = "WithContext"
)

func writeUnionType(file *jen.Group, def *types.UnionType, genAcceptFuncs bool) {
	// Declare exported union struct type
	file.Add(def.CommentLine()).
		Type().
		Id(def.Name).StructFunc(func(g *jen.Group) {
		g.Id("typ").String()
		for _, fieldDef := range def.Fields {
			g.Id(transforms.PrivateFieldName(fieldDef.Name)).Op("*").Add(fieldDef.Type.Code())
		}
	})

	// Declare deserializer struct type
	file.Type().
		Id(unionDeserializerStructName(def.Name)).StructFunc(func(g *jen.Group) {
		g.Id("Type").String().Tag(map[string]string{"json": "type"})
		for _, fieldDef := range def.Fields {
			g.Id(transforms.ExportedFieldName(fieldDef.Name)).
				Op("*").Add(fieldDef.Type.Code()).
				Tag(map[string]string{"json": fieldDef.Name})
		}
	})

	// Declare deserializer toStruct method
	file.Func().
		Params(jen.Id(unionReceiverName).Op("*").Id(unionDeserializerStructName(def.Name))).
		Id("toStruct").
		Params().
		Params(jen.Id(def.Name)).
		Block(jen.Return(jen.Id(def.Name).ValuesFunc(func(g *jen.Group) {
			g.Id("typ").Op(":").Id(unionReceiverName).Dot("Type")
			for _, fieldDef := range def.Fields {
				g.Id(transforms.PrivateFieldName(fieldDef.Name)).
					Op(":").
					Id(unionReceiverName).Dot(transforms.ExportedFieldName(fieldDef.Name))
			}
		})))

	// Declare toSerializer method
	file.Func().
		Params(jen.Id(unionReceiverName).Op("*").Id(def.Name)).
		Id("toSerializer").
		Params().
		Params(jen.Interface(), jen.Error()).
		Block(jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(g *jen.Group) {
			g.Default().Block(jen.Return(
				jen.Nil(), snip.FmtErrorf().Call(jen.Lit("unknown type %s"), jen.Id(unionReceiverName).Dot("typ"))))
			for _, fieldDef := range def.Fields {
				g.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(g *jen.Group) {
					fieldSelector := unionDerefPossibleOptional(g, fieldDef)
					g.Return(
						jen.Struct(
							jen.Id("Type").String().Tag(map[string]string{"json": "type"}),
							jen.Id(transforms.ExportedFieldName(fieldDef.Name)).Add(fieldDef.Type.Code()).Tag(map[string]string{"json": fieldDef.Name}),
						).Values(
							jen.Id("Type").Op(":").Lit(fieldDef.Name),
							jen.Id(transforms.ExportedFieldName(fieldDef.Name)).Op(":").Add(fieldSelector),
						),
						jen.Nil(),
					)
				})
			}
		}))

	// Declare MarshalJSON method
	file.Add(snip.MethodMarshalJSON(unionReceiverName, def.Name).Block(
		jen.List(jen.Id("ser"), jen.Err()).Op(":=").Id(unionReceiverName).Dot("toSerializer").Call(),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Return(snip.SafeJSONMarshal().Call(jen.Id("ser"))),
	))

	// Declare UnmarshalJSON method
	file.Add(snip.MethodUnmarshalJSON(unionReceiverName, def.Name).Block(
		jen.Var().Id("deser").Id(unionDeserializerStructName(def.Name)),
		jen.If(
			jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal().Call(jen.Id(dataVarName), jen.Op("&").Id("deser"))),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		jen.Op("*").Id(unionReceiverName).Op("=").Id("deser").Dot("toStruct").Call(),
		jen.Return(jen.Nil()),
	))

	// Declare yaml methods
	file.Add(snip.MethodMarshalYAML(unionReceiverName, def.Name))
	file.Add(snip.MethodUnmarshalYAML(unionReceiverName, def.Name))

	// Declare AcceptFuncs method & noop helpers
	if genAcceptFuncs {
		file.Func().
			Params(jen.Id(unionReceiverName).Op("*").Id(def.Name)).
			Id("AcceptFuncs").
			ParamsFunc(func(g *jen.Group) {
				for _, fieldDef := range def.Fields {
					g.Id(transforms.PrivateFieldName(fieldDef.Name) + "Func").Func().Params(fieldDef.Type.Code()).Params(jen.Error())
				}
				g.Id("unknownFunc").Func().Params(jen.String()).Params(jen.Error())
			}).
			Params(jen.Error()).
			Block(
				jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(g *jen.Group) {
					g.Default().Block(
						jen.If(jen.Id(unionReceiverName).Dot("typ").Op("==").Lit("")).Block(
							jen.Return(snip.FmtErrorf().Call(jen.Lit("invalid value in union type"))),
						),
						jen.Return(jen.Id("unknownFunc").Call(jen.Id(unionReceiverName).Dot("typ"))),
					)
					for _, fieldDef := range def.Fields {
						g.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(g *jen.Group) {
							selector := unionDerefPossibleOptional(g, fieldDef)
							g.Return(jen.Id(transforms.PrivateFieldName(fieldDef.Name) + "Func").Call(selector))
						})
					}
				}),
			)
		for _, fieldDef := range def.Fields {
			file.Func().
				Params(jen.Id(unionReceiverName).Op("*").Id(def.Name)).
				Id(transforms.ExportedFieldName(fieldDef.Name) + "NoopSuccess").
				Params(fieldDef.Type.Code()).
				Params(jen.Error()).
				Block(jen.Return(jen.Nil()))
		}
		file.Func().
			Params(jen.Id(unionReceiverName).Op("*").Id(def.Name)).
			Id("ErrorOnUnknown").
			Params(jen.Id("typeName").String()).
			Params(jen.Error()).
			Block(jen.Return(snip.FmtErrorf().Call(
				jen.Lit("invalid value in union type. Type name: %s"),
				jen.Id("typeName")),
			))
	}

	// Declare Accept/AcceptWithContext methods & visitor interfaces
	for _, withCtx := range []bool{false, true} {
		suffix := ""
		if withCtx {
			suffix = withContextSuffix
		}
		paramWithCtx := func(param *jen.Statement) func(*jen.Group) {
			return func(g *jen.Group) {
				if withCtx {
					g.Add(snip.ContextVar())
				}
				g.Add(param)
			}
		}
		// Accept method
		file.Func().
			Params(jen.Id(unionReceiverName).Op("*").Id(def.Name)).
			Id("Accept" + suffix).
			ParamsFunc(paramWithCtx(jen.Id("v").Id(def.Name + "Visitor" + suffix))).
			Params(jen.Error()).
			Block(jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(g *jen.Group) {
				g.Default().Block(
					jen.If(jen.Id(unionReceiverName).Dot("typ").Op("==").Lit("")).Block(
						jen.Return(snip.FmtErrorf().Call(jen.Lit("invalid value in union type"))),
					),
					jen.Return(jen.Id("v").Dot("VisitUnknown"+suffix).CallFunc(func(g *jen.Group) {
						if withCtx {
							g.Id("ctx")
						}
						g.Id(unionReceiverName).Dot("typ")
					})),
				)
				for _, fieldDef := range def.Fields {
					g.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(g *jen.Group) {
						fieldSelector := unionDerefPossibleOptional(g, fieldDef)
						g.Return(jen.Id("v").Dot("Visit" + transforms.ExportedFieldName(fieldDef.Name) + suffix).CallFunc(func(g *jen.Group) {
							if withCtx {
								g.Id("ctx")
							}
							g.Add(fieldSelector)
						}))
					})
				}
			}))
		// Visitor Interface
		file.Type().Id(def.Name + "Visitor" + suffix).InterfaceFunc(func(g *jen.Group) {
			for _, fieldDef := range def.Fields {
				g.Id("Visit" + transforms.ExportedFieldName(fieldDef.Name) + suffix).
					ParamsFunc(paramWithCtx(jen.Id("v").Add(fieldDef.Type.Code()))).
					Params(jen.Error())
			}
			g.Id("VisitUnknown" + suffix).
				ParamsFunc(func(g *jen.Group) {
					if withCtx {
						g.Add(snip.ContextVar())
					}
					g.Id("typeName").String()
				}).
				Params(jen.Error())
		})
	}

	// Declare New*From* constructor functions
	for _, fieldDef := range def.Fields {
		file.Func().
			Id("New" + def.Name + "From" + transforms.ExportedFieldName(fieldDef.Name)).
			Params(jen.Id("v").Add(fieldDef.Type.Code())).
			Params(jen.Id(def.Name)).
			Block(
				jen.Return(jen.Id(def.Name).Values(
					jen.Id("typ").Op(":").Lit(fieldDef.Name),
					jen.Id(transforms.PrivateFieldName(fieldDef.Name)).Op(":").Op("&").Id("v"),
				)),
			)
	}
}

func unionDerefPossibleOptional(g *jen.Group, fieldDef *types.Field) *jen.Statement {
	privateName := transforms.PrivateFieldName(fieldDef.Name)
	fieldSelector := jen.Op("*").Id(unionReceiverName).Dot(privateName)
	if fieldDef.Type.IsOptional() {
		// if the type is an optional and is nil, the value should not be dereferenced
		fieldSelector = jen.Id(privateName)
		g.Var().Id(privateName).Add(fieldDef.Type.Code())
		g.If(jen.Id(unionReceiverName).Dot(privateName).Op("!=").Nil()).Block(
			jen.Id(privateName).Op("=").Op("*").Id(unionReceiverName).Dot(privateName),
		)
	}
	return fieldSelector
}

func unionDeserializerStructName(unionTypeName string) string {
	return transforms.Private(transforms.ExportedFieldName(unionTypeName) + "Deserializer")
}
