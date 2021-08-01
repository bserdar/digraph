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
	Connect(n1, n2, NewBasicEdge("label", nil))

	edges := n1.GetAllOutgoingEdges().All()
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, %d", len(edges))
	}
	if edges[0].GetTo() != n2 {
		t.Error("Wrong end")
	}
	edges[0].GetFrom().RemoveOutgoingEdge(edges[0])
	if len(n1.GetAllOutgoingEdges().All()) != 0 {
		t.Error("There are still edges")
	}
}
