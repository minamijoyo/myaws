package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Start starts EC2 instances.
// If wait flag is true, wait until instance is in running state.
func Start(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	client := newEC2Client()

	params := &ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	}

	response, err := client.StartInstances(params)
	if err != nil {
		return errors.Wrap(err, "StartInstances failed:")
	}

	fmt.Println(response)

	if viper.GetBool("ec2.start.wait") {
		fmt.Println("Wait until instance running...")
		err := client.WaitUntilInstanceRunning(
			&ec2.DescribeInstancesInput{
				InstanceIds: instanceIds,
			},
		)
		if err != nil {
			return errors.Wrap(err, "WaitUntilInstanceRunning failed:")
		}
	}

	return nil
}
