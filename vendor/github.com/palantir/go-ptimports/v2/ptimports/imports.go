// Copyright 2016 Palantir Technologies, Inc.
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

// Based on golang.org/x/tools/imports which bears the following license:
//
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ptimports

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"golang.org/x/tools/imports"
)

// ProcessFileFromInput processes the provided file from the provider reader. If the reader is nil, then the file
// described by filename is opened and used as the reader.
func ProcessFileFromInput(filename string, in io.Reader, list, write bool, options *Options, stdout io.Writer) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := Process(filename, src, options)
	if err != nil {
		return err
	}

	if list {
		if !bytes.Equal(src, res) {
			_, _ = fmt.Fprintln(stdout, filename)
		}
		return nil
	}

	if write {
		// only write when file changed
		if !bytes.Equal(src, res) {
			return ioutil.WriteFile(filename, res, 0)
		}
	} else {
		// print regardless of whether they are equal
		_, _ = fmt.Fprint(stdout, string(res))
	}
	return nil
}

type Options struct {
	// if true, converts single-line imports into import blocks
	Refactor bool
	// if true, runs the "gofmt simplify" operation on code
	Simplify bool
	// if true, does not add or remove imports
	FormatOnly bool
	// prefixes to use for goimports operation
	LocalPrefixes []string
}

// Process formats and adjusts imports for the provided file.
func Process(filename string, src []byte, options *Options) ([]byte, error) {
	if options == nil {
		options = &Options{}
	}
	importsOptions := &imports.Options{
		// these values are the default for imports.Process
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
		// use provided formatOnly value
		FormatOnly: options.FormatOnly,
	}

	// run goimports on output. Do this before refactoring so that the refactor operation has the most up-to-date
	// imports (the Process operation may add or remove imports).
	imports.LocalPrefix = strings.Join(options.LocalPrefixes, ",")
	out, err := imports.Process(filename, src, importsOptions)
	if err != nil {
		return nil, err
	}

	// if "simplify" is true, simplify the source
	if options.Simplify {
		out, err = simplifyFile(filename, out)
		if err != nil {
			return nil, err
		}
	}

	// if refactor is true, group import statements
	if options.Refactor {
		out, err = groupImports(filename, out)
		if err != nil {
			return nil, err
		}
		// run goimports on output after grouping imports
		out, err = imports.Process(filename, out, importsOptions)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func simplifyFile(filename string, src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, adjust, err := parse(fset, filename, src)
	if err != nil {
		return nil, err
	}
	simplify(file)

	printerMode := printer.UseSpaces | printer.TabIndent
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, fset, file)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	return out, nil
}

func groupImports(filename string, src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, adjust, err := parse(fset, filename, src)
	if err != nil {
		return nil, err
	}

	cImportsDocs, err := fixImports(fset, file)
	if err != nil {
		return nil, err
	}
	printerMode := printer.UseSpaces | printer.TabIndent
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, fset, file)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	out = addImportSpaces(out)

	cImportCommentIdx := 0
	out = regexp.MustCompile(`\nimport "C"`).ReplaceAllFunc(out, func(match []byte) []byte {
		if cImportCommentIdx >= len(cImportsDocs) {
			return []byte(string(match))
		}
		var commentLines []string
		for _, comment := range cImportsDocs[cImportCommentIdx].List {
			commentLines = append(commentLines, comment.Text)
		}
		val := []byte("\n" + strings.Join(commentLines, "\n") + string(match) + "\n")
		cImportCommentIdx++
		return val
	})
	return out, nil
}

// parse parses src, which was read from filename,
// as a Go source file or statement list.
func parse(fset *token.FileSet, filename string, src []byte) (*ast.File, func(orig, src []byte) []byte, error) {
	parserMode := parser.ParseComments

	// Try as whole source file.
	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	// If the error is that the source file didn't begin with a
	// package line, fall through to try as a source fragment.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	// If this is a declaration list, make it a source file
	// by inserting a package clause.
	// Insert using a ;, not a newline, so that the line numbers
	// in psrc match the ones in src.
	psrc := append([]byte("package main;"), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		// If a main function exists, we will assume this is a main
		// package and leave the file.
		if containsMainFunc(file) {
			return file, nil, nil
		}

		adjust := func(orig, src []byte) []byte {
			// Remove the package clause.
			// Gofmt has turned the ; into a \n.
			src = src[len("package main\n"):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	// If the error is that the source file didn't begin with a
	// declaration, fall through to try as a statement list.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	// If this is a statement list, make it a source file
	// by inserting a package clause and turning the list
	// into a function body.  This handles expressions too.
	// Insert using a ;, not a newline, so that the line numbers
	// in fsrc match the ones in src.
	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the wrapping.
			// Gofmt has turned the ; into a \n\n.
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			// Gofmt has also indented the function body one level.
			// Remove that indent.
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	// Failed, and out of options.
	return nil, nil, err
}

// containsMainFunc checks if a file contains a function declaration with the
// function signature 'func main()'
func containsMainFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				continue
			}

			if len(f.Type.Params.List) != 0 {
				continue
			}

			if f.Type.Results != nil && len(f.Type.Results.List) != 0 {
				continue
			}

			return true
		}
	}

	return false
}

func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}

// matchSpace reformats src to use the same space context as orig.
// 1) If orig begins with blank lines, matchSpace inserts them at the beginning of src.
// 2) matchSpace copies the indentation of the first non-blank line in orig
//    to every non-blank line in src.
// 3) matchSpace copies the trailing space from orig and uses it in place
//   of src's trailing space.
func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	_, _ = b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], line[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' { // not blank
			_, _ = b.Write(indent)
		}
		_, _ = b.Write(line)
	}
	_, _ = b.Write(after)
	return b.Bytes()
}

func addImportSpaces(input []byte) []byte {
	var out bytes.Buffer
	inImports := false
	done := false
	for _, currLineBytes := range bytes.Split(input, []byte("\n")) {
		s := string(currLineBytes)

		if !inImports && !done && strings.HasPrefix(s, "import") {
			inImports = true
		}
		if inImports && (s == ")" ||
			strings.HasPrefix(s, "var") ||
			strings.HasPrefix(s, "func") ||
			strings.HasPrefix(s, "const") ||
			strings.HasPrefix(s, "type")) {
			done = true
			inImports = false
		}
		if !inImports || s != "" {
			fmt.Fprintln(&out, s)
		}
	}
	return out.Bytes()
}
