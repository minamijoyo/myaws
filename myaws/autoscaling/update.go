package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Update(cmd *cobra.Command, args []string) {
	client := newAutoScalingClient()

	asgName := aws.String(args[0])
	desiredCapacity := viper.GetInt64("autoscaling.update.desired-capacity")
	if desiredCapacity == -1 {
		panic("desired-capacity must be set.")
	}

	params := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: asgName,
		DesiredCapacity:      &desiredCapacity,
	}

	_, err := client.SetDesiredCapacity(params)
	if err != nil {
		panic(err)
	}
}
