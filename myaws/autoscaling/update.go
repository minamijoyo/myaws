package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Update(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "AUTO_SCALING_GROUP_NAME is required.")
	}
	asgName := aws.String(args[0])

	desiredCapacity := viper.GetInt64("autoscaling.update.desired-capacity")
	if desiredCapacity == -1 {
		myaws.UsageError(cmd, "--desired-capacity is required.")
	}

	client := newAutoScalingClient()

	params := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: asgName,
		DesiredCapacity:      &desiredCapacity,
	}

	_, err := client.SetDesiredCapacity(params)
	if err != nil {
		panic(err)
	}
}
