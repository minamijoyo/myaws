package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Attach attaches instances or load balancers from autoscaling group.
func Attach(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("AUTO_SCALING_GROUP_NAME is required")
	}
	asgName := args[0]

	if viper.GetString("autoscaling.attach.instance-ids") != "" {
		if err := attachInstances(asgName); err != nil {
			return err
		}
	}

	if viper.GetString("autoscaling.attach.load-balancer-names") != "" {
		if err := attachLoadBalancers(asgName); err != nil {
			return err
		}
	}

	return nil
}

func attachInstances(asgName string) error {
	client := newAutoScalingClient()

	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.instance-ids"))
	params := &autoscaling.AttachInstancesInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          instanceIds,
	}

	if _, err := client.AttachInstances(params); err != nil {
		return errors.Wrap(err, "AttachInstances failed:")
	}

	return nil
}

func attachLoadBalancers(asgName string) error {
	client := newAutoScalingClient()

	loadBalancerNames := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.load-balancer-names"))
	params := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	if _, err := client.AttachLoadBalancers(params); err != nil {
		return errors.Wrap(err, "AttachLoadBalancers failed:")
	}

	return nil
}
