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

	autoscalingAttachCmd.Flags().StringP("instance-ids", "i", "", "One or more instance IDs")
	autoscalingAttachCmd.Flags().StringP("load-balancer-names", "l", "", "One or more load balancer names")

	viper.BindPFlag("autoscaling.attach.instance-ids", autoscalingAttachCmd.Flags().Lookup("instance-ids"))
	viper.BindPFlag("autoscaling.attach.load-balancer-names", autoscalingAttachCmd.Flags().Lookup("load-balancer-names"))
}
