package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ec2StopCmd = &cobra.Command{
	Use:   "stop INSTANCE_ID [...]",
	Short: "Stop EC2 instances",
	Run:   ec2.Stop,
}

func init() {
	ec2Cmd.AddCommand(ec2StopCmd)
	ec2StopCmd.Flags().BoolP("wait", "w", false, "Wait until instance stopped")

	viper.BindPFlag("ec2.stop.wait", ec2StopCmd.Flags().Lookup("wait"))
}
