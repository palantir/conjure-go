package graph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertGraphsAreEqual[T comparable](t *testing.T, expected, actual *Graph[T]) {
	assert.Equalf(t, expected.numNodes(), actual.numNodes(), "graphs do not have the same amount of nodes")
	assert.Equalf(t, expected.numEdges(), actual.numEdges(), "graphs do not have the same amount of edges")
	for _, u1 := range expected.Nodes {
		u2, ok := actual.NodesByID[u1.ID]
		assert.Truef(t, ok, "node %#v does not exist in graph", u1.ID)
		if !ok {
			continue
		}
		assert.Equalf(t, u1.ID, u2.ID, "node %#v in graph has ID %#v", u1.ID, u2.ID)
		assert.Equalf(t, u1.numEdges(), u2.numEdges(), "node %#v does not have expected number of outgoing edges", u2.ID)
		for _, v1 := range u1.Edges {
			v2, ok := u2.Edges[v1.ID]
			assert.Truef(t, ok, "node %#v does not have edge to %#v", u2.ID, v1.ID)
			if !ok {
				continue
			}
			assert.Equalf(t, v1.ID, v2.ID, "node %#v has edge at key %#v mapped to %#v", u2.ID, v1.ID, v2.ID)
		}
	}
}

func printGraph[T comparable](g *Graph[T]) {
	for _, u := range g.Nodes {
		fmt.Printf("Node %#v:\n", u.ID)
		for _, v := range u.Edges {
			fmt.Printf("  -> Node %#v\n", v.ID)
		}
	}
}
