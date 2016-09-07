package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
)

var ec2StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start EC2 instances",
	Run:   ec2.Start,
}

func init() {
	ec2Cmd.AddCommand(ec2StartCmd)
}
