package digraph

// Sinks returns all nodes that have no outgoing edges.
func Sinks(index *Index) []Node {
	ret := make([]Node, 0)
	nodes := index.Nodes()
	for nodes.HasNext() {
		node := nodes.Next()
		if !node.HasOut() {
			ret = append(ret, node)
		}
	}
	return ret
}

// Sources returns all nodes that have no incoming edges.
func Sources(index *Index) []Node {
	nodeMap := make(map[Node]struct{})
	for _, node := range index.NodesSlice() {
		nodeMap[node] = struct{}{}
	}
	for _, node := range index.NodesSlice() {
		targets := node.Out().Targets()
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
func Copy(index *Index, target *Graph, copyNode func(Node) Node, copyEdge func(Edge) Edge) map[Node]Node {
	nodeMap := make(map[Node]Node)
	for _, oldNode := range index.NodesSlice() {
		newNode := copyNode(oldNode)
		if newNode != nil {
			target.AddNode(newNode)
			nodeMap[oldNode] = newNode
		}
	}
	for _, oldNode := range index.NodesSlice() {
		newNode := nodeMap[oldNode]
		if newNode == nil {
			continue
		}
		for edges := oldNode.Out(); edges.HasNext(); {
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

// CopyGraph creates a copy of source graph in target. If target is an
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
func CopyGraph(target, source *Graph, copyNode func(Node) Node, copyEdge func(Edge) Edge) map[Node]Node {
	ix := target.GetIndex()
	return Copy(ix, source, copyNode, copyEdge)
}

// IterateGraph iterates all nodes and edges of the graph until one of the functions returns false
func IterateGraph(g *Graph, nodeFunc func(Node) bool, edgeFunc func(Edge) bool) bool {
	seen := make(map[Node]struct{})
	for node := range g.nodes {
		if !IterateUnique(node, nodeFunc, edgeFunc, seen) {
			return false
		}
	}
	return true
}

// Iterate all nodes and edges of the graph until one of the functions returns false
func Iterate(root Node, nodeFunc func(Node) bool, edgeFunc func(Edge) bool) bool {
	return IterateUnique(root, nodeFunc, edgeFunc, map[Node]struct{}{})
}

// IterateUnique iterates all nodes and edges until one of the
// functions returns false. It skips the nodes in the seen map
func IterateUnique(root Node, nodeFunc func(Node) bool, edgeFunc func(Edge) bool, seen map[Node]struct{}) bool {
	if _, exists := seen[root]; exists {
		return true
	}
	seen[root] = struct{}{}
	if !nodeFunc(root) {
		return false
	}
	for edges := root.Out(); edges.HasNext(); {
		edge := edges.Next()
		if !edgeFunc(edge) {
			return false
		}
		if !IterateUnique(edge.GetTo(), nodeFunc, edgeFunc, seen) {
			return false
		}
	}
	return true
}
