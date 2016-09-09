package cmd

import (
	"github.com/minamijoyo/myaws/myaws/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var autoscalingDetachCmd = &cobra.Command{
	Use:   "detach AUTO_SCALING_GROUP_NAME",
	Short: "Detach instances from autoscaling group",
	Run:   autoscaling.Detach,
}

func init() {
	autoscalingCmd.AddCommand(autoscalingDetachCmd)

	autoscalingDetachCmd.Flags().StringP("instance-ids", "i", "", "List of Instance IDs to detach")

	viper.BindPFlag("autoscaling.detach.instance-ids", autoscalingDetachCmd.Flags().Lookup("instance-ids"))
}
