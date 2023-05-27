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
	"sort"
)

// node represents a node or vertex of a graph and is uniquely identified within the graph by a comparable id.
type node[T comparable] struct {
	id    T
	edges map[T]*node[T]
}

func (u *node[T]) addEdge(v *node[T]) *node[T] {
	u.edges[v.id] = v
	return u
}

func (u *node[T]) numEdges() int {
	return len(u.edges)
}

// sortedEdges returns the edges sorted by a comparator function.
// It is required to keep the results of the graph algorithms stable for tests.
func (u *node[T]) sortedEdges(less func(t1, t2 T) bool) []*node[T] {
	ret := make([]*node[T], 0, len(u.edges))
	for _, v := range u.edges {
		ret = append(ret, v)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return less(ret[i].id, ret[j].id)
	})
	return ret
}

type graph[T comparable] struct {
	nodes     []*node[T]
	nodesByID map[T]*node[T]
}

func newGraph[T comparable](size int) *graph[T] {
	if size > 0 {
		return &graph[T]{
			nodes:     make([]*node[T], 0, size),
			nodesByID: make(map[T]*node[T], size),
		}
	}
	return &graph[T]{
		nodesByID: make(map[T]*node[T]),
	}
}

func (g *graph[T]) addNode(id T) *graph[T] {
	u := &node[T]{
		id:    id,
		edges: make(map[T]*node[T]),
	}
	g.nodes = append(g.nodes, u)
	g.nodesByID[id] = u
	return g
}

func (g *graph[T]) addEdges(u *node[T], vs ...*node[T]) *graph[T] {
	for _, v := range vs {
		u.addEdge(v)
	}
	return g
}

func (g *graph[T]) addEdgesByID(idU T, idVs ...T) *graph[T] {
	if _, hasNode := g.nodesByID[idU]; !hasNode {
		g.addNode(idU)
	}
	u := g.nodesByID[idU]
	for _, idV := range idVs {
		if _, hasNode := g.nodesByID[idV]; !hasNode {
			g.addNode(idV)
		}
		v := g.nodesByID[idV]
		g.addEdges(u, v)
	}
	return g
}

func (g *graph[T]) numNodes() int {
	return len(g.nodes)
}

func (g *graph[T]) numEdges() int {
	cnt := 0
	for _, u := range g.nodes {
		cnt += u.numEdges()
	}
	return cnt
}

// reverseGraph builds the reverse graph of g. In other words, if there is an edge u->v in the original graph
// then there is v->u in the reverse graph.
func reverseGraph[T comparable](g *graph[T]) *graph[T] {
	revG := newGraph[T](len(g.nodes))

	for _, u := range g.nodes {
		revG.addNode(u.id)
	}

	for _, u := range g.nodes {
		for _, v := range u.sortedEdges(func(t1, t2 T) bool { return false }) {
			revG.addEdgesByID(v.id, u.id)
		}
	}

	return revG
}
