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

//go:build ignore
// +build ignore

package old

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
	Type          spec.Type
}

func TypeJSONMethods(receiverName string, def spec.TypeDefinition, info types.PkgInfo) ([]astgen.ASTDecl, error) {
	if !EnableDirectJSONMethods {
		return reflectJSONMethods(receiverName, def, info)
	}
	return literalJSONMethods(receiverName, def, info)
}
