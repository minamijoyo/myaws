package cmd

import (
	"github.com/spf13/cobra"
)

var ecrCmd = &cobra.Command{
	Use:   "ecr",
	Short: "Manage ECR resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(ecrCmd)
}
