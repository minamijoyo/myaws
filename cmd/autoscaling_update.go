package cmd

import (
	"github.com/minamijoyo/myaws/myaws/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var autoscalingUpdateCmd = &cobra.Command{
	Use:   "update AUTO_SCALING_GROUP_NAME",
	Short: "Update autoscaling group",
	Run:   autoscaling.Update,
}

func init() {
	autoscalingCmd.AddCommand(autoscalingUpdateCmd)

	autoscalingUpdateCmd.Flags().Int64P("desired-capacity", "c", -1, "The number of EC2 instances that should be running in the Auto Scaling group.")

	viper.BindPFlag("autoscaling.update.desired-capacity", autoscalingUpdateCmd.Flags().Lookup("desired-capacity"))
}
