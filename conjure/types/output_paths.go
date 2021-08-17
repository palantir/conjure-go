// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/palantir/conjure-go/v6/conjure/transforms"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func newPathTranslator(outputDir string) (pathTranslator, error) {
	outputDirAbsPath, err := filepath.Abs(outputDir)
	if err != nil {
		return pathTranslator{}, errors.Wrapf(err, "failed to convert %s to absolute path", outputDir)
	}
	outputPkgBasePath, err := outputPackageBasePath(outputDirAbsPath)
	if err != nil {
		return pathTranslator{}, errors.Wrapf(err, "failed to determine base import path for generated code")
	}
	return pathTranslator{outputDir: outputDir, outputPkgBasePath: outputPkgBasePath}, nil
}

type pathTranslator struct {
	outputDir         string
	outputPkgBasePath string
}

func (p pathTranslator) conjurePkgToGoPkg(conjurePkg string) string {
	return path.Join(p.outputPkgBasePath, transforms.PackagePath(conjurePkg))
}

// transforms the provided goPkgPath to the absolute on-disk path where files for the package should be generated.
// This abstraction is needed because the generated location may be different between $GOPATH projects and module
// projects.
func (p pathTranslator) goPkgToFilePath(goPkgPath string) string {
	return path.Join(p.outputDir, strings.TrimPrefix(goPkgPath, p.outputPkgBasePath+"/"))
}

func (p pathTranslator) conjurePkgToFilePath(conjurePkg string) string {
	return p.goPkgToFilePath(p.conjurePkgToGoPkg(conjurePkg))
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
