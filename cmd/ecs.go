package cmd

import (
	"github.com/minamijoyo/myaws/myaws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newECSCmd())
}

func newECSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecs",
		Short: "Manage ECS resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newECSNodeCmd(),
	)

	return cmd
}

func newECSNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage ECS node resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newECSNodeLsCmd(),
	)

	return cmd
}

func newECSNodeLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls CLUSTER",
		Short: "List ECS nodes",
		RunE:  runECSNodeLsCmd,
	}

	return cmd
}

func runECSNodeLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	options := myaws.ECSNodeLsOptions{
		Cluster: args[0],
	}
	return client.ECSNodeLs(options)
}
