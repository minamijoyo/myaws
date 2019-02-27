package myaws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// WaitUntilAutoScalingGroupStable is a helper function which waits until
// the AutoScaling Group converges to the desired state. We only check the
// status of AutoScaling Group. If the ASG has an ELB, the health check status
// of ELB can link with the health status of ASG, so we don't check the status
// of ELB here.
// Due to the current limitation of the implementation of `request.Waiter`,
// we need to wait it in two steps.
// 1. Wait until the number of instances equals `DesiredCapacity`.
// 2. Wait until all instances are InService.
func (client *Client) WaitUntilAutoScalingGroupStable(asgName string) error {
	desiredCapacity, err := client.getAutoScalingGroupDesiredCapacity(asgName)
	if err != nil {
		return err
	}

	ctx := aws.BackgroundContext()
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asgName},
	}

	// make sure instances are created or terminated.
	err = client.waitUntilAutoScalingGroupNumberOfInstancesEqualsDesiredCapacityWithContext(
		ctx,
		desiredCapacity,
		input,
	)

	if err != nil {
		return err
	}

	// if the desired state is no instance, we just return here.
	if desiredCapacity == 0 {
		return nil
	}

	// check all instances are InService state.
	return client.waitUntilAutoScalingGroupAllInstancesAreInServiceWithContext(ctx, input)
}

// waitUntilAutoScalingGroupNumberOfInstancesEqualsDesiredCapacityWithContext
// waits the number of instances equals DesiredCapacity.
func (client *Client) waitUntilAutoScalingGroupNumberOfInstancesEqualsDesiredCapacityWithContext(ctx aws.Context, desiredCapacity int64, input *autoscaling.DescribeAutoScalingGroupsInput, opts ...request.WaiterOption) error {
	// We implicitly assume that the number of AutoScalingGroup is only one to
	// simplify checking desiredCapacity. In our case, multiple AutoScalingGroup
	// doesn't pass this function.
	// Properties in the response returned by aws-sdk-go are reference types and
	// not primitive. Thus we cannot be directly compared on JMESPath.
	matcher := fmt.Sprintf("AutoScalingGroups[].[length(Instances) == `%d`][]", desiredCapacity)

	w := request.Waiter{
		Name:        "WaitUntilAutoScalingGroupNumberOfInstancesEqualsDesiredCapacity",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: matcher,
				Expected: true,
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *autoscaling.DescribeAutoScalingGroupsInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.AutoScaling.DescribeAutoScalingGroupsRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// getAutoScalingGroupDesiredCapacity is a helper function which returns
// DesiredCapacity of the specific AutoScalingGroup.
func (client *Client) getAutoScalingGroupDesiredCapacity(asgName string) (int64, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asgName},
	}

	response, err := client.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return 0, errors.Wrap(err, "getAutoScalingGroupDesiredCapacity failed:")
	}

	desiredCapacity := response.AutoScalingGroups[0].DesiredCapacity

	return *desiredCapacity, nil
}

// waitUntilAutoScalingGroupAllInstancesAreInService waits until all instances
// are in service.  Since the official `WaitUntilGroupInServiceWithContext` in
// aws-sdk-go checks as follow:
// contains(AutoScalingGroups[].[length(Instances[?LifecycleState=='InService']) >= MinSize][], `false`)
// But we found this doesn't work as expected. Properties in the response
// returned by aws-sdk-go are reference type and not primitive. Thus we can not
// be directly compared on JMESPath. So we implement a customized waiter here.
// When the number of desired instances increase or decrease, the affected
// instances are in states other than InService until the operation completes.
// So we should check that all the states of instances are InService.
func (client *Client) waitUntilAutoScalingGroupAllInstancesAreInServiceWithContext(ctx aws.Context, input *autoscaling.DescribeAutoScalingGroupsInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilAutoScalingGroupAllInstancesAreInService",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "AutoScalingGroups[].Instances[].LifecycleState",
				Expected: "InService",
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *autoscaling.DescribeAutoScalingGroupsInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.AutoScaling.DescribeAutoScalingGroupsRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
