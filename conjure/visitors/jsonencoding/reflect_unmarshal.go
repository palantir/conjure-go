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

func reflectAliasTypeUnmarshalMethods(
	receiverName string,
	receiverType string,
	aliasType spec.Type,
	info types.PkgInfo,
) ([]astgen.ASTDecl, error) {
	panic("implement me")
}

func reflectStructFieldsUnmarshalMethods(
	receiverName string,
	receiverType string,
	fields []JSONField,
	info types.PkgInfo,
) ([]astgen.ASTDecl, error) {
	panic("implement me")
}
