package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newSTSCmd())
}

func newSTSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sts",
		Short: "Manage STS resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // nolint: errcheck
		},
	}

	cmd.AddCommand(
		newSTSIDCmd(),
	)

	return cmd
}

func newSTSIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "id",
		Short: "Get caller identity",
		RunE:  runSTSIDCmd,
	}

	return cmd
}

func runSTSIDCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	return client.STSID()
}
