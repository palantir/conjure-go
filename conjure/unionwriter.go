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
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	unionReceiverName = "u"
	withContextSuffix = "WithContext"
)

func writeUnionType(file *jen.Group, unionDef *types.UnionType, genAcceptFuncs bool) {
	// Declare exported union struct type
	file.Add(unionDef.CommentLine()).
		Type().
		Id(unionDef.Name).StructFunc(func(structFields *jen.Group) {
		structFields.Id("typ").String()
		for _, fieldDef := range unionDef.Fields {
			structFields.Id(transforms.PrivateFieldName(fieldDef.Name)).Op("*").Add(fieldDef.Type.Code())
		}
	})

	// Declare deserializer struct type
	file.Type().
		Id(unionDeserializerStructName(unionDef.Name)).StructFunc(func(structFields *jen.Group) {
		structFields.Id("Type").String().Tag(map[string]string{"json": "type"})
		for _, fieldDef := range unionDef.Fields {
			structFields.Id(transforms.ExportedFieldName(fieldDef.Name)).
				Op("*").Add(fieldDef.Type.Code()).
				Tag(map[string]string{"json": fieldDef.Name})
		}
	})

	// Declare deserializer toStruct method
	file.Func().
		Params(jen.Id(unionReceiverName).Op("*").Id(unionDeserializerStructName(unionDef.Name))).
		Id("toStruct").
		Params().
		Params(jen.Id(unionDef.Name)).
		Block(jen.Return(jen.Id(unionDef.Name).ValuesFunc(func(values *jen.Group) {
			values.Id("typ").Op(":").Id(unionReceiverName).Dot("Type")
			for _, fieldDef := range unionDef.Fields {
				values.Id(transforms.PrivateFieldName(fieldDef.Name)).
					Op(":").
					Id(unionReceiverName).Dot(transforms.ExportedFieldName(fieldDef.Name))
			}
		})))

	// Declare toSerializer method
	file.Func().
		Params(jen.Id(unionReceiverName).Op("*").Id(unionDef.Name)).
		Id("toSerializer").
		Params().
		Params(jen.Interface(), jen.Error()).
		Block(jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
			cases.Default().Block(jen.Return(
				jen.Nil(), snip.FmtErrorf().Call(jen.Lit("unknown type %s"), jen.Id(unionReceiverName).Dot("typ"))))
			for _, fieldDef := range unionDef.Fields {
				cases.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(caseBody *jen.Group) {
					fieldSelector := unionDerefPossibleOptional(caseBody, fieldDef, jen.Nil())
					caseBody.Return(
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
	file.Add(snip.MethodMarshalJSON(unionReceiverName, unionDef.Name).Block(
		jen.List(jen.Id("ser"), jen.Err()).Op(":=").Id(unionReceiverName).Dot("toSerializer").Call(),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Return(snip.SafeJSONMarshal().Call(jen.Id("ser"))),
	))

	// Declare UnmarshalJSON method
	file.Add(snip.MethodUnmarshalJSON(unionReceiverName, unionDef.Name).Block(
		jen.Var().Id("deser").Id(unionDeserializerStructName(unionDef.Name)),
		jen.If(
			jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal().Call(jen.Id(dataVarName), jen.Op("&").Id("deser"))),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		jen.Op("*").Id(unionReceiverName).Op("=").Id("deser").Dot("toStruct").Call(),
		jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
			cases.Default().Block(jen.Return(
				snip.FmtErrorf().Call(jen.Lit("unknown type %s"), jen.Id(unionReceiverName).Dot("typ"))))
			for _, fieldDef := range unionDef.Fields {
				cases.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(caseBody *jen.Group) {
					if !fieldDef.Type.IsOptional() {
						selector := jen.Id(unionReceiverName).Dot(transforms.PrivateFieldName(fieldDef.Name))
						caseBody.If(selector.Op("==").Nil()).Block(
							jen.Return(snip.FmtErrorf().Call(jen.Lit(fmt.Sprintf("field %s is required", fieldDef.Name)))),
						)
					}
				})
			}
		}),
		jen.Return(jen.Nil()),
	))

	// Declare yaml methods
	file.Add(snip.MethodMarshalYAML(unionReceiverName, unionDef.Name))
	file.Add(snip.MethodUnmarshalYAML(unionReceiverName, unionDef.Name))

	// Declare AcceptFuncs method & noop helpers
	if genAcceptFuncs {
		file.Func().
			Params(jen.Id(unionReceiverName).Op("*").Id(unionDef.Name)).
			Id("AcceptFuncs").
			ParamsFunc(func(args *jen.Group) {
				for _, fieldDef := range unionDef.Fields {
					args.Id(transforms.PrivateFieldName(fieldDef.Name) + "Func").Func().Params(fieldDef.Type.Code()).Params(jen.Error())
				}
				args.Id("unknownFunc").Func().Params(jen.String()).Params(jen.Error())
			}).
			Params(jen.Error()).
			Block(
				jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
					cases.Default().Block(
						jen.If(jen.Id(unionReceiverName).Dot("typ").Op("==").Lit("")).Block(
							jen.Return(snip.FmtErrorf().Call(jen.Lit("invalid value in union type"))),
						),
						jen.Return(jen.Id("unknownFunc").Call(jen.Id(unionReceiverName).Dot("typ"))),
					)
					for _, fieldDef := range unionDef.Fields {
						cases.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(caseBody *jen.Group) {
							selector := unionDerefPossibleOptional(caseBody, fieldDef, nil)
							caseBody.Return(jen.Id(transforms.PrivateFieldName(fieldDef.Name) + "Func").Call(selector))
						})
					}
				}),
			)
		for _, fieldDef := range unionDef.Fields {
			file.Func().
				Params(jen.Id(unionReceiverName).Op("*").Id(unionDef.Name)).
				Id(transforms.ExportedFieldName(fieldDef.Name) + "NoopSuccess").
				Params(fieldDef.Type.Code()).
				Params(jen.Error()).
				Block(jen.Return(jen.Nil()))
		}
		file.Func().
			Params(jen.Id(unionReceiverName).Op("*").Id(unionDef.Name)).
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
			return func(args *jen.Group) {
				if withCtx {
					args.Add(snip.ContextVar())
				}
				args.Add(param)
			}
		}
		// Accept method
		file.Func().
			Params(jen.Id(unionReceiverName).Op("*").Id(unionDef.Name)).
			Id("Accept" + suffix).
			ParamsFunc(paramWithCtx(jen.Id("v").Id(unionDef.Name + "Visitor" + suffix))).
			Params(jen.Error()).
			Block(jen.Switch(jen.Id(unionReceiverName).Dot("typ")).BlockFunc(func(cases *jen.Group) {
				cases.Default().Block(
					jen.If(jen.Id(unionReceiverName).Dot("typ").Op("==").Lit("")).Block(
						jen.Return(snip.FmtErrorf().Call(jen.Lit("invalid value in union type"))),
					),
					jen.Return(jen.Id("v").Dot("VisitUnknown"+suffix).CallFunc(func(args *jen.Group) {
						if withCtx {
							args.Id("ctx")
						}
						args.Id(unionReceiverName).Dot("typ")
					})),
				)
				for _, fieldDef := range unionDef.Fields {
					cases.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(caseBody *jen.Group) {
						fieldSelector := unionDerefPossibleOptional(caseBody, fieldDef, nil)
						caseBody.Return(jen.Id("v").Dot("Visit" + transforms.ExportedFieldName(fieldDef.Name) + suffix).CallFunc(func(args *jen.Group) {
							if withCtx {
								args.Id("ctx")
							}
							args.Add(fieldSelector)
						}))
					})
				}
			}))
		// Visitor Interface
		file.Type().Id(unionDef.Name + "Visitor" + suffix).InterfaceFunc(func(methods *jen.Group) {
			for _, fieldDef := range unionDef.Fields {
				methods.Id("Visit" + transforms.ExportedFieldName(fieldDef.Name) + suffix).
					ParamsFunc(paramWithCtx(jen.Id("v").Add(fieldDef.Type.Code()))).
					Params(jen.Error())
			}
			methods.Id("VisitUnknown" + suffix).
				ParamsFunc(func(args *jen.Group) {
					if withCtx {
						args.Add(snip.ContextVar())
					}
					args.Id("typeName").String()
				}).
				Params(jen.Error())
		})
	}

	// Declare New*From* constructor functions
	for _, fieldDef := range unionDef.Fields {
		file.Func().
			Id("New" + unionDef.Name + "From" + transforms.ExportedFieldName(fieldDef.Name)).
			Params(jen.Id("v").Add(fieldDef.Type.Code())).
			Params(jen.Id(unionDef.Name)).
			Block(
				jen.Return(jen.Id(unionDef.Name).Values(
					jen.Id("typ").Op(":").Lit(fieldDef.Name),
					jen.Id(transforms.PrivateFieldName(fieldDef.Name)).Op(":").Op("&").Id("v"),
				)),
			)
	}
}

func unionDerefPossibleOptional(caseBody *jen.Group, fieldDef *types.Field, returnVal jen.Code) *jen.Statement {
	privateName := transforms.PrivateFieldName(fieldDef.Name)
	var fieldSelector *jen.Statement
	if fieldDef.Type.IsOptional() {
		// if the type is an optional and is nil, the value should not be dereferenced
		fieldSelector = jen.Id(privateName)
		caseBody.Var().Id(privateName).Add(fieldDef.Type.Code())
		caseBody.If(jen.Id(unionReceiverName).Dot(privateName).Op("!=").Nil()).Block(
			jen.Id(privateName).Op("=").Op("*").Id(unionReceiverName).Dot(privateName),
		)
	} else {
		caseBody.If(jen.Id(unionReceiverName).Dot(privateName).Op("==").Nil()).Block(
			jen.ReturnFunc(func(returns *jen.Group) {
				if returnVal != nil {
					returns.Add(returnVal)
				}
				returns.Add(snip.FmtErrorf().Call(jen.Lit(fmt.Sprintf("field %s is required", fieldDef.Name))))
			}),
		)
		fieldSelector = jen.Op("*").Id(unionReceiverName).Dot(privateName)
	}
	return fieldSelector
}

func unionDeserializerStructName(unionTypeName string) string {
	return transforms.Private(transforms.ExportedFieldName(unionTypeName) + "Deserializer")
}

func writeUnionTypeWithGenerics(file *jen.Group, unionType *types.UnionType) {
	unionTypeWithT(file, unionType)
	unionTypeWithTAccept(file, unionType)
	unionVisitorWithT(file, unionType)
}

func unionTypeWithT(file *jen.Group, unionType *types.UnionType) {
	file.Type().
		Id(unionType.Name + "WithT").
		Add(snip.TAny()).
		Add(unionType.Code())
}

func unionTypeWithTAccept(file *jen.Group, unionType *types.UnionType) {
	file.
		Func().
		Params(jen.Id(unionReceiverName).Op("*").Id(unionType.Name+"WithT").Op("[").Id("T").Op("]")).
		Id("Accept").
		Params(snip.ContextVar(), jen.Id("v").Id(unionType.Name+"VisitorWithT").Op("[").Id("T").Op("]")).
		Params(jen.Id("T"), jen.Error()).
		Block(
			jen.Var().Id("result").Id("T"),
			jen.Switch(jen.Id(unionReceiverName).Dot("typ")).
				BlockFunc(func(cases *jen.Group) {
					cases.Default().Block(
						jen.If(jen.Id(unionReceiverName).Dot("typ").Op("==").Lit("")).Block(
							jen.Return(jen.Id("result"), snip.FmtErrorf().Call(jen.Lit("invalid value in union type"))),
						),
						jen.Return(jen.Id("v").Dot("VisitUnknown").Call(jen.Id("ctx"), jen.Id(unionReceiverName).Dot("typ"))),
					)
					for _, fieldDef := range unionType.Fields {
						cases.Case(jen.Lit(fieldDef.Name)).BlockFunc(func(caseBody *jen.Group) {
							caseBody.
								Return(jen.Id("v").
									Dot("Visit"+transforms.ExportedFieldName(fieldDef.Name)).
									Call(jen.Id("ctx"), unionDerefPossibleOptional(caseBody, fieldDef, jen.Id("result"))))
						})
					}
				}),
		)
}

func unionVisitorWithT(file *jen.Group, union *types.UnionType) {
	file.
		Type().
		Id(union.Name + "VisitorWithT").
		Add(snip.TAny()).
		InterfaceFunc(func(methods *jen.Group) {
			for _, fieldDef := range union.Fields {
				methods.Id("Visit"+transforms.ExportedFieldName(fieldDef.Name)).
					Params(snip.ContextVar(), jen.Id("v").Add(fieldDef.Type.Code())).
					Params(jen.Id("T"), jen.Error())
			}
			methods.Id("VisitUnknown").
				Params(snip.ContextVar(), jen.Id("typ").Add(jen.String())).
				Params(jen.Id("T"), jen.Error())
		})
}
