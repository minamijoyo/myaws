package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Detach detaches instances or load balancers from autoscaling group.
func Detach(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("AUTO_SCALING_GROUP_NAME is required")
	}
	asgName := args[0]

	if viper.GetString("autoscaling.detach.instance-ids") != "" {
		if err := detachInstances(asgName); err != nil {
			return err
		}
	}

	if viper.GetString("autoscaling.detach.load-balancer-names") != "" {
		if err := detachLoadBalancers(asgName); err != nil {
			return err
		}
	}

	return nil
}

func detachInstances(asgName string) error {
	client := newAutoScalingClient()

	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.instance-ids"))
	decrementCapacity := true
	params := &autoscaling.DetachInstancesInput{
		AutoScalingGroupName:           &asgName,
		InstanceIds:                    instanceIds,
		ShouldDecrementDesiredCapacity: &decrementCapacity,
	}

	if _, err := client.DetachInstances(params); err != nil {
		return errors.Wrap(err, "DetachInstances failed:")
	}

	return nil
}

func detachLoadBalancers(asgName string) error {
	client := newAutoScalingClient()

	loadBalancerNames := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.load-balancer-names"))
	params := &autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	if _, err := client.DetachLoadBalancers(params); err != nil {
		return errors.Wrap(err, "DetachLoadBalancers failed:")
	}

	return nil
}
