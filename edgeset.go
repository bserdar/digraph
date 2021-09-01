package digraph

// Most nodes have few edges. For those nodes, a slice is a better
// storage than a map. The edgeSet interface allows usings both slices
// and maps as needed
type edgeSet interface {
	addEdge(Edge)
	hasEdges() bool
	removeEdge(Edge)
	getEdges() Edges
	getEdgesWith(interface{}) Edges
	length() int
	next() []Node
	nextWith(interface{}) []Node
}

type sliceEdgeSet []Edge

func newSliceEdgeSet() *sliceEdgeSet { return &sliceEdgeSet{} }

func (set sliceEdgeSet) length() int { return len(set) }

func (set *sliceEdgeSet) addEdge(e Edge) {
	*set = append(*set, e)
}

func (set sliceEdgeSet) hasEdges() bool { return len(set) > 0 }

func (set *sliceEdgeSet) removeEdge(e Edge) {
	if set == nil {
		return
	}
	w := 0
	for _, edge := range *set {
		if edge != e {
			(*set)[w] = edge
			w++
		}
	}
	*set = (*set)[:w]
}

func (set sliceEdgeSet) getEdges() Edges { return Edges{&edgeSliceIterator{Edges: set}} }

func (set sliceEdgeSet) getEdgesWith(label interface{}) Edges {
	return Edges{&edgeSliceIterator{Edges: set}}.Select(func(e Edge) bool { return e.GetLabel() == label })
}

func (set sliceEdgeSet) toMap() *mapEdgeSet {
	ret := newMapEdgeSet()
	ret.s = set
	for _, edge := range set {
		ret.m[edge.GetLabel()] = append(ret.m[edge.GetLabel()], edge)
	}
	return ret
}

func (set sliceEdgeSet) next() []Node {
	switch len(set) {
	case 0:
		return nil
	case 1:
		return []Node{set[0].GetTo()}
	default:
		ret := make([]Node, 0, len(set))
		seen := make(map[Node]struct{})
		for _, e := range set {
			to := e.GetTo()
			if _, ok := seen[to]; !ok {
				ret = append(ret, to)
				seen[to] = struct{}{}
			}
		}
		return ret
	}
}

func (set sliceEdgeSet) nextWith(label interface{}) []Node {
	switch len(set) {
	case 0:
		return nil
	case 1:
		if set[0].GetLabel() == label {
			return []Node{set[0].GetTo()}
		}
		return nil
	default:
		ret := make([]Node, 0, len(set))
		seen := make(map[Node]struct{})
		for _, e := range set {
			if e.GetLabel() == label {
				to := e.GetTo()
				if _, ok := seen[to]; !ok {
					ret = append(ret, to)
					seen[to] = struct{}{}
				}
			}
		}
		return ret
	}
}

type mapEdgeSet struct {
	m map[interface{}]sliceEdgeSet
	s sliceEdgeSet
}

func newMapEdgeSet() *mapEdgeSet {
	return &mapEdgeSet{m: make(map[interface{}]sliceEdgeSet)}
}

func (set *mapEdgeSet) addEdge(e Edge) {
	label := e.GetLabel()
	set.m[label] = append(set.m[label], e)
	set.s.addEdge(e)
}

func (set mapEdgeSet) hasEdges() bool { return len(set.s) > 0 }

func (set *mapEdgeSet) removeEdge(e Edge) {
	label := e.GetLabel()
	x := set.m[label]
	if x == nil {
		return
	}
	x.removeEdge(e)
	if len(x) == 0 {
		delete(set.m, label)
	} else {
		set.m[label] = x
	}
	set.s.removeEdge(e)
}

func (set mapEdgeSet) getEdgesWith(label interface{}) Edges {
	return set.m[label].getEdges()
}

func (set mapEdgeSet) getEdges() Edges {
	return Edges{&edgeSliceIterator{Edges: set.s}}
}

func (set mapEdgeSet) length() int { return len(set.s) }

func (set mapEdgeSet) next() []Node {
	return set.s.next()
}

func (set mapEdgeSet) nextWith(label interface{}) []Node {
	m := set.m[label]
	if m == nil {
		return nil
	}
	return m.nextWith(label)
}
