package digraph

import (
	"testing"
)

func TestBasicGraph(t *testing.T) {
	g := New()

	n1 := NewBasicNode("1", nil)
	g.AddNode(n1)
	n2 := NewBasicNode("2", nil)
	g.AddNode(n2)
	g.AddEdge(n1, n2, NewBasicEdge("label", nil))

	edges := n1.AllOutgoingEdges().All()
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, %d", len(edges))
	}
	if edges[0].From() != n1 {
		t.Error("Wrong start")
	}
	if edges[0].To() != n2 {
		t.Error("Wrong end")
	}
	edges = n2.AllIncomingEdges().All()
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, %d", len(edges))
	}
	if edges[0].From() != n1 {
		t.Error("Wrong start")
	}
	if edges[0].To() != n2 {
		t.Error("Wrong end")
	}
	edges[0].Remove()
	if len(n1.AllOutgoingEdges().All()) != 0 {
		t.Error("There are still edges")
	}
	if len(n2.AllIncomingEdges().All()) != 0 {
		t.Error("There are still edges")
	}
}
