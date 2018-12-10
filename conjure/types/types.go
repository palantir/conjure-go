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
	"path"
	"strings"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
)

type Typer interface {
	// GoType returns the string that can be used as the Go type declaration for this type. currPkg and pkgAliases
	// are used to determine the import paths that should be used to qualify the usage of type names. If a type
	// occurs in a package that matches currPkg, the package will not be referenced. Otherwise, if the package path
	// matches a key in the pkgAliases map, the package name of the value (which is the last portion of its import
	// path) will be used.
	//
	// As an example, if type "t" is a type alias named "ExampleAlias" defined in "github.com/palantir/generated/alias":
	//
	// * t.GoType("github.com/palantir/generated/alias", nil) -> ExampleAlias
	// * t.GoType("github.com/project", nil) -> alias.ExampleAlias
	// * t.GoType("github.com/project", map[string]string{"github.com/palantir/generated/alias": "pkgalias" }) -> pkgalias.ExampleAlias
	GoType(currPkgPath string, pkgAliases map[string]string) string

	// ImportPath returns the strings that can be used as the Go import path for this type. Returns an empty string
	// if the type is a primitive and does not require an import. We must return a list for all collection types
	ImportPaths() []string
}

const (
	conjureTypesImport = "github.com/palantir/conjure-go/conjure/types/conjuretype"
	uuidTypeImport     = "github.com/google/uuid"
)

var (
	String Typer = &simpleType{
		goType: "string",
	}
	Integer Typer = &simpleType{
		goType: "int",
	}
	Double Typer = &simpleType{
		goType: "float64",
	}
	Boolean Typer = &simpleType{
		goType: "bool",
	}
	BinaryType Typer = &simpleType{
		goType: "[]byte",
	}
	SafeLongType Typer = &goType{
		name:       "SafeLong",
		importPath: conjureTypesImport,
	}
	Rid Typer = &goType{
		name:       "Rid",
		importPath: conjureTypesImport,
	}
	Bearertoken Typer = &goType{
		name:       "Bearertoken",
		importPath: conjureTypesImport,
	}
	DateTimeType Typer = &goType{
		name:       "DateTime",
		importPath: conjureTypesImport,
	}
	Any Typer = &simpleType{
		goType: "interface{}",
	}
	UUIDType Typer = &goType{
		name:       "UUID",
		importPath: uuidTypeImport,
	}
)

type simpleType struct {
	goType string
}

func (t *simpleType) GoType(currPkgPath string, pkgAliases map[string]string) string {
	return t.goType
}

func (t *simpleType) ImportPaths() []string {
	return nil
}

type mapType struct {
	conjureType string
	keyType     Typer
	valType     Typer
}

func NewMapType(keyType, valType Typer) Typer {
	return &mapType{
		keyType: keyType,
		valType: valType,
	}
}

func (t *mapType) GoType(currPkgPath string, pkgAliases map[string]string) string {
	return fmt.Sprintf("map[%s]%s", t.keyType.GoType(currPkgPath, pkgAliases), t.valType.GoType(currPkgPath, pkgAliases))
}

func (t *mapType) ImportPaths() []string {
	return append(t.keyType.ImportPaths(), t.valType.ImportPaths()...)
}

type singleGenericValType struct {
	valType   Typer
	fmtString string
}

func (t *singleGenericValType) GoType(currPkgPath string, pkgAliases map[string]string) string {
	return fmt.Sprintf(t.fmtString, t.valType.GoType(currPkgPath, pkgAliases))
}

func (t *singleGenericValType) ImportPaths() []string {
	return t.valType.ImportPaths()
}

func NewListType(valType Typer) Typer {
	return &singleGenericValType{
		valType:   valType,
		fmtString: "[]%s",
	}
}

// NewSetType creates a new Typer for a set type.
//
// TODO: currently, sets and lists are treated identically. If we want to be more semantically precise, then the proper
// approach would be to define a Set as a map type with the provided key type and an empty struct as the value type.
// Because Go doesn't support generics, this would require generating a different Set type ("IntSet", "TestTypeSet", etc.)
// for each different set type that is required. This type would also need to implement custom JSON serialization and
// deserialization to translate to a JSON list, since that's the underlying representation required by the spec.
func NewSetType(valType Typer) Typer {
	return NewListType(valType)
}

func NewOptionalType(valType Typer) Typer {
	return &singleGenericValType{
		valType:   valType,
		fmtString: "*%s",
	}
}

func NewGoType(name, importPath string) Typer {
	return &goType{
		name:       name,
		importPath: importPath,
	}
}

func NewGoTypeFromExternalType(externalType spec.ExternalReference) Typer {
	pathAndName := strings.Split(externalType.ExternalReference.Name, ":")
	return &goType{
		name:       pathAndName[1],
		importPath: externalType.ExternalReference.Package + "." + pathAndName[0],
	}
}

// goType represents a type that is defined in a Go package.
type goType struct {
	name string
	// full import path to the type (including package)
	importPath string
}

func (t *goType) GoType(currPkgPath string, pkgAliases map[string]string) string {
	// if name is fully qualified, only use the last component
	name := t.name
	if lastDotIdx := strings.LastIndex(name, "."); lastDotIdx != -1 {
		name = name[lastDotIdx+1:]
	}

	if currPkgPath == t.importPath {
		// if current package is the same as the import path, no need to qualify type
		return name
	}

	// start package name as final component of the import path
	_, pkgName := path.Split(t.importPath)
	if alias := pkgAliases[t.importPath]; alias != "" {
		// if non-empty alias exists for full import path, use that instead
		pkgName = alias
	}
	return fmt.Sprintf("%s.%s", pkgName, name)
}

func (t *goType) ImportPaths() []string {
	return []string{t.importPath}
}
