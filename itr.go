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
	nodes := make(map[*Node]struct{})
	ret := &arrNodes{}
	for x := l.at; x != nil; x = x.Next() {
		node := x.Value.(*Edge).To()
		if _, exists := nodes[node]; !exists {
			nodes[node] = struct{}{}
			ret.arr = append(ret.arr, node)
		}
	}
	return ret
}

func (l *listEdges) Sources() Nodes {
	nodes := make(map[*Node]struct{})
	ret := &arrNodes{}
	for x := l.at; x != nil; x = x.Next() {
		node := x.Value.(*Edge).From()
		if _, exists := nodes[node]; !exists {
			nodes[node] = struct{}{}
			ret.arr = append(ret.arr, node)
		}
	}
	return ret
}
