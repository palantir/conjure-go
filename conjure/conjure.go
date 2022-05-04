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
	"path/filepath"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/snip"
	"github.com/palantir/conjure-go/v6/conjure/types"
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
	for _, pkg := range def.Packages {
		if len(pkg.Aliases) > 0 {
			aliasFile := newJenFile(pkg, def)
			for _, alias := range pkg.Aliases {
				writeAliasType(aliasFile.Group, alias)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "aliases.conjure.go"), aliasFile))
		}
		if len(pkg.Enums) > 0 {
			enumFile := newJenFile(pkg, def)
			for _, enum := range pkg.Enums {
				writeEnumType(enumFile.Group, enum)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "enums.conjure.go"), enumFile))
		}
		if len(pkg.Objects) > 0 {
			objectFile := newJenFile(pkg, def)
			for _, object := range pkg.Objects {
				writeObjectType(objectFile.Group, object)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "structs.conjure.go"), objectFile))
		}
		if len(pkg.Unions) > 0 {
			unionFile := newJenFile(pkg, def)
			go1Dot18UnionFile := newJenFile(pkg, def)
			go1Dot18UnionFile.Comment("//go:build go1.18")
			for _, union := range pkg.Unions {
				writeUnionType(unionFile.Group, union, cfg.GenerateFuncsVisitor)
				writeUnionTypeWithGenerics(go1Dot18UnionFile.Group, union)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "unions.conjure.go"), unionFile))
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "unions.conjure_go1.18.go"), go1Dot18UnionFile))
		}
		if len(pkg.Errors) > 0 {
			errorFile := newJenFile(pkg, def)
			for _, errorDef := range pkg.Errors {
				writeErrorType(errorFile.Group, errorDef)
			}
			astErrorInitFunc(errorFile.Group, pkg.Errors)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "errors.conjure.go"), errorFile))
		}
		if len(pkg.Services) > 0 {
			serviceFile := newJenFile(pkg, def)
			for _, service := range pkg.Services {
				writeServiceType(serviceFile.Group, service)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "services.conjure.go"), serviceFile))
		}
		if len(pkg.Services) > 0 && cfg.GenerateCLI {
			cliFile := newJenFile(pkg, def)
			writeCLIType(cliFile.Group, pkg.Services)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "cli.conjure.go"), cliFile))
		}
		if len(pkg.Services) > 0 && cfg.GenerateServer {
			serverFile := newJenFile(pkg, def)
			for _, server := range pkg.Services {
				writeServerType(serverFile.Group, server)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "servers.conjure.go"), serverFile))
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

func newJenFile(pkg types.ConjurePackage, def *types.ConjureDefinition) *jen.File {
	f := jen.NewFilePathName(pkg.ImportPath, pkg.PackageName)
	f.ImportNames(snip.DefaultImportsToPackageNames)
	for _, conjurePackage := range def.Packages {
		f.ImportName(conjurePackage.ImportPath, conjurePackage.PackageName)
	}
	return f
}

func newGoFile(filePath string, file *jen.File) *OutputFile {
	return &OutputFile{
		absPath: filePath,
		file:    file,
	}
}
