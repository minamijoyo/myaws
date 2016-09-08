package cmd

import (
	"github.com/spf13/cobra"
)

var autoscalingCmd = &cobra.Command{
	Use:   "autoscaling",
	Short: "Manage autoscaling resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(autoscalingCmd)
}
