package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/minamijoyo/myaws/myaws/ecr"
)

func init() {
	RootCmd.AddCommand(newECRCmd())
}

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
		RunE:  runECRGetLoginCmd,
	}

	return cmd
}

func runECRGetLoginCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	return ecr.GetLogin(client)
}
