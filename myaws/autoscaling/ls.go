package autoscaling

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
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
	for _, asg := range response.AutoScalingGroups {
		fmt.Println(formatAutoScalingGroup(asg))
	}
}

func formatAutoScalingGroup(asg *autoscaling.Group) string {
	var instanceLength int
	if asg.Instances != nil {
		instanceLength = len(asg.Instances)
	}
	output := []string{
		strconv.Itoa(instanceLength),
		*asg.AutoScalingGroupName,
	}

	if asg.LoadBalancerNames != nil {
		loadBalancerNames := aws.StringValueSlice(asg.LoadBalancerNames)
		output = append(output, loadBalancerNames...)
	} else {
		output = append(output, "")
	}

	var instanceIds []string
	if instanceLength > 0 {
		instanceIds = lookupInstanceIds(asg.Instances)
		output = append(output, instanceIds...)
	}
	return strings.Join(output[:], "\t")
}

func lookupInstanceIds(instances []*autoscaling.Instance) []string {
	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, *instance.InstanceId)
	}
	return instanceIds
}
