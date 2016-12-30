package myaws

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// AutoscalingLsOptions customize the behavior of the Ls command.
type AutoscalingLsOptions struct {
	All bool
}

// AutoscalingLs describes autoscaling groups.
func (client *Client) AutoscalingLs(options AutoscalingLsOptions) error {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}

	response, err := client.AutoScaling.DescribeAutoScalingGroups(params)
	if err != nil {
		return errors.Wrap(err, "DescribeAutoScalingGroups failed:")
	}

	for _, asg := range response.AutoScalingGroups {
		if options.All || len(asg.Instances) > 0 {
			fmt.Fprintln(client.stdout, formatAutoscalingGroup(asg))
		}
	}

	return nil
}

func formatAutoscalingGroup(asg *autoscaling.Group) string {
	output := []string{
		formatAutoscalingInstacesLen(asg.Instances),
		*asg.AutoScalingGroupName,
		formatAutoscalingInstanceIds(asg.Instances),
		formatAutoscalingLoadBalancerNames(asg.LoadBalancerNames),
	}

	return strings.Join(output[:], "\t")
}

func formatAutoscalingInstacesLen(instances []*autoscaling.Instance) string {
	if instances == nil {
		return "0"
	}
	return strconv.Itoa(len(instances))
}

func formatAutoscalingInstanceIds(instances []*autoscaling.Instance) string {
	if instances == nil {
		return ""
	}
	instanceIds := lookupAutoscalingInstanceIds(instances)
	return strings.Join(instanceIds[:], " ")
}

func lookupAutoscalingInstanceIds(instances []*autoscaling.Instance) []string {
	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, *instance.InstanceId)
	}
	return instanceIds
}

func formatAutoscalingLoadBalancerNames(lbNames []*string) string {
	if lbNames == nil {
		return ""
	}
	return strings.Join(aws.StringValueSlice(lbNames)[:], " ")
}
