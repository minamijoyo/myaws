package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// AutoscalingAttachOptions customize the behavior of the Attach command.
type AutoscalingAttachOptions struct {
	AsgName           string
	InstanceIds       []*string
	LoadBalancerNames []*string
	Wait              bool
}

// AutoscalingAttach attaches instances or load balancers from autoscaling group.
func (client *Client) AutoscalingAttach(options AutoscalingAttachOptions) error {
	if len(options.InstanceIds) > 0 {
		if err := client.autoscalingAttachInstances(options.AsgName, options.InstanceIds); err != nil {
			return err
		}
	}

	if len(options.LoadBalancerNames) > 0 {
		if err := client.autoscalingAttachLoadBalancers(options.AsgName, options.LoadBalancerNames); err != nil {
			return err
		}
	}

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until desired capacity instances are InService...")
		return client.WaitUntilAutoScalingGroupStable(options.AsgName)
	}

	return nil
}

func (client *Client) autoscalingAttachInstances(asgName string, instanceIds []*string) error {
	params := &autoscaling.AttachInstancesInput{
		AutoScalingGroupName: &asgName,
		InstanceIds:          instanceIds,
	}

	if _, err := client.AutoScaling.AttachInstances(params); err != nil {
		return errors.Wrap(err, "AttachInstances failed:")
	}

	return nil
}

func (client *Client) autoscalingAttachLoadBalancers(asgName string, loadBalancerNames []*string) error {
	params := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: &asgName,
		LoadBalancerNames:    loadBalancerNames,
	}

	if _, err := client.AutoScaling.AttachLoadBalancers(params); err != nil {
		return errors.Wrap(err, "AttachLoadBalancers failed:")
	}

	return nil
}
