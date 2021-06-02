package digraph

import (
	"container/list"
)

// Graph represents a labeled or unlabeled directed graph. Zero value
// for Graph is an empty graph ready to use.
//
// The nodes of the graph are instances of the Node interface. Each node
// is a pointer to a struct that includes a NodeHeader.
//
// The edges between nodes are instances of the Edge interface. Each
// edge is a pointer to a struct that includes an EdgeHeader
type Graph struct {
	nodesByLabel map[interface{}]*list.List
	allNodes     list.List
}

// Node is a node of a directed graph.
//
// Node interface provides a way to get the node header from the underlying struct
//
// A Node may have a label. Multiple nodes can have the same label in
// a graph.
//
// Node keeps a list/map of all outgoing and incoming edges. These
// edges may or may not be labeled. Same label can be used for
// outgoing and incoming edges.
//
// A Node object can belong to at most one graph. The node must be a
// pointer to a struct.
//
type Node interface {
	GetNodeHeader() *NodeHeader

	Label() interface{}
	Remove()
	SetLabel(label interface{})
	NextNode(label interface{}) Node
	PrevNode(label interface{}) Node
	AllOutgoingEdges() Edges
	AllIncomingEdges() Edges
	AllOutgoingEdgesWithLabel(label interface{}) Edges
	AllIncomingEdgesWithLabel(label interface{}) Edges
}

// Edge represents a labeled or unlabeled directed edge between two
// nodes of a graph
//
// An edge has from and to nodes, which are defined when edge is
// created, and connot be changed.
//
// An edge may have a label. Same label can be reused on multiple
// edges outgoing from a node.
//
type Edge interface {
	GetEdgeHeader() *EdgeHeader

	Label() interface{}
	SetLabel(interface{})
	From() Node
	To() Node
	Remove()
	SetTo(node Node)
	SetFrom(node Node)
}

// NodeHeader keeps the links of the node within the graph. Every node
// of the graph must include a NodeHeader
type NodeHeader struct {
	label interface{}

	nodesByLabelEl *list.Element
	allNodesEl     *list.Element
	graph          *Graph
	node           Node

	out    map[interface{}]*list.List
	allOut list.List
	in     map[interface{}]*list.List
	allIn  list.List
}

// Label returns the node label. Label may be nil
func (node *NodeHeader) Label() interface{} { return node.label }

// GetNodeHeader returns the node header
func (node *NodeHeader) GetNodeHeader() *NodeHeader { return node }

func (node *NodeHeader) init() {
	node.out = make(map[interface{}]*list.List)
	node.allOut = list.List{}
	node.in = make(map[interface{}]*list.List)
	node.allIn = list.List{}
}

// BasicNode contains an application defined payload
type BasicNode struct {
	NodeHeader
	Payload interface{}
}

// EdgeHeader keeps the links of an edge within the graph. Every edge
// of the graph must include an EdgeHeader
type EdgeHeader struct {
	label interface{}
	from  Node
	to    Node
	edge  Edge

	outEl    *list.Element
	allOutEl *list.Element
	inEl     *list.Element
	allInEl  *list.Element
}

// Label returns the edge label. Label may be nil
func (edge *EdgeHeader) Label() interface{} { return edge.label }

// From returns the source node for the edge. This cannot be nil.
func (edge *EdgeHeader) From() Node { return edge.from }

// To returns the target node for the edge. This cannot be nil.
func (edge *EdgeHeader) To() Node { return edge.to }

func (edge *EdgeHeader) GetEdgeHeader() *EdgeHeader { return edge }

// BasicEdge contains an application-defined payload
type BasicEdge struct {
	EdgeHeader
	Payload interface{}
}

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
	g.init()
	lst := g.nodesByLabel[label]
	if lst == nil {
		return emptyNodes{}
	}
	return &listNodes{at: lst.Front()}
}

// NewBasicNode creates a new BasicNode with the given label and payload. Both
// parameters are optional and can be nil. Returns the new node
func NewBasicNode(label, payload interface{}) *BasicNode {
	node := &BasicNode{Payload: payload}
	node.label = label
	return node
}

// AddNode adds the node to the graph. The node must not belong to another graph
func (g *Graph) AddNode(node Node) {
	g.init()
	nh := node.GetNodeHeader()
	if nh.graph != nil {
		panic("Node belongs to a graph already")
	}
	nh.init()
	nh.graph = g
	nh.node = node
	llist := g.nodesByLabel[nh.label]
	if llist == nil {
		llist = list.New()
		g.nodesByLabel[nh.label] = llist
	}
	nh.nodesByLabelEl = llist.PushBack(node)
	nh.allNodesEl = g.allNodes.PushBack(node)
}

// NewBasicEdge creates a new basic edge with a label and
// payload. Both are optional, and can be nil.
//
// Returns the new edge.
func NewBasicEdge(label, payload interface{}) *BasicEdge {
	edge := &BasicEdge{Payload: payload}
	edge.label = label
	return edge
}

// AddEdge adds the given edge to the graph.
func (g *Graph) AddEdge(from, to Node, edge Edge) {
	g.init()
	if from.GetNodeHeader().graph != g {
		panic("from is not in the graph")
	}
	if to.GetNodeHeader().graph != g {
		panic("to is not in the graph")
	}
	edge.GetEdgeHeader().edge = edge
	edge.GetEdgeHeader().attachFrom(from)
	edge.GetEdgeHeader().attachTo(to)
}

func (edgehdr *EdgeHeader) attachFrom(from Node) {
	fromhdr := from.GetNodeHeader()
	lst := fromhdr.out[edgehdr.label]
	if lst == nil {
		lst = list.New()
		fromhdr.out[edgehdr.label] = lst
	}
	edgehdr.from = from
	edgehdr.outEl = lst.PushBack(edgehdr.edge)
	edgehdr.allOutEl = fromhdr.allOut.PushBack(edgehdr.edge)
}

func (edgehdr *EdgeHeader) attachTo(to Node) {
	tohdr := to.GetNodeHeader()
	lst := tohdr.in[edgehdr.label]
	if lst == nil {
		lst = list.New()
		tohdr.in[edgehdr.label] = lst
	}
	edgehdr.to = to
	edgehdr.inEl = lst.PushBack(edgehdr.edge)
	edgehdr.allInEl = tohdr.allIn.PushBack(edgehdr.edge)
}

// Remove the edge. The edge is removed from the source and target
// nodes. This method runs in constant-time.
func (edgehdr *EdgeHeader) Remove() {
	edgehdr.detachFrom()
	edgehdr.detachTo()
}

func (edgehdr *EdgeHeader) detachFrom() {
	if edgehdr.from != nil {
		lst := edgehdr.from.GetNodeHeader().out[edgehdr.label]
		lst.Remove(edgehdr.outEl)
		if lst.Len() == 0 {
			delete(edgehdr.from.GetNodeHeader().out, edgehdr.label)
		}
		edgehdr.from.GetNodeHeader().allOut.Remove(edgehdr.allOutEl)
		edgehdr.from = nil
	}
}

func (edgehdr *EdgeHeader) detachTo() {
	if edgehdr.to != nil {
		lst := edgehdr.to.GetNodeHeader().in[edgehdr.label]
		lst.Remove(edgehdr.inEl)
		if lst.Len() == 0 {
			delete(edgehdr.to.GetNodeHeader().in, edgehdr.label)
		}
		edgehdr.to.GetNodeHeader().allIn.Remove(edgehdr.allInEl)
		edgehdr.to = nil
	}
}

// SetTo redirects the target node of the edge
func (edgehdr *EdgeHeader) SetTo(node Node) {
	if node.GetNodeHeader().graph != edgehdr.from.GetNodeHeader().graph {
		panic("Not in same graph")
	}
	edgehdr.detachTo()
	edgehdr.attachTo(node)
}

// SetFrom sets the source node of the edge
func (edgehdr *EdgeHeader) SetFrom(node Node) {
	if node.GetNodeHeader().graph != edgehdr.to.GetNodeHeader().graph {
		panic("Not in same graph")
	}
	edgehdr.detachFrom()
	edgehdr.attachFrom(node)
}

// SetLabel sets the edge label
func (edgehdr *EdgeHeader) SetLabel(label interface{}) {
	from := edgehdr.from
	to := edgehdr.to
	edgehdr.detachFrom()
	edgehdr.detachTo()
	edgehdr.label = label
	if from != nil {
		edgehdr.attachFrom(from)
	}
	if to != nil {
		edgehdr.attachTo(to)
	}
}

// Remove the node. The node and all the edges incoming and outgoing
// from this node are also removed. This method runs in O(n) time
// where n is the number of adjacent edges.
func (nodehdr *NodeHeader) Remove() {
	if nodehdr.graph == nil {
		panic("Node is not in a graph")
	}
	for edge := nodehdr.allOut.Front(); edge != nil; edge = nodehdr.allOut.Front() {
		edge.Value.(Edge).GetEdgeHeader().Remove()
	}
	for edge := nodehdr.allIn.Front(); edge != nil; edge = nodehdr.allIn.Front() {
		edge.Value.(Edge).GetEdgeHeader().Remove()
	}
	lst := nodehdr.graph.nodesByLabel[nodehdr.label]
	lst.Remove(nodehdr.nodesByLabelEl)
	if lst.Len() == 0 {
		delete(nodehdr.graph.nodesByLabel, nodehdr.label)
	}
	nodehdr.graph.allNodes.Remove(nodehdr.allNodesEl)
	nodehdr.graph = nil
}

// SetLabel sets the label of the node
func (nodehdr *NodeHeader) SetLabel(label interface{}) {
	if nodehdr.graph != nil {
		lst := nodehdr.graph.nodesByLabel[nodehdr.label]
		lst.Remove(nodehdr.nodesByLabelEl)
	}
	nodehdr.label = label
	if nodehdr.graph != nil {
		lst := nodehdr.graph.nodesByLabel[nodehdr.label]
		if lst == nil {
			lst = list.New()
			nodehdr.graph.nodesByLabel[nodehdr.label] = lst
		}
		nodehdr.nodesByLabelEl = lst.PushBack(nodehdr.node)
	}
}

// NextNode returns the next node reached following the edge with the
// given label. If there is no such node, returns nil. If there are
// multiple, panics. This runs in constant time.
func (nodehdr *NodeHeader) NextNode(label interface{}) Node {
	nxt := nodehdr.out[label]
	if nxt == nil {
		return nil
	}
	if nxt.Len() > 1 {
		panic("Multiple edges with given label")
	}
	return nxt.Front().Value.(Edge).GetEdgeHeader().to
}

// PrevNode returns the node that reaches this node following the edge
// with the given label. If there is no such node, returns nil. If
// there are multiple, panics. This runs in constant-time.
func (nodehdr *NodeHeader) PrevNode(label interface{}) Node {
	prv := nodehdr.in[label]
	if prv == nil {
		return nil
	}
	if prv.Len() > 1 {
		panic("Multiple edges with given label")
	}
	return prv.Front().Value.(Edge).GetEdgeHeader().from
}

// AllOutgoingEdges returns an iterator over all outgoing edges of the
// node. Never returns nil.
func (nodehdr *NodeHeader) AllOutgoingEdges() Edges {
	return &listEdges{at: nodehdr.allOut.Front()}
}

// AllIncomingEdges returns an iterator over all the incoming edges of
// the node. Never returns nil.
func (nodehdr *NodeHeader) AllIncomingEdges() Edges {
	return &listEdges{at: nodehdr.allIn.Front()}
}

// AllOutgoingEdgesWithLabel returns an iterator over all outgoing
// edges with the given label. Never returns nil.
func (nodehdr *NodeHeader) AllOutgoingEdgesWithLabel(label interface{}) Edges {
	lst := nodehdr.out[label]
	if lst == nil {
		return emptyEdges{}
	}
	return &listEdges{at: lst.Front()}
}

// AllIncomingEdgesWithLabel returns an iterator over all incoming
// edges with the given label. Never returns nil.
func (nodehdr *NodeHeader) AllIncomingEdgesWithLabel(label interface{}) Edges {
	lst := nodehdr.in[label]
	if lst == nil {
		return emptyEdges{}
	}
	return &listEdges{at: lst.Front()}
}
