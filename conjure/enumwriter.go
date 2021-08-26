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
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	enumReceiverName    = "e"
	enumUpperVarName    = "v"
	enumUnknownValue    = "UNKNOWN"
	enumStructFieldName = "val"
)

func writeEnumType(file *jen.Group, def *types.EnumType, cfg OutputConfiguration) {
	file.Add(def.CommentLine()).Add(astForEnumTypeDecls(def.Name))
	file.Add(astForEnumValueConstants(def.Name, def.Values))
	file.Add(astForEnumValuesFunction(def.Name, def.Values))
	file.Add(astForEnumConstructor(def.Name))
	file.Add(astForEnumIsUnknown(def.Name, def.Values))
	file.Add(astForEnumValueMethod(def.Name))
	file.Add(astForEnumStringMethod(def.Name))
	if cfg.LiteralJSONMethods {
		file.Add(astForEnumLiteralMarshalJSON(def.Name))
		file.Add(astForEnumLiteralAppendJSON(def.Name))
	} else {
		file.Add(astForEnumMarshalText(def.Name))
	}
	file.Add(astForEnumUnmarshalText(def.Name, def.Values))
}

func astForEnumTypeDecls(typeName string) *jen.Statement {
	return jen.Type().Id(typeName).Struct(jen.Id(enumStructFieldName).Id(typeName + "_Value")).
		Line().Line().
		Type().Id(typeName + "_Value").String()
}

func astForEnumValueConstants(typeName string, values []*types.Field) *jen.Statement {
	return jen.Const().DefsFunc(func(g *jen.Group) {
		for _, valDef := range values {
			g.Add(valDef.CommentLine()).
				Id(typeName + "_" + valDef.Name).Id(typeName + "_Value").Op("=").Lit(valDef.Name)
		}
		g.Id(typeName + "_" + enumUnknownValue).Id(typeName + "_Value").Op("=").Lit(enumUnknownValue)
	})
}

func astForEnumValuesFunction(typeName string, values []*types.Field) *jen.Statement {
	return jen.Commentf("%s_Values returns all known variants of %s.", typeName, typeName).
		Line().
		Func().
		Id(typeName + "_Values").
		Params().
		Params(jen.Op("[]").Id(typeName + "_Value")).
		Block(
			jen.Return(jen.Op("[]").Id(typeName + "_Value").ValuesFunc(func(g *jen.Group) {
				for _, valDef := range values {
					g.Id(typeName + "_" + valDef.Name)
				}
			})),
		)
}

func astForEnumConstructor(typeName string) *jen.Statement {
	return jen.Func().Id("New_" + typeName).Params(jen.Id("value").Id(typeName + "_Value")).Params(jen.Id(typeName)).Block(
		jen.Return(jen.Id(typeName).Values(jen.Id(enumStructFieldName).Op(":").Id("value"))),
	)
}

func astForEnumIsUnknown(typeName string, values []*types.Field) *jen.Statement {
	return jen.Commentf("IsUnknown returns false for all known variants of %s and true otherwise.", typeName).
		Line().
		Func().
		Params(jen.Id(enumReceiverName).Id(typeName)).Id("IsUnknown").Params().Params(jen.Bool()).Block(
		jen.Switch(jen.Id(enumReceiverName).Dot(enumStructFieldName)).Block(
			jen.CaseFunc(func(g *jen.Group) {
				for _, valDef := range values {
					g.Id(typeName + "_" + valDef.Name)
				}
			}).
				Block(jen.Return(jen.False())),
		),
		jen.Return(jen.True()),
	)
}

func astForEnumValueMethod(typeName string) *jen.Statement {
	return jen.Func().
		Params(jen.Id(enumReceiverName).Id(typeName)).
		Id("Value").
		Params().
		Params(jen.Id(typeName+"_Value")).
		Block(
			jen.If(jen.Id(enumReceiverName).Dot("IsUnknown").Call()).Block(
				jen.Return(jen.Id(typeName+"_"+enumUnknownValue)),
			),
			jen.Return(jen.Id(enumReceiverName).Dot(enumStructFieldName)),
		)
}

func astForEnumStringMethod(typeName string) *jen.Statement {
	return snip.MethodString(enumReceiverName, typeName).Block(
		jen.Return(jen.String().Call(jen.Id(enumReceiverName).Dot(enumStructFieldName))),
	)
}

func astForEnumMarshalText(typeName string) *jen.Statement {
	return snip.MethodMarshalText(enumReceiverName, typeName).Block(
		jen.Return(jen.Id("[]byte").Call(jen.Id(enumReceiverName).Dot(enumStructFieldName)), jen.Nil()),
	)
}

func astForEnumLiteralMarshalJSON(typeName string) *jen.Statement {
	return snip.MethodMarshalJSON(enumReceiverName, typeName).Block(
		encoding.MarshalJSONMethodBody(enumReceiverName),
	)
}

func astForEnumLiteralAppendJSON(typeName string) *jen.Statement {
	return snip.MethodAppendJSON(enumReceiverName, typeName).BlockFunc(func(g *jen.Group) {
		encoding.EnumMethodBodyAppendJSON(g, enumReceiverName)
	})
}

func astForEnumUnmarshalText(typeName string, values []*types.Field) *jen.Statement {
	return snip.MethodUnmarshalText(enumReceiverName, typeName).Block(
		jen.Switch(
			jen.Id(enumUpperVarName).Op(":=").Add(snip.StringsToUpper()).Call(jen.String().Call(jen.Id(dataVarName))),
			jen.Id(enumUpperVarName),
		).BlockFunc(func(g *jen.Group) {
			assign := func(val jen.Code) *jen.Statement {
				return jen.Op("*").Add(jen.Id(enumReceiverName)).Op("=").Id("New_" + typeName).Call(val)
			}
			g.Default().Block(assign(jen.Id(typeName + "_Value").Call(jen.Id(enumUpperVarName))))
			for _, valDef := range values {
				g.Case(jen.Lit(valDef.Name)).Block(assign(jen.Id(typeName + "_" + valDef.Name)))
			}
		}),
		jen.Return(jen.Nil()),
	)
}