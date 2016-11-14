package autoscaling

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Ls describes autoscaling groups.
func Ls(*cobra.Command, []string) error {
	client := newAutoScalingClient()
	params := &autoscaling.DescribeAutoScalingGroupsInput{}

	response, err := client.DescribeAutoScalingGroups(params)
	if err != nil {
		return errors.Wrap(err, "DescribeAutoScalingGroups failed:")
	}

	for _, asg := range response.AutoScalingGroups {
		if viper.GetBool("autoscaling.ls.all") || len(asg.Instances) > 0 {
			fmt.Println(formatAutoScalingGroup(asg))
		}
	}

	return nil
}

func formatAutoScalingGroup(asg *autoscaling.Group) string {
	output := []string{
		formatInstacesLen(asg.Instances),
		*asg.AutoScalingGroupName,
		formatInstanceIds(asg.Instances),
		formatLoadBalancerNames(asg.LoadBalancerNames),
	}

	return strings.Join(output[:], "\t")
}

func formatInstacesLen(instances []*autoscaling.Instance) string {
	if instances == nil {
		return "0"
	}
	return strconv.Itoa(len(instances))
}

func formatInstanceIds(instances []*autoscaling.Instance) string {
	if instances == nil {
		return ""
	}
	instanceIds := lookupInstanceIds(instances)
	return strings.Join(instanceIds[:], " ")
}

func lookupInstanceIds(instances []*autoscaling.Instance) []string {
	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, *instance.InstanceId)
	}
	return instanceIds
}

func formatLoadBalancerNames(lbNames []*string) string {
	if lbNames == nil {
		return ""
	}
	return strings.Join(aws.StringValueSlice(lbNames)[:], " ")
}
