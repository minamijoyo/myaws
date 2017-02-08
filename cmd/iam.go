package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newIAMResetPasswordCmd(),
	)

	return cmd
}

func newIAMResetPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-password",
		Short: "Reset login password for IAM user",
		RunE:  runIAMResetPasswordCmd,
	}

	flags := cmd.Flags()
	flags.StringP("username", "u", "", "Username whose password is to be changed")

	viper.BindPFlag("iam.reset-password.username", flags.Lookup("username"))
	return cmd
}

func runIAMResetPasswordCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	username := viper.GetString("iam.reset-password.username")
	if username == "" {
		return errors.New("username is required")
	}

	options := myaws.IAMResetPasswordOptions{
		UserName: username,
	}
	return client.IAMResetPassword(options)
}
