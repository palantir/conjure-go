// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package jsonencoding

import (
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/goastwriter/astgen"
)

var EnableDirectJSONMethods = true

type JSONField struct {
	// FieldSelector is the name of the Go field in the struct.
	FieldSelector string
	JSONKey       string
	ValueType     spec.Type
}

func AliasTypeJSONMethods(
	receiverName string,
	receiverType string,
	aliasType spec.Type,
	info types.PkgInfo,
) ([]astgen.ASTDecl, error) {
	var marshal, unmarshal []astgen.ASTDecl
	var err error
	if EnableDirectJSONMethods {
		marshal, err = literalAliasTypeMarshalMethods(receiverName, receiverType, aliasType, info)
		if err != nil {
			return nil, err
		}
		unmarshal, err = literalAliasTypeUnmarshalMethods(receiverName, receiverType, aliasType, info)
		if err != nil {
			return nil, err
		}
	} else {
		marshal, err = reflectAliasTypeMarshalMethods(receiverName, receiverType, aliasType, info)
		if err != nil {
			return nil, err
		}
		unmarshal, err = reflectAliasTypeUnmarshalMethods(receiverName, receiverType, aliasType, info)
		if err != nil {
			return nil, err
		}
	}
	return append(marshal, unmarshal...), nil
}

func StructFieldsJSONMethods(
	receiverName string,
	receiverType string,
	fields []JSONField,
	info types.PkgInfo,
) ([]astgen.ASTDecl, error) {
	var marshal, unmarshal []astgen.ASTDecl
	var err error
	if EnableDirectJSONMethods {
		marshal, err = literalStructFieldsMarshalMethods(receiverName, receiverType, fields, info)
		if err != nil {
			return nil, err
		}
		unmarshal, err = literalStructFieldsUnmarshalMethods(receiverName, receiverType, fields, info)
		if err != nil {
			return nil, err
		}
	} else {
		marshal, err = reflectStructFieldsMarshalMethods(receiverName, receiverType, fields, info)
		if err != nil {
			return nil, err
		}
		unmarshal, err = reflectStructFieldsUnmarshalMethods(receiverName, receiverType, fields, info)
		if err != nil {
			return nil, err
		}
	}
	return append(marshal, unmarshal...), nil
}
