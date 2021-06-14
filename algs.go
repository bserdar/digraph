package digraph

// Sinks returns all nodes that have no outgoing edges. If
// includeDisconnected is true, includes nodes that are not connected
// to the graph
func Sinks(g *Graph, includeDisconnected bool) []Node {
	ret := make([]Node, 0)
	nodes := g.AllNodes()
	for nodes.HasNext() {
		node := nodes.Next()
		if !node.GetNodeHeader().AllOutgoingEdges().HasNext() {
			if includeDisconnected || node.GetNodeHeader().AllIncomingEdges().HasNext() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}

// SinksAmongNodes finds the sink nodes among the given nodes
func SinksAmongNodes(nodes []Node, includeDisconnected bool) []Node {
	ret := make([]Node, 0)
	for _, node := range nodes {
		if !node.GetNodeHeader().AllOutgoingEdges().HasNext() {
			if includeDisconnected || node.GetNodeHeader().AllIncomingEdges().HasNext() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}

// Sources returns all nodes that have no incoming edges. If
// includeDisconnected is true, includes nodes that are not connected
// to the graph
func Sources(g *Graph, includeDisconnected bool) []Node {
	ret := make([]Node, 0)
	nodes := g.AllNodes()
	for nodes.HasNext() {
		node := nodes.Next()
		if !node.GetNodeHeader().AllIncomingEdges().HasNext() {
			if includeDisconnected || node.GetNodeHeader().AllOutgoingEdges().HasNext() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}

// SourcesAmongNodes returns all nodes that have no incoming edges
// among the given nodes.
func SourcesAmongNodes(nodes []Node, includeDisconnected bool) []Node {
	ret := make([]Node, 0)
	for _, node := range nodes {
		if !node.GetNodeHeader().AllIncomingEdges().HasNext() {
			if includeDisconnected || node.GetNodeHeader().AllOutgoingEdges().HasNext() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}

// Copy creates a copy of source graph in target. If target is an
// empty graph, the result is a clone of the source graph. If target
// is not empty, after this operation target gets a subgraph that is a
// copy of the source
//
// copyNode function copies the given node, and returns the new
// node. The node must not be inserted into the target graph. copyEdge
// function does the same, creates a copy of the given edge without
// connecting the edges to any of the nodes. The returned edge will be
// connected to the matching nodes.
//
// Returns a map of nodes where the key is the node in the source
// graph, and value is the corresponding node in the target graph
func Copy(target, source *Graph, copyNode func(Node) Node, copyEdge func(Edge) Edge) map[Node]Node {
	nodeMap := make(map[Node]Node)
	for nodes := source.AllNodes(); nodes.HasNext(); {
		oldNode := nodes.Next()
		newNode := copyNode(oldNode)
		target.AddNode(newNode)
		nodeMap[oldNode] = newNode
	}
	for nodes := source.AllNodes(); nodes.HasNext(); {
		oldNode := nodes.Next()
		for edges := oldNode.AllOutgoingEdges(); edges.HasNext(); {
			edge := edges.Next()
			target.AddEdge(nodeMap[edge.From()], nodeMap[edge.To()], copyEdge(edge))
		}
	}
	return nodeMap
}
