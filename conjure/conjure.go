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

// Generate generates the Conjure output files specified by the provided conjureDefinition and outputConfiguration and
// writes the result to disk. Effectively combines the functionality of the GenerateOutputFiles and WriteOutputFiles
// functions.
func Generate(conjureDefinition spec.ConjureDefinition, outputConfiguration OutputConfiguration) error {
	files, err := GenerateOutputFiles(conjureDefinition, outputConfiguration)
	if err != nil {
		return errors.Wrapf(err, "failed to generate Conjure output files to write")
	}
	if err := WriteOutputFiles(files); err != nil {
		return errors.Wrapf(err, "failed to write generated Conjure output files")
	}
	return nil
}

// WriteOutputFiles writes the OutputFile structs in the provided slice by calling the Write() function on each file.
// Skips nil elements. Returns an error if any of the Write() operations return a non-nil error.
func WriteOutputFiles(files []*OutputFile) error {
	for _, file := range files {
		if file == nil {
			// skip any files that are nil. It would be better for the files parameter to just be []OutputFile,
			// but keeping this construction instead because the existing GenerateOutputFiles function returns
			// []*OutputFile and this function is expected to be most commonly be used on the output of that function.
			continue
		}
		if err := file.Write(); err != nil {
			return errors.Wrapf(err, "failed to write output file %q", file.AbsPath())
		}
	}
	return nil
}

// GenerateOutputFiles returns a slice of OutputFile structs that represents the files that would be generated if the
// provided spec.ConjureDefinition was generated using the provided OutputConfiguration. Does not modify any on-disk
// state. The returned slice is ordered based on the absolute path of the file. If the function does not return an
// error, then the slice is guaranteed not to contain any nil elements (given that, it would be better for this function
// to return []OutputFile, but the signature was []*OutputFile when it was originally written and given that this is an
// exported function
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
		if len(pkg.Services) > 0 {
			serviceFile := newJenFile(pkg, def, errorRegistryImportPath)
			for _, service := range pkg.Services {
				writeServiceType(serviceFile.Group, service, errorRegistryImportPath)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "services.conjure.go"), serviceFile))

			if cfg.GenerateCLI {
				cliFile := newJenFile(pkg, def, errorRegistryImportPath)
				writeCLIType(cliFile.Group, pkg.Services)
				files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "cli.conjure.go"), cliFile))
			}

			if cfg.GenerateServer {
				serverFile := newJenFile(pkg, def, errorRegistryImportPath)
				for _, server := range pkg.Services {
					writeServerType(serverFile.Group, server)
				}
				files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "servers.conjure.go"), serverFile))
			}

			const recommendedProductDependencies = "recommended-product-dependencies"
			if v, ok := def.Extensions[recommendedProductDependencies]; ok {
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
