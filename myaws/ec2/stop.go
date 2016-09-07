package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

func Stop(cmd *cobra.Command, args []string) {
	client := newEC2Client()

	params := &ec2.StopInstancesInput{
		InstanceIds: aws.StringSlice(args),
	}

	response, err := client.StopInstances(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
