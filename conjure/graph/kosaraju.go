package graph

import (
	"sort"
)

// ComponentID is an identifier for strongly connected components.
type ComponentID int

type StronglyConnectedComponents[T comparable] struct {
	Components      map[ComponentID][]T
	ComponentByItem map[T]ComponentID
	ComponentGraph  *Graph[ComponentID]
}

// calculateStronglyConnectedComponents takes in a graph and groups nodes into "strongly connected components" (SCCs).
// Two nodes are in the same SCC if and only if they are reachable from each other (i.e. they are part of a cycle).
// The graph of SCCs is a directed acyclic graph (DAG).
func calculateStronglyConnectedComponents[T comparable](g *Graph[T]) *StronglyConnectedComponents[T] {
	return (&kosaraju[T]{g: g}).run()
}

type kosaraju[T comparable] struct {
	// input to kosaraju's algorithm
	g *Graph[T]

	// output of kosaraju's algorithm
	componentByItem map[T]ComponentID

	// variables internal to the algorithm
	revG          *Graph[T]
	visited       map[T]struct{}
	curComponent  ComponentID
	numComponents int
	revToposort   []T
	comparator    func(t1, t2 T) bool
}

// run kosaraju's algorithm for finding strongly connected components
// https://www.geeksforgeeks.org/strongly-connected-components/
// https://en.wikipedia.org/wiki/Kosaraju%27s_algorithm
func (k *kosaraju[T]) run() *StronglyConnectedComponents[T] {
	n := len(k.g.Nodes)

	index := make(map[T]int, n)
	for i, u := range k.g.Nodes {
		index[u.ID] = i
	}
	k.comparator = func(t1, t2 T) bool {
		return index[t1] < index[t2]
	}

	k.revG = k.reverseGraph(k.g)
	k.componentByItem = make(map[T]ComponentID, n)

	k.revToposort = make([]T, 0, n)
	k.visited = make(map[T]struct{})
	for _, u := range k.revG.Nodes {
		if _, visited := k.visited[u.ID]; !visited {
			k.revDfs(u)
		}
	}

	k.visited = make(map[T]struct{})
	k.numComponents = 0
	for i := n - 1; i >= 0; i-- {
		u := k.g.NodesByID[k.revToposort[i]]
		if _, visited := k.visited[u.ID]; !visited {
			k.curComponent = ComponentID(k.numComponents)
			k.numComponents++
			k.dfs(u)
		}
	}

	components := make(map[ComponentID][]T)
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

	return &StronglyConnectedComponents[T]{
		Components:      components,
		ComponentByItem: k.componentByItem,
		ComponentGraph:  k.buildComponentGraph(components),
	}
}

// reverseGraph builds the reverse graph of g. In other words, if there is an edge u->v in the original graph
// if and only if there is v->u in the reverse graph.
func (k *kosaraju[T]) reverseGraph(g *Graph[T]) *Graph[T] {
	revG := newGraph[T](len(g.Nodes))

	for _, u := range g.Nodes {
		revG.addNode(u.ID)
	}

	for _, u := range g.Nodes {
		for _, v := range u.sortedEdges(k.comparator) {
			revG.addEdgesByID(v.ID, u.ID)
		}
	}

	return revG
}

func (k *kosaraju[T]) revDfs(u *Node[T]) {
	k.visited[u.ID] = struct{}{}
	for _, v := range u.sortedEdges(k.comparator) {
		if _, visited := k.visited[v.ID]; !visited {
			k.revDfs(v)
		}
	}
	k.revToposort = append(k.revToposort, u.ID)
}

func (k *kosaraju[T]) dfs(u *Node[T]) {
	k.visited[u.ID] = struct{}{}
	k.componentByItem[u.ID] = k.curComponent
	for _, v := range u.sortedEdges(k.comparator) {
		if _, visited := k.visited[v.ID]; !visited {
			k.dfs(v)
		}
	}
}

func (k *kosaraju[T]) buildComponentGraph(components map[ComponentID][]T) *Graph[ComponentID] {
	compG := newGraph[ComponentID](k.numComponents)

	for compID := range components {
		compG.addNode(compID)
	}
	sort.SliceStable(compG.Nodes, func(i, j int) bool {
		return compG.Nodes[i].ID < compG.Nodes[j].ID
	})

	for _, u := range k.g.Nodes {
		uCompID := k.componentByItem[u.ID]
		for _, v := range u.sortedEdges(k.comparator) {
			vCompID := k.componentByItem[v.ID]
			if uCompID != vCompID {
				compG.addEdgesByID(uCompID, vCompID)
			}
		}
	}

	return compG
}
