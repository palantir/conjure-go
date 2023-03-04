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

// mergeColors takes in a directed acyclic graph (nodes of type T) where each node has a color (color of type V).
// It attempts to merge as many nodes as possible in a stable manner such that two nodes can be merged
// if and only if they have the same color and merging them won't introduce a cycle.
func mergeColors[T, V comparable](g *graph[T], colorByID map[T]V) map[V][][]T {
	return (&merger[T, V]{g: g, colorByID: colorByID}).run()
}

type merger[T, V comparable] struct {
	// merger input
	g         *graph[T]
	colorByID map[T]V

	// merger internal variables
	numNodes       int
	idToBit        map[T]bitID
	disallowedBits map[bitID]bitset // For each node, the nodes that, if merged, would cause a cycle.
	dependencies   map[bitID]bitset // For each node u, the nodes that can be reached by u.
	visited        map[T]struct{}   // Visit each node once in the dfs
	comparator     func(id1, id2 T) bool
	revToposort    []T
}

func (m *merger[T, V]) run() map[V][][]T {
	m.numNodes = len(m.g.nodes)

	// Step 1: Represent all nodes as a bitID. We'll use bitsets to efficiently represent dependencies in O(n^2/64).
	m.idToBit = make(map[T]bitID, m.numNodes)
	for i, u := range m.g.nodes {
		m.idToBit[u.id] = bitID(i)
	}
	m.comparator = func(id1, id2 T) bool {
		return m.idToBit[id1] < m.idToBit[id2]
	}

	// Step 2: For each node, calculate from its dependency subgraph the set of nodes that cannot be merged in because:
	// - Dependency has a different color
	// - At least one node of a different color in the path from the node to the dependency.
	// Also while we're at it, calculate a reversed topological sort.
	m.visited = make(map[T]struct{}, m.numNodes)
	m.disallowedBits = make(map[bitID]bitset, m.numNodes)
	m.dependencies = make(map[bitID]bitset, m.numNodes)
	m.revToposort = make([]T, 0, m.numNodes)
	for _, u := range m.g.nodes {
		m.dfs(u)
	}

	// Step 3: Traverse nodes and try to merge into one of the groupings or create a new grouping if not possible.
	// Do it in reverse topological order to prioritize groups that have less dependencies.
	groupsByColor := make(map[V][][]T)
	for _, id := range m.revToposort {
		color := m.colorByID[id]
		merged := false
		for groupIdx, group := range groupsByColor[color] {
			// Since we merge all restrictions of the group into the first one, we can check if we can merge
			// by just checking the first one.
			leader := group[0]
			if m.canMerge(id, leader) {
				m.merge(id, leader)
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

func (m *merger[T, V]) dfs(u *node[T]) {
	if _, alreadyVisited := m.visited[u.id]; alreadyVisited {
		return
	}
	m.visited[u.id] = struct{}{}

	disallowed := newBitset(m.numNodes)
	dependencies := newBitset(m.numNodes)
	for _, v := range u.sortedEdges(m.comparator) {
		m.dfs(v)

		vBitID := m.idToBit[v.id]
		dependencies = dependencies.merge(m.dependencies[vBitID])

		// If the node v (dependency of u) can't merge with node w (dependency of v),
		// then node u also can't merge with w (indirect dependency of u) because either:
		// - u and w have different colors
		// - u and w have the same color and both have different colors than v: in this case there is a color cycle
		// - u and w are the same color which is also the same as v: since w is disallowed to v there exists at least
		//   one path from v to w that has a node of a different color and this path is a sub-path of u to w.
		// Therefore, we can just merge the v's disallowed set into u's.
		disallowed = disallowed.merge(m.disallowedBits[vBitID])
		// If color of u and v have different colors, then u and v can't be merged.
		if m.colorByID[u.id] != m.colorByID[v.id] {
			disallowed = disallowed.merge(m.dependencies[vBitID])
		}
	}

	uBitID := m.idToBit[u.id]
	m.disallowedBits[uBitID] = disallowed
	dependencies.turnBitOn(uBitID)
	m.dependencies[uBitID] = dependencies
	m.revToposort = append(m.revToposort, u.id)
}

func (m *merger[T, V]) canMerge(id1, id2 T) bool {
	// Can't merge two nodes if the colors are different
	if m.colorByID[id1] != m.colorByID[id2] {
		return false
	}
	bit1 := m.idToBit[id1]
	bit2 := m.idToBit[id2]
	// Check if one is a disallowed dependency of the other.
	return !m.disallowedBits[bit1].getBit(bit2) && !m.disallowedBits[bit2].getBit(bit1)
}

func (m *merger[T, V]) merge(id1, id2 T) {
	bit1 := m.idToBit[id1]
	bit2 := m.idToBit[id2]
	m.disallowedBits[bit1] = m.disallowedBits[bit1].merge(m.disallowedBits[bit2])
	m.disallowedBits[bit2] = m.disallowedBits[bit2].merge(m.disallowedBits[bit1])
}
