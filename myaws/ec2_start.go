package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// EC2StartOptions customize the behavior of the Start command.
type EC2StartOptions struct {
	InstanceIds []*string
	Wait        bool
}

// EC2Start starts EC2 instances.
// If wait flag is true, wait until instance is in running state.
func (client *Client) EC2Start(options EC2StartOptions) error {
	params := &ec2.StartInstancesInput{
		InstanceIds: options.InstanceIds,
	}

	response, err := client.EC2.StartInstances(params)
	if err != nil {
		return errors.Wrap(err, "StartInstances failed:")
	}

	fmt.Fprintln(client.stdout, response)

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until instance running...")
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
