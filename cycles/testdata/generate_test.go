// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	spec_old "github.com/palantir/conjure-go/v6/conjure-api/spec_old"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecs(t *testing.T) {
	t.Skip()
	for _, conjureDir := range []string{
		"no-cycles",
		"cycle-within-pkg",
		"pkg-cycle",
		"pkg-cycle-disconnected",
		"type-cycle",
	} {
		t.Run(conjureDir, func(t *testing.T) {
			bytes, err := os.ReadFile(filepath.Join(conjureDir, "in.conjure.json"))
			require.NoError(t, err)
			var def1 spec.ConjureDefinition
			err = json.Unmarshal(bytes, &def1)
			require.NoError(t, err)
			var def2 spec_old.ConjureDefinition
			err = json.Unmarshal(bytes, &def2)
			require.NoError(t, err)
			//assert.EqualValues(t, def1, def2)
			//yaml1, err := yaml.Marshal(def1)
			//require.NoError(t, err)
			//yaml2, err := yaml.Marshal(def2)
			//require.NoError(t, err)
			//assert.Equal(t, string(yaml1), string(yaml2))
			json1, err := json.MarshalIndent(def1, "", "  ")
			require.NoError(t, err)
			json2, err := json.MarshalIndent(def2, "", "  ")
			require.NoError(t, err)
			assert.Equal(t, string(json1), string(json2))

			//ir1, err := conjure.FromIRFile(filepath.Join(conjureDir, "in.conjure.json"))
			//if err != nil {
			//	panic(err)
			//}
			//outputDir := filepath.Join(conjureDir, "conjure")
			//if err := os.RemoveAll(outputDir); err != nil {
			//	panic(err)
			//}
			//if err := conjure.Generate(ir, conjure.OutputConfiguration{
			//	GenerateFuncsVisitor: true,
			//	OutputDir:            outputDir,
			//}); err != nil {
			//	panic(err)
			//}
		})
	}
}
