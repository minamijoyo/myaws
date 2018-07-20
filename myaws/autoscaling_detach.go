package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// AutoscalingDetachOptions customize the behavior of the Detach command.
type AutoscalingDetachOptions struct {
	AsgName           string
	InstanceIds       []*string
	LoadBalancerNames []*string
	Wait              bool
}

// AutoscalingDetach detaches instances or load balancers from autoscaling group.
func (client *Client) AutoscalingDetach(options AutoscalingDetachOptions) error {
	if len(options.InstanceIds) > 0 {
		if err := client.autoscalingDetachInstances(options.AsgName, options.InstanceIds); err != nil {
			return err
		}
	}

	if len(options.LoadBalancerNames) > 0 {
		if err := client.autoscalingDetachLoadBalancers(options.AsgName, options.LoadBalancerNames); err != nil {
			return err
		}
	}

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until the desired capacity instances are InService...")
		return client.waitUntilAutoScalingGroupDesiredState(options.AsgName)
	}

	return nil
}

func (client *Client) autoscalingDetachInstances(asgName string, instanceIds []*string) error {
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

func (client *Client) autoscalingDetachLoadBalancers(asgName string, loadBalancerNames []*string) error {
	params := &autoscaling.DetachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	if _, err := client.AutoScaling.DetachLoadBalancers(params); err != nil {
		return errors.Wrap(err, "DetachLoadBalancers failed:")
	}

	return nil
}
