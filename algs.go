package digraph

// NodeIndex provides indexed access to all nodes of an underlying
// graph. The nodes of the graph are discovered when the index is
// constructed. Thus, a node index may not include nodes that are
// added to the graph after it is constructed,
type NodeIndex struct {
	nodes        []Node
	nodesByLabel map[interface{}][]Node
}

// GetNodeIndex builds a node index from the graph for quickly
// accessing all accessible nodes
func (g *Graph) GetNodeIndex() NodeIndex {
	seen := make(map[Node]struct{})
	arr := make([]Node, 0, len(g.nodes))
	for n := range g.nodes {
		arr = append(arr, n)
		seen[n] = struct{}{}
	}
	for i := 0; i < len(arr); i++ {
		edges := arr[i].GetAllOutgoingEdges()
		for edges.HasNext() {
			to := edges.Next().GetTo()
			if _, ok := seen[to]; ok {
				continue
			}
			seen[to] = struct{}{}
			arr = append(arr, to)
		}
	}
	return NodeIndex{nodes: arr}
}

// Len returns the number of nodes in the index
func (n NodeIndex) Len() int {
	return len(n.nodes)
}

// Slice returns a slice of all nodes
func (n NodeIndex) Slice() []Node {
	return n.nodes
}

// Nodes returns an iterator over all nodes
func (n NodeIndex) Nodes() Nodes {
	return Nodes{&NodeArrayIterator{n.nodes}}
}

// NodesByLabel returns an iterator of nodes with the given label
func (n *NodeIndex) NodesByLabel(label interface{}) Nodes {
	if n.nodesByLabel == nil {
		n.buildNodesByLabel()
	}
	return Nodes{&NodeArrayIterator{n.nodesByLabel[label]}}
}

func (n *NodeIndex) buildNodesByLabel() {
	n.nodesByLabel = make(map[interface{}][]Node)
	for _, node := range n.nodes {
		label := node.GetLabel()
		n.nodesByLabel[label] = append(n.nodesByLabel[label], node)
	}
}

// Sinks returns all nodes that have no outgoing edges.
func (n *NodeIndex) Sinks() []Node {
	ret := make([]Node, 0)
	nodes := n.Nodes()
	for nodes.HasNext() {
		node := nodes.Next()
		if !node.HasOutgoingEdges() {
			ret = append(ret, node)
		}
	}
	return ret
}

// Sources returns all nodes that have no incoming edges.
func (n *NodeIndex) Sources() []Node {
	nodeMap := make(map[Node]struct{})
	for _, node := range n.nodes {
		nodeMap[node] = struct{}{}
	}
	for _, node := range n.nodes {
		targets := node.GetAllOutgoingEdges().Targets()
		for targets.HasNext() {
			delete(nodeMap, targets.Next())
		}
	}
	ret := make([]Node, 0, len(nodeMap))
	for node := range nodeMap {
		ret = append(ret, node)
	}
	return ret
}

// Copy creates a copy of source graph in target. If target is an
// empty graph, the result is a clone of the source graph. If target
// is not empty, after this operation target gets a subgraph that is a
// copy of the source
//
// copyNode function copies the given node, and returns the new
// node. If it returns nil, the node is not copied.  copyEdge function
// does the same, creates a copy of the given edge without connecting
// the edges to any of the nodes. The returned edge will be connected
// to the matching nodes.
//
// Returns a map of nodes where the key is the node in the source
// graph, and value is the corresponding node in the target graph
func (n *NodeIndex) Copy(target *Graph, copyNode func(Node) Node, copyEdge func(Edge) Edge) map[Node]Node {
	nodeMap := make(map[Node]Node)
	for _, oldNode := range n.nodes {
		newNode := copyNode(oldNode)
		if newNode != nil {
			target.AddNode(newNode)
			nodeMap[oldNode] = newNode
		}
	}
	for _, oldNode := range n.nodes {
		newNode := nodeMap[oldNode]
		if newNode == nil {
			continue
		}
		for edges := oldNode.GetAllOutgoingEdges(); edges.HasNext(); {
			edge := edges.Next()
			newTarget := nodeMap[edge.GetTo()]
			if newTarget == nil {
				continue
			}
			newEdge := copyEdge(edge)
			if newEdge != nil {
				Connect(newNode, newTarget, copyEdge(edge))
			}
		}
	}
	return nodeMap
}

// Copy creates a copy of source graph in target. If target is an
// empty graph, the result is a clone of the source graph. If target
// is not empty, after this operation target gets a subgraph that is a
// copy of the source
//
// copyNode function copies the given node, and returns the new
// node. If it returns nil, the node is not copied.  copyEdge function
// does the same, creates a copy of the given edge without connecting
// the edges to any of the nodes. The returned edge will be connected
// to the matching nodes.
//
// Returns a map of nodes where the key is the node in the source
// graph, and value is the corresponding node in the target graph
func Copy(target, source *Graph, copyNode func(Node) Node, copyEdge func(Edge) Edge) map[Node]Node {
	ix := target.GetNodeIndex()
	return ix.Copy(source, copyNode, copyEdge)
}
