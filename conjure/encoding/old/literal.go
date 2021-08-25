//go:build ignore
// +build ignore

package old

import (
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
)

func literalJSONMethods(receiverName string, receiverType string, def spec.TypeDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	info.AddImports(
		"context",
		"encoding/base64",
		"math",
		"strconv",
		"github.com/palantir/pkg/safejson",
		"github.com/tidwall/gjson")
	info.SetImports("wparams", "github.com/palantir/witchcraft-go-params")
	info.SetImports("werror", "github.com/palantir/witchcraft-go-error")

	var decls []astgen.ASTDecl
	if err := def.AcceptFuncs(
		func(def spec.AliasDefinition) error {
			marshalBody, err := visitAliasMarshalGJSONMethodBody(receiverName, def.Alias, info)
			if err != nil {
				return err
			}
			decls = publicMarshalJSONMethods(receiverName, receiverType, marshalBody)

			decls = append(decls, publicUnmarshalJSONMethods(receiverName, receiverType)...)
			unmarshalBody, err := visitAliasUnmarshalGJSONMethodBody(receiverName, receiverType, def.Alias, info)
			if err != nil {
				return err
			}
			decls = append(decls, unmarshalGJSONMethod(receiverName, receiverType, unmarshalBody))
			return nil
		},
		func(def spec.EnumDefinition) error {
			decls = publicMarshalJSONMethods(receiverName, receiverType, []astgen.ASTStmt{
				appendMarshalBufferQuotedString(expression.NewCallExpression(expression.StringType, expression.VariableVal(receiverName))),
				statement.NewReturn(marshalBufVar, expression.Nil),
			})
			panic("add unmarshal methods!")
			return nil
		},
		func(def spec.ObjectDefinition) error {
			fields := make([]JSONField, len(def.Fields))
			for i, field := range def.Fields {
				fields[i] = JSONField{
					FieldSelector: transforms.ExportedFieldName(string(field.FieldName)),
					JSONKey:       string(field.FieldName),
					Type:          field.Type,
				}
			}
			marshalBody, err := visitStructFieldsMarshalGJSONMethodBody(receiverName, receiverType, fields, info)
			if err != nil {
				return err
			}
			decls = publicMarshalJSONMethods(receiverName, receiverType, marshalBody)

			decls = append(decls, publicUnmarshalJSONMethods(receiverName, receiverType)...)
			unmarshalBody, err := visitStructFieldsUnmarshalGJSONMethodBody(receiverName, receiverType, fields, info)
			if err != nil {
				return err
			}
			decls = append(decls, unmarshalGJSONMethod(receiverName, receiverType, unmarshalBody))
			return nil
		},
		func(def spec.UnionDefinition) error {
			fields := []JSONField{{
				FieldSelector: "typ",
				JSONKey:       "type",
				Type:          spec.NewTypeFromPrimitive(spec.New_PrimitiveType(spec.PrimitiveType_STRING)),
			}}
			for _, field := range def.Union {
				fields = append(fields, JSONField{
					FieldSelector: transforms.PrivateFieldName(string(field.FieldName)),
					JSONKey:       string(field.FieldName),
					Type:          spec.NewTypeFromOptional(spec.OptionalType{ItemType: field.Type}),
				})
			}
			marshalBody, err := visitStructFieldsMarshalGJSONMethodBody(receiverName, receiverType, fields, info)
			if err != nil {
				return err
			}
			decls = publicMarshalJSONMethods(receiverName, receiverType, marshalBody)

			decls = append(decls, publicUnmarshalJSONMethods(receiverName, receiverType)...)
			unmarshalBody, err := visitStructFieldsUnmarshalGJSONMethodBody(receiverName, receiverType, fields, info)
			if err != nil {
				return err
			}
			decls = append(decls, unmarshalGJSONMethod(receiverName, receiverType, unmarshalBody))
			return nil
		},
		def.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return decls, nil
}
