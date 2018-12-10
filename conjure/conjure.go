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
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/transforms"
	"github.com/palantir/conjure-go/conjure/types"
	"github.com/palantir/conjure-go/conjure/visitors"
)

type Value struct {
	Value string
	Docs  string
}

type StringSet map[string]struct{}

func NewStringSet(vals ...string) StringSet {
	s := make(StringSet)
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

func (s StringSet) AddAll(other StringSet) {
	for k := range other {
		s[k] = struct{}{}
	}
}

func (s StringSet) Sorted() []string {
	var sorted []string
	for k := range s {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

func Generate(conjureDefinition spec.ConjureDefinition, outputDir string) error {
	files, err := GenerateOutputFiles(conjureDefinition, outputDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := file.Write(); err != nil {
			return err
		}
	}
	return nil
}

func createMappingFunction(outputDir string) (func(string) string, error) {
	outputDirGoPkgPather, err := toGoPkgPather(outputDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine absolute Go package path for output directory %s", outputDir)
	}
	outputDirGoPkgSrcRel, err := outputDirGoPkgPather.GoPathSrcRel()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine $GOPATH-relative Go package path for output directory %s", outputDir)
	}
	conjurePkgToGoPkg := func(conjurePkg string) string {
		return path.Join(outputDirGoPkgSrcRel, transforms.PackagePath(conjurePkg))
	}
	return conjurePkgToGoPkg, nil
}

func GenerateOutputFiles(conjureDefinition spec.ConjureDefinition, outputDir string) ([]*OutputFile, error) {
	conjurePkgToGoPkg, err := createMappingFunction(outputDir)
	if err != nil {
		return nil, err
	}
	customTypes, err := types.GetCustomConjureTypes(conjureDefinition.Types, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid configuration in types block")
	}
	newConjureTypeFilterVisitor := types.NewConjureTypeFilterVisitor()
	for _, typeDefinition := range conjureDefinition.Types {
		if err := typeDefinition.Accept(newConjureTypeFilterVisitor); err != nil {
			return nil, errors.Wrapf(err, "illegal recursive object type definition")
		}
	}
	var files []*OutputFile
	enumFiles, err := collectEnumFiles(newConjureTypeFilterVisitor.EnumDefinitions, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for enums")
	}
	files = append(files, enumFiles...)

	aliasFiles, err := collectAliasFiles(newConjureTypeFilterVisitor.AliasDefinitions, customTypes, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for aliases")
	}
	files = append(files, aliasFiles...)

	objectFiles, err := collectObjectFiles(newConjureTypeFilterVisitor.ObjectDefinitions, customTypes, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for objects")
	}
	files = append(files, objectFiles...)

	unionFiles, err := collectUnionFiles(newConjureTypeFilterVisitor.UnionDefinitions, customTypes, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for unions")
	}
	files = append(files, unionFiles...)

	errorFiles, err := collectErrorFiles(conjureDefinition.Errors, customTypes, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for errors")
	}
	files = append(files, errorFiles...)

	serviceFiles, err := collectServiceFiles(conjureDefinition.Services, customTypes, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for services")
	}
	files = append(files, serviceFiles...)

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

func collectEnumFiles(enums []spec.EnumDefinition, conjurePkgToGoPk func(string) string) ([]*OutputFile, error) {
	// group enums by Go package
	var sortedPkgNames []string
	goPkgToEnums := make(map[string][]astgen.ASTDecl)
	goPkgToImports := make(map[string]StringSet)

	for _, enumDefinition := range enums {
		goPkgName := conjurePkgToGoPk(enumDefinition.TypeName.Package)
		sortedPkgNames = append(sortedPkgNames, goPkgName)
		var values []Value
		for _, enumValueDefinition := range enumDefinition.Values {

			values = append(values, Value{
				Docs:  transforms.Documentation(enumValueDefinition.Docs),
				Value: enumValueDefinition.Value,
			})
		}

		enumDecl := &Enum{
			Name:    enumDefinition.TypeName.Name,
			Values:  values,
			Comment: transforms.Documentation(enumDefinition.Docs),
		}

		declers, imports := enumDecl.ASTDeclers()

		goPkgToEnums[goPkgName] = append(goPkgToEnums[goPkgName], declers...)

		if goPkgToImports[goPkgName] == nil {
			goPkgToImports[goPkgName] = NewStringSet()
		}
		goPkgToImports[goPkgName].AddAll(imports)
	}
	sort.Strings(sortedPkgNames)

	var files []*OutputFile
	for _, goPkgImportPath := range sortedPkgNames {
		importToAlias := createAliasMap(goPkgToImports[goPkgImportPath])
		file, err := newGoFile("enums", goPkgImportPath, importToAlias, goPkgToEnums[goPkgImportPath])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go enums for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectAliasFiles(aliasDefinitions []spec.AliasDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPk func(string) string) ([]*OutputFile, error) {
	// group aliases by Go package
	var sortedPkgNames []string
	goPkgToAliases := make(map[string][]astgen.ASTDecl)
	goPkgToUniqueImports := make(map[string]map[string]struct{})
	for _, aliasDefinition := range aliasDefinitions {
		goPkgName := conjurePkgToGoPk(aliasDefinition.TypeName.Package)
		sortedPkgNames = append(sortedPkgNames, goPkgName)
		conjureTypeProvider, err := visitors.NewConjureTypeProvider(aliasDefinition.Alias)
		if err != nil {
			return nil, err
		}
		aliasTyper, err := conjureTypeProvider.ParseType(customTypes)
		if err != nil {
			return nil, errors.Wrapf(err, "alias type %s specifies unrecognized type", aliasDefinition.TypeName.Name)
		}
		for _, importPath := range aliasTyper.ImportPaths() {
			if goPkgToUniqueImports[goPkgName] == nil {
				goPkgToUniqueImports[goPkgName] = make(map[string]struct{})
			}
			goPkgToUniqueImports[goPkgName][importPath] = struct{}{}
		}
		goPkgToAliases[goPkgName] = append(goPkgToAliases[goPkgName], &decl.Alias{
			Name:    aliasDefinition.TypeName.Name,
			Type:    expression.Type(aliasTyper.GoType(goPkgName, nil)),
			Comment: transforms.Documentation(aliasDefinition.Docs),
		})
	}
	sort.Strings(sortedPkgNames)

	var files []*OutputFile
	for _, goPkgImportPath := range sortedPkgNames {
		importToAlias := createAliasMap(goPkgToUniqueImports[goPkgImportPath])
		file, err := newGoFile("aliases", goPkgImportPath, importToAlias, goPkgToAliases[goPkgImportPath])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go aliases for %s", goPkgImportPath)
		}
		files = append(files, file)
	}

	return files, nil
}

func collectObjectFiles(objects []spec.ObjectDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPk func(string) string) ([]*OutputFile, error) {
	var files []*OutputFile
	// group objects by Go package
	packageNameToObjects := make(map[string][]spec.ObjectDefinition)
	for _, obj := range objects {
		goPkgName := conjurePkgToGoPk(obj.TypeName.Package)
		packageNameToObjects[goPkgName] = append(packageNameToObjects[goPkgName], obj)
	}
	for goPkgImportPath, objectList := range packageNameToObjects {

		uniqueGoImports := make(map[string]struct{})

		for _, object := range objectList {
			uniqueGoPkgs, err := getImportPathsFromFields(object.Fields, customTypes)
			if err != nil {
				return nil, err
			}
			for _, k := range uniqueGoPkgs.Sorted() {
				if k != goPkgImportPath {
					// if package required by type is not the current package, track as import
					uniqueGoImports[k] = struct{}{}
				}
			}
		}
		importToAlias := createAliasMap(uniqueGoImports)
		var objDefs []astgen.ASTDecl
		for _, object := range objectList {
			decl, imports, err := astForObject(object, customTypes, goPkgImportPath, importToAlias)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate AST for object %s", object.TypeName.Name)
			}
			for _, k := range imports.Sorted() {
				if _, ok := importToAlias[k]; !ok {
					importToAlias[k] = ""
				}
			}
			objDefs = append(objDefs, decl...)
		}
		file, err := newGoFile("structs", goPkgImportPath, importToAlias, objDefs)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go objects for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func getImportPathsFromFields(fields []spec.FieldDefinition, customTypes types.CustomConjureTypes) (StringSet, error) {
	uniqueGoPkgs := NewStringSet()
	for _, field := range fields {
		typer, err := getTyperFromType(field.Type, customTypes)
		if err != nil {
			return nil, err
		}
		for _, importPath := range typer.ImportPaths() {
			uniqueGoPkgs[importPath] = struct{}{}
		}
	}
	return uniqueGoPkgs, nil
}

func getTyperFromType(specType spec.Type, customTypes types.CustomConjureTypes) (types.Typer, error) {
	conjureTypeProvider, err := visitors.NewConjureTypeProvider(specType)
	if err != nil {
		return nil, err
	}
	return conjureTypeProvider.ParseType(customTypes)
}

func collectUnionFiles(unionDefinitions []spec.UnionDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPkg func(string) string) ([]*OutputFile, error) {
	var files []*OutputFile
	// group objects by Go package
	packageNameToObjects := make(map[string][]spec.UnionDefinition)
	for _, unionDefinition := range unionDefinitions {
		goPkgName := conjurePkgToGoPkg(unionDefinition.TypeName.Package)
		packageNameToObjects[goPkgName] = append(packageNameToObjects[goPkgName], unionDefinition)
	}

	for goPkgImportPath, unionDefinitionList := range packageNameToObjects {

		uniqueGoImports := NewStringSet()
		for _, unionDefinition := range unionDefinitionList {
			uniqueGoPkgs, err := getImportPathsFromFields(unionDefinition.Union, customTypes)
			if err != nil {
				return nil, err
			}
			for _, k := range uniqueGoPkgs.Sorted() {
				if k != goPkgImportPath {
					// if package required by type is not the current package, track as import
					uniqueGoImports[k] = struct{}{}
				}
			}
		}
		importToAlias := createAliasMap(uniqueGoImports)
		var unionDefs []astgen.ASTDecl
		for _, unionDefinition := range unionDefinitionList {
			declers, imports, err := astForUnion(unionDefinition, customTypes, goPkgImportPath, importToAlias)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate AST for union type %q", unionDefinition.TypeName.Name)
			}
			for _, k := range imports.Sorted() {
				if _, ok := importToAlias[k]; !ok {
					importToAlias[k] = ""
				}
			}
			unionDefs = append(unionDefs, declers...)
		}

		file, err := newGoFile("unions", goPkgImportPath, importToAlias, unionDefs)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go unions for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectErrorFiles(errorDefinitions []spec.ErrorDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPk func(string) string) ([]*OutputFile, error) {
	var files []*OutputFile
	// group errors by Go package
	packageNameToErrors := make(map[string][]spec.ErrorDefinition)
	for _, errorDefinition := range errorDefinitions {
		goPkgName := conjurePkgToGoPk(errorDefinition.ErrorName.Package)
		packageNameToErrors[goPkgName] = append(packageNameToErrors[goPkgName], errorDefinition)
	}
	for goPkgImportPath, errorList := range packageNameToErrors {

		uniqueGoImports := make(map[string]struct{})

		for _, errorDefinition := range errorList {
			allArgs := make([]spec.FieldDefinition, 0, len(errorDefinition.SafeArgs)+len(errorDefinition.UnsafeArgs))
			allArgs = append(allArgs, errorDefinition.SafeArgs...)
			allArgs = append(allArgs, errorDefinition.UnsafeArgs...)
			uniqueGoPkgs, err := getImportPathsFromFields(allArgs, customTypes)
			if err != nil {
				return nil, err
			}
			for _, k := range uniqueGoPkgs.Sorted() {
				if k != goPkgImportPath {
					// if package required by type is not the current package, track as import
					uniqueGoImports[k] = struct{}{}
				}
			}
		}
		importToAlias := createAliasMap(uniqueGoImports)
		var decls []astgen.ASTDecl
		for _, errorDefinition := range errorList {
			decl, imports, err := astForError(errorDefinition, customTypes, goPkgImportPath, importToAlias)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate AST for error %s", errorDefinition.ErrorName.Name)
			}
			for _, k := range imports.Sorted() {
				if _, ok := importToAlias[k]; !ok {
					importToAlias[k] = ""
				}
			}
			decls = append(decls, decl...)
		}
		file, err := newGoFile("errors", goPkgImportPath, importToAlias, decls)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go errors for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectServiceFiles(services []spec.ServiceDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPkg func(string) string) ([]*OutputFile, error) {
	var files []*OutputFile
	pkgToServiceDefinitions := make(map[string][]spec.ServiceDefinition)
	for _, serviceDefinition := range services {
		goPkgName := conjurePkgToGoPkg(serviceDefinition.ServiceName.Package)
		pkgToServiceDefinitions[goPkgName] = append(pkgToServiceDefinitions[goPkgName], serviceDefinition)
	}
	for goPkgImportPath, serviceDefinitionList := range pkgToServiceDefinitions {
		uniqueGoImports := NewStringSet()
		for _, serviceDefinition := range serviceDefinitionList {
			for _, endpointDefinition := range serviceDefinition.Endpoints {
				for _, endpointArg := range endpointDefinition.Args {
					typer, err := getTyperFromType(endpointArg.Type, customTypes)
					if err != nil {
						return nil, err
					}

					for _, importPath := range typer.ImportPaths() {
						if importPath != goPkgImportPath {
							uniqueGoImports[importPath] = struct{}{}
						}
					}
				}
				if endpointDefinition.Returns != nil {
					typer, err := getTyperFromType(*endpointDefinition.Returns, customTypes)
					if err != nil {
						return nil, err
					}
					for _, importPath := range typer.ImportPaths() {
						if importPath != goPkgImportPath {
							uniqueGoImports[importPath] = struct{}{}
						}
					}
				}
			}
		}
		importToAlias := createAliasMap(uniqueGoImports)

		var decls []astgen.ASTDecl
		for _, serviceDefinition := range serviceDefinitionList {
			declers, imports, err := astForService(serviceDefinition, customTypes, goPkgImportPath, importToAlias)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate AST for service %s", serviceDefinition.ServiceName.Name)
			}
			for _, k := range imports.Sorted() {
				if _, ok := importToAlias[k]; !ok {
					importToAlias[k] = ""
				}
			}
			decls = append(decls, declers...)
		}
		file, err := newGoFile("services", goPkgImportPath, importToAlias, decls)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go services for %s", goPkgImportPath)
		}
		files = append(files, file)

	}
	return files, nil
}

func createAliasMap(goImports StringSet) map[string]string {
	importToAlias := make(map[string]string)
	used := make(map[string]struct{})

	for _, currImportPath := range goImports.Sorted() {
		_, pkgName := path.Split(currImportPath)
		if _, ok := used[pkgName]; !ok {
			// package name has not been used yet -- no need for alias
			used[pkgName] = struct{}{}
			importToAlias[currImportPath] = ""
			continue
		}

		// package name has been used before -- need to find a unique alias and record
		currIdx := 1
		// append number to package name to make it unique. Increment counter until unique identifier is found.
		for {
			if _, ok := used[pkgName]; !ok {
				// package name is available
				break
			}
			pkgName = fmt.Sprintf("%s_%d", pkgName, currIdx)
		}

		// add entry to alias map
		used[pkgName] = struct{}{}
		importToAlias[currImportPath] = pkgName
	}

	return importToAlias
}

func newGoFile(fileName, goImportPath string, importsToAliases map[string]string, goTypeObjs []astgen.ASTDecl) (*OutputFile, error) {
	fileName += ".conjure.go"
	_, pkgName := path.Split(goImportPath)

	pkgDir := pkgpath.NewGoPathSrcRelPkgPath(goImportPath).Abs()
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return nil, errors.Wrapf(err, "failed to create parent directory for Go file output")
	}

	var components []astgen.ASTDecl
	if len(importsToAliases) > 0 {
		components = append(components, decl.NewImports(importsToAliases))
	}
	components = append(components, goTypeObjs...)

	file := OutputFile{
		pkgName:    pkgName,
		absPath:    path.Join(pkgDir, fileName),
		goTypeObjs: components,
	}

	return &file, nil
}

func toGoPkgPather(dir string) (pkgpath.PkgPather, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert %s to absolute path", dir)
	}
	return pkgpath.NewAbsPkgPath(absPath), nil
}
