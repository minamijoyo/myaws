package myaws

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	funk "github.com/thoas/go-funk"
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
	chunks := (funk.Chunk(options.InstanceIds, maxInstanceIDCount)).([][]*string)
	for _, c := range chunks {
		if err := client.autoScalingSetInstanceProtectionInstances(options.AsgName, c, options.ProtectedFromScaleIn); err != nil {
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
