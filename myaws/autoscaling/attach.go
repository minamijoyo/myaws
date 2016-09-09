package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Attach(cmd *cobra.Command, args []string) {
	client := newAutoScalingClient()

	asgName := aws.String(args[0])
	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.instance-ids"))
	params := &autoscaling.AttachInstancesInput{
		AutoScalingGroupName: asgName,
		InstanceIds:          instanceIds,
	}

	response, err := client.AttachInstances(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
