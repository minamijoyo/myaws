package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

func Start(cmd *cobra.Command, args []string) {
	client := newEC2Client()

	fmt.Println(args)
	params := &ec2.StartInstancesInput{
		InstanceIds: aws.StringSlice(args),
	}

	response, err := client.StartInstances(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
