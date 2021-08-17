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
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	aliasReceiverName   = "a"
	aliasValueFieldName = "Value"
)

func aliasDotValue() *jen.Statement { return jen.Id(aliasReceiverName).Dot(aliasValueFieldName) }

func writeAliasType(file *jen.Group, def *types.AliasType) {
	if def.IsOptional() {
		file.Add(def.Docs.CommentLine()).Type().Id(def.Name).Struct(
			jen.Id("Value").Add(def.Item.Code()),
		)
		// Optionals have special method ASTs.
		valueInit := def.Make()
		if valueInit == nil {
			valueInit = jen.New(jen.Id(strings.TrimPrefix(def.Item.Code().GoString(), "*")))
		}
		if def.IsBinary() {
			file.Add(astForAliasOptionalBinaryTextMarshal(def.Name))
			file.Add(astForAliasOptionalBinaryTextUnmarshal(def.Name))
		} else if def.IsString() {
			file.Add(astForAliasOptionalStringTextMarshal(def.Name))
			file.Add(astForAliasOptionalStringTextUnmarshal(def.Name))
		} else if def.IsText() {
			file.Add(astForAliasOptionalTextMarshal(def.Name))
			file.Add(astForAliasOptionalTextUnmarshal(def.Name, valueInit))
		} else {
			file.Add(astForAliasOptionalJSONMarshal(def.Name))
			file.Add(astForAliasOptionalJSONUnmarshal(def.Name, valueInit))
		}
		file.Add(snip.MethodMarshalYAML(aliasReceiverName, def.Name))
		file.Add(snip.MethodUnmarshalYAML(aliasReceiverName, def.Name))
	} else {
		file.Add(def.Docs.CommentLine()).Type().Id(def.Name).Add(def.Item.Code())
		if !isSimpleAliasType(def.Item) {
			// Everything else gets MarshalJSON/UnmarshalJSON that delegate to the aliased type
			if _, isBinary := def.Item.(types.Binary); isBinary {
				file.Add(astForAliasString(def.Name, snip.BinaryNew()))
				file.Add(astForAliasTextMarshal(def.Name, snip.BinaryNew()))
				file.Add(astForAliasBinaryTextUnmarshal(def.Name))
			} else if def.IsText() {
				// If we have gotten here, we have a non-go-builtin text type that implements MarshalText/UnmarshalText.
				file.Add(astForAliasString(def.Name, def.Item.Code()))
				file.Add(astForAliasTextMarshal(def.Name, def.Item.Code()))
				file.Add(astForAliasTextUnmarshal(def.Name, def.Item.Code()))
			} else {
				// By default, we delegate json/yaml encoding to the aliased type.
				file.Add(astForAliasJSONMarshal(def.Name, def.Item.Code()))
				file.Add(astForAliasJSONUnmarshal(def.Name, def.Item.Code()))
			}
			file.Add(snip.MethodMarshalYAML(aliasReceiverName, def.Name))
			file.Add(snip.MethodUnmarshalYAML(aliasReceiverName, def.Name))
		}
	}
}

func isSimpleAliasType(t types.Type) bool {
	switch v := t.(type) {
	case types.Any, types.Boolean, types.Double, types.Integer, types.String:
		// Plain builtins do not need encoding methods; do nothing.
		return true
	case *types.List:
		return isSimpleAliasType(v.Item)
	case *types.Map:
		return isSimpleAliasType(v.Key) && isSimpleAliasType(v.Val)
	case *types.Optional:
		return isSimpleAliasType(v.Item)
	case *types.AliasType:
		return isSimpleAliasType(v.Item)
	default:
		return false
	}
}

func astForAliasString(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodString(aliasReceiverName, typeName).Block(
		jen.Return(aliasGoType.Call(jen.Id(aliasReceiverName)).Dot("String").Call()),
	)
}

func astForAliasTextMarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodMarshalText(aliasReceiverName, typeName).Block(
		jen.Return(aliasGoType.Call(jen.Id(aliasReceiverName)).Dot("MarshalText").Call()),
	)
}

func astForAliasOptionalTextMarshal(typeName string) *jen.Statement {
	return snip.MethodMarshalText(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil().Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)),
		jen.Return(aliasDotValue().Dot("MarshalText").Call()),
	)
}

func astForAliasOptionalStringTextMarshal(typeName string) *jen.Statement {
	return snip.MethodMarshalText(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil().Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)),
		jen.Return(jen.Id("[]byte").Call(jen.Op("*").Id(aliasReceiverName).Dot("Value")), jen.Nil()),
	)
}

func astForAliasOptionalBinaryTextMarshal(typeName string) *jen.Statement {
	return snip.MethodMarshalText(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil().Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)),
		jen.Return(snip.BinaryNew().Call(jen.Op("*").Add(aliasDotValue())).Dot("MarshalText").Call()),
	)
}

func astForAliasTextUnmarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	rawVarName := "raw" + typeName
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.Var().Id(rawVarName).Add(aliasGoType),
		jen.If(
			jen.Err().Op(":=").Id(rawVarName).Dot("UnmarshalText").Call(jen.Id(dataVarName)),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		jen.Op("*").Id(aliasReceiverName).Op("=").Id(typeName).Call(jen.Id(rawVarName)),
		jen.Return(jen.Nil()),
	)
}

func astForAliasBinaryTextUnmarshal(typeName string) *jen.Statement {
	rawVarName := "raw" + typeName
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.List(jen.Id(rawVarName), jen.Err()).Op(":=").
			Add(snip.BinaryBinary()).Call(jen.Id(dataVarName)).Dot("Bytes").Call(),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Op("*").Id(aliasReceiverName).Op("=").Id(typeName).Call(jen.Id(rawVarName)),
		jen.Return(jen.Nil()),
	)
}

func astForAliasOptionalTextUnmarshal(typeName string, aliasValueInit *jen.Statement) *jen.Statement {
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil()).Block(
			aliasDotValue().Op("=").Add(aliasValueInit),
		),
		jen.Return(aliasDotValue().Dot("UnmarshalText").Call(jen.Id(dataVarName))),
	)
}

func astForAliasOptionalStringTextUnmarshal(typeName string) *jen.Statement {
	rawVarName := "raw" + typeName
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.Id(rawVarName).Op(":=").String().Call(jen.Id(dataVarName)),
		aliasDotValue().Op("=").Op("&").Id(rawVarName),
		jen.Return(jen.Nil()),
	)
}

func astForAliasOptionalBinaryTextUnmarshal(typeName string) *jen.Statement {
	rawVarName := "raw" + typeName
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.List(jen.Id(rawVarName), jen.Err()).Op(":=").
			Add(snip.BinaryBinary()).Call(jen.Id(dataVarName)).Dot("Bytes").Call(),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Op("*").Add(aliasDotValue()).Op("=").Id(rawVarName),
		jen.Return(jen.Nil()),
	)
}

func astForAliasJSONMarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodMarshalJSON(aliasReceiverName, typeName).Block(
		jen.Return(snip.SafeJSONMarshal().Call(aliasGoType.Call(jen.Id(aliasReceiverName)))),
	)
}

func astForAliasOptionalJSONMarshal(typeName string) *jen.Statement {
	return snip.MethodMarshalJSON(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil()).Block(
			jen.Return(jen.Nil(), jen.Nil()),
		),
		jen.Return(snip.SafeJSONMarshal().Call(aliasDotValue())),
	)
}

func astForAliasJSONUnmarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	rawVarName := "raw" + typeName
	return snip.MethodUnmarshalJSON(aliasReceiverName, typeName).Block(
		jen.Var().Id(rawVarName).Add(aliasGoType),
		jen.If(
			jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal()).Call(jen.Id(dataVarName), jen.Op("&").Id(rawVarName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		),
		jen.Op("*").Add(jen.Id(aliasReceiverName)).Op("=").Id(typeName).Call(jen.Id(rawVarName)),
		jen.Return(jen.Nil()),
	)
}

func astForAliasOptionalJSONUnmarshal(typeName string, aliasValueInit *jen.Statement) *jen.Statement {
	return snip.MethodUnmarshalJSON(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil()).Block(
			aliasDotValue().Op("=").Add(aliasValueInit),
		),
		jen.Return(snip.SafeJSONUnmarshal().Call(jen.Id(dataVarName), aliasDotValue())),
	)
}
