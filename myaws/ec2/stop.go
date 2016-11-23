package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// StopOptions customize the behavior of the Stop command.
type StopOptions struct {
	InstanceIds []*string
	Wait        bool
}

// Stop stops EC2 instances.
// If wait flag is true, wait until instance is in stopped state.
func Stop(client *myaws.Client, options StopOptions) error {
	params := &ec2.StopInstancesInput{
		InstanceIds: options.InstanceIds,
	}

	response, err := client.EC2.StopInstances(params)
	if err != nil {
		return errors.Wrap(err, "StopInstances failed:")
	}

	fmt.Println(response)

	if options.Wait {
		fmt.Println("Wait until instance stopped...")
		err := client.EC2.WaitUntilInstanceStopped(
			&ec2.DescribeInstancesInput{
				InstanceIds: options.InstanceIds,
			},
		)
		if err != nil {
			return errors.Wrap(err, "WaitUntilInstanceStopped failed:")
		}
	}

	return nil
}
