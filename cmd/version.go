package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is a version number.
var version = "0.4.8"

func init() {
	RootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("%s\n", version)
		},
	}

	return cmd
}
