package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ec2lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List EC2 instances",
	Run:   ec2.Ls,
}

func init() {
	ec2Cmd.AddCommand(ec2lsCmd)

	ec2lsCmd.Flags().BoolP("all", "a", false, "List all instances (default: false)")

	viper.BindPFlag("ec2.ls.all", ec2lsCmd.Flags().Lookup("all"))
}
