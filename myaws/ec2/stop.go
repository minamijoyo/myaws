package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

// Stop stops EC2 instances.
// If wait flag is true, wait until instance is in stopped state.
func Stop(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "INSTANCE_ID is required.")
	}
	instanceIds := aws.StringSlice(args)

	client := newEC2Client()

	params := &ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	}

	response, err := client.StopInstances(params)
	if err != nil {
		panic(err)
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
			panic(err)
		}
	}
}
