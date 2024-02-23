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

// partition takes in a directed acyclic graph (nodes of type T) where each node has a color (color of type V).
// It attempts to merge as many nodes as possible in a stable manner such that two nodes can be merged
// if and only if they have the same color and merging them won't introduce a cycle.
func partition[T, V comparable](g *graph[T], colorByID map[T]V) map[V][][]T {
	return (&partitioner[T, V]{g: g, colorByID: colorByID}).run()
}

type partitioner[T, V comparable] struct {
	// partitioner input
	g         *graph[T]
	colorByID map[T]V

	// partitioner internal variables
	numNodes     int
	revG         *graph[T]
	idToBit      map[T]bitID
	disallowed   map[T]bitset // For each node, the nodes that, if merged, would cause a cycle.
	dependencies map[T]bitset // For each node u, the nodes that can be reached by u. Contains u.
	group        map[T]bitset // For each node u, the nodes that are in the same grouping as u. Contains u.
	revToposort  []T

	// used for sorting edges to keep algorithm stable.
	comparator func(t1, t2 T) bool
}

func (p *partitioner[T, V]) run() map[V][][]T {
	p.numNodes = len(p.g.nodes)
	p.revG = reverseGraph(p.g)

	// Step 1: Represent all nodes as a bitID. We'll use bitsets to efficiently represent dependencies in O(n^2/64).
	p.idToBit = make(map[T]bitID, p.numNodes)
	for i, u := range p.g.nodes {
		p.idToBit[u.id] = bitID(i)
	}
	p.comparator = func(t1, t2 T) bool {
		return p.idToBit[t1] < p.idToBit[t2]
	}

	// Step 2: For each node, calculate its total dependencies and disallowed dependencies.
	// A dependency is disallowed if and only if:
	// - Dependency has a different color
	// - At least one node of a different color exists in a the path from the node to the dependency.
	// Also while we're at it, calculate a reversed topological sort.
	visited := make(map[T]struct{}, p.numNodes)
	p.disallowed = make(map[T]bitset, p.numNodes)
	p.dependencies = make(map[T]bitset, p.numNodes)
	p.group = make(map[T]bitset, p.numNodes)
	p.revToposort = make([]T, 0, p.numNodes)
	for _, u := range p.g.nodes {
		p.dfs(u, visited)
	}

	// Step 3: Traverse nodes and try to merge into one of the groups or create a new group if not possible.
	// Do it in reverse topological order to prioritize groups that have less dependencies.
	// This is the same as grouping nodes by their depth (biggest path starting at the node).
	groupsByColor := make(map[V][][]T)
	for _, id := range p.revToposort {
		color := p.colorByID[id]
		merged := false
		for groupIdx, group := range groupsByColor[color] {
			// Since we merge all restrictions of the group into the first one node of it, we can check if we can merge
			// by just checking the first one.
			leader := group[0]
			if p.canMerge(id, leader) {
				p.merge(id, leader)
				group = append(group, id)
				groupsByColor[color][groupIdx] = group
				merged = true
				break
			}
		}
		if !merged {
			groupsByColor[color] = append(groupsByColor[color], []T{id})
		}
	}
	return groupsByColor
}

func (p *partitioner[T, V]) dfs(u *node[T], visited map[T]struct{}) {
	if _, alreadyVisited := visited[u.id]; alreadyVisited {
		return
	}
	visited[u.id] = struct{}{}

	p.disallowed[u.id] = newBitset(p.numNodes)
	dependencies := newBitset(p.numNodes)
	dependencies.add(p.idToBit[u.id])
	p.dependencies[u.id] = dependencies
	group := newBitset(p.numNodes)
	group.add(p.idToBit[u.id])
	p.group[u.id] = group

	for _, v := range u.sortedEdges(p.comparator) {
		p.dfs(v, visited)
		p.processDependency(u.id, v.id)
	}

	p.revToposort = append(p.revToposort, u.id)
}

func (p *partitioner[T, V]) processDependency(uID, vID T) {
	// fmt.Printf("Processing dependency %v -> %v\n", uID, vID)
	p.dependencies[uID] = p.dependencies[uID].merge(p.dependencies[vID])

	// If the node v (dependency of u) can't merge with node w (dependency of v),
	// then node u also can't merge with w (indirect dependency of u) because either:
	// - u and w have different colors
	// - u and w have the same color and both have different colors than v: in this case there is a color cycle
	// - u and w are the same color which is also the same as v: since w is disallowed to v there exists at least
	//   one path from v to w that has a node of a different color and this path is a sub-path of u to w.
	// Therefore, we can just merge the v's disallowed set into u's.
	p.disallowed[uID] = p.disallowed[uID].merge(p.disallowed[vID])
	// If color of u and v have different colors, then u and v can't be merged.
	if p.colorByID[uID] != p.colorByID[vID] {
		p.disallowed[uID] = p.disallowed[uID].merge(p.dependencies[vID])
	}
}

func (p *partitioner[T, V]) canMerge(id1, id2 T) bool {
	// Can't merge two nodes if the colors are different
	if p.colorByID[id1] != p.colorByID[id2] {
		return false
	}
	// Check if one is a disallowed dependency of the other.
	return !p.disallowed[id1].intersects(p.group[id2]) && !p.disallowed[id2].intersects(p.group[id1])
}

func (p *partitioner[T, V]) merge(id1, id2 T) {
	u := p.revG.nodesByID[id1]
	v := p.revG.nodesByID[id2]
	u.addEdge(v)
	v.addEdge(u)
	p.revDfs(u, v, make(map[T]struct{}))
	p.revDfs(v, u, make(map[T]struct{}))
	p.group[u.id] = p.group[u.id].merge(p.group[v.id])
	p.group[v.id] = p.group[u.id]
}

func (p *partitioner[T, V]) revDfs(u *node[T], newDep *node[T], visited map[T]struct{}) {
	if _, alreadyVisited := visited[u.id]; alreadyVisited {
		return
	}
	visited[u.id] = struct{}{}

	if p.dependencies[u.id].has(p.idToBit[newDep.id]) {
		return
	}
	p.processDependency(u.id, newDep.id)

	for _, v := range u.sortedEdges(p.comparator) {
		p.revDfs(v, newDep, visited)
	}
}
