package digraph

import (
	"fmt"
	"io"
)

// DOTRenderer renders a graph in Graphviz dot format
type DOTRenderer struct {
	NodeRenderer func(string, Node, io.Writer) error
	EdgeRenderer func(string, string, Edge, io.Writer) error
}

// RenderNode renders a node. If node renderer is not set, calls the default renderer
func (d DOTRenderer) RenderNode(ID string, node Node, w io.Writer) error {
	if d.NodeRenderer == nil {
		return DefaultDOTNodeRender(ID, node, w)
	}
	return d.NodeRenderer(ID, node, w)
}

// RenderEdge renders an edge. If edge renderer is not set, call the default rendeded
func (d DOTRenderer) RenderEdge(fromID, toID string, edge Edge, w io.Writer) error {
	if d.EdgeRenderer == nil {
		return DefaultDOTEdgeRender(fromID, toID, edge, w)
	}
	return d.EdgeRenderer(fromID, toID, edge, w)
}

// DefaultDOTNodeRender renders the node with the given ID. If the
// node has a label, it uses that label, otherwise node is not
// labeled.
func DefaultDOTNodeRender(ID string, node Node, w io.Writer) error {
	if node.GetLabel() != nil {
		_, err := fmt.Fprintf(w, "  %s [label=\"%v\"];\n", ID, node.GetLabel())
		return err
	}
	_, err := fmt.Fprintf(w, "  %s;\n", ID)
	return err
}

// DefaultDOTEdgeRender renders the edge with a label if there is
// one, or without a label if there is not a label.
func DefaultDOTEdgeRender(fromNode, toNode string, edge Edge, w io.Writer) error {
	lbl := edge.GetLabel()
	if lbl != nil {
		if _, err := fmt.Fprintf(w, "  %s -> %s [label=\"%s\"];\n", fromNode, toNode, lbl); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w, "  %s -> %s;\n", fromNode, toNode); err != nil {
			return err
		}
	}
	return nil
}

// Render writes a DOT graph with the given name
func (d DOTRenderer) Render(g *Graph, graphName string, out io.Writer) error {
	if _, err := fmt.Fprintf(out, "digraph %s {\n", graphName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "rankdir=\"LR\";\n"); err != nil {
		return err
	}

	// Give nodes unique IDs for the graph
	nodeMap := map[Node]string{}
	x := 0
	for itr := g.GetAllNodes(); itr.HasNext(); {
		node := itr.Next()
		nodeId := fmt.Sprintf("n%d", x)
		nodeMap[node] = nodeId
		x++
		if err := d.RenderNode(nodeId, node, out); err != nil {
			return err
		}
	}
	for itr := g.GetAllNodes(); itr.HasNext(); {
		node := itr.Next()
		for edgeItr := node.GetAllOutgoingEdges(); edgeItr.HasNext(); {
			edge := edgeItr.Next()
			fromNodeId := nodeMap[edge.GetFrom()]
			toNodeId := nodeMap[edge.GetTo()]
			if err := d.RenderEdge(fromNodeId, toNodeId, edge, out); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintf(out, "}\n"); err != nil {
		return err
	}
	return nil
}
