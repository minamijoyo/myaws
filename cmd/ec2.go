package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Manage EC2 resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ec2 called")
	},
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
}
