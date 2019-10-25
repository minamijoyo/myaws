package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
	funk "github.com/thoas/go-funk"
)

// ECSNodeDrainOptions customize the behavior of the Drain command.
type ECSNodeDrainOptions struct {
	Cluster            string
	ContainerInstances []*string
	Wait               bool
}

// ECSNodeDrain Drain ECS container instances.
// We want to wait until drain action is completed, but the ECSNodeUpdate
// method is general purpose, so we implement a wait option to specialized
// method for draining.
func (client *Client) ECSNodeDrain(options ECSNodeDrainOptions) error {
	// We can specify up to 10 container instances to update state in a single operation.
	// This constraint is not specified in the API reference, but it returns the following error:
	//   InvalidParameterException: instanceIds can have at most 10 items.
	// So we need to divide the list by 10.
	chunks := (funk.Chunk(options.ContainerInstances, 10)).([][]*string)
	for _, c := range chunks {
		_, err := client.ECS.UpdateContainerInstancesState(
			&ecs.UpdateContainerInstancesStateInput{
				Cluster:            &options.Cluster,
				ContainerInstances: c,
				Status:             aws.String("DRAINING"),
			},
		)
		if err != nil {
			return errors.Wrapf(err, "UpdateContainerInstancesState failed")
		}
	}

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until container instances are drained...")
		return client.WaitUntilECSContainerInstancesAreDrained(options.Cluster, options.ContainerInstances)
	}

	return nil
}
