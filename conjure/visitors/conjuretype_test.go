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

package visitors

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/palantir/conjure-go/v4/conjure-api/conjure/spec"
	"github.com/palantir/conjure-go/v4/conjure/types"
)

func TestConjureTypeVisitor(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Type     spec.Type
		Expected string
	}{
		{
			Name:     "string",
			Type:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
			Expected: "string",
		},
		{
			Name:     "binary",
			Type:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeBinary),
			Expected: "[]byte",
		},
		{
			Name:     "boolean",
			Type:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeBoolean),
			Expected: "bool",
		},
		{
			Name:     "double",
			Type:     spec.NewTypeFromPrimitive(spec.PrimitiveTypeDouble),
			Expected: "float64",
		},
		{
			Name: "list<string>",
			Type: spec.NewTypeFromList(spec.ListType{
				ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
			}),
			Expected: "[]string",
		},
		{
			Name: "list<binary>",
			Type: spec.NewTypeFromList(spec.ListType{
				ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeBinary),
			}),
			Expected: "[][]byte",
		},
		{
			Name: "list<boolean>",
			Type: spec.NewTypeFromList(spec.ListType{
				ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeBoolean),
			}),
			Expected: "[]bool",
		},
		{
			Name: "list<double>",
			Type: spec.NewTypeFromList(spec.ListType{
				ItemType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeDouble),
			}),
			Expected: "[]float64",
		},
		{
			Name: "map<string, string>",
			Type: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
				ValueType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeString),
			}),
			Expected: "map[string]string",
		},
		{
			Name: "map<boolean, boolean>",
			Type: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.PrimitiveTypeBoolean),
				ValueType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeBoolean),
			}),
			Expected: "map[bool]bool",
		},
		{
			Name: "map<double, double>",
			Type: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.PrimitiveTypeDouble),
				ValueType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeDouble),
			}),
			Expected: "map[float64]float64",
		},
		{
			Name: "map<binary, binary>",
			Type: spec.NewTypeFromMap(spec.MapType{
				KeyType:   spec.NewTypeFromPrimitive(spec.PrimitiveTypeBinary),
				ValueType: spec.NewTypeFromPrimitive(spec.PrimitiveTypeBinary),
			}),
			Expected: "map[binary.Binary][]byte",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			info := types.NewPkgInfo("", nil)
			typer, err := NewConjureTypeProviderTyper(test.Type, info)
			require.NoError(t, err)
			require.Equal(t, test.Expected, typer.GoType(info))
		})
	}
}
