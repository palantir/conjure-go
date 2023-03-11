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

	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		graph     *graph[int]
		colorByID map[int]string
		expected  map[string][][]int
	}{
		{
			name: "all nodes of the same color are connected and can be merged",
			graph: newGraph[int](8).
				addNode(1).
				addNode(2).
				addNode(3).
				addNode(4).
				addNode(5).
				addNode(6).
				addNode(7).
				addNode(8).
				addEdgesByID(1, 2, 4).
				addEdgesByID(2, 3, 7).
				addEdgesByID(3, 8).
				addEdgesByID(4, 6, 7).
				addEdgesByID(6, 5, 8).
				addEdgesByID(7, 8),
			colorByID: map[int]string{
				1: "red",
				2: "red",
				3: "red",
				4: "blue",
				5: "blue",
				6: "blue",
				7: "green",
				8: "green",
			},
			expected: map[string][][]int{
				"red": {
					{3, 2, 1},
				},
				"blue": {
					{5, 6, 4},
				},
				"green": {
					{8, 7},
				},
			},
		},
		{
			name: "nodes of the same color may be disconnected but they still can be merged",
			graph: newGraph[int](8).
				addNode(1).
				addNode(2).
				addNode(3).
				addNode(4).
				addNode(5).
				addNode(6).
				addNode(7).
				addNode(8).
				addEdgesByID(1, 3, 4).
				addEdgesByID(3, 6, 7).
				addEdgesByID(4, 7).
				addEdgesByID(5, 8).
				addEdgesByID(7, 8),
			colorByID: map[int]string{
				1: "red",
				2: "red",
				3: "red",
				4: "blue",
				5: "blue",
				6: "blue",
				7: "green",
				8: "green",
			},
			expected: map[string][][]int{
				"red": {
					{3, 1, 2},
				},
				"blue": {
					{6, 4, 5},
				},
				"green": {
					{8, 7},
				},
			},
		},
		{
			name: "nodes of the same color can't be merged due to indirect dependencies",
			graph: newGraph[int](8).
				addNode(1).
				addNode(2).
				addNode(3).
				addNode(4).
				addNode(5).
				addNode(6).
				addNode(7).
				addNode(8).
				addEdgesByID(1, 2, 3).
				addEdgesByID(2, 4, 6, 7).
				addEdgesByID(3, 7).
				addEdgesByID(4, 5, 8).
				addEdgesByID(6, 7).
				addEdgesByID(7, 8),
			colorByID: map[int]string{
				1: "red",
				2: "red",
				3: "blue",
				4: "green",
				5: "green",
				6: "red",
				7: "red",
				8: "red",
			},
			expected: map[string][][]int{
				"red": {
					{8, 7, 6},
					{2, 1},
				},
				"blue": {
					{3},
				},
				"green": {
					{5, 4},
				},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual := partition(testCase.graph, testCase.colorByID)
			assert.Equal(t, len(testCase.expected), len(actual))
			for k := range testCase.expected {
				_, ok := actual[k]
				assert.True(t, ok)
				assert.ElementsMatch(t, testCase.expected[k], actual[k])
			}
		})
	}
}
