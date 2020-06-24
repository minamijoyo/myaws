package cmd

import (
	"time"

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
		newECSStatusCmd(),
		newECSNodeCmd(),
		newECSServiceCmd(),
	)

	return cmd
}

func newECSStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status CLUSTER",
		Short: "Print ECS status",
		RunE:  runECSStatusCmd,
	}

	return cmd
}

func runECSStatusCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	options := myaws.ECSStatusOptions{
		Cluster: args[0],
	}
	return client.ECSStatus(options)
}

func newECSNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage ECS node resources (container instances)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newECSNodeLsCmd(),
		newECSNodeUpdateCmd(),
		newECSNodeDrainCmd(),
		newECSNodeRenewCmd(),
	)

	return cmd
}

func newECSNodeLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls CLUSTER",
		Short: "List ECS nodes (container instances)",
		RunE:  runECSNodeLsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("print-header", "H", false, "Print Header")

	viper.BindPFlag("ecs.node.ls.print-header", flags.Lookup("print-header"))

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
		Cluster:     args[0],
		PrintHeader: viper.GetBool("ecs.node.ls.print-header"),
	}
	return client.ECSNodeLs(options)
}

func newECSNodeUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update CLUSTER",
		Short: "Update ECS nodes (container instances)",
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

func newECSNodeDrainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drain CLUSTER",
		Short: "Drain ECS nodes (container instances)",
		RunE:  runECSNodeDrainCmd,
	}

	flags := cmd.Flags()
	flags.StringP("container-instances", "i", "", "A list of container instance IDs or full ARN entries separated by space")
	flags.BoolP("wait", "w", false, "Wait until container instances are drained")
	flags.Int64P("timeout", "t", 600, "Number of secconds to wait before timeout")

	viper.BindPFlag("ecs.node.drain.container-instances", flags.Lookup("container-instances"))
	viper.BindPFlag("ecs.node.drain.wait", flags.Lookup("wait"))
	viper.BindPFlag("ecs.node.drain.timeout", flags.Lookup("timeout"))

	return cmd
}

func runECSNodeDrainCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	containerInstances := aws.StringSlice(viper.GetStringSlice("ecs.node.drain.container-instances"))
	if len(containerInstances) == 0 {
		return errors.New("container-instances is required")
	}

	timeout := time.Duration(viper.GetInt64("ecs.node.drain.timeout")) * time.Second

	options := myaws.ECSNodeDrainOptions{
		Cluster:            args[0],
		ContainerInstances: containerInstances,
		Wait:               viper.GetBool("ecs.node.drain.wait"),
		Timeout:            timeout,
	}

	return client.ECSNodeDrain(options)
}

func newECSNodeRenewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew CLUSTER",
		Short: "Renew ECS nodes (container instances) with blue-grean deployment",
		RunE:  runECSNodeRenewCmd,
	}

	flags := cmd.Flags()
	flags.StringP("asg-name", "a", "", "A name of AutoScalingGroup to which the ECS container instances belong")

	// Note that this is a total timeout, and indivisual wait operations can
	// timeout in shorter amount of time.
	flags.Int64P("timeout", "t", 3600, "Number of secconds to wait before timeout")

	viper.BindPFlag("ecs.node.renew.asg-name", flags.Lookup("asg-name"))
	viper.BindPFlag("ecs.node.renew.timeout", flags.Lookup("timeout"))

	return cmd
}

func runECSNodeRenewCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	asgName := viper.GetString("ecs.node.renew.asg-name")
	if len(asgName) == 0 {
		return errors.New("asg-name is required")
	}

	timeout := time.Duration(viper.GetInt64("ecs.node.renew.timeout")) * time.Second

	options := myaws.ECSNodeRenewOptions{
		Cluster: args[0],
		AsgName: asgName,
		Timeout: timeout,
	}

	return client.ECSNodeRenew(options)
}

func newECSServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage ECS service resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newECSServiceLsCmd(),
		newECSServiceUpdateCmd(),
	)

	return cmd
}

func newECSServiceLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls CLUSTER",
		Short: "List ECS services",
		RunE:  runECSServiceLsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("print-header", "H", false, "Print Header")

	viper.BindPFlag("ecs.service.ls.print-header", flags.Lookup("print-header"))

	return cmd
}

func runECSServiceLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	options := myaws.ECSServiceLsOptions{
		Cluster:     args[0],
		PrintHeader: viper.GetBool("ecs.service.ls.print-header"),
	}
	return client.ECSServiceLs(options)
}

func newECSServiceUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update CLUSTER",
		Short: "Update ECS services",
		RunE:  runECSServiceUpdateCmd,
	}

	flags := cmd.Flags()
	flags.StringP("service", "s", "", "Name of service to be updated")
	flags.Int64P("desired-capacity", "c", -1, "Number of task to place and keep running")
	flags.BoolP("wait", "w", false, "Wait until desired capacity tasks are InService")

	// We may use time.Duration directly here via flags.Duration,
	// but time.Duration is unfamiliar for non-Gopher
	// so we use simple int64 as seconds for CLI interface.
	flags.Int64P("timeout", "t", 600, "Number of secconds to wait before timeout")

	flags.BoolP("force", "f", false, "Force new deployment")

	viper.BindPFlag("ecs.service.update.service", flags.Lookup("service"))
	viper.BindPFlag("ecs.service.update.desired-capacity", flags.Lookup("desired-capacity"))
	viper.BindPFlag("ecs.service.update.wait", flags.Lookup("wait"))
	viper.BindPFlag("ecs.service.update.timeout", flags.Lookup("timeout"))
	viper.BindPFlag("ecs.service.update.force", flags.Lookup("force"))
	return cmd
}

func runECSServiceUpdateCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("CLUSTER is required")
	}

	service := viper.GetString("ecs.service.update.service")
	if len(service) == 0 {
		return errors.New("--service is required")
	}

	// For desiredCapacity, 0 is valid value.
	// So we use -1 as a default value which indicates unset.
	// In this case, ECSServiceUpdateOptions.DesiredCount should be nil to allow us force deploy
	desiredCapacity := viper.GetInt64("ecs.service.update.desired-capacity")
	var desiredCapacityP *int64
	if desiredCapacity != -1 {
		desiredCapacityP = &desiredCapacity
	}

	timeout := time.Duration(viper.GetInt64("ecs.service.update.timeout")) * time.Second

	options := myaws.ECSServiceUpdateOptions{
		Cluster:      args[0],
		Service:      service,
		DesiredCount: desiredCapacityP,
		Wait:         viper.GetBool("ecs.service.update.wait"),
		Timeout:      timeout,
		Force:        viper.GetBool("ecs.service.update.force"),
	}
	return client.ECSServiceUpdate(options)
}
