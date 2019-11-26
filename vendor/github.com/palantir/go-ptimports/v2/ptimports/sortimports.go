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
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

func fixImports(fset *token.FileSet, f *ast.File) (cImportDocs []*ast.CommentGroup, rErr error) {
	imports, cImports, cImportsDocs := takeImports(f)
	if imports == nil || len(imports.Specs) == 0 {
		return
	}

	fixParens(imports)
	f.Decls = append(cImports, append([]ast.Decl{imports}, f.Decls...)...)

	var comments []*ast.CommentGroup
	for _, fileComment := range f.Comments {
		skip := false
		for _, cImportComment := range cImportsDocs {
			if fileComment == cImportComment {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		comments = append(comments, fileComment)
	}
	f.Comments = comments

	buf := &bytes.Buffer{}
	if err := printer.Fprint(buf, fset, f); err != nil {
		return nil, err
	}
	newF, err := parser.ParseFile(fset, f.Name.Name, buf, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	*f = *newF
	return cImportsDocs, nil
}

func takeImports(f *ast.File) (imports *ast.GenDecl, cImports []ast.Decl, cImportDocs []*ast.CommentGroup) {
	for len(f.Decls) > 0 {
		d, ok := f.Decls[0].(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Import decls are always first.
			break
		}

		cImport := false
		if len(d.Specs) > 0 {
			for _, spec := range d.Specs {
				impSpec := spec.(*ast.ImportSpec)
				if impSpec.Path.Value == `"C"` {
					cImport = true
					cImportDocs = append(cImportDocs, d.Doc)
					d.Doc = nil
				}
			}
		}

		if cImport {
			cImports = append(cImports, d)
		} else if imports == nil {
			imports = d
		} else {
			if imports.Doc == nil {
				imports.Doc = d.Doc
			} else if d.Doc != nil {
				imports.Doc.List = append(imports.Doc.List, d.Doc.List...)
			}
			imports.Specs = append(imports.Specs, d.Specs...)
		}

		// Put back later in a single decl
		f.Decls = f.Decls[1:]
	}
	return imports, cImports, cImportDocs
}

// All import decls require parens, even with only a single import.
func fixParens(d *ast.GenDecl) {
	if !d.Lparen.IsValid() {
		d.Lparen = d.Specs[0].Pos()
	}
}
