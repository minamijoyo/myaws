package cmd

import (
	"github.com/spf13/cobra"
)

var elbCmd = &cobra.Command{
	Use:   "elb",
	Short: "Manage ELB resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(elbCmd)
}
