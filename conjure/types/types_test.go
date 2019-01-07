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
