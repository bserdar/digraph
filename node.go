package digraph

// Node of a directed graph. The node can contain a label
type Node interface {
	GetLabel() interface{}

	HasOutgoingEdges() bool
	GetAllOutgoingEdges() Edges
	GetAllOutgoingEdgesWithLabel(interface{}) Edges

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
