package digraph

// CheckIsomoprhism checks to see if graphs whose node indexes are
// given are equal as defined by the edge equivalence and node
// equivalence functions. The nodeEquivalenceFunction will be called
// for nodes whose labels are the same. The edgeEquivalenceFunction
// will be called for edges connecting equivalent nodes with the same
// labels.
//
// Node isomorphism check will fail if one node is equivalent to
// multiple nodes
func CheckIsomorphism(nodes1, nodes2 *Index, nodeEquivalenceFunc func(n1, n2 Node) bool, edgeEquivalenceFunc func(e1, e2 Edge) bool) bool {
	// Map of nodes1 -> nodes2
	nodeMapping1_2 := make(map[Node]Node)
	// Map of nodes2 -> nodes1
	nodeMapping2_1 := make(map[Node]Node)

	if len(nodes1.NodesSlice()) != len(nodes2.NodesSlice()) {
		return false
	}

	for nodes := nodes1.Nodes(); nodes.HasNext(); {
		node1 := nodes.Next()
		for x := nodes2.NodesByLabel(node1.GetLabel()); x.HasNext(); {
			node2 := x.Next()
			if nodeEquivalenceFunc(node1, node2) {
				if _, ok := nodeMapping1_2[node1]; ok {
					return false
				}
				nodeMapping1_2[node1] = node2
				if _, ok := nodeMapping2_1[node2]; ok {
					return false
				}
				nodeMapping2_1[node2] = node1
			}
		}
	}

	if len(nodeMapping1_2) != len(nodes1.NodesSlice()) {
		return false
	}
	// No need to check the other map, adding one will add to the other

	// Node equivalences are established, now check edge equivalences for each node
	for node1, node2 := range nodeMapping1_2 {
		// node1 and node2 are equivalent. Now we check if equivalent edges go to equivalent nodes
		edges1 := node1.Out().All()
		edges2 := node2.Out().All()
		// There must be same number of edges
		if len(edges1) != len(edges2) {
			return false
		}
		// Find equivalent edges
		edgeMap := make(map[Edge]Edge)
		for _, edge1 := range edges1 {
			found := false
			for _, edge2 := range edges2 {
				if edge1.GetLabel() == edge2.GetLabel() &&
					nodeMapping1_2[edge1.GetTo()] == edge2.GetTo() &&
					edgeEquivalenceFunc(edge1, edge2) {
					if found {
						// Multiple edges match
						return false
					}
					edgeMap[edge1] = edge2
					found = true
				}
			}
			if !found {
				return false
			}
		}
		if len(edgeMap) != len(edges1) {
			return false
		}
	}
	return true
}
