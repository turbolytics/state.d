package internal

import (
	"errors"
	"testing"

	"github.com/awalterschulze/gographviz"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		graphStr string
		expected bool
		err      error
	}{
		{
			name: "Valid graph with s_init and s_end",
			graphStr: `
				digraph G {
					s_init [label="start"];
					s_end [label="end"];
					s_init -> s_end;
				}`,
			expected: true,
			err:      nil,
		},
		{
			name: "Graph without s_init node",
			graphStr: `
				digraph G {
					s_end [label="end"];
				}`,
			expected: false,
			err:      errors.New("graph must have both 's_init' and 's_end' nodes"),
		},
		{
			name: "Graph without s_end node",
			graphStr: `
				digraph G {
					s_init [label="start"];
				}`,
			expected: false,
			err:      errors.New("graph must have both 's_init' and 's_end' nodes"),
		},
		{
			name: "Graph without s_init and s_end nodes",
			graphStr: `
				digraph G {
					node1 [label="node1"];
					node2 [label="node2"];
					node1 -> node2;
				}`,
			expected: false,
			err:      errors.New("graph must have both 's_init' and 's_end' nodes"),
		},
		{
			name: "Graph without s_init and s_end nodes",
			graphStr: `
			digraph G {
				node1 [label="node1"];
				node2 [label="node2"];
				node1 -> node2;
			}`,
			expected: false,
			err:      errors.New("graph must have both 's_init' and 's_end' nodes"),
		},
		{
			name: "Graph with incoming edge to s_init",
			graphStr: `
			digraph G {
				s_init [label="start"];
				s_end [label="end"];
				node1 [label="node1"];
				node1 -> s_init;
				s_init -> s_end;
			}`,
			expected: false,
			err:      errors.New("the 's_init' node must have no incoming edges"),
		},
		{
			name: "Graph with outgoing edge from s_end",
			graphStr: `
			digraph G {
				s_init [label="start"];
				s_end [label="end"];
				s_init -> s_end;
				s_end -> node1;
			}`,
			expected: false,
			err:      errors.New("the 's_end' node must have no outgoing edges"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphAst, _ := gographviz.ParseString(tt.graphStr)
			graph := gographviz.NewGraph()
			if err := gographviz.Analyse(graphAst, graph); err != nil {
				t.Fatalf("failed to parse graph: %v", err)
			}

			valid, err := Validate(graph)
			if valid != tt.expected || (err != nil && err.Error() != tt.err.Error()) {
				t.Errorf("Validate() = %v, %v; want %v, %v", valid, err, tt.expected, tt.err)
			}
		})
	}
}
