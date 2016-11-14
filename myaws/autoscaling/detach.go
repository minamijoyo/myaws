package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

// Detach detaches instances or load balancers from autoscaling group.
func Detach(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "AUTO_SCALING_GROUP_NAME is required.")
	}
	asgName := args[0]

	if viper.GetString("autoscaling.detach.instance-ids") != "" {
		detachInstances(asgName)
	}

	if viper.GetString("autoscaling.detach.load-balancer-names") != "" {
		detachLoadBalancers(asgName)
	}
}

func detachInstances(asgName string) {
	client := newAutoScalingClient()

	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.instance-ids"))
	decrementCapacity := true
	params := &autoscaling.DetachInstancesInput{
		AutoScalingGroupName:           &asgName,
		InstanceIds:                    instanceIds,
		ShouldDecrementDesiredCapacity: &decrementCapacity,
	}

	_, err := client.DetachInstances(params)
	if err != nil {
		panic(err)
	}
}

func detachLoadBalancers(asgName string) {
	client := newAutoScalingClient()

	loadBalancerNames := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.load-balancer-names"))
	params := &autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	_, err := client.DetachLoadBalancers(params)
	if err != nil {
		panic(err)
	}
}
