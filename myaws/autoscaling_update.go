package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// AutoscalingUpdateOptions customize the behavior of the Update command.
type AutoscalingUpdateOptions struct {
	AsgName         string
	DesiredCapacity int64
	Wait            bool
}

// AutoscalingUpdate updates autoscaling group setting.
// Available param is currently desired-capacity only.
func (client *Client) AutoscalingUpdate(options AutoscalingUpdateOptions) error {
	params := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: &options.AsgName,
		DesiredCapacity:      &options.DesiredCapacity,
	}

	if _, err := client.AutoScaling.SetDesiredCapacity(params); err != nil {
		return errors.Wrap(err, "SetDesiredCapacity failed:")
	}

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until the desired capacity instances are InService...")
		return client.waitUntilAutoScalingGroupDesiredState(options.AsgName)
	}

	return nil
}
