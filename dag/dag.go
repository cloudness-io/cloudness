package dag

import (
	"github.com/cloudness-io/cloudness/app/usererror"
)

// Graph represents a directed acyclic graph using adjacency list
type Graph[T comparable] struct {
	adjList  map[T][]T
	inDegree map[T]int
	vertices map[T]bool
}

// NewGraph initializes a new graph
func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		adjList:  make(map[T][]T),
		inDegree: make(map[T]int),
		vertices: make(map[T]bool),
	}
}

// AddVertex ensures the vertex exists in the graph
func (g *Graph[T]) AddVertex(v T) {
	if _, exits := g.vertices[v]; !exits {
		g.vertices[v] = true
		g.inDegree[v] = 0
		g.adjList[v] = make([]T, 0)
	}
}

// AddEdge adds a directed edge from u to v
func (g *Graph[T]) AddEdge(from, to T) {
	g.AddVertex(from)
	g.AddVertex(to)
	g.adjList[from] = append(g.adjList[from], to)
	g.inDegree[to]++
}

// TopoSort performs a topological sort of the graph
func (g *Graph[T]) TopoSort() ([]T, error) {
	queue := make([]T, 0)
	for v := range g.vertices {
		if g.inDegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	var sorted []T
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		sorted = append(sorted, curr)

		for _, neighbor := range g.adjList[curr] {
			g.inDegree[neighbor]--
			if g.inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(sorted) != len(g.vertices) {
		return nil, usererror.ErrCyclicHierarchy
	}

	return sorted, nil
}
