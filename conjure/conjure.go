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

func (s StringSet) Add(vals ...string) {
	for _, v := range vals {
		s[v] = struct{}{}
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
	customTypes, err := visitors.GetCustomConjureTypes(conjureDefinition.Types, conjurePkgToGoPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid configuration in types block")
	}
	outputByPackage, err := visitors.ConjureDefinitionsByPackage(conjureDefinition)
	if err != nil {
		return nil, err
	}

	var files []*OutputFile
	for packageName, conjureDef := range outputByPackage {
		importPath := conjurePkgToGoPkg(packageName)
		goPkgDir := goPkgToFilePath(importPath)
		collector := &outputFileCollector{
			aliases:  fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			enums:    fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			objects:  fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			unions:   fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			errors:   fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			services: fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
		}
		if err := visitors.VisitConjureDefinition(conjureDef, collector); err != nil {
			return nil, err
		}

		for filename, ast := range map[string]fileASTCollector{
			"aliases.conjure.go":  collector.aliases,
			"enums.conjure.go":    collector.enums,
			"structs.conjure.go":  collector.objects,
			"unions.conjure.go":   collector.unions,
			"errors.conjure.go":   collector.errors,
			"services.conjure.go": collector.services,
		} {
			if len(ast.Decls) == 0 {
				continue
			}
			file, err := newGoFile(path.Join(goPkgDir, filename), importPath, ast.Info, ast.Decls)
			if err != nil {
				return nil, err
			}
			files = append(files, file)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

type outputFileCollector struct {
	// Track outputs (imports and decls) per-file
	aliases  fileASTCollector
	enums    fileASTCollector
	objects  fileASTCollector
	unions   fileASTCollector
	errors   fileASTCollector
	services fileASTCollector
}

type fileASTCollector struct {
	Info  types.PkgInfo
	Decls []astgen.ASTDecl
}

func (c *outputFileCollector) VisitAlias(aliasDefinition spec.AliasDefinition) error {
	info := c.aliases.Info
	conjureTypeProvider, err := visitors.NewConjureTypeProvider(aliasDefinition.Alias)
	if err != nil {
		return err
	}
	aliasTyper, err := conjureTypeProvider.ParseType(info)
	if err != nil {
		return errors.Wrapf(err, "alias type %s specifies unrecognized type", aliasDefinition.TypeName.Name)
	}
	info.AddImports(aliasTyper.ImportPaths()...)
	decls, imports, err := astForAlias(aliasDefinition, info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for alias %s", aliasDefinition.TypeName.Name)
	}

	info.AddImports(imports.Sorted()...)
	c.aliases.Decls = append(c.aliases.Decls, decls...)
	return nil
}

func (c *outputFileCollector) VisitEnum(enumDefinition spec.EnumDefinition) error {
	decls, imports := astForEnum(enumDefinition)
	c.enums.Info.AddImports(imports.Sorted()...)
	c.enums.Decls = append(c.enums.Decls, decls...)
	return nil
}

func (c *outputFileCollector) VisitObject(objectDefinition spec.ObjectDefinition) error {
	info := c.objects.Info
	uniqueGoPkgs, err := getImportPathsFromFields(objectDefinition.Fields, info)
	if err != nil {
		return err
	}
	info.AddImports(uniqueGoPkgs.Sorted()...)

	objDecls, imports, err := astForObject(objectDefinition, info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for object %s", objectDefinition.TypeName.Name)
	}

	info.AddImports(imports.Sorted()...)
	c.objects.Decls = append(c.objects.Decls, objDecls...)
	return nil
}

func (c *outputFileCollector) VisitUnion(unionDefinition spec.UnionDefinition) error {
	info := c.unions.Info
	uniqueGoPkgs, err := getImportPathsFromFields(unionDefinition.Union, info)
	if err != nil {
		return err
	}
	info.AddImports(uniqueGoPkgs.Sorted()...)

	declers, imports, err := astForUnion(unionDefinition, info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for union type %q", unionDefinition.TypeName.Name)
	}
	info.AddImports(imports.Sorted()...)
	c.unions.Decls = append(c.unions.Decls, declers...)
	return nil
}

func (c *outputFileCollector) VisitError(errorDefinition spec.ErrorDefinition) error {
	info := c.errors.Info
	allArgs := make([]spec.FieldDefinition, 0, len(errorDefinition.SafeArgs)+len(errorDefinition.UnsafeArgs))
	allArgs = append(allArgs, errorDefinition.SafeArgs...)
	allArgs = append(allArgs, errorDefinition.UnsafeArgs...)
	uniqueGoPkgs, err := getImportPathsFromFields(allArgs, info)
	if err != nil {
		return err
	}
	info.AddImports(uniqueGoPkgs.Sorted()...)

	errorDecls, imports, err := astForError(errorDefinition, info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for error %s", errorDefinition.ErrorName.Name)
	}
	info.AddImports(imports.Sorted()...)
	c.errors.Decls = append(c.errors.Decls, errorDecls...)
	return nil
}

func (c *outputFileCollector) VisitService(serviceDefinition spec.ServiceDefinition) error {
	info := c.services.Info
	for _, endpointDefinition := range serviceDefinition.Endpoints {
		for _, endpointArg := range endpointDefinition.Args {
			typer, err := getTyperFromType(endpointArg.Type, info)
			if err != nil {
				return err
			}
			info.AddImports(typer.ImportPaths()...)
		}
		if endpointDefinition.Returns != nil {
			typer, err := getTyperFromType(*endpointDefinition.Returns, info)
			if err != nil {
				return err
			}
			info.AddImports(typer.ImportPaths()...)
		}
	}

	declers, imports, err := astForService(serviceDefinition, info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for service %s", serviceDefinition.ServiceName.Name)
	}
	info.AddImports(imports.Sorted()...)
	c.services.Decls = append(c.services.Decls, declers...)
	return nil
}

func (c *outputFileCollector) VisitUnknown(typeName string) error {
	return errors.New("Unknown Type found " + typeName)
}

func getImportPathsFromFields(fields []spec.FieldDefinition, info types.PkgInfo) (StringSet, error) {
	uniqueGoPkgs := NewStringSet()
	for _, field := range fields {
		typer, err := getTyperFromType(field.Type, info)
		if err != nil {
			return nil, err
		}
		for _, importPath := range typer.ImportPaths() {
			uniqueGoPkgs[importPath] = struct{}{}
		}
	}
	return uniqueGoPkgs, nil
}

func getTyperFromType(specType spec.Type, info types.PkgInfo) (types.Typer, error) {
	conjureTypeProvider, err := visitors.NewConjureTypeProvider(specType)
	if err != nil {
		return nil, err
	}
	return conjureTypeProvider.ParseType(info)
}

func newGoFile(filePath, goImportPath string, info types.PkgInfo, goTypeObjs []astgen.ASTDecl) (*OutputFile, error) {
	var components []astgen.ASTDecl
	imports := info.ImportAliases()
	if len(imports) > 0 {
		components = append(components, decl.NewImports(imports))
	}
	components = append(components, goTypeObjs...)

	_, pkgName := path.Split(goImportPath)
	return &OutputFile{
		pkgName:    pkgName,
		absPath:    filePath,
		goTypeObjs: components,
	}, nil
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
