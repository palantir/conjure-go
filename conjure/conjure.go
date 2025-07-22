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
	"path"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/palantir/pkg/safejson"
	"github.com/pkg/errors"
)

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

func GenerateOutputFiles(conjureDefinition spec.ConjureDefinition, cfg OutputConfiguration) ([]*OutputFile, error) {
	def, err := types.NewConjureDefinition(cfg.OutputDir, conjureDefinition)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid configuration")
	}

	var files []*OutputFile

	var errorRegistryImportPath string
	if len(conjureDefinition.Errors) > 0 {
		errorRegistryImportPath, err = types.GetGoPackageForInternalErrors(cfg.OutputDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine import path for error registry package")
		}
		errorRegistryJenFile := jen.NewFilePathName(errorRegistryImportPath, path.Base(errorRegistryImportPath))
		errorRegistryJenFile.ImportNames(snip.DefaultImportsToPackageNames)
		writeErrorRegistryFile(errorRegistryJenFile.Group)
		errorRegistryOutputDir, err := types.GetOutputDirectoryForGoPackage(cfg.OutputDir, errorRegistryImportPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine output directory for error registry package")
		}
		files = append(files, newGoFile(filepath.Join(errorRegistryOutputDir, "error_registry.conjure.go"), errorRegistryJenFile))
	}

	for _, pkg := range def.Packages {
		if len(pkg.Aliases) > 0 {
			aliasFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, alias := range pkg.Aliases {
				writeAliasType(aliasFile.Group, alias)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "aliases.conjure.go"), aliasFile))
		}
		if len(pkg.Enums) > 0 {
			enumFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, enum := range pkg.Enums {
				writeEnumType(enumFile.Group, enum)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "enums.conjure.go"), enumFile))
		}
		if len(pkg.Objects) > 0 {
			objectFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, object := range pkg.Objects {
				writeObjectType(objectFile.Group, object)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "structs.conjure.go"), objectFile))
		}
		if len(pkg.Unions) > 0 {
			unionFile := newJenFile(pkg, def, errorRegistryImportPath)
			goUnionGenericsFile := newJenFile(pkg, def, errorRegistryImportPath)
			goUnionGenericsFile.Comment("//go:build go1.18")
			for _, union := range pkg.Unions {
				writeUnionType(unionFile.Group, union, cfg.GenerateFuncsVisitor)
				writeUnionTypeWithGenerics(goUnionGenericsFile.Group, union)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "unions.conjure.go"), unionFile))
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "unions_generics.conjure.go"), goUnionGenericsFile))
		}
		if len(pkg.Errors) > 0 {
			errorFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, errorDef := range pkg.Errors {
				writeErrorType(errorFile.Group, errorDef)
			}
			astErrorInitFunc(errorFile.Group, pkg.Errors, errorRegistryImportPath)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "errors.conjure.go"), errorFile))
		}
		if len(pkg.Services) > 0 {
			serviceFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, service := range pkg.Services {
				writeServiceType(serviceFile.Group, service, errorRegistryImportPath)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "services.conjure.go"), serviceFile))
		}
		if len(pkg.Services) > 0 && cfg.GenerateCLI {
			cliFile := newJenFile(pkg, def, errorRegistryImportPath)
			writeCLIType(cliFile.Group, pkg.Services)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "cli.conjure.go"), cliFile))
		}
		if len(pkg.Services) > 0 && cfg.GenerateServer {
			serverFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, server := range pkg.Services {
				writeServerType(serverFile.Group, server)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "servers.conjure.go"), serverFile))
		}
		if len(def.Extensions) > 0 {
			const extensions = "extensions.conjure.json"

			extensionsContent, err := safejson.MarshalIndent(def.Extensions, "", "\t")
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal the conjure IR `extensions` field")
			}
			files = append(files, newRawFile(filepath.Join(pkg.OutputDir, extensions), extensionsContent))

			embedFile := newJenFile(pkg, def, errorRegistryImportPath)
			embedFileAsBlankIdentifierByteSlice(embedFile, extensions)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "embed.conjure.go"), embedFile))
		}

	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

func newJenFile(pkg types.ConjurePackage, def *types.ConjureDefinition, errorRegistryImportPath string) *jen.File {
	f := jen.NewFilePathName(pkg.ImportPath, pkg.PackageName)
	f.ImportNames(snip.DefaultImportsToPackageNames)
	for _, conjurePackage := range def.Packages {
		if packageSuffixRequiresAlias(conjurePackage.ImportPath) {
			f.ImportAlias(conjurePackage.ImportPath, conjurePackage.PackageName)
		} else {
			f.ImportName(conjurePackage.ImportPath, conjurePackage.PackageName)
		}
	}
	if errorRegistryImportPath != "" {
		f.ImportName(errorRegistryImportPath, path.Base(errorRegistryImportPath))
	}
	return f
}

func newGoFile(filePath string, file *jen.File) *OutputFile {
	return &OutputFile{
		absPath: filePath,
		render:  func() ([]byte, error) { return renderJenFile(file) },
	}
}

func newRawFile(filePath string, bytes []byte) *OutputFile {
	return &OutputFile{
		absPath: filePath,
		render:  func() ([]byte, error) { return bytes, nil },
	}
}

func packageSuffixRequiresAlias(importPath string) bool {
	return regexp.MustCompile(`/v[0-9]+$`).MatchString(importPath)
}
