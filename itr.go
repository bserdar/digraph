package digraph

import (
	"container/list"
)

type Nodes interface {
	HasNext() bool
	Next() *Node
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

type Edges interface {
	HasNext() bool
	Next() *Edge
	All() []*Edge
}

type emptyEdges struct{}

func (emptyEdges) HasNext() bool { return false }
func (emptyEdges) Next() *Edge   { panic("No more edges") }
func (emptyEdges) All() []*Edge  { return nil }

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
