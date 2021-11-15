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
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	enumReceiverName    = "e"
	enumUpperVarName    = "v"
	enumUnknownValue    = "UNKNOWN"
	enumStructFieldName = "val"
)

func writeEnumType(file *jen.Group, enumDef *types.EnumType) {
	file.Add(enumDef.CommentLine()).Add(astForEnumTypeDecls(enumDef.Name))
	file.Add(astForEnumValueConstants(enumDef.Name, enumDef.Values))
	file.Add(astForEnumValuesFunction(enumDef.Name, enumDef.Values))
	file.Add(astForEnumConstructor(enumDef.Name))
	file.Add(astForEnumIsUnknown(enumDef.Name, enumDef.Values))
	file.Add(astForEnumValueMethod(enumDef.Name))
	file.Add(astForEnumStringMethod(enumDef.Name))
	file.Add(astForEnumMarshalText(enumDef.Name))
	file.Add(astForEnumUnmarshalText(enumDef.Name, enumDef.Values))
}

func astForEnumTypeDecls(typeName string) *jen.Statement {
	return jen.Type().Id(typeName).Struct(jen.Id(enumStructFieldName).Id(typeName + "_Value")).
		Line().Line().
		Type().Id(typeName + "_Value").String()
}

func astForEnumValueConstants(typeName string, values []*types.Field) *jen.Statement {
	return jen.Const().DefsFunc(func(consts *jen.Group) {
		for _, valDef := range values {
			consts.Add(valDef.CommentLine()).
				Id(typeName + "_" + valDef.Name).Id(typeName + "_Value").Op("=").Lit(valDef.Name)
		}
		consts.Id(typeName + "_" + enumUnknownValue).Id(typeName + "_Value").Op("=").Lit(enumUnknownValue)
	})
}

func astForEnumValuesFunction(typeName string, enumValues []*types.Field) *jen.Statement {
	return jen.Commentf("%s_Values returns all known variants of %s.", typeName, typeName).
		Line().
		Func().
		Id(typeName + "_Values").
		Params().
		Params(jen.Op("[]").Id(typeName + "_Value")).
		Block(
			jen.Return(jen.Op("[]").Id(typeName + "_Value").ValuesFunc(func(values *jen.Group) {
				for _, valDef := range enumValues {
					values.Id(typeName + "_" + valDef.Name)
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
		Params(jen.Id(enumReceiverName).Id(typeName)).Id("IsUnknown").Params().Params(jen.Bool()).BlockFunc(func(methodBody *jen.Group) {
		methodBody.Switch(jen.Id(enumReceiverName).Dot(enumStructFieldName)).BlockFunc(func(switchBlock *jen.Group) {
			if len(values) == 0 {
				switchBlock.Default().Return(jen.False())
			} else {
				switchBlock.CaseFunc(func(conds *jen.Group) {
					for _, valDef := range values {
						conds.Id(typeName + "_" + valDef.Name)
					}
				}).Block(jen.Return(jen.False()))
			}
		methodBody.Return(jen.True())
		})})
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

func astForEnumUnmarshalText(typeName string, values []*types.Field) *jen.Statement {
	return snip.MethodUnmarshalText(enumReceiverName, typeName).Block(
		jen.Switch(
			jen.Id(enumUpperVarName).Op(":=").Add(snip.StringsToUpper()).Call(jen.String().Call(jen.Id(dataVarName))),
			jen.Id(enumUpperVarName),
		).BlockFunc(func(cases *jen.Group) {
			assign := func(val jen.Code) *jen.Statement {
				return jen.Op("*").Add(jen.Id(enumReceiverName)).Op("=").Id("New_" + typeName).Call(val)
			}
			cases.Default().Block(assign(jen.Id(typeName + "_Value").Call(jen.Id(enumUpperVarName))))
			for _, valDef := range values {
				cases.Case(jen.Lit(valDef.Name)).Block(assign(jen.Id(typeName + "_" + valDef.Name)))
			}
		}),
		jen.Return(jen.Nil()),
	)
}
