package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Start(cmd *cobra.Command, args []string) {
	client := newEC2Client()

	instanceIds := aws.StringSlice(args)
	params := &ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	}

	response, err := client.StartInstances(params)
	if err != nil {
		panic(err)
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
			panic(err)
		}
	}
}
