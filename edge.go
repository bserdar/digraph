package digraph

// Edge represents a labeled or unlabeled directed edge between two
// nodes of a graph
type Edge interface {
	GetLabel() interface{}
	GetTo() Node
	GetFrom() Node

	getEdgeHeader() *EdgeHeader
}

// EdgeHeader must be embedded into every edge implementation
type EdgeHeader struct {
	to    Node
	from  Node
	label interface{}
	edge  Edge
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

// Connect two nodes with the given edge. The edge must not be connected before
func Connect(from, to Node, edge Edge) {
	hdr := edge.getEdgeHeader()
	if hdr.edge != nil {
		panic("Edge is already connected")
	}
	hdr.edge = edge
	hdr.to = to
	hdr.from = from
	from.AddOutgoingEdge(edge)
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
