package digraph

// Index provides indexed access to all nodes and edges of an
// underlying graph. The index is lazily constructed to include all
// nodes and edges. Index is constructed by accessible nodes and
// edges, thus the underlying graph should not be modified.
type Index struct {
	g *Graph

	allNodes        []Node
	allNodesByLabel map[interface{}][]Node

	incomingEdges        map[Node][]Edge
	incomingEdgesByLabel map[Node]map[interface{}][]Edge
}

// GetIndex returns an uninitialized index for the graph
func (g *Graph) GetIndex() *Index {
	return &Index{g: g}
}

// NodesSlice returns all accessible nodes as a slice
func (index *Index) NodesSlice() []Node {
	if index.allNodes == nil {
		index.allNodes = make([]Node, 0)
		seen := map[Node]struct{}{}
		for node := range index.g.nodes {
			IterateUnique(node, func(n Node) bool {
				index.allNodes = append(index.allNodes, n)
				return true
			}, func(e Edge) bool { return true }, seen)
		}
	}
	return index.allNodes
}

// Nodes returns an iterator over all nodes
func (index *Index) Nodes() Nodes {
	return Nodes{&nodeSliceIterator{index.NodesSlice()}}
}

// NodesByLabel returns an iterator of nodes with the given label
func (index *Index) NodesByLabelSlice(label interface{}) []Node {
	if index.allNodesByLabel == nil {
		index.allNodesByLabel = make(map[interface{}][]Node)
		seen := map[Node]struct{}{}
		for node := range index.g.nodes {
			IterateUnique(node, func(n Node) bool {
				index.allNodesByLabel[n.GetLabel()] = append(index.allNodesByLabel[n.GetLabel()], n)
				return true
			}, func(e Edge) bool { return true }, seen)
		}
	}
	return index.allNodesByLabel[label]
}

// NodesByLabel returns nodes by label
func (index *Index) NodesByLabel(label interface{}) Nodes {
	return Nodes{&nodeSliceIterator{index.NodesByLabelSlice(label)}}
}

// Out returns the outgoing edges of a node
func (index *Index) Out(node Node) Edges {
	return node.Out()
}

// Out returns the outgoing edges of a node
func (index *Index) OutSlice(node Node) []Edge {
	return node.Out().All()
}

// OutWith returns the outgoing edges of a node with a label
func (index *Index) OutWith(node Node, label interface{}) Edges {
	return node.OutWith(label)
}

// OutWithSlice returns the outgoing edges of a node with a label
func (index *Index) OutWithSlice(node Node, label interface{}) []Edge {
	return node.OutWith(label).All()
}

// InSlice returns the incoming edges of a node. These will include only
// those edges that are from the nodes included in this graph
func (index *Index) InSlice(node Node) []Edge {
	if index.incomingEdges == nil {
		index.incomingEdges = make(map[Node][]Edge)
		seen := map[Node]struct{}{}
		for node := range index.g.nodes {
			IterateUnique(node, func(n Node) bool {
				return true
			}, func(e Edge) bool {
				index.incomingEdges[e.GetTo()] = append(index.incomingEdges[e.GetTo()], e)
				return true
			},
				seen)
		}
	}
	return index.incomingEdges[node]
}

// In returns the incoming edges of a node. These will include only
// those edges that are from the nodes included in this graph
func (index *Index) In(node Node) Edges {
	return Edges{&edgeSliceIterator{index.InSlice(node)}}
}

// InWithSlice returns the incoming edges of a node by label. These will include only
// those edges that are from the nodes included in this graph
func (index *Index) InWithSlice(node Node, label interface{}) []Edge {
	if index.incomingEdgesByLabel == nil {
		index.incomingEdgesByLabel = make(map[Node]map[interface{}][]Edge)
		seen := map[Node]struct{}{}
		for node := range index.g.nodes {
			IterateUnique(node, func(n Node) bool {
				return true
			}, func(e Edge) bool {
				m := index.incomingEdgesByLabel[e.GetTo()]
				if m == nil {
					m = make(map[interface{}][]Edge)
					index.incomingEdgesByLabel[e.GetTo()] = m
				}
				m[e.GetLabel()] = append(m[e.GetLabel()], e)
				return true
			},
				seen)
		}
	}
	m := index.incomingEdgesByLabel[node]
	if m != nil {
		return m[label]
	}
	return nil
}

// InWith returns the incoming edges of a node by label. These will include only
// those edges that are from the nodes included in this graph
func (index *Index) InWith(node Node, label interface{}) Edges {
	return Edges{&edgeSliceIterator{index.InWithSlice(node, label)}}
}
