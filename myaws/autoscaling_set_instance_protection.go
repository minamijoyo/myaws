package myaws

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// AutoScalingSetInstanceProtectionOptions customize the behavior of the Attach command.
type AutoScalingSetInstanceProtectionOptions struct {
	AsgName              string
	InstanceIds          []*string
	ProtectedFromScaleIn bool
}

// AutoScalingSetInstanceProtection protects from termination when scale in your autoscaling group.
func (client *Client) AutoScalingSetInstanceProtection(options AutoScalingSetInstanceProtectionOptions) error {
	if len(options.InstanceIds) > 0 {
		if err := client.AutoScalingSetInstanceProtectionInstances(options.AsgName, options.InstanceIds, options.ProtectedFromScaleIn); err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) AutoScalingSetInstanceProtectionInstances(asgName string, instanceIds []*string, protectedFromScaleIn bool) error {
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
