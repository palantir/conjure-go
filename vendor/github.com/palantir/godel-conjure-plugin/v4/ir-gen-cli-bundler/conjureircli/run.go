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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/mholt/archiver"
	conjureircli_internal "github.com/palantir/godel-conjure-plugin/v4/ir-gen-cli-bundler/conjureircli/internal"
	"github.com/pkg/errors"
)

func YAMLtoIR(in []byte) (rBytes []byte, rErr error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); rErr == nil && err != nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory")
		}
	}()

	inPath := path.Join(tmpDir, "in.yml")
	if err := ioutil.WriteFile(inPath, in, 0644); err != nil {
		return nil, errors.WithStack(err)
	}
	return InputPathToIR(inPath)
}

func InputPathToIR(inPath string) (rBytes []byte, rErr error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); rErr == nil && err != nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory")
		}
	}()

	outPath := path.Join(tmpDir, "out.json")
	if err := Run(inPath, outPath); err != nil {
		return nil, err
	}
	irBytes, err := ioutil.ReadFile(outPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return irBytes, nil
}

func Run(inPath, outPath string) error {
	cliPath, err := cliCmdPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(cliPath, "compile", inPath, outPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "failed to execute %v\nOutput:\n%s", cmd.Args, string(output))
	}
	return nil
}

func cliCmdPath() (string, error) {
	cacheDirPath := path.Join(os.TempDir(), "_conjureircli")
	dstPath := path.Join(cacheDirPath, fmt.Sprintf("conjure-%v", conjureircli_internal.Version))
	if err := ensureCLIExists(dstPath); err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin", "linux":
		return path.Join(dstPath, "bin", "conjure"), nil
	default:
		return "", errors.Errorf("OS %s not supported", runtime.GOOS)
	}
}

func ensureCLIExists(dstPath string) error {
	if fi, err := os.Stat(dstPath); err == nil && fi.IsDir() {
		// destination already exists
		return nil
	}

	// expand asset into destination
	tgzBytes, err := conjureircli_internal.Asset("conjure.tgz")
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(archiver.TarGz.Read(bytes.NewReader(tgzBytes), path.Dir(dstPath)))
}
