package myaws

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"math"
)

// AutoScalingSetInstanceProtectionOptions customizes the behavior of the Attach command.
type AutoScalingSetInstanceProtectionOptions struct {
	AsgName              string
	InstanceIds          []*string
	ProtectedFromScaleIn bool
}

// AutoScalingSetInstanceProtection protects from termination when scale in your autoscaling group.
func (client *Client) AutoScalingSetInstanceProtection(options AutoScalingSetInstanceProtectionOptions) error {
	// the number of maximum InstanceIds is limited to 19.
	// https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_SetInstanceProtection.html
	maxInstanceIDCount := 19
	var maxIteration = int(math.Ceil(float64(len(options.InstanceIds)) / float64(maxInstanceIDCount)))
	for i := 0; i < maxIteration; i++ { // set AutoScalingInstanceProtection every 19 instances.
		firstIndex := maxInstanceIDCount * i
		lastIndex := maxInstanceIDCount * (i + 1)
		if i == maxIteration-1 { // last iteration
			lastIndex = len(options.InstanceIds)
		}
		instanceIds := options.InstanceIds[firstIndex:lastIndex]
		if err := client.autoScalingSetInstanceProtectionInstances(options.AsgName, instanceIds, options.ProtectedFromScaleIn); err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) autoScalingSetInstanceProtectionInstances(asgName string, instanceIds []*string, protectedFromScaleIn bool) error {
	params := &autoscaling.SetInstanceProtectionInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          instanceIds,
		ProtectedFromScaleIn: &protectedFromScaleIn,
	}

	if _, err := client.AutoScaling.SetInstanceProtection(params); err != nil {
		return errors.Wrap(err, "SetInstanceProtection failed:")
	}

	return nil
}
