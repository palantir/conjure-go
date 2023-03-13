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

// componentID is an identifier for strongly connected components.
type componentID int

type stronglyConnectedComponents[T comparable] struct {
	components      map[componentID][]T
	componentByItem map[T]componentID
	componentGraph  *graph[componentID]
}

// calculateStronglyConnectedComponents takes in a graph and groups nodes into "strongly connected components" (SCCs).
// Two nodes are in the same SCC if and only if they are reachable from each other (i.e. they are part of a cycle).
// The graph of SCCs is a directed acyclic graph (DAG).
func calculateStronglyConnectedComponents[T comparable](g *graph[T]) *stronglyConnectedComponents[T] {
	return (&kosaraju[T]{g: g}).run()
}

type kosaraju[T comparable] struct {
	// input to kosaraju's algorithm
	g *graph[T]

	// output of kosaraju's algorithm
	componentByItem map[T]componentID

	// variables internal to the algorithm
	revG          *graph[T]
	visited       map[T]struct{}
	curComponent  componentID
	numComponents int
	revToposort   []T

	// used for sorting edges to keep algorithm stable for tests.
	comparator func(t1, t2 T) bool
}

// run kosaraju's algorithm for finding strongly connected components
// https://www.geeksforgeeks.org/strongly-connected-components/
// https://en.wikipedia.org/wiki/Kosaraju%27s_algorithm
func (k *kosaraju[T]) run() *stronglyConnectedComponents[T] {
	n := len(k.g.nodes)

	index := make(map[T]int, n)
	for i, u := range k.g.nodes {
		index[u.id] = i
	}
	k.comparator = func(t1, t2 T) bool {
		return index[t1] < index[t2]
	}

	k.revG = k.reverseGraph(k.g)
	k.componentByItem = make(map[T]componentID, n)

	k.revToposort = make([]T, 0, n)
	k.visited = make(map[T]struct{})
	for _, u := range k.revG.nodes {
		if _, visited := k.visited[u.id]; !visited {
			k.revDfs(u)
		}
	}

	k.visited = make(map[T]struct{})
	k.numComponents = 0
	for i := n - 1; i >= 0; i-- {
		u := k.g.nodesByID[k.revToposort[i]]
		if _, visited := k.visited[u.id]; !visited {
			k.curComponent = componentID(k.numComponents)
			k.numComponents++
			k.dfs(u)
		}
	}

	components := make(map[componentID][]T)
	for id, componentID := range k.componentByItem {
		components[componentID] = append(components[componentID], id)
	}
	for componentID := range components {
		ids := components[componentID]
		sort.SliceStable(ids, func(i, j int) bool {
			return k.comparator(ids[i], ids[j])
		})
		components[componentID] = ids
	}

	return &stronglyConnectedComponents[T]{
		components:      components,
		componentByItem: k.componentByItem,
		componentGraph:  k.buildComponentGraph(components),
	}
}

// reverseGraph builds the reverse graph of g. In other words, if there is an edge u->v in the original graph
// if and only if there is v->u in the reverse graph.
func (k *kosaraju[T]) reverseGraph(g *graph[T]) *graph[T] {
	revG := newGraph[T](len(g.nodes))

	for _, u := range g.nodes {
		revG.addNode(u.id)
	}

	for _, u := range g.nodes {
		for _, v := range u.sortedEdges(k.comparator) {
			revG.addEdgesByID(v.id, u.id)
		}
	}

	return revG
}

func (k *kosaraju[T]) revDfs(u *node[T]) {
	k.visited[u.id] = struct{}{}
	for _, v := range u.sortedEdges(k.comparator) {
		if _, visited := k.visited[v.id]; !visited {
			k.revDfs(v)
		}
	}
	k.revToposort = append(k.revToposort, u.id)
}

func (k *kosaraju[T]) dfs(u *node[T]) {
	k.visited[u.id] = struct{}{}
	k.componentByItem[u.id] = k.curComponent
	for _, v := range u.sortedEdges(k.comparator) {
		if _, visited := k.visited[v.id]; !visited {
			k.dfs(v)
		}
	}
}

func (k *kosaraju[T]) buildComponentGraph(components map[componentID][]T) *graph[componentID] {
	compG := newGraph[componentID](k.numComponents)

	for compID := range components {
		compG.addNode(compID)
	}
	sort.SliceStable(compG.nodes, func(i, j int) bool {
		return compG.nodes[i].id < compG.nodes[j].id
	})

	for _, u := range k.g.nodes {
		uCompID := k.componentByItem[u.id]
		for _, v := range u.sortedEdges(k.comparator) {
			vCompID := k.componentByItem[v.id]
			if uCompID != vCompID {
				compG.addEdgesByID(uCompID, vCompID)
			}
		}
	}

	return compG
}
