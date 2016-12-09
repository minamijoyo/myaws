package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// EC2StopOptions customize the behavior of the Stop command.
type EC2StopOptions struct {
	InstanceIds []*string
	Wait        bool
}

// EC2Stop stops EC2 instances.
// If wait flag is true, wait until instance is in stopped state.
func (client *Client) EC2Stop(options EC2StopOptions) error {
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
