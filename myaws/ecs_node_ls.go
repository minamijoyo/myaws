package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSNodeLsOptions customize the behavior of the Ls command.
type ECSNodeLsOptions struct {
	Cluster     string
	PrintHeader bool
}

// ECSNodeLs describes ECS container instances.
func (client *Client) ECSNodeLs(options ECSNodeLsOptions) error {
	instances, err := client.findECSNodes(options.Cluster)
	if err != nil {
		return err
	}

	if options.PrintHeader {
		header := fmt.Sprintf("%-32s\t%-19s\t%-10s\t%s\t%s\t%s",
			"ContainerInstanceId",
			"Ec2InstanceId",
			"Status",
			"Running",
			"Pending",
			"RegisteredAt",
		)
		fmt.Fprintln(client.stdout, header)
	}

	for _, instance := range instances {
		fmt.Fprintln(client.stdout, formatECSNode(client, instance))
	}

	return nil
}

func formatECSNode(client *Client, instance *ecs.ContainerInstance) string {
	arn := strings.Split(*instance.ContainerInstanceArn, "/")
	// To fix misalignment, we use the width of state is 10 characters here,
	// because 8 characters + 2 characters as future margin of change.
	// The valid values of status are ACTIVE, INACTIVE, or DRAINING.
	return fmt.Sprintf("%s\t%s\t%-10s\t%d\t%d\t%s",
		arn[2],
		*instance.Ec2InstanceId,
		*instance.Status,
		*instance.RunningTasksCount,
		*instance.PendingTasksCount,
		client.FormatTime(instance.RegisteredAt),
	)
}
