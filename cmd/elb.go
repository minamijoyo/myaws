package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/minamijoyo/myaws/myaws"
)

func init() {
	RootCmd.AddCommand(newELBCmd())
}

func newELBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elb",
		Short: "Manage ELB resources",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help() // nolint: errcheck
		},
	}

	cmd.AddCommand(
		newELBLsCmd(),
		newELBPsCmd(),
	)

	return cmd
}

func newELBLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List ELB instances",
		RunE:  runELBLsCmd,
	}

	return cmd
}

func runELBLsCmd(_ *cobra.Command, _ []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	return client.ELBLs()
}

func newELBPsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps ELB_NAME",
		Short: "Show ELB instances",
		RunE:  runELBPsCmd,
	}

	return cmd
}

func runELBPsCmd(_ *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("ELB_NAME is required")
	}

	options := myaws.ELBPsOptions{
		LoadBalancerName: args[0],
	}

	return client.ELBPs(options)
}
