package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Stop stops EC2 instances.
// If wait flag is true, wait until instance is in stopped state.
func Stop(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	client := newEC2Client()

	params := &ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	}

	response, err := client.StopInstances(params)
	if err != nil {
		return errors.Wrap(err, "StopInstances failed:")
	}

	fmt.Println(response)
	if viper.GetBool("ec2.stop.wait") {
		fmt.Println("Wait until instance stopped...")
		err := client.WaitUntilInstanceStopped(
			&ec2.DescribeInstancesInput{
				InstanceIds: instanceIds,
			},
		)
		if err != nil {
			return errors.Wrap(err, "WaitUntilInstanceStopped failed:")
		}
	}

	return nil
}
