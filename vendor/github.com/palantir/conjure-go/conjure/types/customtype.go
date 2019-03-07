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

package types

import (
	"fmt"
	"strings"
)

type CustomConjureType struct {
	Name string
	// Pkg is the Go package that needs to be imported to use this type from an external package. Empty value
	// indicates that no import is needed (only true in cases where the Typer returns a Go primitive or built-in).
	Pkg string
	Typer
}

type CustomConjureTypes interface {
	Add(name, pkg string, typer Typer) error
	Get(name string) (CustomConjureType, bool)
}

type customConjureTypes struct {
	typeDecls map[string]CustomConjureType
}

func NewCustomConjureTypes() CustomConjureTypes {
	return &customConjureTypes{
		typeDecls: make(map[string]CustomConjureType),
	}
}

func (t *customConjureTypes) Add(name, pkg string, typer Typer) error {
	// normalize for purpose of comparison and keying, but use value provided by user as "Name" field for type and
	// for errors displayed to the user.
	lowercaseName := strings.ToLower(name)

	if existing, ok := t.typeDecls[lowercaseName]; ok {
		return fmt.Errorf("%q has already been defined as a custom Conjure type", existing.Name)
	}
	t.typeDecls[lowercaseName] = CustomConjureType{
		Name:  name,
		Pkg:   pkg,
		Typer: typer,
	}
	return nil
}

func (t *customConjureTypes) Get(name string) (CustomConjureType, bool) {
	typer, ok := t.typeDecls[strings.ToLower(name)]
	return typer, ok
}
