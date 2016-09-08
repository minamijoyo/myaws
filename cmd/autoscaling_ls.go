package cmd

import (
	"github.com/minamijoyo/myaws/myaws/autoscaling"
	"github.com/spf13/cobra"
)

var autoscalingLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List autoscaling groups",
	Run:   autoscaling.Ls,
}

func init() {
	autoscalingCmd.AddCommand(autoscalingLsCmd)
}
