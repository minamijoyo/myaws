package cmd

import (
	"github.com/spf13/cobra"

	"github.com/minamijoyo/myaws/myaws/elb"
)

func newELBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elb",
		Short: "Manage ELB resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
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
		RunE:  elb.Ls,
	}

	return cmd
}

func newELBPsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps ELB_NAME",
		Short: "Show ELB instances",
		RunE:  elb.Ps,
	}

	return cmd
}
