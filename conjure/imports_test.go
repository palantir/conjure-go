// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/palantir/conjure-go/v7/conjure/snip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportPathToAssumedName(t *testing.T) {
	for _, tc := range []struct {
		importPath string
		want       string
	}{
		{"github.com/palantir/pkg/safejson", "safejson"},
		{"github.com/palantir/witchcraft-go-error", "witchcraft"},
		{"github.com/palantir/witchcraft-go-logging/wlog-zap", "wlog"},
		{"github.com/palantir/conjure-go-runtime/v3/conjure-go-client/httpclient", "httpclient"},
		{"github.com/palantir/api/v3", "api"},
		{"github.com/hashicorp/go-multierror", "multierror"},
		{"example.com/conjure/com/palantir/bar_foo", "bar_foo"},
		{"gopkg.in/yaml.v3", "yaml"},
	} {
		assert.Equal(t, tc.want, importPathToAssumedName(tc.importPath), "import path %s", tc.importPath)
	}
}

// TestSetImportNameRendersExplicitAlias verifies that an import whose package name is not inferable from its path is
// rendered with an explicit alias. Generated files are formatted without goimports' import-resolution pass, so
// jennifer is the only thing that can supply the alias.
func TestSetImportNameRendersExplicitAlias(t *testing.T) {
	const werrorPath = "github.com/palantir/witchcraft-go-error"

	for _, tc := range []struct {
		name       string
		importPath string
		wantImport string
	}{
		{
			name:       "package name not inferable from path",
			importPath: werrorPath,
			wantImport: `werror "` + werrorPath + `"`,
		},
		{
			name:       "package name inferable from path",
			importPath: "github.com/palantir/pkg/safejson",
			wantImport: `"github.com/palantir/pkg/safejson"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := jen.NewFilePathName("example.com/testpkg", "testpkg")
			setImportNames(f, snip.ImportsToPackageNames())
			f.Var().Id("_").Op("=").Qual(tc.importPath, "Something")

			rendered, err := newGoFile("/tmp/test.go", f).Render()
			require.NoError(t, err)
			assert.Contains(t, string(rendered), tc.wantImport)
		})
	}
}
