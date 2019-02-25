package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// ECSNodeLsOptions customize the behavior of the Ls command.
type ECSNodeLsOptions struct {
	Cluster string
}

// ECSNodeLs describes ECS container instances.
func (client *Client) ECSNodeLs(options ECSNodeLsOptions) error {
	arns, err := client.ECS.ListContainerInstances(
		&ecs.ListContainerInstancesInput{
			Cluster: &options.Cluster,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "ListContainerInstances failed")
	}

	if len(arns.ContainerInstanceArns) == 0 {
		return errors.New("container instances not found")
	}

	instances, err := client.ECS.DescribeContainerInstances(
		&ecs.DescribeContainerInstancesInput{
			Cluster:            &options.Cluster,
			ContainerInstances: arns.ContainerInstanceArns,
		},
	)

	if len(instances.ContainerInstances) == 0 {
		return errors.New("ListContainerInstances succeed, but DescribeContainerInstances returns no instances")
	}

	for _, instance := range instances.ContainerInstances {
		fmt.Fprintln(client.stdout, formatECSNode(client, options, instance))
	}

	return nil
}

func formatECSNode(client *Client, options ECSNodeLsOptions, instance *ecs.ContainerInstance) string {
	arn := strings.Split(*instance.ContainerInstanceArn, "/")

	return fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t%s",
		arn[1],
		*instance.Ec2InstanceId,
		*instance.Status,
		*instance.RunningTasksCount,
		*instance.PendingTasksCount,
		client.FormatTime(instance.RegisteredAt),
	)
}
