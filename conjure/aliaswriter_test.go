// Copyright (c) 2019 Palantir Technologies. All rights reserved.
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
	"strings"
	"testing"

	"github.com/palantir/goastwriter"
	"github.com/palantir/goastwriter/astgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/conjure/types"
)

func TestAliasWriter(t *testing.T) {
	for caseNum, currCase := range []struct {
		pkg     string
		name    string
		aliases []spec.AliasDefinition
		want    string
	}{
		{
			pkg:  "testpkg",
			name: "single string alias",
			aliases: []spec.AliasDefinition{
				{
					TypeName: spec.TypeName{
						Name:    "Month",
						Package: "api",
					},
					Docs:  docPtr("These represent months"),
					Alias: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
				},
			},
			want: `package testpkg

// These represent months
type Month string
`,
		},
		{
			pkg:  "testpkg",
			name: "single optional string alias",
			aliases: []spec.AliasDefinition{
				{
					TypeName: spec.TypeName{
						Name:    "Month",
						Package: "api",
					},
					Docs: docPtr("These represent months"),
					Alias: spec.NewTypeFromOptional(spec.OptionalType{
						ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
					}),
				},
			},
			want: `package testpkg

// These represent months
type Month struct {
	Value *string
}

func (a Month) MarshalText() ([]byte, error) {
	if a.Value == nil {
		return nil, nil
	}
	return []byte(*a.Value), nil
}
func (a *Month) UnmarshalText(data []byte) error {
	rawMonth := string(data)
	a.Value = &rawMonth
	return nil
}
func (a Month) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (a *Month) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
`,
		},
		{
			pkg:  "testpkg",
			name: "single object alias",
			aliases: []spec.AliasDefinition{
				{
					TypeName: spec.TypeName{
						Name:    "Map",
						Package: "api",
					},
					Alias: spec.NewTypeFromMap(spec.MapType{
						KeyType:   spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
						ValueType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeSafelong),
					}),
				},
			},
			want: `package testpkg

type Map map[string]safelong.SafeLong

func (a Map) MarshalJSON() ([]byte, error) {
	return safejson.Marshal(map[string]safelong.SafeLong(a))
}
func (a *Map) UnmarshalJSON(data []byte) error {
	var rawMap map[string]safelong.SafeLong
	if err := safejson.Unmarshal(data, &rawMap); err != nil {
		return err
	}
	*a = Map(rawMap)
	return nil
}
func (a Map) MarshalYAML() (interface{}, error) {
	jsonBytes, err := safejson.Marshal(a)
	if err != nil {
		return nil, err
	}
	return safeyaml.JSONtoYAMLMapSlice(jsonBytes)
}
func (a *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	jsonBytes, err := safeyaml.UnmarshalerToJSONBytes(unmarshal)
	if err != nil {
		return err
	}
	return safejson.Unmarshal(jsonBytes, *&a)
}
`,
		},
	} {
		t.Run(currCase.name, func(t *testing.T) {
			info := types.NewPkgInfo("", nil)
			var components []astgen.ASTDecl
			for _, a := range currCase.aliases {
				declers, err := astForAlias(a, info)
				require.NoError(t, err)
				components = append(components, declers...)
			}

			got, err := goastwriter.Write(currCase.pkg, components...)
			require.NoError(t, err, "Case %d: %s", caseNum, currCase.name)

			assert.Equal(t, strings.Split(currCase.want, "\n"), strings.Split(string(got), "\n"))
		})
	}
}
