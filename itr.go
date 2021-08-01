package digraph

// Nodes is a convenience wrapper aroung NodeItr for chained methods
type Nodes struct {
	NodeIterator
}

// All returns all remaining nodes
func (n Nodes) All() []Node {
	ret := make([]Node, 0)
	for n.HasNext() {
		ret = append(ret, n.Next())
	}
	return ret
}

// NodeIterator iterates through a list of nodes
type NodeIterator interface {
	// Returns if there are more nodes to go through
	HasNext() bool
	// If HasNext is true, returns the next node and advances. Otherwise, panics
	Next() Node
}

type NodeArrayIterator struct {
	Nodes []Node
}

func (a *NodeArrayIterator) HasNext() bool { return len(a.Nodes) > 0 }
func (a *NodeArrayIterator) Next() Node    { ret := a.Nodes[0]; a.Nodes = a.Nodes[1:]; return ret }

type filterNodes struct {
	source NodeIterator
	filter func(Node) bool

	nextReady bool
	next      Node
}

func (a *filterNodes) adv() {
	if a.nextReady {
		return
	}
	a.nextReady = true
	a.next = nil
	for a.source.HasNext() {
		node := a.source.Next()
		if a.filter(node) {
			a.next = node
			return
		}
	}
}

func (a *filterNodes) HasNext() bool {
	a.adv()
	return a.next != nil
}

func (a *filterNodes) Next() Node {
	a.adv()
	if a.next == nil {
		panic("Next node not available")
	}
	a.nextReady = false
	return a.next
}

// Select returns a subset of the given nodes containing only those nodes selected by the predicate
func (n Nodes) Select(predicate func(Node) bool) Nodes {
	return Nodes{&filterNodes{source: n, filter: predicate}}
}

// NewNodes returns a Nodes for the given array of nodes
func NewNodes(nodes ...Node) Nodes { return Nodes{&NodeArrayIterator{Nodes: nodes}} }

// Unique filters the nodes so only unique nodes are returned
func (n Nodes) Unique() Nodes {
	seen := make(map[Node]struct{})
	return n.Select(func(node Node) bool {
		_, ok := seen[node]
		if !ok {
			seen[node] = struct{}{}
		}
		return !ok
	})
}

type Edges struct {
	EdgeIterator
}

// All returns all remaining edges
func (e Edges) All() []Edge {
	ret := make([]Edge, 0)
	for e.HasNext() {
		ret = append(ret, e.Next())
	}
	return ret
}

type edgeNodeSelector struct {
	source     EdgeIterator
	selectNode func(Edge) Node
}

func (e *edgeNodeSelector) HasNext() bool { return e.source.HasNext() }
func (e *edgeNodeSelector) Next() Node    { return e.selectNode(e.source.Next()) }

// Targets returns a node iterator that goes through the target nodes
func (e Edges) Targets() Nodes {
	return Nodes{&edgeNodeSelector{source: e, selectNode: func(e Edge) Node { return e.GetTo() }}}
}

// EdgeIterator iterates through a list of edges
type EdgeIterator interface {
	// Returns if there are more edges to go through
	HasNext() bool
	// If HasNext is true, returns the next edge and advances. Otherwise panics
	Next() Edge
}

type EdgeArrayIterator struct {
	Edges []Edge
}

func (e *EdgeArrayIterator) HasNext() bool { return len(e.Edges) > 0 }
func (e *EdgeArrayIterator) Next() Edge    { ret := e.Edges[0]; e.Edges = e.Edges[1:]; return ret }
