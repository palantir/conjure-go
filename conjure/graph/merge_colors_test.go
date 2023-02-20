package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeColors(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		graph     *Graph[int]
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
			actual := mergeColors(testCase.graph, testCase.colorByID)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
