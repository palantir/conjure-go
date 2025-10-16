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
	"github.com/palantir/conjure-go/v6/conjure/jsonv2"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	aliasReceiverName   = "a"
	aliasValueFieldName = "Value"
)

func aliasDotValue() *jen.Statement { return jen.Id(aliasReceiverName).Dot(aliasValueFieldName) }

func writeAliasType(file *jen.Group, aliasDef *types.AliasType, cfg OutputConfiguration) {
	if aliasDef.IsOptional() {
		writeOptionalAliasType(file, aliasDef, cfg)
	} else {
		writeNonOptionalAliasType(file, aliasDef, cfg)
	}
}

func writeOptionalAliasType(file *jen.Group, aliasDef *types.AliasType, cfg OutputConfiguration) {
	typeName := aliasDef.Name
	// Define the type

	safety := aliasDef.Safety()
	if aliasDef.Docs != "" {
		file.Comment(string(aliasDef.Docs))
	}
	if !safety.IsUnknown() {
		// Add safety comment if type has safety annotation
		file.Comment("safelogging:" + logSafetyToAnnotation(safety.Value()))
	}
	file.Type().Id(typeName).Struct(
		jen.Id("Value").Add(aliasDef.Item.Code()),
	)
	opt := aliasDef.Item.(*types.Optional)
	// text methods
	switch {
	case opt.Item.IsString():
		file.Add(astForAliasStringOptional(typeName))
		file.Add(astForAliasOptionalStringTextMarshal(typeName))
		file.Add(astForAliasOptionalStringTextUnmarshal(typeName, opt.Item.Code()))
	case opt.Item.IsBinary():
		file.Add(astForAliasOptionalBinaryTextMarshal(typeName))
		file.Add(astForAliasOptionalBinaryTextUnmarshal(typeName))
	case opt.Item.IsText():
		valueInit := aliasDef.Make()
		if valueInit == nil {
			valueInit = jen.New(opt.Item.Code())
		}
		file.Add(astForAliasOptionalTextMarshal(typeName))
		file.Add(astForAliasOptionalTextUnmarshal(typeName, valueInit))
	}
	// json methods
	for _, c := range astForAliasTypeMarshalJSON(aliasDef, cfg) {
		file.Add(c)
	}
	for _, c := range astForAliasTypeUnmarshalJSON(aliasDef, cfg) {
		file.Add(c)
	}
	// yaml methods
	if cfg.JSONv2 {
		file.Add(snip.MethodJSONV2MarshalYAML(aliasReceiverName, aliasDef.Name))
		file.Add(snip.MethodJSONV2UnmarshalYAML(aliasReceiverName, aliasDef.Name))
	} else {
		file.Add(snip.MethodMarshalYAML(aliasReceiverName, aliasDef.Name))
		file.Add(snip.MethodUnmarshalYAML(aliasReceiverName, aliasDef.Name))
	}
}

func writeNonOptionalAliasType(file *jen.Group, aliasDef *types.AliasType, cfg OutputConfiguration) {
	typeName := aliasDef.Name
	// Define the type

	// Add safety inline comment if type has safety annotation
	safety := aliasDef.Safety()
	typeStatement := file.Add(aliasDef.Docs.CommentLine()).Type().Id(typeName).Add(aliasDef.Item.Code())
	if !safety.IsUnknown() {
		typeStatement.Comment("safelogging:" + logSafetyToAnnotation(safety.Value()))
		file.Line()
	}
	if aliasDef.IsInterface() {
		// no methods allowed on interface alias
		return
	}

	if !cfg.JSONv2 && isSimpleAliasType(aliasDef) {
		// no methods allowed on interface alias
		return
	}
	// text methods
	switch {
	case aliasDef.IsString():
		file.Add(astForAliasString(typeName))
		file.Add(astForAliasStringTextMarshal(typeName))
		file.Add(astForAliasStringTextUnmarshal(typeName))
	case aliasDef.IsBinary():
		file.Add(astForAliasStringer(typeName, snip.BinaryNew()))
		file.Add(astForAliasTextMarshal(typeName, snip.BinaryNew()))
		file.Add(astForAliasBinaryTextUnmarshal(typeName))
	case aliasDef.IsText():
		file.Add(astForAliasStringer(typeName, aliasDef.Item.Code()))
		// If we have gotten here, we have a non-go-builtin text type that implements MarshalText/UnmarshalText.
		file.Add(astForAliasTextMarshal(typeName, aliasDef.Item.Code()))
		file.Add(astForAliasTextUnmarshal(typeName, aliasDef.Item.Code()))
	}
	// json methods
	for _, c := range astForAliasTypeMarshalJSON(aliasDef, cfg) {
		file.Add(c)
	}
	for _, c := range astForAliasTypeUnmarshalJSON(aliasDef, cfg) {
		file.Add(c)
	}
	// yaml methods
	if cfg.JSONv2 {
		file.Add(snip.MethodJSONV2MarshalYAML(aliasReceiverName, aliasDef.Name))
		file.Add(snip.MethodJSONV2UnmarshalYAML(aliasReceiverName, aliasDef.Name))
	} else {
		file.Add(snip.MethodMarshalYAML(aliasReceiverName, aliasDef.Name))
		file.Add(snip.MethodUnmarshalYAML(aliasReceiverName, aliasDef.Name))
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
	case *types.External:
		return isSimpleAliasType(v.Fallback)
	default:
		return false
	}
}

func astForAliasString(typeName string) *jen.Statement {
	return snip.MethodString(aliasReceiverName, typeName).Block(
		jen.Return(jen.String().Call(jen.Id(aliasReceiverName))),
	)
}

func astForAliasStringer(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodString(aliasReceiverName, typeName).Block(
		jen.Return(aliasGoType.Call(jen.Id(aliasReceiverName)).Dot("String").Call()),
	)
}

func astForAliasStringOptional(typeName string) *jen.Statement {
	return snip.MethodString(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil().Block(
			jen.Return(jen.Lit("")),
		)),
		jen.Return(jen.String().Call(jen.Op("*").Add(aliasDotValue()))),
	)
}

func astForAliasTextMarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodMarshalText(aliasReceiverName, typeName).Block(
		jen.Return(aliasGoType.Call(jen.Id(aliasReceiverName)).Dot("MarshalText").Call()),
	)
}

func astForAliasStringTextMarshal(typeName string) *jen.Statement {
	return snip.MethodMarshalText(aliasReceiverName, typeName).Block(
		jen.Return(jen.Index().Byte().Call(jen.Id(aliasReceiverName)), jen.Nil()),
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
		jen.Return(jen.Index().Byte().Call(jen.Op("*").Add(aliasDotValue())), jen.Nil()),
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
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.Return(jen.Parens(jen.Op("*").Add(aliasGoType)).Call(jen.Id(aliasReceiverName)).Dot("UnmarshalText").Call(jen.Id(dataVarName))),
	)
}

func astForAliasStringTextUnmarshal(typeName string) *jen.Statement {
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.Op("*").Id(aliasReceiverName).Op("=").Id(typeName).Call(jen.Id(dataVarName)),
		jen.Return(jen.Nil()),
	)
}

func astForAliasBinaryTextUnmarshal(typeName string) *jen.Statement {
	rawVarName := "raw" + typeName
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.List(jen.Id(rawVarName), jen.Err()).Op(":=").
			Add(snip.BinaryBinary()).Call(jen.Id(dataVarName)).Dot("Bytes").Call(),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Op("*").Id(aliasReceiverName).Op("=").Id(rawVarName),
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

func astForAliasOptionalStringTextUnmarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodUnmarshalText(aliasReceiverName, typeName).Block(
		jen.Op("*").Add(aliasDotValue()).Op("=").Add(aliasGoType).Call(jen.Id(dataVarName)),
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

func astForAliasTypeMarshalJSON(aliasDef *types.AliasType, cfg OutputConfiguration) []jen.Code {
	if aliasDef.IsOptional() {
		if cfg.JSONv2 {
			return []jen.Code{
				snip.MethodMarshalJSONV2(aliasReceiverName, aliasDef.Name),
				jsonv2.MethodMarshalJSONTo(aliasReceiverName, aliasDef.Name, aliasDef),
			}
		}
		return []jen.Code{astForAliasOptionalJSONMarshal(aliasDef.Name)}
	}
	if cfg.JSONv2 {
		return []jen.Code{
			snip.MethodMarshalJSONV2(aliasReceiverName, aliasDef.Name),
			jsonv2.MethodMarshalJSONTo(aliasReceiverName, aliasDef.Name, aliasDef),
		}
	}
	return []jen.Code{astForAliasJSONMarshal(aliasDef.Name, aliasDef.Item.Code())}
}

func astForAliasTypeUnmarshalJSON(aliasDef *types.AliasType, cfg OutputConfiguration) []jen.Code {
	typeName := aliasDef.Name
	if aliasDef.IsOptional() {
		if cfg.JSONv2 {
			return []jen.Code{
				snip.MethodUnmarshalJSONV2(aliasReceiverName, aliasDef.Name),
				jsonv2.MethodUnmarshalJSONFrom(aliasReceiverName, aliasDef.Name, aliasDef),
			}
		}
		opt := aliasDef.Item.(*types.Optional)
		valueInit := aliasDef.Make()
		if valueInit == nil {
			valueInit = jen.New(opt.Item.Code())
		}
		return []jen.Code{astForAliasOptionalJSONUnmarshal(typeName, valueInit)}
	}
	if cfg.JSONv2 {
		return []jen.Code{
			snip.MethodUnmarshalJSONV2(aliasReceiverName, aliasDef.Name),
			jsonv2.MethodUnmarshalJSONFrom(aliasReceiverName, aliasDef.Name, aliasDef),
		}
	}
	return []jen.Code{astForAliasJSONUnmarshal(typeName, aliasDef.Item.Code())}
}

func astForAliasJSONMarshal(typeName string, aliasGoType *jen.Statement) *jen.Statement {
	return snip.MethodMarshalJSON(aliasReceiverName, typeName).Block(
		jen.Return(snip.SafeJSONMarshal().Call(aliasGoType.Call(jen.Id(aliasReceiverName)))),
	)
}

func astForAliasOptionalJSONMarshal(typeName string) *jen.Statement {
	return snip.MethodMarshalJSON(aliasReceiverName, typeName).Block(
		jen.If(aliasDotValue().Op("==").Nil()).Block(
			jen.Return(jen.Index().Byte().Call(jen.Lit("null")), jen.Nil()),
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
		jen.Op("*").Id(aliasReceiverName).Op("=").Id(typeName).Call(jen.Id(rawVarName)),
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
