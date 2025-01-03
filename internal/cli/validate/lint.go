package validate

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/spf13/cobra"
	"io"
	"os"
	"turbolytics/state.d/internal"
)

func NewLintCommand() *cobra.Command {
	var filePath string

	var cmd = &cobra.Command{
		Use:   "lint",
		Short: "Lint a graphviz file",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			bs, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			graphAst, _ := gographviz.ParseString(string(bs))
			graph := gographviz.NewGraph()
			if err := gographviz.Analyse(graphAst, graph); err != nil {
				return err
			}

			valid, err := internal.Validate(graph)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", graph)

			fmt.Printf("Graph is valid: %t\n", valid)

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to graph file")
	cmd.MarkFlagRequired("file")

	return cmd
}
