// Copyright (c) 2023 Palantir Technologies. All rights reserved.
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

package cycles

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemovePackageCycles(t *testing.T) {
	for _, testCase := range []struct {
		name              string
		conjureInputFile  string
		conjureOutputFile string
	}{
		{
			name:              "no cycles",
			conjureInputFile:  "testdata/no-cycles/in.conjure.json",
			conjureOutputFile: "testdata/no-cycles/out.conjure.json",
		},
		{
			name:              "cycle within package",
			conjureInputFile:  "testdata/cycle-within-pkg/in.conjure.json",
			conjureOutputFile: "testdata/cycle-within-pkg/out.conjure.json",
		},
		{
			name:              "pkg cycle",
			conjureInputFile:  "testdata/pkg-cycle/in.conjure.json",
			conjureOutputFile: "testdata/pkg-cycle/out.conjure.json",
		},
		{
			name:              "pkg cycle disconnected",
			conjureInputFile:  "testdata/pkg-cycle-disconnected/in.conjure.json",
			conjureOutputFile: "testdata/pkg-cycle-disconnected/out.conjure.json",
		},
		{
			name:              "type cycle",
			conjureInputFile:  "testdata/type-cycle/in.conjure.json",
			conjureOutputFile: "testdata/type-cycle/out.conjure.json",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			inData, err := os.ReadFile(testCase.conjureInputFile)
			require.NoError(t, err)
			var def spec.ConjureDefinition
			err = json.Unmarshal(inData, &def)
			require.NoError(t, err)

			outData, err := os.ReadFile(testCase.conjureOutputFile)
			require.NoError(t, err)
			var expected spec.ConjureDefinition
			err = json.Unmarshal(outData, &expected)
			require.NoError(t, err)

			actual, err := RemovePackageCycles(def)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestMergeTypesIntoSamePackage(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		types    []spec.TypeName
		expected map[spec.TypeName]spec.TypeName
	}{
		{
			name: "no op when all in same package",
			types: []spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"},
			},
			expected: map[spec.TypeName]spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"}: {Package: "com.palantir.foo", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type2"}: {Package: "com.palantir.foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"}: {Package: "com.palantir.foo", Name: "Type3"},
			},
		},
		{
			name: "merges packages without name conflicts",
			types: []spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"},
				{Package: "com.palantir.bar", Name: "Type4"},
				{Package: "com.palantir.bar", Name: "Type5"},
			},
			expected: map[spec.TypeName]spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"}: {Package: "com.palantir.bar_foo", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type2"}: {Package: "com.palantir.bar_foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"}: {Package: "com.palantir.bar_foo", Name: "Type3"},
				{Package: "com.palantir.bar", Name: "Type4"}: {Package: "com.palantir.bar_foo", Name: "Type4"},
				{Package: "com.palantir.bar", Name: "Type5"}: {Package: "com.palantir.bar_foo", Name: "Type5"},
			},
		},
		{
			name: "merges packages with name conflicts",
			types: []spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"},
				{Package: "com.palantir.bar", Name: "Type3"},
				{Package: "com.palantir.bar", Name: "Type5"},
				{Package: "com.palantir.baz", Name: "Type3"},
			},
			expected: map[spec.TypeName]spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"}: {Package: "com.palantir.bar_baz_foo", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type2"}: {Package: "com.palantir.bar_baz_foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"}: {Package: "com.palantir.bar_baz_foo", Name: "FooType3"},
				{Package: "com.palantir.bar", Name: "Type3"}: {Package: "com.palantir.bar_baz_foo", Name: "BarType3"},
				{Package: "com.palantir.bar", Name: "Type5"}: {Package: "com.palantir.bar_baz_foo", Name: "Type5"},
				{Package: "com.palantir.baz", Name: "Type3"}: {Package: "com.palantir.bar_baz_foo", Name: "BazType3"},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual := make(map[spec.TypeName]spec.TypeName)
			mergeTypesIntoSamePackage(testCase.types, actual, 0)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestMergePackages(t *testing.T) {
	for _, testCase := range []struct {
		name                 string
		packageSet           packageSet
		numSimilarPackageSet int
		expected             string
	}{
		{
			name: "single package",
			packageSet: packageSet{
				"com.palantir.foo": {},
			},
			numSimilarPackageSet: 0,
			expected:             "com.palantir.foo",
		},
		{
			name: "single package but name already used",
			packageSet: packageSet{
				"com.palantir.foo": {},
			},
			numSimilarPackageSet: 2,
			expected:             "com.palantir.foo2",
		},
		{
			name: "two packages",
			packageSet: packageSet{
				"com.palantir.foo": {},
				"com.palantir.bar": {},
			},
			numSimilarPackageSet: 0,
			expected:             "com.palantir.bar_foo",
		},
		{
			name: "two packages but name already used",
			packageSet: packageSet{
				"com.palantir.foo": {},
				"com.palantir.bar": {},
			},
			numSimilarPackageSet: 2,
			expected:             "com.palantir.bar_foo2",
		},
		{
			name: "packages with nested dependencies",
			packageSet: packageSet{
				"com.palantir.foo":      {},
				"com.palantir.root.bar": {},
				"com.palantir.root.baz": {},
			},
			numSimilarPackageSet: 0,
			expected:             "com.palantir.foo_rootbar_rootbaz",
		},
		{
			name: "packages with nested dependencies but name already used",
			packageSet: packageSet{
				"com.palantir.foo":      {},
				"com.palantir.root.bar": {},
				"com.palantir.root.baz": {},
			},
			numSimilarPackageSet: 2,
			expected:             "com.palantir.foo_rootbar_rootbaz2",
		},
		{
			name: "packages where one is a prefix of the other",
			packageSet: packageSet{
				"com.palantir.root":     {},
				"com.palantir.root.bar": {},
				"com.palantir.root.baz": {},
			},
			numSimilarPackageSet: 0,
			expected:             "com.palantir.root_rootbar_rootbaz",
		},
		{
			name: "packages where one is a prefix of the other but name already used",
			packageSet: packageSet{
				"com.palantir.root":     {},
				"com.palantir.root.bar": {},
				"com.palantir.root.baz": {},
			},
			numSimilarPackageSet: 2,
			expected:             "com.palantir.root_rootbar_rootbaz2",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual := mergePackages(testCase.packageSet, testCase.numSimilarPackageSet)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestApplyTransformToDef(t *testing.T) {
	for _, testCase := range []struct {
		name              string
		conjureInputFile  string
		conjureOutputFile string
		typeTransform     map[spec.TypeName]spec.TypeName
	}{
		{
			name:              "no cycles",
			conjureInputFile:  "testdata/no-cycles/in.conjure.json",
			conjureOutputFile: "testdata/no-cycles/out.conjure.json",
			typeTransform:     map[spec.TypeName]spec.TypeName{},
		},
		{
			name:              "cycle within pkg",
			conjureInputFile:  "testdata/cycle-within-pkg/in.conjure.json",
			conjureOutputFile: "testdata/cycle-within-pkg/out.conjure.json",
			typeTransform:     map[spec.TypeName]spec.TypeName{},
		},
		{
			name:              "pkg cycle",
			conjureInputFile:  "testdata/pkg-cycle/in.conjure.json",
			conjureOutputFile: "testdata/pkg-cycle/out.conjure.json",
			typeTransform: map[spec.TypeName]spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type1"}: {Package: "com.palantir.foo1", Name: "Type1"},
				{Package: "com.palantir.foo", Name: "Type3"}: {Package: "com.palantir.foo1", Name: "Type3"},
			},
		},
		{
			name:              "type cycle",
			conjureInputFile:  "testdata/type-cycle/in.conjure.json",
			conjureOutputFile: "testdata/type-cycle/out.conjure.json",
			typeTransform: map[spec.TypeName]spec.TypeName{
				{Package: "com.palantir.foo", Name: "Type2"}: {Package: "com.palantir.bar_foo", Name: "Type2"},
				{Package: "com.palantir.foo", Name: "Type3"}: {Package: "com.palantir.bar_foo", Name: "FooType3"},
				{Package: "com.palantir.foo", Name: "Type4"}: {Package: "com.palantir.bar_foo", Name: "Type4"},
				{Package: "com.palantir.bar", Name: "Type1"}: {Package: "com.palantir.bar_foo", Name: "Type1"},
				{Package: "com.palantir.bar", Name: "Type3"}: {Package: "com.palantir.bar_foo", Name: "BarType3"},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			inData, err := os.ReadFile(testCase.conjureInputFile)
			require.NoError(t, err)
			var def spec.ConjureDefinition
			err = json.Unmarshal(inData, &def)
			require.NoError(t, err)

			outData, err := os.ReadFile(testCase.conjureOutputFile)
			require.NoError(t, err)
			var expected spec.ConjureDefinition
			err = json.Unmarshal(outData, &expected)
			require.NoError(t, err)

			actual, err := applyTypeTransformToDef(def, func(typeName spec.TypeName) (spec.TypeName, error) {
				newTypeName, ok := testCase.typeTransform[typeName]
				if !ok {
					return typeName, nil
				}
				return newTypeName, nil
			})
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestApplyTransformToType(t *testing.T) {
	typeTransform := func(typeName spec.TypeName) (spec.TypeName, error) {
		return spec.TypeName{
			Name:    typeName.Name + "Extra",
			Package: typeName.Package + ".extra",
		}, nil
	}

	for _, testCase := range []struct {
		name     string
		typ      []byte
		expected []byte
	}{
		{
			name: "primitive",
			typ: []byte(`
{
  "type" : "primitive",
  "primitive" : "STRING"
}`),
			expected: []byte(`
{
  "type" : "primitive",
  "primitive" : "STRING"
}`),
		},
		{
			name: "reference",
			typ: []byte(`
{
  "type" : "reference",
  "reference" : {
    "name" : "Foo",
    "package" : "com.palantir.package"
  }
}`),
			expected: []byte(`
{
  "type" : "reference",
  "reference" : {
    "name" : "FooExtra",
    "package" : "com.palantir.package.extra"
  }
}`),
		},
		{
			name: "optional",
			typ: []byte(`
{
  "type" : "optional",
  "optional" : {
    "itemType" : {
      "type" : "reference",
      "reference" : {
        "name" : "Foo",
        "package" : "com.palantir.package"
      }
    }
  }
}`),
			expected: []byte(`
{
  "type" : "optional",
  "optional" : {
    "itemType" : {
      "type" : "reference",
      "reference" : {
        "name" : "FooExtra",
        "package" : "com.palantir.package.extra"
      }
    }
  }
}`),
		},
		{
			name: "list",
			typ: []byte(`
{
  "type" : "list",
  "list" : {
    "itemType" : {
      "type" : "reference",
      "reference" : {
        "name" : "Foo",
        "package" : "com.palantir.package"
      }
    }
  }
}`),
			expected: []byte(`
{
  "type" : "list",
  "list" : {
    "itemType" : {
      "type" : "reference",
      "reference" : {
        "name" : "FooExtra",
        "package" : "com.palantir.package.extra"
      }
    }
  }
}`),
		},
		{
			name: "set",
			typ: []byte(`
{
  "type" : "set",
  "set" : {
    "itemType" : {
      "type" : "reference",
      "reference" : {
        "name" : "Foo",
        "package" : "com.palantir.package"
      }
    }
  }
}`),
			expected: []byte(`
{
  "type" : "set",
  "set" : {
    "itemType" : {
      "type" : "reference",
      "reference" : {
        "name" : "FooExtra",
        "package" : "com.palantir.package.extra"
      }
    }
  }
}`),
		},
		{
			name: "map",
			typ: []byte(`
{
  "type" : "map",
  "map" : {
    "keyType" : {
      "type" : "reference",
      "reference" : {
        "name" : "Key",
        "package" : "com.palantir.map"
      }
    },
    "valueType" : {
      "type" : "reference",
      "reference" : {
        "name" : "Value",
        "package" : "com.palantir.map"
      }
    }
  }
}`),
			expected: []byte(`
{
  "type" : "map",
  "map" : {
    "keyType" : {
      "type" : "reference",
      "reference" : {
        "name" : "KeyExtra",
        "package" : "com.palantir.map.extra"
      }
    },
    "valueType" : {
      "type" : "reference",
      "reference" : {
        "name" : "ValueExtra",
        "package" : "com.palantir.map.extra"
      }
    }
  }
}`),
		},
		{
			name: "external",
			typ: []byte(`
{
  "type" : "external",
  "external" : {
    "externalReference" : {
      "name" : "Foo",
      "package" : "com.palantir.external.package"
    },
    "fallback" : {
      "type" : "reference",
      "reference" : {
        "name" : "Foo",
        "package" : "com.palantir.package"
      }
    }
  }
}`),
			expected: []byte(`
{
  "type" : "external",
  "external" : {
    "externalReference" : {
      "name" : "FooExtra",
      "package" : "com.palantir.external.package.extra"
    },
    "fallback" : {
      "type" : "reference",
      "reference" : {
        "name" : "FooExtra",
        "package" : "com.palantir.package.extra"
      }
    }
  }
}`),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var typ spec.Type
			err := json.Unmarshal(testCase.typ, &typ)
			require.NoError(t, err)
			var expected spec.Type
			err = json.Unmarshal(testCase.expected, &expected)
			require.NoError(t, err)
			actual, err := applyTypeTransformToType(typ, typeTransform)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}
