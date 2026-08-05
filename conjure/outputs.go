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
	"bytes"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/go-ptimports/v2/ptimports"
	"github.com/pkg/errors"
)

type OutputFile struct {
	absPath string
	render  func() ([]byte, error)
}

func (f *OutputFile) AbsPath() string {
	return f.absPath
}

func (f *OutputFile) Write() error {
	body, err := f.Render()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(f.absPath), 0755); err != nil {
		return errors.Wrapf(err, "failed to create parent directory for Go file output %s", f.absPath)
	}
	if err := os.WriteFile(f.absPath, body, 0644); err != nil {
		return errors.Wrapf(err, "failed to write Go file output to %s", f.absPath)
	}
	return nil
}

func (f *OutputFile) Render() ([]byte, error) {
	bytes, err := f.render()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to render bytes for file %s", f.absPath)
	}
	return bytes, nil
}

// newlineAfterBrace matches a closing brace at the start of a line that is immediately followed by another declaration.
var newlineAfterBrace = regexp.MustCompile(`(\n}\n)(\S)`)

// ptimportsOptions formats the rendered source without letting goimports add or remove imports.
//
// jennifer emits an import for every qualified reference it renders and for nothing else, so there is never anything
// for goimports to resolve. Leaving resolution enabled is not just redundant, it is expensive: each call builds a fresh
// resolver that shells out to "go env" and reads the whole module-cache index (tens of MB) before it can conclude the
// file is already complete, which dominates generation time for large IRs.
var ptimportsOptions = &ptimports.Options{Refactor: true, FormatOnly: true}

func renderJenFile(f *jen.File) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return nil, errors.Wrapf(err, "failed to generate Go source")
	}

	goFileSrc, err := ptimports.Process("", buf.Bytes(), ptimportsOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run ptimports on generated Go source")
	}
	// add extra newline after braces
	goFileSrc = newlineAfterBrace.ReplaceAll(goFileSrc, []byte("$1\n$2"))

	return goFileSrc, nil
}
