package lib

import (
	"errors"
)

type Graph[T comparable] struct {
	Sources    *Set[T]
	Vertices   *Set[T]
	Indegree   map[T]int
	AdjList    map[T]*Set[T]
	Undirected bool

	SortLevel     *Set[T]
	SortDegrees   map[T]int
	SortRemaining int
}

func NewGraph[T comparable]() *Graph[T] {
	g := new(Graph[T])
	g.Sources = NewSet[T]()
	g.Vertices = NewSet[T]()
	g.Indegree = map[T]int{}
	g.AdjList = map[T]*Set[T]{}

	g.SortLevel = NewSet[T]()
	g.SortDegrees = map[T]int{}
	return g
}

func (g *Graph[T]) AddEdge(u, v T) {
	uExists := g.Vertices.Has(u)
	vExists := g.Vertices.Has(v)
	if uExists && vExists && g.AdjList[u].Has(v) {
		return
	}
	if uExists && vExists && g.AdjList[v].Has(u) {
		g.Undirected = true
	}

	if !uExists {
		g.Sources.Add(u)
		g.Vertices.Add(u)
		g.AdjList[u] = NewSet[T]()
	}
	if !vExists {
		g.Vertices.Add(v)
		g.AdjList[v] = NewSet[T]()
	}

	g.Indegree[v]++
	g.AdjList[u].Add(v)
	g.Sources.Delete(v)

	g.SortLevel = g.Sources
	g.SortDegrees[v]++
	g.SortRemaining = g.Vertices.Size
}

func (g *Graph[T]) HasNextLevel() bool {
	return g.SortRemaining > 0
}

func (g *Graph[T]) GetLevel() ([]T, error) {
	var level []T
	nextLevel := NewSet[T]()
	for u := range g.SortLevel.Iterator() {
		level = append(level, u)
		g.SortRemaining--
		for v := range g.AdjList[u].Iterator() {
			g.SortDegrees[v]--
			if g.SortDegrees[v] == 0 {
				nextLevel.Add(v)
			}
		}
	}

	g.SortLevel = nextLevel
	if g.SortLevel.Size == 0 && g.SortRemaining > 0 {
		return nil, errors.New("cycle detected")
	} else {
		return level, nil
	}
}

func (g *Graph[T]) ResetSort() {
	g.SortRemaining = g.Vertices.Size
	g.SortLevel = g.Sources
	for k, v := range g.Indegree {
		g.SortDegrees[k] = v
	}
}
