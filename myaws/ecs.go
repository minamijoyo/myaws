package myaws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// findECSNodes finds ECS container instances
func (client *Client) findECSNodes(cluster string) ([]*ecs.ContainerInstance, error) {
	arns, err := client.ECS.ListContainerInstances(
		&ecs.ListContainerInstancesInput{
			Cluster: &cluster,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "ListContainerInstances failed")
	}

	if len(arns.ContainerInstanceArns) == 0 {
		return nil, errors.New("container instances not found")
	}

	instances, err := client.ECS.DescribeContainerInstances(
		&ecs.DescribeContainerInstancesInput{
			Cluster:            &cluster,
			ContainerInstances: arns.ContainerInstanceArns,
		},
	)

	if len(instances.ContainerInstances) == 0 {
		return nil, errors.New("ListContainerInstances succeed, but DescribeContainerInstances returns no instances")
	}

	return instances.ContainerInstances, nil
}
