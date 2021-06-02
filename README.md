# Digraph

[![GoDoc](https://godoc.org/github.com/bserdar/digraph?status.svg)](https://godoc.org/github.com/bserdar/digraph)
[![Go Report](https://goreportcard.com/badge/github.com/bserdar/digraph)](https://goreportcard.com/report/github.com/bserdar/digraph)

This package provides the `digraph` package that implements a directed
graph data structure. Nodes and edges may be labeled. The graph
supports application-defined structs as nodes and edges.

Node and edge addition, accessor functions, iterators run in
constant-time. Node deletion is O(n) where n is the number of adjacent
edges, because deleting a node requires deleting all adjacent
edges. Deleting an edge is constant-time.

Digraph is not thread-safe. 

## Example

Construct a graph, and add nodes:

```
g:=digraph.New()
node1:=digraph.NewBasicNode("node1",nil)
node2:=digraph.NewBasicNode("node2",nil)
```

Connect the nodes with edges:

```
// edge node1 --edge1--> node2
edge:=g.NewBasicEdge(node1,node2,"edge1",nil)
```

Get the nodes accessible via a label:

```
n:=node1.NextNode("edge1")
```

If there are multiple, iterate:

```
edges:=node1.AllOutgoingEdgesWithLabel("edge1")
for edges.HasNext() {
   edge:=edges.Next()
}
```

Removing a node removes all adjacent edges.

```
node1.Remove()
```

The `BasicNode` and `BasicEdge` contains an application-define
`Payload` field. This allows for `container/list` style graph
container. It is also possible to use application-specific node and
edge objects.

```
type CustomNode struct {
  digraph.NodeHeader
  // other fields
}

type CustomEdge struct {
  digraph.EdgeHeader
  // other fields
}
```

The `*CustomNode` and `*CustomEdge` instances can be added to the
graph. The embeded `NodeHeader` and `EdgeHeader` connects the object
to the underlying graph.
