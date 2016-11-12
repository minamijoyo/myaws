package cmd

import (
	"github.com/spf13/cobra"

	"github.com/minamijoyo/myaws/myaws/ecr"
)

func newECRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecr",
		Short: "Manage ECR resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newECRGetLoginCmd(),
	)

	return cmd
}

func newECRGetLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-login",
		Short: "Get docker login command for ECR",
		Run:   ecr.GetLogin,
	}

	return cmd
}
