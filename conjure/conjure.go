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
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/conjure-go/v6/conjure/visitors"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
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

func Generate(conjureDefinition spec.ConjureDefinition, outputConfiguration OutputConfiguration) error {
	files, err := GenerateOutputFiles(conjureDefinition, outputConfiguration)
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

func GenerateOutputFiles(conjureDefinition spec.ConjureDefinition, outputConfiguration OutputConfiguration) ([]*OutputFile, error) {
	conjurePkgToGoPkg, goPkgToFilePath, err := createMappingFunctions(outputConfiguration.OutputDir)
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
			cfg:      outputConfiguration,
			aliases:  fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			enums:    fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			objects:  fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			unions:   fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			errors:   fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			servers:  fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
			services: fileASTCollector{Info: types.NewPkgInfo(importPath, customTypes)},
		}
		if err := visitors.VisitConjureDefinition(conjureDef, collector); err != nil {
			return nil, err
		}
		collector.collectInitFuncs(conjureDef)

		for filename, ast := range map[string]fileASTCollector{
			"aliases.conjure.go":  collector.aliases,
			"enums.conjure.go":    collector.enums,
			"structs.conjure.go":  collector.objects,
			"unions.conjure.go":   collector.unions,
			"errors.conjure.go":   collector.errors,
			"servers.conjure.go":  collector.servers,
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
	cfg OutputConfiguration
	// Track outputs (imports and decls) per-file
	aliases  fileASTCollector
	enums    fileASTCollector
	objects  fileASTCollector
	unions   fileASTCollector
	errors   fileASTCollector
	servers  fileASTCollector
	services fileASTCollector
}

type fileASTCollector struct {
	Info  types.PkgInfo
	Decls []astgen.ASTDecl
}

func (c *outputFileCollector) VisitAlias(aliasDefinition spec.AliasDefinition) error {
	decls, err := astForAlias(aliasDefinition, c.aliases.Info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for alias %s", aliasDefinition.TypeName.Name)
	}
	c.aliases.Decls = append(c.aliases.Decls, decls...)
	return nil
}

func (c *outputFileCollector) VisitEnum(enumDefinition spec.EnumDefinition) error {
	decls := astForEnum(enumDefinition, c.enums.Info)
	c.enums.Decls = append(c.enums.Decls, decls...)
	return nil
}

func (c *outputFileCollector) VisitObject(objectDefinition spec.ObjectDefinition) error {
	objDecls, err := astForObject(objectDefinition, c.objects.Info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for object %s", objectDefinition.TypeName.Name)
	}
	c.objects.Decls = append(c.objects.Decls, objDecls...)
	return nil
}

func (c *outputFileCollector) VisitUnion(unionDefinition spec.UnionDefinition) error {
	declers, err := astForUnion(unionDefinition, c.unions.Info, c.cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for union type %q", unionDefinition.TypeName.Name)
	}
	c.unions.Decls = append(c.unions.Decls, declers...)
	return nil
}

func (c *outputFileCollector) VisitError(errorDefinition spec.ErrorDefinition) error {
	errorDecls, err := astForError(errorDefinition, c.errors.Info)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for error %s", errorDefinition.ErrorName.Name)
	}
	c.errors.Decls = append(c.errors.Decls, errorDecls...)
	return nil
}

func (c *outputFileCollector) VisitService(serviceDefinition spec.ServiceDefinition) error {
	servicesInfo := c.services.Info
	if err := addImportsToPkgInfoFromServiceDefinitionEndpoints(&servicesInfo, serviceDefinition); err != nil {
		return err
	}

	declers, imports, err := astForService(serviceDefinition, servicesInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to generate AST for service %s", serviceDefinition.ServiceName.Name)
	}
	servicesInfo.AddImports(imports.Sorted()...)
	c.services.Decls = append(c.services.Decls, declers...)

	// Generate server code if GenerateServer configuration is enabled
	if c.cfg.GenerateServer {
		serversInfo := c.servers.Info
		if err := addImportsToPkgInfoFromServiceDefinitionEndpoints(&serversInfo, serviceDefinition); err != nil {
			return err
		}

		serverInterfaces, err := AstForServerInterface(serviceDefinition, serversInfo)
		if err != nil {
			return errors.Wrapf(err, "failed to generate AST for service %s", serviceDefinition.ServiceName.Name)
		}
		c.servers.Decls = append(c.servers.Decls, serverInterfaces...)

		routeReg, err := ASTForServerRouteRegistration(serviceDefinition, serversInfo)
		if err != nil {
			return errors.Wrapf(err, "failed to generate AST for service %s", serviceDefinition.ServiceName.Name)
		}
		c.servers.Decls = append(c.servers.Decls, routeReg...)
		c.servers.Info.AddImports(imports.Sorted()...)

		handlers, err := AstForServerFunctionHandler(serviceDefinition, serversInfo)
		if err != nil {
			return errors.Wrapf(err, "failed to generate AST for service %s", serviceDefinition.ServiceName.Name)
		}
		c.servers.Decls = append(c.servers.Decls, handlers...)
		c.servers.Info.AddImports(imports.Sorted()...)
	}
	return nil
}

func addImportsToPkgInfoFromServiceDefinitionEndpoints(info *types.PkgInfo, serviceDefinition spec.ServiceDefinition) error {
	for _, endpointDefinition := range serviceDefinition.Endpoints {
		for _, endpointArg := range endpointDefinition.Args {
			typer, err := visitors.NewConjureTypeProviderTyper(endpointArg.Type, *info)
			if err != nil {
				return err
			}
			info.AddImports(typer.ImportPaths()...)
		}
		if endpointDefinition.Returns != nil {
			typer, err := visitors.NewConjureTypeProviderTyper(*endpointDefinition.Returns, *info)
			if err != nil {
				return err
			}
			info.AddImports(typer.ImportPaths()...)
		}
	}
	return nil
}

func (c *outputFileCollector) VisitUnknown(typeName string) error {
	return errors.New("Unknown Type found " + typeName)
}

func (c *outputFileCollector) collectInitFuncs(definition spec.ConjureDefinition) {
	if len(definition.Errors) > 0 {
		c.errors.Decls = append(c.errors.Decls, astErrorInitFunc(definition.Errors, c.errors.Info))
	}
}

func addImportPathsFromFields(fields []spec.FieldDefinition, info types.PkgInfo) error {
	for _, field := range fields {
		typer, err := visitors.NewConjureTypeProviderTyper(field.Type, info)
		if err != nil {
			return err
		}
		for _, importPath := range typer.ImportPaths() {
			info.AddImports(importPath)
		}
	}
	return nil
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

// outputPackageBasePath returns the Go package path to the base output directory.
//
// If the output directory is within a module (that is, if the command "go list -m" completes successfully in the
// directory), the returned value is the module path joined with the path between the module root and the output
// directory. For example, if the project module is "github.com/org/project" and the output directory is "outDir", then
// the returned path is "github.com/org/project/outDir".
//
// If the output directory is not within a module, then the returned package path is the path as determined by running
// "packages.Load" in the output directory. For example, if the project is in package "github.com/org/project" and the
// output directory is the "outDir" directory within that package, the returned path is "github.com/org/project/outDir".
// Any conjure-generated package paths should be appended to this path.
func outputPackageBasePath(outputDirAbsPath string) (string, error) {
	// ensure that output directory exists, as "packages.Load" may require this
	if _, err := os.Stat(outputDirAbsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDirAbsPath, 0755); err != nil {
			return "", errors.Wrapf(err, "failed to create directory")
		}
	}

	modName, modBaseDir, err := goModulePath(outputDirAbsPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to determine if output directory is in a module")
	}
	if modName != "" {
		normalizedOutputPath, err := filepath.EvalSymlinks(outputDirAbsPath)
		if err != nil {
			return "", errors.Wrapf(err, "failed to resolve sym links")
		}
		normalizedModBaseDir, err := filepath.EvalSymlinks(modBaseDir)
		if err != nil {
			return "", errors.Wrapf(err, "failed to resolve sym links")
		}

		relPath, err := filepath.Rel(normalizedModBaseDir, normalizedOutputPath)
		if err != nil {
			return "", errors.Wrapf(err, "failed to determine relative path for module directory")
		}
		return filepath.Join(modName, relPath), nil
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

// Returns the module name and the path to the base module for the specified directory. Returns an empty string if the
// provided directory is not within a module. Returns an error if any errors are encountered in making this
// determination.
func goModulePath(dir string) (modName string, modBaseDir string, rErr error) {
	cmd := exec.Command("go", "list", "-m", "-mod=readonly", "-json")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		if string(output) == "go list -m: not using modules\n" {
			// directory is not a module: return empty and no error so that package mode is used. Note that this can
			// happen even if a go.mod exists if GO111MODULE=off is set, but this is correct behavior (if module mode
			// is explicitly set to off, conjure should not consider modules either).
			return "", "", nil
		}
		// "go list -m" failed for a reason other than not being a module: return error
		return "", "", errors.Wrapf(err, "%v failed with output: %q", cmd.Args, string(output))
	}
	modJSON := struct {
		Path string
		Dir  string
	}{}
	if err := json.Unmarshal(output, &modJSON); err != nil {
		return "", "", errors.Wrapf(err, "failed to unmarshal output of %v as JSON", cmd.Args)
	}
	return modJSON.Path, modJSON.Dir, nil
}
