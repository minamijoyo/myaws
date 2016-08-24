package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
)

var ec2LsFlag ec2.LsFlag

// ec2lsCmd represents the ec2ls command
var ec2lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List EC2 instances",
	Run: func(*cobra.Command, []string) {
		ec2.Ls(&ec2LsFlag)
	},
}

func init() {
	ec2Cmd.AddCommand(ec2lsCmd)
	ec2lsCmd.Flags().BoolVarP(&ec2LsFlag.All, "all", "a", false, "List all instances")
}
