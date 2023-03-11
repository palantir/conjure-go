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
	"testing"

	"github.com/stretchr/testify/require"
)

// testGraph is the graph from the wikipedia animation at https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm
var testGraph = newGraph[int](8).
	addNode(1).
	addNode(2).
	addNode(3).
	addNode(4).
	addNode(5).
	addNode(6).
	addNode(7).
	addNode(8).
	addEdgesByID(1, 2).
	addEdgesByID(2, 3).
	addEdgesByID(3, 1).
	addEdgesByID(4, 2, 3, 5).
	addEdgesByID(5, 4, 6).
	addEdgesByID(6, 3, 7).
	addEdgesByID(7, 6).
	addEdgesByID(8, 5, 7, 8)

func TestCalculateStronglyConnectedComponents(t *testing.T) {
	expectedComponents := map[componentID][]int{
		0: {1, 2, 3},
		1: {6, 7},
		2: {4, 5},
		3: {8},
	}
	expectedComponentByItem := map[int]componentID{
		1: 0,
		2: 0,
		3: 0,
		4: 2,
		5: 2,
		6: 1,
		7: 1,
		8: 3,
	}
	expectedComponentGraph := newGraph[componentID](4).
		addNode(0).
		addNode(1).
		addNode(2).
		addNode(3).
		addEdgesByID(3, 1, 2).
		addEdgesByID(2, 0, 1).
		addEdgesByID(1, 0)

	sccs := calculateStronglyConnectedComponents(testGraph)
	require.Equal(t, expectedComponents, sccs.components)
	require.Equal(t, expectedComponentByItem, sccs.componentByItem)
	assertGraphsAreEqual(t, expectedComponentGraph, sccs.componentGraph)
}

func TestReverseGraph(t *testing.T) {
	expected := newGraph[int](8).
		addNode(1).
		addNode(2).
		addNode(3).
		addNode(4).
		addNode(5).
		addNode(6).
		addNode(7).
		addNode(8).
		addEdgesByID(1, 3).
		addEdgesByID(2, 1, 4).
		addEdgesByID(3, 2, 4, 6).
		addEdgesByID(4, 5).
		addEdgesByID(5, 4, 8).
		addEdgesByID(6, 5, 7).
		addEdgesByID(7, 6, 8).
		addEdgesByID(8, 8)
	actual := (&kosaraju[int]{
		comparator: func(t1, t2 int) bool {
			return t1 < t2
		},
	}).reverseGraph(testGraph)
	assertGraphsAreEqual(t, expected, actual)
}
