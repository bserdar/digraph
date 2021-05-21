package digraph

// Sinks returns all nodes that have no outgoing edges. If
// includeDisconnected is true, includes nodes that are not connected
// to the graph
func Sinks(g *Graph, includeDisconnected bool) []*Node {
	ret := make([]*Node, 0)
	nodes := g.AllNodes()
	for nodes.HasNext() {
		node := nodes.Next()
		if !node.AllOutgoingEdges().HasNext() {
			if includeDisconnected || node.AllIncomingEdges().HasNext() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}

// Sources returns all nodes that have no incoming edges. If
// includeDisconnected is true, includes nodes that are not connected
// to the graph
func Sources(g *Graph, includeDisconnected bool) []*Node {
	ret := make([]*Node, 0)
	nodes := g.AllNodes()
	for nodes.HasNext() {
		node := nodes.Next()
		if !node.AllIncomingEdges().HasNext() {
			if includeDisconnected || node.AllOutgoingEdges().HasNext() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}
