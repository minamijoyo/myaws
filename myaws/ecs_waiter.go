package myaws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// WaitUntilECSContainerInstancesAreDrained is a helper function which waits until
// the ECS container instances are drained.
// Due to the current limitation of the implementation of `request.Waiter`,
// we need to wait it in two steps.
// 1. Wait until container instances are DRAINING state.
// 2. Wait until no running tasks on the container instances.
func (client *Client) WaitUntilECSContainerInstancesAreDrained(cluster string, containerInstances []*string) error {
	ctx := aws.BackgroundContext()

	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            &cluster,
		ContainerInstances: containerInstances,
	}

	// make sure container instances are DRAINING state
	if err := client.waitUntilECSContainerInstancesDrainingStateWithContext(ctx, input); err != nil {
		return err
	}

	// wait until no running tasks on the container instances
	if err := client.waitUntilECSContainerInstancesNoRunningTaskWithContext(ctx, input); err != nil {
		return err
	}

	return nil
}

func (client *Client) waitUntilECSContainerInstancesDrainingStateWithContext(ctx aws.Context, input *ecs.DescribeContainerInstancesInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilECSContainerInstancesDrainingState",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "ContainerInstances[].Status",
				Expected: "DRAINING",
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *ecs.DescribeContainerInstancesInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.ECS.DescribeContainerInstancesRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

func (client *Client) waitUntilECSContainerInstancesNoRunningTaskWithContext(ctx aws.Context, input *ecs.DescribeContainerInstancesInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilECSContainerInstancesNoRunningTask",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "ContainerInstances[].RunningTasksCount",
				Expected: aws.Int64(0),
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *ecs.DescribeContainerInstancesInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.ECS.DescribeContainerInstancesRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
