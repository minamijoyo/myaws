package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
)

// ec2lsCmd represents the ec2ls command
var ec2lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List EC2 instances",
	Long:  `List EC2 instances`,
	Run:   ec2.Ls,
}

func init() {
	ec2Cmd.AddCommand(ec2lsCmd)
}
