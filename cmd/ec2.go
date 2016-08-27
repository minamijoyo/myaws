package cmd

import (
	"github.com/spf13/cobra"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Manage EC2 resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
}
