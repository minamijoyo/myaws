package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSNodeLsOptions customize the behavior of the Ls command.
type ECSNodeLsOptions struct {
	Cluster string
}

// ECSNodeLs describes ECS container instances.
func (client *Client) ECSNodeLs(options ECSNodeLsOptions) error {
	instances, err := client.findECSNodes(options.Cluster)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		fmt.Fprintln(client.stdout, formatECSNode(client, options, instance))
	}

	return nil
}

func formatECSNode(client *Client, options ECSNodeLsOptions, instance *ecs.ContainerInstance) string {
	arn := strings.Split(*instance.ContainerInstanceArn, "/")

	// To fix misalignment, we use the width of state is 10 characters here,
	// because 8 characters + 2 characters as future margin of change.
	// The valid values of status are ACTIVE, INACTIVE, or DRAINING.
	return fmt.Sprintf("%s\t%s\t%-10s\t%d\t%d\t%s",
		arn[1],
		*instance.Ec2InstanceId,
		*instance.Status,
		*instance.RunningTasksCount,
		*instance.PendingTasksCount,
		client.FormatTime(instance.RegisteredAt),
	)
}
