package cmd

import (
	"github.com/minamijoyo/myaws/myaws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newELBV2Cmd())
}

func newELBV2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elbv2",
		Short: "Manage ELBV2 resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // nolint: errcheck
		},
	}

	cmd.AddCommand(
		newELBV2LsCmd(),
		newELBV2PsCmd(),
	)

	return cmd
}

func newELBV2LsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List ELBV2 instances",
		RunE:  runELBV2LsCmd,
	}

	return cmd
}

func runELBV2LsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	return client.ELBV2Ls()
}

func newELBV2PsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps TARGET_GROUP_NAME",
		Short: "Show ELBV2 target group health",
		RunE:  runELBV2PsCmd,
	}

	return cmd
}

func runELBV2PsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("TARGET_GROUP_NAME is required")
	}

	options := myaws.ELBV2PsOptions{
		TargetGroupName: args[0],
	}

	return client.ELBV2Ps(options)
}
