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
	"testing"

	"github.com/stretchr/testify/require"
)

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
