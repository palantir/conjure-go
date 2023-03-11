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

package cycles

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/palantir/conjure-go/v6/conjure-api/conjure/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertGraphsAreEqual[T comparable](t *testing.T, expected, actual *graph[T]) {
	assert.Equalf(t, expected.numNodes(), actual.numNodes(), "graphs do not have the same amount of nodes")
	assert.Equalf(t, expected.numEdges(), actual.numEdges(), "graphs do not have the same amount of edges")
	for _, u1 := range expected.nodes {
		u2, ok := actual.nodesByID[u1.id]
		assert.Truef(t, ok, "node %#v does not exist in graph", u1.id)
		if !ok {
			continue
		}
		assert.Equalf(t, u1.id, u2.id, "node %#v in graph has ID %#v", u1.id, u2.id)
		assert.Equalf(t, u1.numEdges(), u2.numEdges(), "node %#v does not have expected number of outgoing edges", u2.id)
		for _, v1 := range u1.edges {
			v2, ok := u2.edges[v1.id]
			assert.Truef(t, ok, "node %#v does not have edge to %#v", u2.id, v1.id)
			if !ok {
				continue
			}
			assert.Equalf(t, v1.id, v2.id, "node %#v has edge at key %#v mapped to %#v", u2.id, v1.id, v2.id)
		}
	}
}

func TestBuildTypeGraph(t *testing.T) {
	for _, testCase := range []struct {
		name             string
		conjureInputFile string
		expectedGraph    *graph[spec.TypeName]
	}{
		{
			name:             "no cycles",
			conjureInputFile: "testdata/no-cycles/in.conjure.json",
			expectedGraph: newGraph[spec.TypeName](11).
				addNode(spec.TypeName{Package: "com.palantir.errors", Name: "MyError"}).
				addNode(spec.TypeName{Package: "com.palantir.services", Name: "MyService"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type4"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.errors", Name: "MyError"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.services", Name: "MyService"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}),
		},
		{
			name:             "cycle within pkg",
			conjureInputFile: "testdata/cycle-within-pkg/in.conjure.json",
			expectedGraph: newGraph[spec.TypeName](11).
				addNode(spec.TypeName{Package: "com.palantir.errors", Name: "MyError"}).
				addNode(spec.TypeName{Package: "com.palantir.services", Name: "MyService"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type4"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.errors", Name: "MyError"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.services", Name: "MyService"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}),
		},
		{
			name:             "pkg cycle",
			conjureInputFile: "testdata/pkg-cycle/in.conjure.json",
			expectedGraph: newGraph[spec.TypeName](11).
				addNode(spec.TypeName{Package: "com.palantir.errors", Name: "MyError"}).
				addNode(spec.TypeName{Package: "com.palantir.services", Name: "MyService"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type4"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.errors", Name: "MyError"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.services", Name: "MyService"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}),
		},
		{
			name:             "type cycle",
			conjureInputFile: "testdata/type-cycle/in.conjure.json",
			expectedGraph: newGraph[spec.TypeName](11).
				addNode(spec.TypeName{Package: "com.palantir.errors", Name: "MyError"}).
				addNode(spec.TypeName{Package: "com.palantir.services", Name: "MyService"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.foo", Name: "Type4"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type2"}).
				addNode(spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addNode(spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addNode(spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.errors", Name: "MyError"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.services", Name: "MyService"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.fizz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type2"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type1"},
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"}).
				addEdgesByID(
					spec.TypeName{Package: "com.palantir.bar", Name: "Type3"},
					spec.TypeName{Package: "com.palantir.foo", Name: "Type4"},
					spec.TypeName{Package: "com.palantir.buzz", Name: "Type1"}),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := os.ReadFile(testCase.conjureInputFile)
			require.NoError(t, err)
			var def spec.ConjureDefinition
			err = json.Unmarshal(data, &def)
			require.NoError(t, err)
			actualGraph, err := buildTypeGraph(def)
			require.NoError(t, err)
			assertGraphsAreEqual(t, testCase.expectedGraph, actualGraph)
		})
	}
}

func TestTypeNamesWithinType(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		typ      []byte
		expected []spec.TypeName
	}{
		{
			name: "primitive",
			typ: []byte(`
{
  "type" : "primitive",
  "primitive" : "STRING"
}`),
			expected: []spec.TypeName{},
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
			expected: []spec.TypeName{
				{
					Name:    "Foo",
					Package: "com.palantir.package",
				},
			},
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
			expected: []spec.TypeName{
				{
					Name:    "Foo",
					Package: "com.palantir.package",
				},
			},
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
			expected: []spec.TypeName{
				{
					Name:    "Foo",
					Package: "com.palantir.package",
				},
			},
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
			expected: []spec.TypeName{
				{
					Name:    "Foo",
					Package: "com.palantir.package",
				},
			},
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
			expected: []spec.TypeName{
				{
					Name:    "Key",
					Package: "com.palantir.map",
				},
				{
					Name:    "Value",
					Package: "com.palantir.map",
				},
			},
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
			expected: []spec.TypeName{
				{
					Name:    "Foo",
					Package: "com.palantir.external.package",
				},
				{
					Name:    "Foo",
					Package: "com.palantir.package",
				},
			},
		},
		/*		{
							name: "object",
							typ: []byte(`
				{
				  "type" : "object",
				  "object" : {
				    "typeName" : {
				      "name" : "ObjectName",
				      "package" : "com.palantir.package"
				    },
				    "fields" : [ {
				      "fieldName" : "foo",
				      "type" : {
				        "type" : "primitive",
				        "primitive" : "STRING"
				      }
				    }, {
				      "fieldName" : "bar",
				      "type" : {
				        "type" : "list",
				        "list" : {
				          "itemType" : {
				            "type" : "reference",
				            "reference" : {
				              "name" : "Bar",
				              "package" : "com.palantir.package"
				            }
				          }
				        }
				      }
				    }, {
				      "fieldName" : "map",
				      "type" : {
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
				      }
				    } ]
				  }
				}`),
							expected: []spec.TypeName{
								{
									Name:    "Key",
									Package: "com.palantir.map",
								},
								{
									Name:    "Value",
									Package: "com.palantir.map",
								},
								{
									Name:    "Bar",
									Package: "com.palantir.package",
								},
								{
									Name:    "ObjectName",
									Package: "com.palantir.package",
								},
							},
						},*/
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var typ spec.Type
			err := json.Unmarshal(testCase.typ, &typ)
			require.NoError(t, err)
			actual, err := typeNamesWithinType(typ)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestDedup(t *testing.T) {
	input := []spec.TypeName{
		{
			Package: "com.palantir.spec",
			Name:    "Foo",
		},
		{
			Package: "com.palantir.status",
			Name:    "Foo",
		},
		{
			Package: "com.palantir.spec",
			Name:    "Foo",
		},
		{
			Package: "com.palantir.status",
			Name:    "Bar",
		},
		{
			Package: "com.palantir.status",
			Name:    "Foo",
		},
	}
	expected := []spec.TypeName{
		{
			Package: "com.palantir.spec",
			Name:    "Foo",
		},
		{
			Package: "com.palantir.status",
			Name:    "Bar",
		},
		{
			Package: "com.palantir.status",
			Name:    "Foo",
		},
	}
	actual := dedup(input)
	assert.Equal(t, expected, actual)
}
