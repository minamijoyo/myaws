package cmd

import (
	"github.com/minamijoyo/myaws/myaws/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var autoscalingAttachCmd = &cobra.Command{
	Use:   "attach AUTO_SCALING_GROUP_NAME",
	Short: "Attach instances to autoscaling group",
	Run:   autoscaling.Attach,
}

func init() {
	autoscalingCmd.AddCommand(autoscalingAttachCmd)

	autoscalingAttachCmd.Flags().StringP("instance-ids", "i", "", "List of Instance IDs to attach")

	viper.BindPFlag("autoscaling.attach.instance-ids", autoscalingAttachCmd.Flags().Lookup("instance-ids"))
}
