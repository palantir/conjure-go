package graph

import (
	"sort"
)

type Node[T comparable] struct {
	ID    T
	Edges map[T]*Node[T]
}

func (u *Node[T]) addEdge(v *Node[T]) *Node[T] {
	u.Edges[v.ID] = v
	return u
}

func (u *Node[T]) numEdges() int {
	return len(u.Edges)
}

// sortedEdges returns the edges sorted by a comparator function.
// It is required to keep the results of the graph algorithms stable.
func (u *Node[T]) sortedEdges(less func(t1, t2 T) bool) []*Node[T] {
	ret := make([]*Node[T], 0, len(u.Edges))
	for _, u := range u.Edges {
		ret = append(ret, u)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return less(ret[i].ID, ret[j].ID)
	})
	return ret
}

type Graph[T comparable] struct {
	Nodes     []*Node[T]
	NodesByID map[T]*Node[T]
}

func newGraph[T comparable](sz int) *Graph[T] {
	if sz > 0 {
		return &Graph[T]{
			Nodes:     make([]*Node[T], 0, sz),
			NodesByID: make(map[T]*Node[T], sz),
		}
	}
	return &Graph[T]{
		NodesByID: make(map[T]*Node[T]),
	}
}

func (g *Graph[T]) addNode(id T) *Graph[T] {
	u := &Node[T]{
		ID:    id,
		Edges: make(map[T]*Node[T]),
	}
	g.Nodes = append(g.Nodes, u)
	g.NodesByID[id] = u
	return g
}

func (g *Graph[T]) addEdges(u *Node[T], vs ...*Node[T]) *Graph[T] {
	for _, v := range vs {
		u.addEdge(v)
	}
	return g
}

func (g *Graph[T]) addEdgesByID(idU T, idVs ...T) *Graph[T] {
	if _, hasNode := g.NodesByID[idU]; !hasNode {
		g.addNode(idU)
	}
	u := g.NodesByID[idU]
	for _, idV := range idVs {
		if _, hasNode := g.NodesByID[idV]; !hasNode {
			g.addNode(idV)
		}
		v := g.NodesByID[idV]
		g.addEdges(u, v)
	}
	return g
}

func (g *Graph[T]) numNodes() int {
	return len(g.Nodes)
}

func (g *Graph[T]) numEdges() int {
	cnt := 0
	for _, u := range g.Nodes {
		cnt += u.numEdges()
	}
	return cnt
}
