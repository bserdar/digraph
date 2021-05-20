package digraph

import (
	"container/list"
)

// Nodes iterates through a list of nodes
type Nodes interface {
	// Returns if there are more nodes to go through
	HasNext() bool
	// If HasNext is true, returns the next node and advances. Otherwise, panics
	Next() *Node
	// Returns all remaining nodes
	All() []*Node
}

type emptyNodes struct{}

func (emptyNodes) HasNext() bool { return false }
func (emptyNodes) Next() *Node   { panic("No more nodes") }
func (emptyNodes) All() []*Node  { return nil }

type listNodes struct {
	at *list.Element
}

func (l *listNodes) HasNext() bool {
	return l.at != nil
}

func (l *listNodes) Next() *Node {
	ret := l.at.Value.(*Node)
	l.at = l.at.Next()
	return ret
}

func (l *listNodes) All() []*Node {
	ret := make([]*Node, 0)
	for ; l.at != nil; l.at = l.at.Next() {
		ret = append(ret, l.at.Value.(*Node))
	}
	return ret
}

type arrNodes struct {
	arr []*Node
	at  int
}

func (a *arrNodes) HasNext() bool { return a.at < len(a.arr) }
func (a *arrNodes) Next() *Node   { ret := a.arr[a.at]; a.at++; return ret }
func (a *arrNodes) All() []*Node  { return a.arr[a.at:] }

type edgeNodes struct {
	edge   *list.Element
	seen   map[*Node]struct{}
	node   func(*Edge) *Node
	seeked bool
}

func (a *edgeNodes) HasNext() bool {
	if a.seeked {
		return true
	}
	if a.seen == nil {
		a.seen = make(map[*Node]struct{})
	}
	for {
		if a.edge == nil {
			return false
		}
		node := a.node(a.edge.Value.(*Edge))
		if _, seen := a.seen[node]; !seen {
			a.seeked = true
			return true
		}
		a.edge = a.edge.Next()
	}
}

func (a *edgeNodes) Next() *Node {
	if !a.seeked {
		a.HasNext()
	}
	if a.seen == nil {
		a.seen = make(map[*Node]struct{})
	}
	node := a.node(a.edge.Value.(*Edge))
	a.seen[node] = struct{}{}
	a.seeked = false
	return node
}

func (a *edgeNodes) All() []*Node {
	ret := make([]*Node, 0)
	for a.HasNext() {
		ret = append(ret, a.Next())
	}
	return ret
}

// Edges iterates through a list of edges
type Edges interface {
	// Returns if there are more edges to go through
	HasNext() bool
	// If HasNext is true, returns the next edge and advances. Otherwise panics
	Next() *Edge
	// Returns all remaining edges
	All() []*Edge

	// Returns a node iterator that will go through each target node once
	Targets() Nodes
	// Returns a node iterator that will go through each source node once
	Sources() Nodes
}

type emptyEdges struct{}

func (emptyEdges) HasNext() bool  { return false }
func (emptyEdges) Next() *Edge    { panic("No more edges") }
func (emptyEdges) All() []*Edge   { return nil }
func (emptyEdges) Targets() Nodes { return &emptyNodes{} }
func (emptyEdges) Sources() Nodes { return &emptyNodes{} }

type listEdges struct {
	at *list.Element
}

func (l *listEdges) HasNext() bool {
	return l.at != nil
}

func (l *listEdges) Next() *Edge {
	ret := l.at.Value.(*Edge)
	l.at = l.at.Next()
	return ret
}

func (l *listEdges) All() []*Edge {
	ret := make([]*Edge, 0)
	for ; l.at != nil; l.at = l.at.Next() {
		ret = append(ret, l.at.Value.(*Edge))
	}
	return ret
}

func (l *listEdges) Targets() Nodes {
	return &edgeNodes{edge: l.at, node: func(e *Edge) *Node { return e.To() }, seeked: l.at != nil}
}

func (l *listEdges) Sources() Nodes {
	return &edgeNodes{edge: l.at, node: func(e *Edge) *Node { return e.From() }, seeked: l.at != nil}
}
