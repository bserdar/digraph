# Digraph

[![GoDoc](https://godoc.org/github.com/bserdar/digraph?status.svg)](https://godoc.org/github.com/bserdar/digraph)
[![Go Report](https://goreportcard.com/badge/github.com/bserdar/digraph)](https://goreportcard.com/report/github.com/bserdar/digraph)

This package provides the `digraph` package that implements a directed
graph data structure. Nodes and edges may be labeled. The graph
supports application-defined structs as nodes and edges.


The graph structure is designed so that nodes know the outgoing edges,
and edges know both the source and target nodes. The graph structure
itself knows only "some" of the nodes, so retrieving all the nodes of
the graph or accessing nodes by label requires an intermediate
structure, the `NodeIndex`. A `NodeIndex` discovers all nodes when
requested, and then provides indexes access to the nodes. A
`NodeIndex` only sees the nodes that were accessible from the `Graph`
when it is created, thus it does not provide a dynamic view of the
graph.

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
edge:=g.NewBasicEdge("edge1",nil)
digraph.Connect(node1,node2,edge)
```

Get the nodes accessible via a label:

```
n:=node1.Next("edge1")
```

If there are multiple, iterate:

```
edges:=node1.GetAllOutgoingEdgesWithLabel("edge1")
for edges.HasNext() {
   edge:=edges.Next()
}
```

The `BasicNode` and `BasicEdge` contains an application-defined
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
