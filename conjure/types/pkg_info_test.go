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

package types_test

import (
	"testing"

	"github.com/palantir/conjure-go/conjure/types"
	"github.com/stretchr/testify/assert"
)

func TestPkgInfo_AddImports(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    []string
		Expected map[string]string
	}{
		{
			Name: "no duplicates",
			Input: []string{
				"bytes",
				"encoding/json",
				"io",
			},
			Expected: map[string]string{
				"bytes":         "",
				"encoding/json": "",
				"io":            "",
			},
		},
		{
			Name: "with duplicates",
			Input: []string{
				"bytes",
				"encoding/json",
				"example1/bytes",
				"example2/bytes",
				"io",
				"encoding/json",
				"example1/bytes",
			},
			Expected: map[string]string{
				"bytes":          "",
				"encoding/json":  "",
				"example1/bytes": "bytes_1",
				"example2/bytes": "bytes_2",
				"io":             "",
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("call once", func(t *testing.T) {
				info := types.NewPkgInfo("foo", nil)
				info.AddImports(test.Input...)
				assert.Equal(t, test.Expected, info.ImportAliases())
			})
			t.Run("call per import", func(t *testing.T) {
				info := types.NewPkgInfo("foo", nil)
				for _, input := range test.Input {
					info.AddImports(input)
				}
				assert.Equal(t, test.Expected, info.ImportAliases())
			})
		})
	}
}
