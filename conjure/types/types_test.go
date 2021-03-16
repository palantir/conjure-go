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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimitives(t *testing.T) {

	for _, test := range []struct {
		name       string
		typer      Typer
		currPkg    string
		importsMap map[string]string
		want       string
	}{
		{
			name:  "String type",
			typer: String,
			want:  "string",
		},
		{
			name:  "Integer type",
			typer: Integer,
			want:  "int",
		},
		{
			name:  "Double type",
			typer: Double,
			want:  "float64",
		},
		{
			name:  "Boolean type",
			typer: Boolean,
			want:  "bool",
		},
		{
			name:  "Any type",
			typer: Any,
			want:  "interface{}",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			info := NewPkgInfo(test.currPkg, NewCustomConjureTypes())
			assert.Equal(t, test.want, test.typer.GoType(info))
		})
	}
}

func TestFunctionTypes(t *testing.T) {
	for _, test := range []struct {
		name            string
		ft              funcType
		currPkg         string
		importsMap      map[string]string
		want            string
		expectedImports []string
	}{
		{
			name:            "empty func",
			ft:              funcType{},
			want:            "func()",
			expectedImports: []string{},
		},
		{
			name: "one input",
			ft: funcType{
				inputs: []Typer{
					IOReadCloserType,
				},
			},
			want:            "func(io.ReadCloser)",
			expectedImports: []string{"io"},
		},
		{
			name: "two inputs",
			ft: funcType{
				inputs: []Typer{
					IOReadCloserType,
					Bearertoken,
				},
			},
			want:            "func(io.ReadCloser, bearertoken.Token)",
			expectedImports: []string{"io", "github.com/palantir/pkg/bearertoken"},
		},
		{
			name: "one output",
			ft: funcType{
				outputs: []Typer{
					String,
				},
			},
			want:            "func() string",
			expectedImports: []string{},
		},
		{
			name: "two outputs",
			ft: funcType{
				outputs: []Typer{
					String,
					IOReadCloserType,
				},
			},
			want:            "func() (string, io.ReadCloser)",
			expectedImports: []string{"io"},
		},
		{
			name: "one input, one output",
			ft: funcType{
				inputs: []Typer{
					IOReadCloserType,
				},
				outputs: []Typer{
					String,
				},
			},
			want:            "func(io.ReadCloser) string",
			expectedImports: []string{"io"},
		},
		{
			name: "two inputs, two outputs",
			ft: funcType{
				inputs: []Typer{
					IOReadCloserType,
					Integer,
				},
				outputs: []Typer{
					String,
					Boolean,
				},
			},
			want:            "func(io.ReadCloser, int) (string, bool)",
			expectedImports: []string{"io"},
		},
		{
			name: "function composition",
			ft: funcType{
				inputs: []Typer{
					&funcType{},
					&funcType{},
				},
				outputs: []Typer{
					&funcType{},
					&funcType{},
				},
			},
			want:            "func(func(), func()) (func(), func())",
			expectedImports: []string{},
		},
		{
			name: "no deduplication",
			ft: funcType{
				inputs: []Typer{
					IOReadCloserType,
					IOReadCloserType,
				},
			},
			want:            "func(io.ReadCloser, io.ReadCloser)",
			expectedImports: []string{"io", "io"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			info := NewPkgInfo(test.currPkg, NewCustomConjureTypes())
			assert.Equal(t, test.want, test.ft.GoType(info))
			assert.Equal(t, test.expectedImports, test.ft.ImportPaths())
		})
	}
}

func TestReadCloseMapper(t *testing.T) {

	for _, test := range []struct {
		name         string
		initialTyper Typer
		desiredTyper Typer
		hasImports   bool
	}{
		{
			name:         "String type",
			initialTyper: String,
			desiredTyper: String,
		},
		{
			name:         "SafeLong type",
			initialTyper: SafeLong,
			desiredTyper: SafeLong,
		},
		{
			name:         "Binary type",
			initialTyper: BinaryType,
			desiredTyper: IOReadCloserType,
			hasImports:   true,
		},
		{
			name:         "Optional String type",
			initialTyper: NewOptionalType(String),
			desiredTyper: NewOptionalType(String),
		},
		{
			name:         "Optional SafeLong type",
			initialTyper: NewOptionalType(SafeLong),
			desiredTyper: NewOptionalType(SafeLong),
		},
		{
			name:         "Optional Binary type",
			initialTyper: NewOptionalType(BinaryType),
			desiredTyper: NewOptionalType(IOReadCloserType),
			hasImports:   true,
		},
		{
			name:         "List Binary type",
			initialTyper: NewListType(BinaryType),
			desiredTyper: NewListType(BinaryType),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			typer, imports := MapBinaryTypeToReadCloserType(test.initialTyper)
			assert.Equal(t, typer, test.desiredTyper)
			if test.hasImports {
				assert.Equal(t, imports, IOReadCloserType.ImportPaths())
			} else {
				assert.Nil(t, imports)
			}
		})
	}
}
