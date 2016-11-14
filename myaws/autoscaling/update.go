package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Update updates autoscaling group setting.
// Available param is currently desired-capacity only.
func Update(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("AUTO_SCALING_GROUP_NAME is required")
	}
	asgName := aws.String(args[0])

	desiredCapacity := viper.GetInt64("autoscaling.update.desired-capacity")
	if desiredCapacity == -1 {
		return errors.New("--desired-capacity is required")
	}

	client := newAutoScalingClient()

	params := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: asgName,
		DesiredCapacity:      &desiredCapacity,
	}

	if _, err := client.SetDesiredCapacity(params); err != nil {
		return errors.Wrap(err, "SetDesiredCapacity failed:")
	}

	return nil
}
