package lib

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph[string]()
	if g.Sources.Size != 0 {
		t.Fatal("expecting empty sources set")
	}
	if g.Vertices.Size != 0 {
		t.Fatal("expecting empty vertices set")
	}
	if g.Indegree == nil {
		t.Fatal("expecting initialized indegree map")
	}
	if g.AdjList == nil {
		t.Fatal("expecting initialized adj list map")
	}
	if g.Undirected {
		t.Fatal("expecting undirected flag to be false")
	}

	if g.SortLevel.Size != 0 {
		t.Fatal("expecting empty sort level set")
	}
	if g.SortDegrees == nil {
		t.Fatal("expecting initialized sort degrees map on init")
	}
	if g.SortRemaining != 0 {
		t.Fatal("expecting remaining count of 0 on init")
	}
}

func TestGraph_AddEdge(t *testing.T) {
	g := NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("d", "c")
	if g.Vertices.Size != 4 {
		t.Fatal("expecting size of 4")
	}
	if g.Indegree["b"] != 1 {
		t.Fatal("expecting indegree=1 for b")
	}
	if g.Indegree["c"] != 2 {
		t.Fatal("expecting indegree=2 for c")
	}
	_, exist := g.Indegree["a"]
	if exist {
		t.Fatal("expecting no indegree for a")
	}
	_, exist = g.Indegree["d"]
	if exist {
		t.Fatal("expecting no indegree for d")
	}
	if !g.AdjList["a"].Has("b") {
		t.Fatal("expecting a->b in adj list")
	}
	if g.AdjList["a"].Size != 1 {
		t.Fatal("expecting one neighbor for a")
	}
	if !g.AdjList["b"].Has("c") {
		t.Fatal("expecting b->c in adj list")
	}
	if g.AdjList["b"].Size != 1 {
		t.Fatal("expecting one neighbor for b")
	}
	if !g.AdjList["d"].Has("c") {
		t.Fatal("expecting d->c in adj list")
	}
	if g.AdjList["d"].Size != 1 {
		t.Fatal("expecting one neighbor for d")
	}
	if g.Undirected {
		t.Fatal("expecting directed graph")
	}
}

func TestGraph_Undirected(t *testing.T) {
	g := NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("b", "a")
	if !g.Undirected {
		t.Fatal("expecting undirected graph")
	}
}

func TestGraph_TestCycle(t *testing.T) {
	g := NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("c", "a")
	_, err := g.GetLevel()
	if err == nil {
		t.Fatal("expecting cycle error")
	}
	if err.Error() != "cycle detected" {
		t.Fatal("expecting cycle error")
	}

	g = NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("c", "d")
	g.AddEdge("d", "e")
	g.AddEdge("e", "c")
	_, err = g.GetLevel() //a
	_, err = g.GetLevel() //b
	_, err = g.GetLevel() //c
	_, err = g.GetLevel() //d
	_, err = g.GetLevel() //e
	_, err = g.GetLevel() //c
	if err == nil {
		t.Fatal("expecting cycle error")
	}
	if err.Error() != "cycle detected" {
		t.Fatal("expecting cycle error")
	}
}

func TestGraph_HasNextLevel(t *testing.T) {
	g := NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("c", "d")

	// a
	if !g.HasNextLevel() {
		t.Fatal("expecting next level true")
	}
	_, _ = g.GetLevel()

	//b
	if !g.HasNextLevel() {
		t.Fatal("expecting next level true")
	}
	_, _ = g.GetLevel()

	// c
	if !g.HasNextLevel() {
		t.Fatal("expecting next level true")
	}
	_, _ = g.GetLevel()

	// d
	if !g.HasNextLevel() {
		t.Fatal("expecting next level true")
	}
	_, _ = g.GetLevel()

	// empty
	if g.HasNextLevel() {
		err, level := g.GetLevel()
		fmt.Println(err, level, "@@@")
		t.Fatal("expecting no more levels")
	}
}

func TestGraph_NextLevel(t *testing.T) {
	g := NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("c", "d")

	// a
	level, err := g.GetLevel()
	if err != nil {
		t.Fatal(err)
	}
	if level[0] != "a" && len(level[0]) != 1 {
		t.Fatal("failure on level 1")
	}

	// b
	level, err = g.GetLevel()
	if err != nil {
		t.Fatal(err)
	}
	if level[0] != "b" && len(level[0]) != 1 {
		t.Fatal("failure on level 2")
	}

	// c
	level, err = g.GetLevel()
	if err != nil {
		t.Fatal(err)
	}
	if level[0] != "c" && len(level[0]) != 1 {
		t.Fatal("failure on level 3")
	}

	// d
	level, err = g.GetLevel()
	if err != nil {
		t.Fatal(err)
	}
	if level[0] != "d" && len(level[0]) != 1 {
		t.Fatal("failure on level 4")
	}
}

func TestGraph_ResetSort(t *testing.T) {
	g := NewGraph[string]()
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("c", "d")

	g.SortRemaining = 0
	g.SortLevel = NewSet[string]()
	g.SortDegrees = map[string]int{}
	g.ResetSort()

	if g.SortRemaining != g.Vertices.Size {
		t.Fatal("expecting reset of sort remaining")
	}
	sources := strings.Join(g.Sources.Items(), "")
	level := strings.Join(g.SortLevel.Items(), "")
	if sources != level {
		t.Fatal("expecting sources level on reset")
	}

	expected := map[string]int{
		"b": 1,
		"c": 1,
	}
	if fmt.Sprint(expected) == fmt.Sprint(g.SortDegrees) {
		t.Fatal("expecting a reset of sort indegree map on reset")
	}
}
