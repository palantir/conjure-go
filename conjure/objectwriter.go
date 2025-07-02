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
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
)

const (
	objReceiverName = "o"
	dataVarName     = "data"
)

// logSafetyToAnnotation converts a LogSafety value to its corresponding safelogging annotation string
func logSafetyToAnnotation(safety spec.LogSafety_Value) string {
	switch safety {
	case spec.LogSafety_SAFE:
		return "@Safe"
	case spec.LogSafety_UNSAFE:
		return "@Unsafe"
	case spec.LogSafety_DO_NOT_LOG:
		return "@DoNotLog"
	default:
		panic("unhandled LogSafety value: " + safety)
	}
}

func writeObjectType(file *jen.Group, objectDef *types.ObjectType, safetyCache map[types.Type]spec.LogSafety) {
	// Declare struct type with fields
	containsCollection := false // If contains collection, we need JSON methods to initialize empty values.

	// Compute the overall safety of the struct and add comment-based annotation
	overallSafety := computeObjectSafety(objectDef, safetyCache)
	var safetyComment jen.Code = jen.Null()
	if !overallSafety.IsUnknown() {
		safetyComment = jen.Comment("safelogging:" + logSafetyToAnnotation(overallSafety.Value()))
	}

	structFunc := func(structDecl *jen.Group) {
		for _, fieldDef := range objectDef.Fields {
			jsonTag := fieldDef.Name
			if _, isOptional := fieldDef.Type.(*types.Optional); isOptional {
				jsonTag += ",omitempty"
			}
			fieldTags := map[string]string{"json": jsonTag}

			if fieldDef.Docs != "" {
				// backtick characters ("`") are really painful to deal with in struct tags
				// (which are themselves defined within backtick literals), so replace with
				// double quotes instead.
				fieldTags["conjure-docs"] = strings.Replace(strings.TrimSpace(string(fieldDef.Docs)), "`", `"`, -1)
			}

			// Add safety struct tag based on field's safety annotation or type safety (with recursive computation)
			fieldSafety := getSafetyFromField(fieldDef, safetyCache)
			if !fieldSafety.IsUnknown() {
				fieldTags["safelogging"] = logSafetyToAnnotation(fieldSafety.Value())
			}
			if fieldDef.Type.Make() != nil {
				containsCollection = true
			}
			structDecl.Add(fieldDef.Docs.CommentLineWithDeprecation(fieldDef.Deprecated)).Id(transforms.ExportedFieldName(fieldDef.Name)).Add(fieldDef.Type.Code()).Tag(fieldTags)
		}
	}

	// If there are docs, add them without the trailing line since the safety comment will alreaady add the trailing line
	if objectDef.Docs != "" {
		file.Comment(string(objectDef.Docs))
	}
	file.Add(safetyComment)
	file.Type().Id(objectDef.Name).StructFunc(structFunc)

	// If there are no collections, we can defer to the default json behavior
	// Otherwise we need to override MarshalJSON and UnmarshalJSON
	if containsCollection {
		// We use this prefix to ensure that the resulting type alias does not conflict with any of the types in the object's fields, which will always be exported.
		tmpAliasName := "_tmp" + objectDef.Name
		// Declare MarshalJSON
		file.Add(snip.MethodMarshalJSON(objReceiverName, objectDef.Name).BlockFunc(func(methodBody *jen.Group) {
			writeStructMarshalInitDecls(methodBody, objectDef.Fields, objReceiverName)
			methodBody.Type().Id(tmpAliasName).Id(objectDef.Name)
			methodBody.Return(snip.SafeJSONMarshal().Call(jen.Id(tmpAliasName).Call(jen.Id(objReceiverName))))
		}))
		// Declare UnmarshalJSON
		file.Add(snip.MethodUnmarshalJSON(objReceiverName, objectDef.Name).BlockFunc(func(methodBody *jen.Group) {
			rawVarName := "raw" + objectDef.Name
			methodBody.Type().Id(tmpAliasName).Id(objectDef.Name)
			methodBody.Var().Id(rawVarName).Id(tmpAliasName)
			methodBody.If(jen.Err().Op(":=").Add(snip.SafeJSONUnmarshal()).Call(jen.Id(dataVarName), jen.Op("&").Id(rawVarName)),
				jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			)
			writeStructMarshalInitDecls(methodBody, objectDef.Fields, rawVarName)
			methodBody.Op("*").Id(objReceiverName).Op("=").Id(objectDef.Name).Call(jen.Id(rawVarName))
			methodBody.Return(jen.Nil())
		}))
	}

	file.Add(snip.MethodMarshalYAML(objReceiverName, objectDef.Name))
	file.Add(snip.MethodUnmarshalYAML(objReceiverName, objectDef.Name))
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

func getSafetyFromType(fieldType types.Type, safetyCache map[types.Type]spec.LogSafety) spec.LogSafety {
	if safetyCache == nil {
		safetyCache = make(map[types.Type]spec.LogSafety)
	}

	// Check cache first
	if cached, exists := safetyCache[fieldType]; exists {
		return cached
	}

	// Insert placeholder to detect cycles - use UNKNOWN as safe default
	safetyCache[fieldType] = spec.New_LogSafety(spec.LogSafety_UNKNOWN)

	// First check if the field type itself has safety annotations
	fieldSafety := fieldType.Safety()
	if !fieldSafety.IsUnknown() {
		// Update cache with real result
		safetyCache[fieldType] = fieldSafety
		return fieldSafety
	}

	// If it's an ObjectType, recursively compute safety from its fields
	if objectType, ok := fieldType.(*types.ObjectType); ok {
		result := computeObjectSafety(objectType, safetyCache)
		// Update cache with computed result
		safetyCache[fieldType] = result
		return result
	}

	// For other types without explicit safety, return unknown so no tags are added
	// Cache already contains UNKNOWN, so just return it
	return spec.New_LogSafety(spec.LogSafety_UNKNOWN)
}

// getSafetyFromField returns the safety of a field, considering both field-level and type-level safety annotations
func getSafetyFromField(field *types.Field, safetyCache map[types.Type]spec.LogSafety) spec.LogSafety {
	// First check if the field has an explicit safety annotation
	if field.Safety != nil {
		return *field.Safety
	}
	// Fall back to type-based safety
	return getSafetyFromType(field.Type, safetyCache)
}

func computeObjectSafety(object *types.ObjectType, typeSafetyCache map[types.Type]spec.LogSafety) spec.LogSafety {
	// Empty struct has no information to determine safety, so it should be unknown
	if len(object.Fields) == 0 {
		return spec.New_LogSafety(spec.LogSafety_UNKNOWN)
	}

	// Default to SAFE, then find the most restrictive level
	// Hierarchy: safe -> unannotated -> unsafe -> do-not-log
	overallSafety := spec.New_LogSafety(spec.LogSafety_SAFE)

	for _, fieldDef := range object.Fields {
		fieldSafety := getSafetyFromField(fieldDef, typeSafetyCache)

		if fieldSafety.IsUnknown() {
			// Field for which safety cannot be determined "contaminates" the struct - it becomes unannotated
			// unless we already found something more restrictive
			if overallSafety.Value() == spec.LogSafety_SAFE {
				overallSafety = spec.New_LogSafety(spec.LogSafety_UNKNOWN)
			}
		} else {
			// Specific safety value determined for field
			switch fieldSafety.Value() {
			case spec.LogSafety_DO_NOT_LOG:
				return fieldSafety // Most restrictive, return immediately
			case spec.LogSafety_UNSAFE:
				if overallSafety.Value() == spec.LogSafety_SAFE || overallSafety.IsUnknown() {
					overallSafety = fieldSafety
				}
			case spec.LogSafety_SAFE:
				// SAFE is least restrictive, only set if we haven't found anything else
				// (overallSafety is already SAFE by default)
			}
		}
	}

	return overallSafety
}
