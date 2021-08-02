package digraph

// A Graph is a labeled directed graph. The labels can be nil. Zero
// value of a Graph is ready to use.
//
// A graph knows some of the nodes of the graph. The remaining nodes
// are discovered when needd. This allows a graph to have multiple
// disjoint components.
//
// Two graphs can be merged simply by adding an edge between their two
// nods. Then the graph containing the source node of the edge will
// include all the accessible nodes of the second graph.
type Graph struct {
	// nodes keeps some of the nodes of the graph
	nodes map[Node]struct{}
}

func (g *Graph) init() {
	if g.nodes == nil {
		g.nodes = make(map[Node]struct{})
	}
}

// New returns a new empty graph
func New() *Graph {
	g := new(Graph)
	g.init()
	return g
}

// AddNode adds the node to the graph. The node can be the node of a
// disconnected graph.
func (g *Graph) AddNode(node Node) {
	g.init()
	g.nodes[node] = struct{}{}
}

// GetAllNodes returns an iterator over all nodes of a graph
func (g *Graph) GetAllNodes() Nodes {
	return Nodes{&NodeArrayIterator{g.GetNodeIndex().Slice()}}
}
