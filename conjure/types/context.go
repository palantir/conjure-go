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
)

type TypeContext struct {
	currPkgPath   string
	customTypes   CustomConjureTypes
	importAliases map[string]string
}

func NewTypeContext(currPkgPath string, customTypes CustomConjureTypes) TypeContext {
	return TypeContext{
		currPkgPath:   currPkgPath,
		customTypes:   customTypes,
		importAliases: make(map[string]string),
	}
}

func (ctx *TypeContext) CustomTypes() CustomConjureTypes {
	return ctx.customTypes
}

// ImportAliases returns a copy of the map of import paths to aliases.
// Modifications to the returned map will not be written to the context.
func (ctx *TypeContext) ImportAliases() map[string]string {
	m := make(map[string]string, len(ctx.importAliases))
	for k, v := range ctx.importAliases {
		m[k] = v
	}
	return m
}

// AddImports adds imports to the internal mapping tracking import paths
// and package aliases in the event of conflicts.
// Typer.GoType uses this map to correctly build the selector for an imported declaration.
func (ctx *TypeContext) AddImports(imports ...string) {
	usedPkgNames := make(map[string]struct{})
	for usedImport, usedPkgName := range ctx.importAliases {
		if usedPkgName == "" {
			_, pkgName := path.Split(usedImport)
			usedPkgNames[pkgName] = struct{}{}
		} else {
			usedPkgNames[usedPkgName] = struct{}{}
		}
	}

	for _, importName := range imports {
		if importName == ctx.currPkgPath {
			// skip local package
			continue
		}
		if _, ok := ctx.importAliases[importName]; ok {
			// package is already imported
			continue
		}

		// TODO(bmoylan): Use golang.org/x/tools/go/packages to get more correct package names? This would create more dependency on generating into a correctly-formed project, maybe we just do best effort.
		_, pkgName := path.Split(importName)

		if _, ok := usedPkgNames[pkgName]; !ok {
			// package name has not been used yet -- no need for alias
			ctx.importAliases[importName] = ""
			usedPkgNames[pkgName] = struct{}{}
			continue
		}

		// package name has been used before -- need to find a unique alias and record
		currIdx := 1
		// append number to package name to make it unique. Increment counter until unique identifier is found.
		pkgShortName := pkgName
		for {
			pkgName = fmt.Sprintf("%s_%d", pkgShortName, currIdx)
			if _, ok := usedPkgNames[pkgName]; !ok {
				// package name is available
				break
			}
			currIdx++
		}
		// add entry to alias map
		ctx.importAliases[importName] = pkgName
		usedPkgNames[pkgName] = struct{}{}
	}
}
