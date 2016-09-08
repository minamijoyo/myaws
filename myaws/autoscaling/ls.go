package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
)

func Ls(*cobra.Command, []string) {
	client := newAutoScalingClient()
	params := &autoscaling.DescribeAutoScalingGroupsInput{}

	response, err := client.DescribeAutoScalingGroups(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
