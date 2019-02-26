package myaws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// ECSNodeUpdateOptions customize the behavior of the Update command.
type ECSNodeUpdateOptions struct {
	Cluster            string
	ContainerInstances []*string
	Status             string
}

// ECSNodeUpdate Update ECS container instances.
func (client *Client) ECSNodeUpdate(options ECSNodeUpdateOptions) error {
	_, err := client.ECS.UpdateContainerInstancesState(
		&ecs.UpdateContainerInstancesStateInput{
			Cluster:            &options.Cluster,
			ContainerInstances: options.ContainerInstances,
			Status:             &options.Status,
		},
	)

	if err != nil {
		return errors.Wrapf(err, "UpdateContainerInstancesState failed")
	}

	return nil
}
