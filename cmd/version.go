package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version string = "v0.0.3"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("myaws version: %s\n", version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
