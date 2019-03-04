package myaws

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// getAutoScalingGroupDesiredCapacity is a helper function which returns
// DesiredCapacity of the specific AutoScalingGroup.
func (client *Client) getAutoScalingGroupDesiredCapacity(asgName string) (int64, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asgName},
	}

	response, err := client.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return 0, errors.Wrap(err, "getAutoScalingGroupDesiredCapacity failed:")
	}

	desiredCapacity := response.AutoScalingGroups[0].DesiredCapacity

	return *desiredCapacity, nil
}
