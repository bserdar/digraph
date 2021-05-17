# Digraph

[![GoDoc](https://godoc.org/github.com/bserdar/digraph?status.svg)](https://godoc.org/github.com/bserdar/digraph)
[![Go Report](https://goreportcard.com/badge/github.com/bserdar/digraph)](https://goreportcard.com/report/github.com/bserdar/digraph)

This package provides the `digraph` package that implements a directed
graph data structure. Nodes and edges may be labeled, and they may
have application-specific payloads.

Node and edge creations, accessor functions, iterators run in
constant-time. Node deletion is O(n) where n is the number of adjacent
edges, because deleting a node requires deleting all adjacent
edges. Deleting an edge is constant-time.

Digraph is not thread-safe. If you need thread-safety, you have to
explicitly implement it yourself.

## Example

Construct a graph, and add nodes:

```
g:=digraph.New()
node1:=digraph.NewNode("node1",nil)
node2:=digraph.NewNode("node2",nil)
```

You can connect the nodes with edges:

```
// edge node1 --edge1--> node2
edge:=g.NewEdge(node1,node2,"edge1",nil)
```

You can get the nodes accessible via a label:

```
n:=node1.NextNode("edge1")
```

If there are multiple, you can iterate:

```
edges:=node1.AllOutgoingEdgesWithLabel("edge1")
for edges.HasNext() {
   edge:=edges.Next()
}
```

You can remove a node or an edge. When you remove a node, all edges
adjacent to that node are also removed from the graph.

```
node1.Remove()
```

