package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// StartOptions customize the behavior of the Start command.
type StartOptions struct {
	InstanceIds []*string
	Wait        bool
}

// Start starts EC2 instances.
// If wait flag is true, wait until instance is in running state.
func Start(client *myaws.Client, options StartOptions) error {
	params := &ec2.StartInstancesInput{
		InstanceIds: options.InstanceIds,
	}

	response, err := client.EC2.StartInstances(params)
	if err != nil {
		return errors.Wrap(err, "StartInstances failed:")
	}

	fmt.Println(response)

	if options.Wait {
		fmt.Println("Wait until instance running...")
		err := client.EC2.WaitUntilInstanceRunning(
			&ec2.DescribeInstancesInput{
				InstanceIds: options.InstanceIds,
			},
		)
		if err != nil {
			return errors.Wrap(err, "WaitUntilInstanceRunning failed:")
		}
	}

	return nil
}
