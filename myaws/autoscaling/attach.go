package autoscaling

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// AttachOptions customize the behavior of the Attach command.
type AttachOptions struct {
	AsgName           string
	InstanceIds       []*string
	LoadBalancerNames []*string
}

// Attach attaches instances or load balancers from autoscaling group.
func Attach(client *myaws.Client, options AttachOptions) error {
	if len(options.InstanceIds) > 0 {
		if err := attachInstances(client, options.AsgName, options.InstanceIds); err != nil {
			return err
		}
	}

	if len(options.LoadBalancerNames) > 0 {
		if err := attachLoadBalancers(client, options.AsgName, options.LoadBalancerNames); err != nil {
			return err
		}
	}

	return nil
}

func attachInstances(client *myaws.Client, asgName string, instanceIds []*string) error {
	params := &autoscaling.AttachInstancesInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          instanceIds,
	}

	if _, err := client.AutoScaling.AttachInstances(params); err != nil {
		return errors.Wrap(err, "AttachInstances failed:")
	}

	return nil
}

func attachLoadBalancers(client *myaws.Client, asgName string, loadBalancerNames []*string) error {
	params := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	if _, err := client.AutoScaling.AttachLoadBalancers(params); err != nil {
		return errors.Wrap(err, "AttachLoadBalancers failed:")
	}

	return nil
}
