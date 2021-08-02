package digraph

// Node of a directed graph. The node can contain a label
type Node interface {
	GetLabel() interface{}
	SetLabel(interface{})

	// Returns if the node has any outgoing edges
	HasOutgoingEdges() bool
	// Returns all outgoing edges of the node
	GetAllOutgoingEdges() Edges
	// Returns the outgoing edges with the given label
	GetAllOutgoingEdgesWithLabel(interface{}) Edges
	// Returns the node accessed by following the edge with the given
	// label. If there are multiple, panics
	Next(label interface{}) Node

	RemoveOutgoingEdge(Edge)
	AddOutgoingEdge(Edge)
}

// NodeHeader must be embedded into every node
type NodeHeader struct {
	label         interface{}
	outgoingEdges map[interface{}]map[Edge]struct{}
}

func (hdr *NodeHeader) GetLabel() interface{} {
	return hdr.label
}

func (hdr *NodeHeader) SetLabel(label interface{}) {
	hdr.label = label
}

func (hdr *NodeHeader) AddOutgoingEdge(edge Edge) {
	if hdr.outgoingEdges == nil {
		hdr.outgoingEdges = make(map[interface{}]map[Edge]struct{})
	}
	m := hdr.outgoingEdges[edge.GetLabel()]
	if m == nil {
		m = make(map[Edge]struct{})
		hdr.outgoingEdges[edge.GetLabel()] = m
	}
	m[edge] = struct{}{}
}

func (hdr *NodeHeader) HasOutgoingEdges() bool {
	return len(hdr.outgoingEdges) > 0
}

func (hdr *NodeHeader) RemoveOutgoingEdge(edge Edge) {
	if hdr.outgoingEdges == nil {
		return
	}
	m := hdr.outgoingEdges[edge.GetLabel()]
	if m == nil {
		return
	}
	delete(m, edge)
	if len(m) == 0 {
		delete(hdr.outgoingEdges, edge.GetLabel())
	}
}

func (hdr *NodeHeader) GetAllOutgoingEdges() Edges {
	if hdr.outgoingEdges == nil {
		return Edges{&EdgeArrayIterator{}}
	}
	arr := make([]Edge, 0, len(hdr.outgoingEdges))
	for _, m := range hdr.outgoingEdges {
		for edge := range m {
			arr = append(arr, edge)
		}
	}
	return Edges{&EdgeArrayIterator{arr}}
}

func (hdr *NodeHeader) GetAllOutgoingEdgesWithLabel(label interface{}) Edges {
	if hdr.outgoingEdges == nil {
		return Edges{&EdgeArrayIterator{}}
	}
	m := hdr.outgoingEdges[label]
	arr := make([]Edge, 0, len(m))
	for edge := range m {
		arr = append(arr, edge)
	}
	return Edges{&EdgeArrayIterator{arr}}
}

// Next returns the node reached by following the edge with the given
// label. If there are none, returns nil. If there are multiple,
// panics
func (hdr *NodeHeader) Next(label interface{}) Node {
	if hdr.outgoingEdges == nil {
		return nil
	}
	m := hdr.outgoingEdges[label]
	switch len(m) {
	case 0:
		return nil
	case 1:
		for k := range m {
			return k.GetTo()
		}
	}
	panic("Multiple nodes for Next")
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
