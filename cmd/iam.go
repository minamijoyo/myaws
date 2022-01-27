package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/minamijoyo/myaws/myaws"
)

func init() {
	RootCmd.AddCommand(newIAMCmd())
}

func newIAMCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iam",
		Short: "Manage IAM resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // nolint: errcheck
		},
	}

	cmd.AddCommand(
		newIAMUserCmd(),
	)

	return cmd
}

func newIAMUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage IAM user resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // nolint: errcheck
		},
	}

	cmd.AddCommand(
		newIAMUserLsCmd(),
		newIAMUserResetPasswordCmd(),
	)

	return cmd
}

func newIAMUserLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List IAM users",
		RunE:  runIAMUserLsCmd,
	}

	return cmd
}

func runIAMUserLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	return client.IAMUserLs()
}

func newIAMUserResetPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-password USERNAME",
		Short: "Reset login password for IAM user",
		RunE:  runIAMUserResetPasswordCmd,
	}

	return cmd
}

func runIAMUserResetPasswordCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("USERNAME is required")
	}

	options := myaws.IAMUserResetPasswordOptions{
		UserName: args[0],
	}
	return client.IAMUserResetPassword(options)
}
