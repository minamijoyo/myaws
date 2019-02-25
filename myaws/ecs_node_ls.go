package myaws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

// ECSNodeLsOptions customize the behavior of the Ls command.
type ECSNodeLsOptions struct {
	Cluster string
}

// ECSNodeLs describes EC2 instances.
func (client *Client) ECSNodeLs(options ECSNodeLsOptions) error {
	params := &ecs.ListContainerInstancesInput{
		Cluster: &options.Cluster,
	}

	instances, err := client.ECS.ListContainerInstances(params)
	if err != nil {
		return errors.Wrapf(err, "ListContainerInstances failed")
	}

	pp.Println(instances)

	return nil
}
