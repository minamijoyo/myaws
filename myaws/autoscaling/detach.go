package autoscaling

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// DetachOptions customize the behavior of the Detach command.
type DetachOptions struct {
	AsgName           string
	InstanceIds       []*string
	LoadBalancerNames []*string
}

// Detach detaches instances or load balancers from autoscaling group.
func Detach(client *myaws.Client, options DetachOptions) error {
	if len(options.InstanceIds) > 0 {
		if err := detachInstances(client, options.AsgName, options.InstanceIds); err != nil {
			return err
		}
	}

	if len(options.LoadBalancerNames) > 0 {
		if err := detachLoadBalancers(client, options.AsgName, options.LoadBalancerNames); err != nil {
			return err
		}
	}

	return nil
}

func detachInstances(client *myaws.Client, asgName string, instanceIds []*string) error {
	decrementCapacity := true
	params := &autoscaling.DetachInstancesInput{
		AutoScalingGroupName:           &asgName,
		InstanceIds:                    instanceIds,
		ShouldDecrementDesiredCapacity: &decrementCapacity,
	}

	if _, err := client.AutoScaling.DetachInstances(params); err != nil {
		return errors.Wrap(err, "DetachInstances failed:")
	}

	return nil
}

func detachLoadBalancers(client *myaws.Client, asgName string, loadBalancerNames []*string) error {
	params := &autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	if _, err := client.AutoScaling.DetachLoadBalancers(params); err != nil {
		return errors.Wrap(err, "DetachLoadBalancers failed:")
	}

	return nil
}
