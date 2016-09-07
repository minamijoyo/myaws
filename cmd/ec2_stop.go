package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
)

var ec2StopCmd = &cobra.Command{
	Use:   "stop INSTANCE_ID [...]",
	Short: "Stop EC2 instances",
	Run:   ec2.Stop,
}

func init() {
	ec2Cmd.AddCommand(ec2StopCmd)
}
