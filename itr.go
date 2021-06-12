package digraph

import (
	"container/list"
)

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

type listNodes struct {
	at *list.Element
}

func (l *listNodes) HasNext() bool { return l.at != nil }

func (l *listNodes) Next() Node {
	ret := l.at.Value.(Node)
	l.at = l.at.Next()
	return ret
}

type arrNodes struct {
	arr []Node
	at  int
}

func (a *arrNodes) HasNext() bool { return a.at < len(a.arr) }
func (a *arrNodes) Next() Node    { ret := a.arr[a.at]; a.at++; return ret }

type filterNodes struct {
	source Nodes
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
func NewNodes(nodes ...Node) Nodes { return Nodes{&arrNodes{arr: nodes}} }

// Unique filters the nodes so only unique nodes are returned
func (n Nodes) Uniqu() Nodes {
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
	return Nodes{&edgeNodeSelector{source: e, selectNode: func(e Edge) Node { return e.To() }}}
}

// Sources returns a node iterator that goes through the source nodes
func (e Edges) Sources() Nodes {
	return Nodes{&edgeNodeSelector{source: e, selectNode: func(e Edge) Node { return e.From() }}}
}

// EdgeIterator iterates through a list of edges
type EdgeIterator interface {
	// Returns if there are more edges to go through
	HasNext() bool
	// If HasNext is true, returns the next edge and advances. Otherwise panics
	Next() Edge
}

type emptyEdges struct{}

func (emptyEdges) HasNext() bool { return false }
func (emptyEdges) Next() Edge    { panic("No more edges") }

type listEdges struct {
	at *list.Element
}

func (l *listEdges) HasNext() bool { return l.at != nil }

func (l *listEdges) Next() Edge {
	ret := l.at.Value.(Edge)
	l.at = l.at.Next()
	return ret
}
