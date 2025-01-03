package validate

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates state files",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("validate state files!")
			return nil
		},
	}

	cmd.AddCommand(NewLintCommand())
	return cmd
}
