package internal

import (
	"errors"
	"github.com/awalterschulze/gographviz"
)

type AdjacencyList map[*gographviz.Node][]*gographviz.Node

type NodeIdentifier string

const (
	NodeInit NodeIdentifier = "s_init"
	NodeEnd  NodeIdentifier = "s_end"
)

// Validator interface defines a method for validating a graph
type Validator interface {
	Validate(graph *gographviz.Graph, outgoing AdjacencyList, incoming AdjacencyList) (bool, error)
}

type EndIsLeafValidator struct{}

func (v *EndIsLeafValidator) Validate(graph *gographviz.Graph, outgoing AdjacencyList, _ AdjacencyList) (bool, error) {
	n, ok := graph.Nodes.Lookup[string(NodeEnd)]
	if !ok {
		return false, errors.New("the 's_end' node must exist")
	}
	edges, ok := outgoing[n]
	if !ok {
		return false, errors.New("the 's_end' node must exist")
	}

	if len(edges) > 0 {
		return false, errors.New("the 's_end' node must have no outgoing edges")
	}

	return true, nil
}

// NoIncomingEdgesToInitValidator checks if the "start" node has no incoming edges
type NoIncomingEdgesToInitValidator struct{}

func (v *NoIncomingEdgesToInitValidator) Validate(graph *gographviz.Graph, _ AdjacencyList, incoming AdjacencyList) (bool, error) {
	n, ok := graph.Nodes.Lookup[string(NodeInit)]
	if !ok {
		return false, errors.New("the 's_init' node must exist")
	}
	edges, ok := incoming[n]
	if !ok {
		return false, errors.New("the 's_init' node must exist")
	}

	if len(edges) > 0 {
		return false, errors.New("the 's_init' node must have no incoming edges")
	}

	return true, nil
}

// HasInitAndEndValidator checks if the graph has both "start" and "end" nodes
type HasInitAndEndValidator struct{}

func (v *HasInitAndEndValidator) Validate(graph *gographviz.Graph, outgoing AdjacencyList, incoming AdjacencyList) (bool, error) {
	_, hasInit := graph.Nodes.Lookup[string(NodeInit)]
	_, hasEnd := graph.Nodes.Lookup[string(NodeEnd)]

	if !hasInit || !hasEnd {
		return false, errors.New("graph must have both 's_init' and 's_end' nodes")
	}

	return true, nil
}

// Validate runs all validators against the input graph
func Validate(graph *gographviz.Graph) (bool, error) {
	validators := []Validator{
		&HasInitAndEndValidator{},
		&NoIncomingEdgesToInitValidator{},
		&EndIsLeafValidator{},
	}

	outgoing := make(map[*gographviz.Node][]*gographviz.Node)
	incoming := make(map[*gographviz.Node][]*gographviz.Node)
	for _, node := range graph.Nodes.Nodes {
		outgoing[node] = []*gographviz.Node{}
		incoming[node] = []*gographviz.Node{}
	}

	for _, edge := range graph.Edges.Edges {
		src := graph.Nodes.Lookup[edge.Src]
		dst := graph.Nodes.Lookup[edge.Dst]
		outgoing[src] = append(outgoing[src], dst)
		incoming[dst] = append(incoming[dst], src)
	}

	for _, validator := range validators {
		if valid, err := validator.Validate(graph, outgoing, incoming); !valid {
			return false, err
		}
	}

	return true, nil
}
