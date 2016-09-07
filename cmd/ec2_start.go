package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ec2StartCmd = &cobra.Command{
	Use:   "start INSTANCE_ID [...]",
	Short: "Start EC2 instances",
	Run:   ec2.Start,
}

func init() {
	ec2Cmd.AddCommand(ec2StartCmd)
	ec2StartCmd.Flags().BoolP("wait", "w", false, "Wait until instance running")

	viper.BindPFlag("ec2.start.wait", ec2StartCmd.Flags().Lookup("wait"))
}
