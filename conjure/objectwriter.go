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

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
	"github.com/palantir/conjure-go/v6/conjure/visitors/jsonencoding"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/pkg/errors"
)

func astForObject(objectDefinition spec.ObjectDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	if err := addImportPathsFromFields(objectDefinition.Fields, info); err != nil {
		return nil, err
	}

	var structFields []*expression.StructField

	for _, fieldDefinition := range objectDefinition.Fields {
		conjureTypeProvider, err := visitors.NewConjureTypeProvider(fieldDefinition.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create type provider for field %s for object %s",
				fieldDefinition.FieldName,
				objectDefinition.TypeName.Name)
		}
		typer, err := conjureTypeProvider.ParseType(info)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse type field %s for object %s",
				fieldDefinition.FieldName,
				objectDefinition.TypeName.Name)
		}
		info.AddImports(typer.ImportPaths()...)

		structFields = append(structFields, &expression.StructField{
			Name:    transforms.ExportedFieldName(string(fieldDefinition.FieldName)),
			Type:    expression.Type(typer.GoType(info)),
			Tag:     fmt.Sprintf("json:%q", fieldDefinition.FieldName),
			Comment: transforms.Documentation(fieldDefinition.Docs),
		})
	}

	decls := []astgen.ASTDecl{
		decl.NewStruct(objectDefinition.TypeName.Name, structFields, transforms.Documentation(objectDefinition.Docs)),
	}

	var fields []jsonencoding.JSONField
	for _, field := range objectDefinition.Fields {
		fields = append(fields, jsonencoding.JSONField{
			FieldSelector: transforms.ExportedFieldName(string(field.FieldName)),
			JSONKey:       string(field.FieldName),
			ValueType:     field.Type,
		})
	}
	jsonMethods, err := jsonencoding.StructFieldsJSONMethods(objReceiverName, objectDefinition.TypeName.Name, fields, info)
	if err != nil {
		return nil, err
	}
	decls = append(decls, jsonMethods...)
	decls = append(decls, newMarshalYAMLMethod(objReceiverName, objectDefinition.TypeName.Name, info))
	decls = append(decls, newUnmarshalYAMLMethod(objReceiverName, objectDefinition.TypeName.Name, info))

	return decls, nil
}

const (
	objReceiverName = "o"
)
