package myaws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
	funk "github.com/thoas/go-funk"
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
	if err := client.waitUntilECSContainerInstancesStatusWithContext(ctx, input, "DRAINING"); err != nil {
		return err
	}

	// wait until no running tasks on the container instances
	if err := client.waitUntilECSContainerInstancesNoRunningTaskWithContext(ctx, input); err != nil {
		return err
	}

	return nil
}

func (client *Client) waitUntilECSContainerInstancesStatusWithContext(ctx aws.Context, input *ecs.DescribeContainerInstancesInput, status string, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilECSContainerInstancesStatus",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "ContainerInstances[].Status",
				Expected: status,
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

// WaitUntilECSContainerInstancesAreRegistered is a helper function which waits until
// the ECS container instances are registered.
// Due to the current limitation of the implementation of `request.Waiter`,
// we need to wait it in two steps.
// 1. Wait until the number of container instances is targetCapacity.
// 2. Wait until container instances are ACTIVE state.
func (client *Client) WaitUntilECSContainerInstancesAreRegistered(cluster string, targetCapacity int64) error {
	ctx := aws.BackgroundContext()

	listInput := &ecs.ListContainerInstancesInput{
		Cluster: &cluster,
	}

	// Simple count the number of container instances
	if err := client.waitUntilECSContainerInstancesCountWithContext(ctx, listInput, targetCapacity); err != nil {
		return err
	}

	// build descirbe input
	arns, err := client.ECS.ListContainerInstancesWithContext(ctx, listInput)

	if err != nil {
		return errors.Wrapf(err, "ListContainerInstances failed")
	}

	describeInput := &ecs.DescribeContainerInstancesInput{
		Cluster:            &cluster,
		ContainerInstances: arns.ContainerInstanceArns,
	}

	// make sure container instances are ACTIVE state
	if err := client.waitUntilECSContainerInstancesStatusWithContext(ctx, describeInput, "ACTIVE"); err != nil {
		return err
	}

	return nil
}

func (client *Client) waitUntilECSContainerInstancesCountWithContext(ctx aws.Context, input *ecs.ListContainerInstancesInput, targetCapacity int64, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilECSContainerInstancesCount",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "length(ContainerInstanceArns[])",
				Expected: float64(targetCapacity),
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *ecs.ListContainerInstancesInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.ECS.ListContainerInstancesRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

func (client *Client) getECSServiceArns(cluster string) ([]*string, error) {
	serviceArns := []*string{}

	err := client.ECS.ListServicesPages(
		&ecs.ListServicesInput{
			Cluster: &cluster,
		},
		func(p *ecs.ListServicesOutput, lastPage bool) bool {
			serviceArns = append(serviceArns, p.ServiceArns...)
			return true
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "ListServices failed")
	}

	if len(serviceArns) == 0 {
		return nil, errors.New("services not found")
	}

	return serviceArns, nil
}

func (client *Client) getECSTargetGroupArns(cluster string, serviceArn string) ([]string, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []*string{&serviceArn},
	}
	response, err := client.ECS.DescribeServices(input)
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeServices failed")
	}

	if len(response.Services) != 1 {
		return nil, errors.Errorf("ECS.DescribeServices expects to return 1 service, but found %d services", len(response.Services))
	}

	s := response.Services[0]

	targetGroupArns := []string{}
	for _, lb := range s.LoadBalancers {
		targetGroupArns = append(targetGroupArns, *lb.TargetGroupArn)
	}

	return targetGroupArns, nil
}

// WaitUntilECSAllServicesStable is a helper function which wait until all ECS
// servcies are running the desired number of containers.
// The official (*ECS) WaitUntilServicesStable does not support more than 10
// services.
// We need to check 10 services at a time.
func (client *Client) WaitUntilECSAllServicesStable(cluster string) error {
	serviceArns, err := client.getECSServiceArns(cluster)
	if err != nil {
		return err
	}

	// We can specify up to 10 services to describe in a single operation.
	// So we need to divide the list by 10.
	chunks := (funk.Chunk(serviceArns, 10)).([][]*string)
	for _, c := range chunks {
		input := &ecs.DescribeServicesInput{
			Cluster:  &cluster,
			Services: c,
		}
		err := client.ECS.WaitUntilServicesStable(input)
		if err != nil {
			return errors.Wrapf(err, "WaitUntilServicesStable failed")
		}
	}

	return nil
}

// WaitUntilECSAllTargetsInService is a helper function which wait until all
// target related to ECS servcies are healthy.
func (client *Client) WaitUntilECSAllTargetsInService(cluster string) error {
	serviceArns, err := client.getECSServiceArns(cluster)
	if err != nil {
		return err
	}

	for _, s := range serviceArns {
		targetGroupArns, err := client.getECSTargetGroupArns(cluster, *s)
		if err != nil {
			return err
		}
		for _, t := range targetGroupArns {
			input := &elbv2.DescribeTargetHealthInput{
				TargetGroupArn: &t,
			}
			err := client.ELBV2.WaitUntilTargetInService(input)
			if err != nil {
				return errors.Wrapf(err, "WaitUntilTargetInService failed")
			}
		}
	}

	return nil
}
