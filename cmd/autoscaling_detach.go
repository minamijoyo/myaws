package cmd

import (
	"github.com/minamijoyo/myaws/myaws/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var autoscalingDetachCmd = &cobra.Command{
	Use:   "detach AUTO_SCALING_GROUP_NAME",
	Short: "Detach instances/loadbalancers from autoscaling group",
	Run:   autoscaling.Detach,
}

func init() {
	autoscalingCmd.AddCommand(autoscalingDetachCmd)

	autoscalingDetachCmd.Flags().StringP("instance-ids", "i", "", "One or more instance IDs")
	autoscalingDetachCmd.Flags().StringP("load-balancer-names", "l", "", "One or more load balancer names")

	viper.BindPFlag("autoscaling.detach.instance-ids", autoscalingDetachCmd.Flags().Lookup("instance-ids"))
	viper.BindPFlag("autoscaling.detach.load-balancer-names", autoscalingDetachCmd.Flags().Lookup("load-balancer-names"))
}
