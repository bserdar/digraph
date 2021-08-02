package digraph

// Edge represents a labeled or unlabeled directed edge between two
// nodes of a graph
type Edge interface {
	// Return the label of the edge. Label cannot be changed. Remove the
	// edge and add a new one with a different label
	GetLabel() interface{}
	// Return the target node
	GetTo() Node
	// Return the source node
	GetFrom() Node

	// Remove an edge
	Disconnect()

	getEdgeHeader() *EdgeHeader
}

// EdgeHeader must be embedded into every edge implementation
type EdgeHeader struct {
	to    Node
	from  Node
	label interface{}
	edge  Edge
}

// NewEdgeHeader returns a new constructed edge header with the given
// label
func NewEdgeHeader(label interface{}) EdgeHeader {
	return EdgeHeader{label: label}
}

func (hdr *EdgeHeader) getEdgeHeader() *EdgeHeader {
	return hdr
}

// GetLabel returns the edge label. Once set, label cannot be changed
func (hdr *EdgeHeader) GetLabel() interface{} {
	return hdr.label
}

// GetTo returns the target node of the edge
func (hdr *EdgeHeader) GetTo() Node {
	return hdr.to
}

// GetFrom returns the source node of the edge
func (hdr *EdgeHeader) GetFrom() Node {
	return hdr.from
}

// Disconnect an edge
func (hdr *EdgeHeader) Disconnect() {
	if hdr.from != nil && hdr.edge != nil {
		hdr.from.removeOutgoingEdge(hdr.edge)
	}
	hdr.edge = nil
	hdr.from = nil
	hdr.to = nil
}

// Connect two nodes with the given edge. The edge must not be connected before
func Connect(from, to Node, edge Edge) {
	hdr := edge.getEdgeHeader()
	if hdr.edge != nil {
		panic("Edge is already connected")
	}
	hdr.edge = edge
	hdr.to = to
	hdr.from = from
	from.addOutgoingEdge(edge)
}

// BasicEdge contains an application-defined payload
type BasicEdge struct {
	EdgeHeader
	Payload interface{}
}

// NewBasicEdge returns a new unconnected edge with the given label and payload
func NewBasicEdge(label, payload interface{}) *BasicEdge {
	return &BasicEdge{
		EdgeHeader: EdgeHeader{
			label: label,
		},
		Payload: payload,
	}
}
