package digraph

import (
	"container/list"
)

// Graph represents a labeled or unlabeled directed graph. Zero value
// for Graph is an empty graph ready to use.
type Graph struct {
	nodesByLabel map[interface{}]*list.List
	allNodes     list.List
}

// Node is a node of a directed graph.
//
// A Node may have a label. The label is given when the node is
// created, and cannot be changed. Multiple nodes can have the same
// label in a graph.
//
// Node keeps a list/map of all outgoing and incoming edges. These
// edges may or may not be labeled. Same label can be used for
// outgoing and incoming edges.
//
// Each node may have an application-defined payload.
type Node struct {
	Payload interface{}

	label interface{}

	nodesByLabelEl *list.Element
	allNodesEl     *list.Element
	graph          *Graph

	out    map[interface{}]*list.List
	allOut list.List
	in     map[interface{}]*list.List
	allIn  list.List
}

// Edge represents a labeled or unlabeled directed edge between two
// nodes of a graph.
//
// An edge has from and to nodes, which are defined when edge is
// created, and connot be changed.
//
// An edge may have a label. Same label can be reused on multiple
// edges outgoing from a node.
//
// An edge may have an application defined payload
type Edge struct {
	Payload interface{}

	label interface{}
	from  *Node
	to    *Node

	outEl    *list.Element
	allOutEl *list.Element
	inEl     *list.Element
	allInEl  *list.Element
}

// Label returns the edge label. Label may be nil
func (edge *Edge) Label() interface{} { return edge.label }

// From returns the source node for the edge. This cannot be nil.
func (edge *Edge) From() *Node { return edge.from }

// To returns the target node for the edge. This cannot be nil.
func (edge *Edge) To() *Node { return edge.to }

// Label returns the node label. Label may be nil
func (node *Node) Label() interface{} { return node.label }

func (g *Graph) init() *Graph {
	if g.nodesByLabel == nil {
		g.nodesByLabel = make(map[interface{}]*list.List)
	}
	return g
}

// New returns a new empty graph
func New() *Graph {
	return new(Graph).init()
}

// Len returns the number of nodes in the graph
func (g *Graph) Len() int { return g.allNodes.Len() }

// AllNodes returns an iterator over all nodes of a graph
func (g *Graph) AllNodes() Nodes {
	return &listNodes{at: g.allNodes.Front()}
}

// AllNodesWithLabel returns an iterator over all nodes with the given label
func (g *Graph) AllNodesWithLabel(label interface{}) Nodes {
	lst := g.nodesByLabel[label]
	if lst == nil {
		return emptyNodes{}
	}
	return &listNodes{at: lst.Front()}
}

// NewNode creates a new node with the given label and payload. Both
// parameters are optional and can be nil. Returns the new node in the graph.
// This method runs in constant-time.
func (g *Graph) NewNode(label, payload interface{}) *Node {
	g.init()
	node := &Node{Payload: payload,
		label: label,
		graph: g,
		out:   make(map[interface{}]*list.List),
		in:    make(map[interface{}]*list.List),
	}
	llist := g.nodesByLabel[label]
	if llist == nil {
		llist = list.New()
		g.nodesByLabel[label] = llist
	}
	node.nodesByLabelEl = llist.PushBack(node)
	node.allNodesEl = g.allNodes.PushBack(node)
	return node
}

// NewEdge creates a new edge between the two nodes. Both nodes must
// be in the same graph g. From and to nodes can be the same node. An
// edge may have a label and payload. Both are optional, and can be nil.
// This method runs in constant-time.
//
// Returns the new edge.
func (g *Graph) NewEdge(from, to *Node, label, payload interface{}) *Edge {
	g.init()
	if from.graph != g {
		panic("from is not in the graph")
	}
	if to.graph != g {
		panic("to is not in the graph")
	}
	edge := &Edge{Payload: payload,
		label: label,
		from:  from,
		to:    to,
	}
	edge.attachFrom(from)
	edge.attachTo(to)
	return edge
}

func (edge *Edge) attachFrom(from *Node) {
	lst := from.out[edge.label]
	if lst == nil {
		lst = list.New()
		from.out[edge.label] = lst
	}
	edge.outEl = lst.PushBack(edge)
	edge.allOutEl = from.allOut.PushBack(edge)
}

func (edge *Edge) attachTo(to *Node) {
	lst := to.in[edge.label]
	if lst == nil {
		lst = list.New()
		to.in[edge.label] = lst
	}
	edge.inEl = lst.PushBack(edge)
	edge.allInEl = to.allIn.PushBack(edge)
}

// Remove the edge. The edge is removed from the source and target
// nodes. This method runs in constant-time.
func (edge *Edge) Remove() {
	edge.detachFrom()
	edge.detachTo()
}

func (edge *Edge) detachFrom() {
	if edge.from != nil {
		lst := edge.from.out[edge.label]
		lst.Remove(edge.outEl)
		if lst.Len() == 0 {
			delete(edge.from.out, edge.label)
		}
		edge.from.allOut.Remove(edge.allOutEl)
		edge.from = nil
	}
}

func (edge *Edge) detachTo() {
	if edge.to != nil {
		lst := edge.to.in[edge.label]
		lst.Remove(edge.inEl)
		if lst.Len() == 0 {
			delete(edge.to.in, edge.label)
		}
		edge.to.allIn.Remove(edge.allInEl)
		edge.to = nil
	}
}

// SetTo redirects the target node of the edge
func (edge *Edge) SetTo(node *Node) {
	if node.graph != edge.from.graph {
		panic("Not in same graph")
	}
	edge.detachTo()
	edge.attachTo(node)
}

// SetFrom sets the source node of the edge
func (edge *Edge) SetFrom(node *Node) {
	if node.graph != edge.to.graph {
		panic("Not in same graph")
	}
	edge.detachFrom()
	edge.attachFrom(node)
}

// Remove the node. The node and all the edges incoming and outgoing
// from this node are also removed. This method runs in O(n) time
// where n is the number of adjacent edges.
func (node *Node) Remove() {
	if node.graph == nil {
		panic("Node is not in a graph")
	}
	for edge := node.allOut.Front(); edge != nil; edge = node.allOut.Front() {
		edge.Value.(*Edge).Remove()
	}
	for edge := node.allIn.Front(); edge != nil; edge = node.allIn.Front() {
		edge.Value.(*Edge).Remove()
	}
	lst := node.graph.nodesByLabel[node.label]
	lst.Remove(node.nodesByLabelEl)
	if lst.Len() == 0 {
		delete(node.graph.nodesByLabel, node.label)
	}
	node.graph.allNodes.Remove(node.allNodesEl)
	node.graph = nil
}

// NextNode returns the next node reached following the edge with the
// given label. If there is no such node, returns nil. If there are
// multiple, panics. This runs in constant time.
func (node *Node) NextNode(label interface{}) *Node {
	nxt := node.out[label]
	if nxt == nil {
		return nil
	}
	if nxt.Len() > 1 {
		panic("Multiple edges with given label")
	}
	return nxt.Front().Value.(*Edge).to
}

// PrevNode returns the node that reaches this node following the edge
// with the given label. If there is no such node, returns nil. If
// there are multiple, panics. This runs in constant-time.
func (node *Node) PrevNode(label interface{}) *Node {
	prv := node.in[label]
	if prv == nil {
		return nil
	}
	if prv.Len() > 1 {
		panic("Multiple edges with given label")
	}
	return prv.Front().Value.(*Edge).from
}

// AllOutgoingEdges returns an iterator over all outgoing edges of the
// node. Never returns nil.
func (node *Node) AllOutgoingEdges() Edges {
	return &listEdges{at: node.allOut.Front()}
}

// AllIncomingEdges returns an iterator over all the incoming edges of
// the node. Never returns nil.
func (node *Node) AllIncomingEdges() Edges {
	return &listEdges{at: node.allIn.Front()}
}

// AllOutgoingEdgesWithLabel returns an iterator over all outgoing
// edges with the given label. Never returns nil.
func (node *Node) AllOutgoingEdgesWithLabel(label interface{}) Edges {
	lst := node.out[label]
	if lst == nil {
		return emptyEdges{}
	}
	return &listEdges{at: lst.Front()}
}

// AllIncomingEdgesWithLabel returns an iterator over all incoming
// edges with the given label. Never returns nil.
func (node *Node) AllIncomingEdgesWithLabel(label interface{}) Edges {
	lst := node.in[label]
	if lst == nil {
		return emptyEdges{}
	}
	return &listEdges{at: lst.Front()}
}
