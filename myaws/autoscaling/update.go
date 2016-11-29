package autoscaling

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// UpdateOptions customize the behavior of the Update command.
type UpdateOptions struct {
	AsgName         string
	DesiredCapacity int64
}

// Update updates autoscaling group setting.
// Available param is currently desired-capacity only.
func Update(client *myaws.Client, options UpdateOptions) error {
	params := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: &options.AsgName,
		DesiredCapacity:      &options.DesiredCapacity,
	}

	if _, err := client.AutoScaling.SetDesiredCapacity(params); err != nil {
		return errors.Wrap(err, "SetDesiredCapacity failed:")
	}

	return nil
}
