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
	// Set the runtime versions before generating any code
	snip.SetCGRModuleVersion(cfg.CGRModuleVersion)
	snip.SetWGSModuleVersion(cfg.WGSModuleVersion)

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
		errorRegistryJenFile.ImportNames(snip.ImportsToPackageNames())
		writeErrorRegistryFile(errorRegistryJenFile.Group)
		errorRegistryOutputDir, err := types.GetOutputDirectoryForGoPackage(cfg.OutputDir, errorRegistryImportPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine output directory for error registry package")
		}
		files = append(files, newGoFile(filepath.Join(errorRegistryOutputDir, "error_registry.conjure.go"), errorRegistryJenFile))

		if cfg.ExportErrorDecoder {
			errorDecoderImportPath, err := types.GetGoPackageForErrorDecoder(cfg.OutputDir)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to determine import path for error decoder package")
			}
			errorDecoderJenFile := jen.NewFilePathName(errorDecoderImportPath, path.Base(errorDecoderImportPath))
			errorDecoderJenFile.ImportNames(snip.ImportsToPackageNames())
			errorDecoderJenFile.ImportAlias(errorRegistryImportPath, "internalconjureerrors")
			writeErrorDecoderExportFile(errorDecoderJenFile.Group, errorRegistryImportPath)
			errorDecoderOutputDir, err := types.GetOutputDirectoryForGoPackage(cfg.OutputDir, errorDecoderImportPath)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to determine output directory for error decoder package")
			}
			files = append(files, newGoFile(filepath.Join(errorDecoderOutputDir, "decoder.conjure.go"), errorDecoderJenFile))
		}
	}

	extensionsImportPath, err := types.GetGoPackageForEmbedFile(cfg.OutputDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine import path for extensions package")
	}

	extensionsOutputDir, err := types.GetOutputDirectoryForGoPackage(cfg.OutputDir, extensionsImportPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine output directory for extensions package")
	}

	// todo(aviradinsky): delete once rolled out
	for _, fileName := range [...]string{
		"extensions.conjure.json",
		"embed.conjure.go",
	} {
		filePath := filepath.Join(extensionsOutputDir, fileName)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed when attempted to remove file: %s: %w", filePath, err)
		}
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
			safetyCache := make(map[types.Type]spec.LogSafety)
			for _, object := range pkg.Objects {
				writeObjectType(objectFile.Group, object, safetyCache)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "structs.conjure.go"), objectFile))
		}
		if len(pkg.Unions) > 0 {
			unionFile := newJenFile(pkg, def, errorRegistryImportPath)
			goUnionGenericsFile := newJenFile(pkg, def, errorRegistryImportPath)
			goUnionGenericsFile.Comment("//go:build go1.18")
			for _, union := range pkg.Unions {
				writeUnionType(unionFile.Group, union, cfg.GenerateFuncsVisitor)
				writeUnionTypeWithGenerics(goUnionGenericsFile.Group, union, cfg.GenerateFuncsVisitor)
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
		hasServices := len(pkg.Services) > 0
		if hasServices {
			serviceFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, service := range pkg.Services {
				writeServiceType(serviceFile.Group, service, errorRegistryImportPath)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "services.conjure.go"), serviceFile))
		}
		if hasServices && cfg.GenerateCLI {
			cliFile := newJenFile(pkg, def, errorRegistryImportPath)
			writeCLIType(cliFile.Group, pkg.Services)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "cli.conjure.go"), cliFile))
		}
		if hasServices && cfg.GenerateServer {
			serverFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, server := range pkg.Services {
				writeServerType(serverFile.Group, server)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "servers.conjure.go"), serverFile))
		}

		const recommendedProductDependencies = "recommended-product-dependencies"
		if v, ok := def.Extensions[recommendedProductDependencies]; ok && hasServices {
			if vList, ok := v.([]any); ok && len(vList) != 0 {
				const extensions = "extensions.conjure.json"

				extensionsContent, err := safejson.MarshalIndent(map[string]any{
					recommendedProductDependencies: v,
				}, "", "\t")
				if err != nil {
					return nil, errors.Wrapf(err, "failed to marshal the conjure IR `extensions` field")
				}
				files = append(files, newRawFile(filepath.Join(pkg.OutputDir, extensions), extensionsContent))

				embedFile := newJenFile(pkg, def, errorRegistryImportPath)
				embedFileAsBlankIdentifierString(embedFile, extensions)
				files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "embed.conjure.go"), embedFile))
			}
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

func newJenFile(pkg types.ConjurePackage, def *types.ConjureDefinition, errorRegistryImportPath string) *jen.File {
	f := jen.NewFilePathName(pkg.ImportPath, pkg.PackageName)
	f.ImportNames(snip.ImportsToPackageNames())
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
	// add header to indicate that this is generated code using official format for comment indicating generated code as specified in https://github.com/golang/go/issues/13560#issuecomment-288457920
	file.HeaderComment("Code generated by conjure-go. DO NOT EDIT.")

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
