package graph

import (
	"reflect"
	"testing"
)

func setupTestGraph() *Graph {
	g := NewGraph()

	for i := 1; i <= 5; i++ {
		g.AddNode(i)
	}
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(4, 5)

	return g
}

func TestShortestPathSuccess(t *testing.T) {
	g := setupTestGraph()

	path := g.ShortestPath(1, 3)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(path, expected) {
		t.Errorf("Expected %v, got %v", expected, path)
	}
}

func TestShortestPathNoPath(t *testing.T) {
	g := setupTestGraph()

	path := g.ShortestPath(1, 5)

	if path != nil {
		t.Errorf("Expected nil, got %v", path)
	}
}

func TestInfluenceInDegree(t *testing.T) {
	g := NewGraph()

	g.AddNode(1)
	g.AddNode(2)
	g.AddNode(3)

	g.AddEdge(1, 2)
	g.AddEdge(3, 2)

	metrics := g.Influence(2)

	if metrics.InDegree != 2 {
		t.Errorf("Expected in-degree 2, got %d", metrics.InDegree)
	}
}

func TestTopInfluencers(t *testing.T) {
	g := NewGraph()

	g.AddNode(1)
	g.AddNode(2)
	g.AddNode(3)

	g.AddEdge(1, 2)
	g.AddEdge(3, 2)

	top := g.TopInfluencers(1)

	if len(top) != 1 || top[0] != 2 {
		t.Errorf("Expected top influencer 2, got %v", top)
	}
}

func TestConnectedComponents(t *testing.T) {
	g := setupTestGraph()

	components := g.ConnectedComponents()

	if components != 2 {
		t.Errorf("Expected 2 components, got %d", components)
	}
}