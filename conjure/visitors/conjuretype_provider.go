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

package visitors

import (
	"github.com/palantir/goastwriter/expression"

	"github.com/palantir/conjure-go/v5/conjure/types"
)

type TypeCheck string

const (
	IsText     TypeCheck = "TEXT" // anything serialized as a string
	IsOptional TypeCheck = "OPTIONAL"
	IsBinary   TypeCheck = "BINARY"
	IsString   TypeCheck = "STRING"
	IsList     TypeCheck = "LIST"
	IsMap      TypeCheck = "MAP"
	IsSet      TypeCheck = "SET"
	IsBoolean  TypeCheck = "BOOLEAN"
)

type ConjureTypeProvider interface {
	ParseType(info types.PkgInfo) (types.Typer, error)
	CollectionInitializationIfNeeded(info types.PkgInfo) (*expression.CallExpression, error)
	IsSpecificType(typeCheck TypeCheck) bool
}
