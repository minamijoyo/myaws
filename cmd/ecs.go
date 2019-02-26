package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		newECSNodeUpdateCmd(),
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

func newECSNodeUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update CLUSTER",
		Short: "Update ECS nodes",
		RunE:  runECSNodeUpdateCmd,
	}

	flags := cmd.Flags()
	flags.StringP("container-instances", "i", "", "A list of container instance IDs or full ARN entries separated by space")
	flags.StringP("status", "s", "", "container instance state (ACTIVE | DRAINING)")

	viper.BindPFlag("ecs.node.update.container-instances", flags.Lookup("container-instances"))
	viper.BindPFlag("ecs.node.update.status", flags.Lookup("status"))

	return cmd
}

func runECSNodeUpdateCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	containerInstances := aws.StringSlice(viper.GetStringSlice("ecs.node.update.container-instances"))
	if len(containerInstances) == 0 {
		return errors.New("container-instances is required")
	}

	status := viper.GetString("ecs.node.update.status")
	if len(status) == 0 {
		return errors.New("status is required")
	}

	options := myaws.ECSNodeUpdateOptions{
		Cluster:            args[0],
		ContainerInstances: containerInstances,
		Status:             status,
	}

	return client.ECSNodeUpdate(options)
}
