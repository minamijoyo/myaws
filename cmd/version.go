package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is version number which automatically set on build.
	Version string
	// Revision is git commit hash which automatically set `git rev-parse --short HEAD` on build.
	Revision string
)

func init() {
	RootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("myaws version: %s, revision: %s\n", Version, Revision)
		},
	}

	return cmd
}
