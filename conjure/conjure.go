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
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v6/conjure/types"
	"github.com/pkg/errors"
)

func Generate(conjureDefinition spec.ConjureDefinition, outputConfiguration OutputConfiguration) error {
	//TODO(revert!)
	outputConfiguration.LiteralJSONMethods = true

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
			aliasFile := newJenFile(pkg.ImportPath)
			for _, alias := range pkg.Aliases {
				writeAliasType(aliasFile.Group, alias, cfg)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "aliases.conjure.go"), pkg.ImportPath, aliasFile))
		}
		if len(pkg.Enums) > 0 {
			enumFile := newJenFile(pkg.ImportPath)
			for _, enum := range pkg.Enums {
				writeEnumType(enumFile.Group, enum, cfg)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "enums.conjure.go"), pkg.ImportPath, enumFile))
		}
		if len(pkg.Objects) > 0 {
			objectFile := newJenFile(pkg.ImportPath)
			for _, object := range pkg.Objects {
				writeObjectType(objectFile.Group, object, cfg)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "structs.conjure.go"), pkg.ImportPath, objectFile))
		}
		if len(pkg.Unions) > 0 {
			unionFile := newJenFile(pkg.ImportPath)
			for _, union := range pkg.Unions {
				writeUnionType(unionFile.Group, union, cfg)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "unions.conjure.go"), pkg.ImportPath, unionFile))
		}
		if len(pkg.Errors) > 0 {
			errorFile := newJenFile(pkg.ImportPath)
			for _, errorDef := range pkg.Errors {
				writeErrorType(errorFile.Group, errorDef, cfg)
			}
			astErrorInitFunc(errorFile.Group, pkg.Errors)
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "errors.conjure.go"), pkg.ImportPath, errorFile))
		}
		if len(pkg.Services) > 0 {
			serviceFile := newJenFile(pkg.ImportPath)
			for _, service := range pkg.Services {
				writeServiceType(serviceFile.Group, service, cfg)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "services.conjure.go"), pkg.ImportPath, serviceFile))
		}
		if len(pkg.Services) > 0 && cfg.GenerateServer {
			serverFile := newJenFile(pkg.ImportPath)
			for _, server := range pkg.Services {
				writeServerType(serverFile.Group, server)
			}
			files = append(files, newGoFile(filepath.Join(pkg.OutputDir, "servers.conjure.go"), pkg.ImportPath, serverFile))
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].AbsPath() < files[j].AbsPath()
	})

	return files, nil
}

func newJenFile(importPath string) *jen.File {
	f := jen.NewFilePath(importPath)
	f.ImportNames(map[string]string{
		"github.com/palantir/witchcraft-go-params": "wparams",
		"github.com/palantir/witchcraft-go-error":  "werror",
	})
	return f
}

func newGoFile(filePath, goImportPath string, file *jen.File) *OutputFile {
	_, pkgName := path.Split(goImportPath)
	return &OutputFile{
		pkgName: pkgName,
		absPath: filePath,
		file:    file,
	}
}
