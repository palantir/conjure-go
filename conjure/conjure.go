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
	"strings"

	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

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

func createMappingFunctions(outputDir string) (conjurePkgToGoPkg, goPkgToFilePath func(string) string, rErr error) {
	outputDirAbsPath, err := filepath.Abs(outputDir)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to convert %s to absolute path", outputDir)
	}

	outputPkgBasePath, err := outputPackageBasePath(outputDirAbsPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to determine base import path for generated code")
	}
	conjurePkgToGoPkg = func(conjurePkg string) string {
		return path.Join(outputPkgBasePath, transforms.PackagePath(conjurePkg))
	}

	// transforms the provided goPkgPath to the absolute on-disk path where files for the package should be generated.
	// This abstraction is needed because the generated location may be different between $GOPATH projects and module
	// projects.
	goPkgToFilePath = func(goPkgPath string) string {
		return path.Join(outputDir, strings.TrimPrefix(goPkgPath, outputPkgBasePath+"/"))
	}
	return conjurePkgToGoPkg, goPkgToFilePath, nil
}

func GenerateOutputFiles(conjureDefinition spec.ConjureDefinition, outputDir string) ([]*OutputFile, error) {
	conjurePkgToGoPkg, goPkgToFilePath, err := createMappingFunctions(outputDir)
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
	enumFiles, err := collectEnumFiles(newConjureTypeFilterVisitor.EnumDefinitions, conjurePkgToGoPkg, goPkgToFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for enums")
	}
	files = append(files, enumFiles...)

	aliasFiles, err := collectAliasFiles(newConjureTypeFilterVisitor.AliasDefinitions, customTypes, conjurePkgToGoPkg, goPkgToFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for aliases")
	}
	files = append(files, aliasFiles...)

	objectFiles, err := collectObjectFiles(newConjureTypeFilterVisitor.ObjectDefinitions, customTypes, conjurePkgToGoPkg, goPkgToFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for objects")
	}
	files = append(files, objectFiles...)

	unionFiles, err := collectUnionFiles(newConjureTypeFilterVisitor.UnionDefinitions, customTypes, conjurePkgToGoPkg, goPkgToFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for unions")
	}
	files = append(files, unionFiles...)

	errorFiles, err := collectErrorFiles(conjureDefinition.Errors, customTypes, conjurePkgToGoPkg, goPkgToFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for errors")
	}
	files = append(files, errorFiles...)

	serviceFiles, err := collectServiceFiles(conjureDefinition.Services, customTypes, conjurePkgToGoPkg, goPkgToFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write files for services")
	}
	files = append(files, serviceFiles...)

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

func collectEnumFiles(enums []spec.EnumDefinition, conjurePkgToGoPk, goPkgToFilePath func(string) string) ([]*OutputFile, error) {
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
		file, err := newGoFile("enums", goPkgImportPath, goPkgToFilePath, importToAlias, goPkgToEnums[goPkgImportPath])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go enums for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectAliasFiles(aliasDefinitions []spec.AliasDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPk, goPkgToFilePath func(string) string) ([]*OutputFile, error) {
	var files []*OutputFile
	// group objects by Go package
	packageNameToAliases := make(map[string][]spec.AliasDefinition)
	for _, alias := range aliasDefinitions {
		goPkgName := conjurePkgToGoPk(alias.TypeName.Package)
		packageNameToAliases[goPkgName] = append(packageNameToAliases[goPkgName], alias)
	}
	for goPkgImportPath, aliasList := range packageNameToAliases {
		uniqueGoImports := make(map[string]struct{})
		for _, aliasDefinition := range aliasList {
			conjureTypeProvider, err := visitors.NewConjureTypeProvider(aliasDefinition.Alias)
			if err != nil {
				return nil, err
			}
			aliasTyper, err := conjureTypeProvider.ParseType(customTypes)
			if err != nil {
				return nil, errors.Wrapf(err, "alias type %s specifies unrecognized type", aliasDefinition.TypeName.Name)
			}
			for _, k := range aliasTyper.ImportPaths() {
				if k != goPkgImportPath {
					// if package required by type is not the current package, track as import
					uniqueGoImports[k] = struct{}{}
				}
			}
		}
		importToAlias := createAliasMap(uniqueGoImports)
		var aliasDefs []astgen.ASTDecl
		for _, alias := range aliasList {
			decls, imports, err := astForAlias(alias, customTypes, goPkgImportPath, importToAlias)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate AST for alias %s", alias.TypeName.Name)
			}
			for _, k := range imports.Sorted() {
				if _, ok := importToAlias[k]; !ok {
					importToAlias[k] = ""
				}
			}
			aliasDefs = append(aliasDefs, decls...)
		}

		file, err := newGoFile("aliases", goPkgImportPath, goPkgToFilePath, importToAlias, aliasDefs)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go aliases for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectObjectFiles(objects []spec.ObjectDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPk, goPkgToFilePath func(string) string) ([]*OutputFile, error) {
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
		file, err := newGoFile("structs", goPkgImportPath, goPkgToFilePath, importToAlias, objDefs)
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

func collectUnionFiles(unionDefinitions []spec.UnionDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPkg, goPkgToFilePath func(string) string) ([]*OutputFile, error) {
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

		file, err := newGoFile("unions", goPkgImportPath, goPkgToFilePath, importToAlias, unionDefs)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go unions for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectErrorFiles(errorDefinitions []spec.ErrorDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPkg, goPkgToFilePath func(string) string) ([]*OutputFile, error) {
	var files []*OutputFile
	// group errors by Go package
	packageNameToErrors := make(map[string][]spec.ErrorDefinition)
	for _, errorDefinition := range errorDefinitions {
		goPkgName := conjurePkgToGoPkg(errorDefinition.ErrorName.Package)
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
		file, err := newGoFile("errors", goPkgImportPath, goPkgToFilePath, importToAlias, decls)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create Go errors for %s", goPkgImportPath)
		}
		files = append(files, file)
	}
	return files, nil
}

func collectServiceFiles(services []spec.ServiceDefinition, customTypes types.CustomConjureTypes, conjurePkgToGoPkg, goPkgToFilePath func(string) string) ([]*OutputFile, error) {
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
		file, err := newGoFile("services", goPkgImportPath, goPkgToFilePath, importToAlias, decls)
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

func newGoFile(fileName, goImportPath string, goPkgToFilePath func(string) string, importsToAliases map[string]string, goTypeObjs []astgen.ASTDecl) (*OutputFile, error) {
	fileName += ".conjure.go"
	_, pkgName := path.Split(goImportPath)
	pkgDir := goPkgToFilePath(goImportPath)

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

// outputPackageBasePath returns the Go package path to the base output directory. For example, if the project is in
// package "github.com/org/project" and the output directory is the "outDir" directory within that package, the returned
// path is "github.com/org/project/outDir". Any conjure-generated package paths should be appended to this path.
func outputPackageBasePath(outputDirAbsPath string) (string, error) {
	// ensure that output directory exists, as "packages.Load" may require this
	if _, err := os.Stat(outputDirAbsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDirAbsPath, 0755); err != nil {
			return "", errors.Wrapf(err, "failed to create directory")
		}
	}

	pkgs, err := packages.Load(&packages.Config{
		Dir: outputDirAbsPath,
	}, "")
	if err != nil {
		return "", errors.Wrapf(err, "failed to load packages in %s", outputDirAbsPath)
	}
	if len(pkgs) == 0 {
		return "", errors.Errorf("could not determine package of %s", outputDirAbsPath)
	}
	return pkgs[0].PkgPath, nil
}
