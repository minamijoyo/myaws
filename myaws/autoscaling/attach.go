package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Attach(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "AUTO_SCALING_GROUP_NAME is required.")
	}
	asgName := args[0]

	if viper.GetString("autoscaling.attach.instance-ids") != "" {
		attachInstances(asgName)
	}

	fmt.Println(viper.GetString("autoscaling.attach.load-balancer-names"))
	if viper.GetString("autoscaling.attach.load-balancer-names") != "" {
		attachLoadBalancers(asgName)
	}
}

func attachInstances(asgName string) {
	client := newAutoScalingClient()

	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.instance-ids"))
	params := &autoscaling.AttachInstancesInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          instanceIds,
	}

	_, err := client.AttachInstances(params)
	if err != nil {
		panic(err)
	}
}

func attachLoadBalancers(asgName string) {
	client := newAutoScalingClient()

	loadBalancerNames := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.load-balancer-names"))
	params := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	_, err := client.AttachLoadBalancers(params)
	if err != nil {
		panic(err)
	}
}
