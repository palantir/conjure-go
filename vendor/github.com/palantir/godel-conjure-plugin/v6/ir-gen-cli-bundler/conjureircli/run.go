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

package conjureircli

import (
	_ "embed" // required for go:embed directive
	"os"
	"path/filepath"

	"github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli/internal"
	"github.com/palantir/pkg/clipackager"
	"github.com/palantir/pkg/safejson"
	"github.com/pkg/errors"
)

var (
	//go:embed internal/conjure.tgz
	conjureCLITGZ []byte

	// CLI runner that runs the Conjure CLI
	cliRunner = clipackager.NewDefaultPackagedCLIRunner(
		"conjure",
		internal.Version,
		conjureCLITGZ,
		".tgz",
	)
)

func YAMLtoIR(in []byte) (rBytes []byte, rErr error) {
	return YAMLtoIRWithParams(in)
}

func YAMLtoIRWithParams(in []byte, params ...Param) (rBytes []byte, rErr error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); rErr == nil && err != nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory")
		}
	}()

	inPath := filepath.Join(tmpDir, "in.yml")
	if err := os.WriteFile(inPath, in, 0644); err != nil {
		return nil, errors.WithStack(err)
	}
	return InputPathToIRWithParams(inPath, params...)
}

func InputPathToIR(inPath string) (rBytes []byte, rErr error) {
	return InputPathToIRWithParams(inPath)
}

func InputPathToIRWithParams(inPath string, params ...Param) (rBytes []byte, rErr error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); rErr == nil && err != nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory")
		}
	}()

	outPath := filepath.Join(tmpDir, "out.json")
	if err := RunWithParams(inPath, outPath, params...); err != nil {
		return nil, err
	}
	irBytes, err := os.ReadFile(outPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return irBytes, nil
}

// Run invokes the "compile" operation on the Conjure CLI with the provided inPath and outPath as arguments.
func Run(inPath, outPath string) error {
	return RunWithParams(inPath, outPath)
}

type runArgs struct {
	extensionsContent []byte
}

type Param interface {
	apply(*runArgs)
}

type paramFn func(*runArgs)

func (fn paramFn) apply(r *runArgs) {
	fn(r)
}

// ExtensionsParam returns a parameter that sets the extensions of the generated Conjure IR to be the JSON-marshalled
// content of the provided map if it is non-empty. Returns a no-op parameter if the provided map is nil or empty.
func ExtensionsParam(extensionsContent map[string]interface{}) (Param, error) {
	if len(extensionsContent) == 0 {
		return nil, nil
	}
	extensionBytes, err := safejson.Marshal(extensionsContent)
	if err != nil {
		return nil, err
	}
	return paramFn(func(r *runArgs) {
		r.extensionsContent = extensionBytes
	}), nil
}

// RunWithParams invokes the "compile" operation on the Conjure CLI with the provided inPath and outPath as arguments.
// Any arguments or configuration supplied by the provided params are also applied.
func RunWithParams(inPath, outPath string, params ...Param) error {
	// apply provided params
	var runArgCollector runArgs
	for _, param := range params {
		if param == nil {
			continue
		}
		param.apply(&runArgCollector)
	}

	// invoke the "compile" command
	args := []string{"compile"}

	// if extensionsContent is non-empty, add as flag
	if len(runArgCollector.extensionsContent) > 0 {
		args = append(args, "--extensions", string(runArgCollector.extensionsContent))
	}

	// set the inPath and outPath as final arguments
	args = append(args, inPath, outPath)

	if cliPath, output, err := clipackager.RunPackagedCLI(cliRunner, args...); err != nil {
		return errors.Wrapf(err, "failed to execute %v\nOutput:\n%s", append([]string{cliPath}, args...), string(output))
	}
	return nil
}
