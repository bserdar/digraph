package digraph

// A Graph is a labeled directed graph. The labels can be nil. Zero
// value of a Graph is ready to use.
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

// AddNode adds the node to the graph. The node is disconnected
func (g *Graph) AddNode(node Node) {
	g.init()
	g.nodes[node] = struct{}{}
}

// AllNodes returns an iterator over all nodes of a graph
func (g *Graph) AllNodes() Nodes {
	return Nodes{&NodeArrayIterator{g.GetNodeIndex().Slice()}}
}
