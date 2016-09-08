package cmd

import (
	"github.com/minamijoyo/myaws/myaws/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var autoscalingLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List autoscaling groups",
	Run:   autoscaling.Ls,
}

func init() {
	autoscalingCmd.AddCommand(autoscalingLsCmd)

	autoscalingLsCmd.Flags().BoolP("all", "a", false, "List all autoscaling groups (by default, list autoscaling groups only having at least 1 attached instance)")

	viper.BindPFlag("autoscaling.ls.all", autoscalingLsCmd.Flags().Lookup("all"))
}
