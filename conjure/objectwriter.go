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
	"github.com/palantir/conjure-go/v6/conjure/encoding"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	objReceiverName = "o"
	dataVarName     = "data"
)

func writeObjectType(file *jen.Group, def *types.ObjectType, cfg OutputConfiguration) {
	// Declare struct type with fields
	containsCollection := false // If contains collection, we need JSON methods to initialize empty values.
	file.Add(def.Docs.CommentLine()).Type().Id(def.Name).StructFunc(func(structDecl *jen.Group) {
		for _, fieldDef := range def.Fields {
			fieldName := fieldDef.Name
			fieldTags := map[string]string{"json": fieldName}

			if fieldDef.Docs != "" {
				// backtick characters ("`") are really painful to deal with in struct tags
				// (which are themselves defined within backtick literals), so replace with
				// double quotes instead.
				fieldTags["conjure-docs"] = strings.Replace(string(fieldDef.Docs), "`", `"`, -1)
			}
			if fieldDef.Type.Make() != nil {
				containsCollection = true
			}
			structDecl.Add(fieldDef.Docs.CommentLine()).Id(transforms.ExportedFieldName(fieldName)).Add(fieldDef.Type.Code()).Tag(fieldTags)
		}
	})

	if cfg.LiteralJSONMethods {
		file.Add(astForStructLiteralMarshalJSON(def))
		file.Add(astForStructLiteralAppendJSON(def))
	} else {
		// If there are no collections, we can defer to the default json behavior
		// Otherwise we need to override MarshalJSON and UnmarshalJSON
		if containsCollection {
			tmpAliasName := def.Name + "Alias"
			// Declare MarshalJSON
			file.Add(snip.MethodMarshalJSON(objReceiverName, def.Name).BlockFunc(func(methodBody *jen.Group) {
				writeStructMarshalInitDecls(methodBody, def.Fields, objReceiverName)
				methodBody.Type().Id(tmpAliasName).Id(def.Name)
				methodBody.Return(snip.SafeJSONMarshal().Call(jen.Id(tmpAliasName).Call(jen.Id(objReceiverName))))
			}))
			// Declare UnmarshalJSON
			file.Add(snip.MethodUnmarshalJSON(objReceiverName, def.Name).BlockFunc(func(methodBody *jen.Group) {
				rawVarName := "raw" + def.Name
				methodBody.Type().Id(tmpAliasName).Id(def.Name)
				methodBody.Var().Id(rawVarName).Id(tmpAliasName)
				methodBody.If(jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal()).Call(jen.Id(dataVarName), jen.Op("&").Id(rawVarName)),
					jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Err()),
				)
				writeStructMarshalInitDecls(methodBody, def.Fields, rawVarName)
				methodBody.Op("*").Id(objReceiverName).Op("=").Id(def.Name).Call(jen.Id(rawVarName))
				methodBody.Return(jen.Nil())
			}))
		}
	}

	if cfg.GenerateYAMLMethods {
		file.Add(snip.MethodMarshalYAML(objReceiverName, def.Name))
		file.Add(snip.MethodUnmarshalYAML(objReceiverName, def.Name))
	}
}

func writeStructMarshalInitDecls(methodBody *jen.Group, fields []*types.Field, rawVarName string) {
	for _, fieldDef := range fields {
		if collInit := fieldDef.Type.Make(); collInit != nil {
			// if there is a map or slice field, the struct contains a collection.
			fName := transforms.ExportedFieldName(fieldDef.Name)
			methodBody.If(jen.Id(rawVarName).Dot(fName).Op("==").Nil()).Block(
				jen.Id(rawVarName).Dot(fName).Op("=").Add(collInit),
			)
		}
	}
}

func astForStructLiteralMarshalJSON(def *types.ObjectType) *jen.Statement {
	return snip.MethodMarshalJSON(objReceiverName, def.Name).Block(
		encoding.MarshalJSONMethodBody(objReceiverName),
	)
}

func astForStructLiteralAppendJSON(def *types.ObjectType) *jen.Statement {
	return snip.MethodAppendJSON(objReceiverName, def.Name).BlockFunc(func(methodBody *jen.Group) {
		var fields []encoding.JSONStructField
		for _, field := range def.Fields {
			fields = append(fields, encoding.JSONStructField{
				Spec:     field,
				Selector: jen.Id(objReceiverName).Dot(transforms.ExportedFieldName(field.Name)).Clone,
			})
		}
		encoding.StructMethodBodyAppendJSON(methodBody, fields)
	})
}
