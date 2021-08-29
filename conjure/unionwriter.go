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
	"github.com/palantir/conjure-go/v6/conjure/encoding"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	unionReceiverName = "u"
	withContextSuffix = "WithContext"
)

func writeUnionType(file *jen.Group, unionDef *types.UnionType, cfg OutputConfiguration) {
	// Declare exported union struct type
	file.Add(unionDef.CommentLine()).
		Type().
		Id(unionDef.Name).StructFunc(func(structFields *jen.Group) {
		structFields.Id("typ").String()
		for _, fieldDef := range unionDef.Fields {
			structFields.Id(transforms.PrivateFieldName(fieldDef.Name)).Op("*").Add(fieldDef.Type.Code())
		}
	})

	// Declare New*From* constructor functions
	for _, fieldDef := range unionDef.Fields {
		file.Add(astForUnionConstructorFromElemFunc(unionDef.Name, fieldDef.Name, fieldDef.Type))
	}

	// Declare Accept/AcceptWithContext methods & visitor interfaces
	for _, withCtx := range []bool{false, true} {
		file.Add(astforUnionVisitorInterfaceTypeDecl(unionDef, withCtx))
		file.Add(astForUnionAcceptMethod(unionDef, withCtx))
	}

	// Declare AcceptFuncs method & noop helpers
	if cfg.GenerateFuncsVisitor {
		file.Add(astForUnionAcceptFuncsMethod(unionDef))
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

	if cfg.LiteralJSONMethods {
		file.Add(astForUnionLiteralMarshalJSON(unionDef))
		file.Add(astForUnionLiteralAppendJSON(unionDef))
		for _, stmt := range encoding.UnmarshalJSONMethods(unionReceiverName, unionDef.Name, unionDef) {
			file.Add(stmt)
		}
	} else {
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

		astForUnionReflectMarshalJSON(file, unionDef)
		astForUnionReflectUnmarshalJSON(file, unionDef)
	}

	// Declare yaml methods
	if cfg.GenerateYAMLMethods {
		file.Add(snip.MethodMarshalYAML(unionReceiverName, unionDef.Name))
		file.Add(snip.MethodUnmarshalYAML(unionReceiverName, unionDef.Name))
	}
}

func astForUnionLiteralMarshalJSON(unionDef *types.UnionType) *jen.Statement {
	return snip.MethodMarshalJSON(unionReceiverName, unionDef.Name).Block(
		encoding.MarshalJSONMethodBody(unionReceiverName),
	)
}

func astForUnionLiteralAppendJSON(unionDef *types.UnionType) *jen.Statement {
	return snip.MethodAppendJSON(unionReceiverName, unionDef.Name).BlockFunc(func(methodBody *jen.Group) {
		encoding.UnionMethodBodyAppendJSON(
			methodBody,
			jen.Id(unionReceiverName).Dot("typ").Clone,
			unionEncodingFields(unionDef.Fields),
		)
	})
}

func astForUnionReflectUnmarshalJSON(file *jen.Group, unionDef *types.UnionType) {
	file.Add(snip.MethodUnmarshalJSON(unionReceiverName, unionDef.Name).Block(
		jen.Var().Id("deser").Id(unionDeserializerStructName(unionDef.Name)),
		jen.If(
			jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal().Call(jen.Id(dataVarName), jen.Op("&").Id("deser"))),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		jen.Op("*").Id(unionReceiverName).Op("=").Id("deser").Dot("toStruct").Call(),
		jen.Return(jen.Nil()),
	))
}

func astForUnionReflectMarshalJSON(file *jen.Group, unionDef *types.UnionType) {
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
					fieldSelector := unionDerefPossibleOptional(caseBody, fieldDef)
					caseBody.Return(
						jen.Struct(
							jen.Id("Type").String().Tag(map[string]string{"json": "type"}),
							jen.Id(transforms.ExportedFieldName(fieldDef.Name)).
								Add(fieldDef.Type.Code()).
								Tag(map[string]string{"json": fieldDef.Name}),
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
}

func astForUnionAcceptFuncsMethod(unionDef *types.UnionType) *jen.Statement {
	return jen.Func().
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
						selector := unionDerefPossibleOptional(caseBody, fieldDef)
						caseBody.Return(jen.Id(transforms.PrivateFieldName(fieldDef.Name) + "Func").Call(selector))
					})
				}
			}),
		)
}

func astForUnionAcceptMethod(unionDef *types.UnionType, withCtx bool) *jen.Statement {
	suffix := ""
	if withCtx {
		suffix = withContextSuffix
	}
	return jen.Func().
		Params(jen.Id(unionReceiverName).Op("*").Id(unionDef.Name)).
		Id("Accept" + suffix).
		ParamsFunc(func(args *jen.Group) {
			if withCtx {
				args.Add(snip.ContextVar())
			}
			args.Id("v").Id(unionDef.Name + "Visitor" + suffix)
		}).
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
					fieldSelector := unionDerefPossibleOptional(caseBody, fieldDef)
					visitFn := "Visit" + transforms.ExportedFieldName(fieldDef.Name) + suffix
					caseBody.Return(jen.Id("v").Dot(visitFn).CallFunc(func(args *jen.Group) {
						if withCtx {
							args.Id("ctx")
						}
						args.Add(fieldSelector)
					}))
				})
			}
		}))
}

func astforUnionVisitorInterfaceTypeDecl(unionDef *types.UnionType, withCtx bool) *jen.Statement {
	suffix := ""
	if withCtx {
		suffix = withContextSuffix
	}
	return jen.Type().Id(unionDef.Name + "Visitor" + suffix).InterfaceFunc(func(methods *jen.Group) {
		for _, fieldDef := range unionDef.Fields {
			methods.Id("Visit" + transforms.ExportedFieldName(fieldDef.Name) + suffix).
				ParamsFunc(func(args *jen.Group) {
					if withCtx {
						args.Add(snip.Context())
					}
					args.Add(fieldDef.Type.Code())
				}).
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

func astForUnionConstructorFromElemFunc(typeName string, fieldName string, fieldType types.Type) *jen.Statement {
	return jen.Func().
		Id("New" + typeName + "From" + transforms.ExportedFieldName(fieldName)).
		Params(jen.Id("v").Add(fieldType.Code())).
		Params(jen.Id(typeName)).
		Block(
			jen.Return(jen.Id(typeName).ValuesFunc(func(values *jen.Group) {
				values.Id("typ").Op(":").Lit(fieldName)
				values.Id(transforms.PrivateFieldName(fieldName)).Op(":").Op("&").Id("v")
			})),
		)
}

func unionDerefPossibleOptional(caseBody *jen.Group, fieldDef *types.Field) *jen.Statement {
	privateName := transforms.PrivateFieldName(fieldDef.Name)
	fieldSelector := jen.Op("*").Id(unionReceiverName).Dot(privateName)
	if fieldDef.Type.IsOptional() {
		// if the type is an optional and is nil, the value should not be dereferenced
		fieldSelector = jen.Id(privateName)
		caseBody.Var().Id(privateName).Add(fieldDef.Type.Code())
		caseBody.If(jen.Id(unionReceiverName).Dot(privateName).Op("!=").Nil()).Block(
			jen.Id(privateName).Op("=").Op("*").Id(unionReceiverName).Dot(privateName),
		)
	}
	return fieldSelector
}

func unionDeserializerStructName(unionTypeName string) string {
	return transforms.Private(transforms.ExportedFieldName(unionTypeName) + "Deserializer")
}

func unionEncodingFields(spec []*types.Field) []encoding.JSONStructField {
	var fields []encoding.JSONStructField
	for _, field := range spec {
		fields = append(fields, encoding.JSONStructField{
			Key:      field.Name,
			Type:     field.Type,
			Selector: jen.Id(unionReceiverName).Dot(transforms.PrivateFieldName(field.Name)).Clone,
		})
	}
	return fields
}
