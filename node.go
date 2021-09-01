package digraph

// Node of a directed graph. The node can contain a label
type Node interface {
	GetLabel() interface{}
	SetLabel(interface{})

	// Returns if the node has any outgoing edges
	HasOut() bool
	// Returns all outgoing edges of the node
	Out() Edges
	// Returns the outgoing edges with the given label
	OutWith(interface{}) Edges

	// Returns all directly accessible nodes
	Next() []Node

	// Returns all directly accessible nodes with label
	NextWith(interface{}) []Node

	removeOutgoingEdge(Edge)
	addOutgoingEdge(Edge)
}

// NodeHeader must be embedded into every node
type NodeHeader struct {
	label interface{}
	out   edgeSet
}

// GetLabel returns the node label
func (hdr *NodeHeader) GetLabel() interface{} {
	return hdr.label
}

// SetLabel sets the node label
func (hdr *NodeHeader) SetLabel(label interface{}) {
	hdr.label = label
}

// addOutgoingEdge adds a new outgoing edge to this node. The edge must be disconnected.
func (hdr *NodeHeader) addOutgoingEdge(edge Edge) {
	if hdr.out == nil {
		hdr.out = newSliceEdgeSet()
	}
	hdr.out.addEdge(edge)
	if sl, ok := hdr.out.(*sliceEdgeSet); ok {
		if sl.length() > 10 {
			hdr.out = sl.toMap()
		}
	}
}

// Next returns all next nodes
func (hdr *NodeHeader) Next() []Node {
	if hdr.out == nil {
		return nil
	}
	return hdr.out.next()
}

// NextWith returns all next nodes reachable with label
func (hdr *NodeHeader) NextWith(label interface{}) []Node {
	if hdr.out == nil {
		return nil
	}
	return hdr.out.nextWith(label)
}

// HasOut returns true if the node has outgoing edges
func (hdr *NodeHeader) HasOut() bool {
	if hdr.out == nil {
		return false
	}
	return hdr.out.hasEdges()
}

func (hdr *NodeHeader) removeOutgoingEdge(edge Edge) {
	if hdr.out == nil {
		return
	}
	hdr.out.removeEdge(edge)
}

// Out returns all outgoing edges of the node
func (hdr *NodeHeader) Out() Edges {
	if hdr.out == nil {
		return Edges{&edgeSliceIterator{}}
	}
	return hdr.out.getEdges()
}

// OutWith returns all edges with the given label
func (hdr *NodeHeader) OutWith(label interface{}) Edges {
	if hdr.out == nil {
		return Edges{&edgeSliceIterator{}}
	}
	return hdr.out.getEdgesWith(label)
}

// BasicNode contains an application defined payload
type BasicNode struct {
	NodeHeader
	Payload interface{}
}

func NewBasicNode(label, payload interface{}) *BasicNode {
	return &BasicNode{
		NodeHeader: NodeHeader{
			label: label,
		},
		Payload: payload,
	}
}
