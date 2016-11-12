package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version string = "v0.0.6"

// Revision is git commit hash which automatically set `git describe --always` on build
var Revision string

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("myaws version: %s, revision: %s\n", version, Revision)
		},
	}

	return cmd
}
