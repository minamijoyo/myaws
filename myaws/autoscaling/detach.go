package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Detach(cmd *cobra.Command, args []string) {
	client := newAutoScalingClient()

	asgName := aws.String(args[0])
	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.instance-ids"))
	decrementCapacity := true
	params := &autoscaling.DetachInstancesInput{
		AutoScalingGroupName:           asgName,
		InstanceIds:                    instanceIds,
		ShouldDecrementDesiredCapacity: &decrementCapacity,
	}

	response, err := client.DetachInstances(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
