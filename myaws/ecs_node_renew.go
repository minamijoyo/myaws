package myaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
)

// ECSNodeRenewOptions customize the behavior of the Renew command.
type ECSNodeRenewOptions struct {
	Cluster string
	AsgName string
}

// ECSNodeRenew renew ECS container instances with blue-green deployment.
// This method is an automation process to renew your ECS container instances
// if you update the AMI. creates new instances, drains the old instances,
// and discards the old instances.
func (client *Client) ECSNodeRenew(options ECSNodeRenewOptions) error {
	fmt.Fprintf(client.stdout, "start: ecs node renew\noptions: %s\n", awsutil.Prettify(options))

	if err := client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// get the current desired capacity
	desiredCapacity, err := client.getAutoScalingGroupDesiredCapacity(options.AsgName)
	if err != nil {
		return err
	}

	// list the current container instances
	oldNodes, err := client.findECSNodes(options.Cluster)
	if err != nil {
		return err
	}

	if len(oldNodes) != int(desiredCapacity) {
		return errors.Errorf("assertion failed: currentCapacity(%d) != desiredCapacity(%d)", len(oldNodes), desiredCapacity)
	}

	// Update the desired capacity and wait until new instances are InService
	// We simply double the number of instances here.
	// If you need more flexible control, please implement a strategy such as
	// rolling update.
	targetCapacity := desiredCapacity * 2

	fmt.Fprintf(client.stdout, "Update autoscaling group %s (DesiredCapacity: %d => %d)\n", options.AsgName, desiredCapacity, targetCapacity)

	err = client.AutoscalingUpdate(AutoscalingUpdateOptions{
		AsgName:         options.AsgName,
		DesiredCapacity: targetCapacity,
		Wait:            true,
	})
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// A status of instance in autoscaling group is InService doesn't mean the
	// container instance is registered. We should make sure container instances
	// are registered
	fmt.Fprintln(client.stdout, "Wait until ECS container instances are registered...")
	err = client.WaitUntilECSContainerInstancesAreRegistered(options.Cluster, targetCapacity)
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// drain old container instances and wait until no task running
	oldNodeArns := []*string{}
	for _, oldNode := range oldNodes {
		oldNodeArns = append(oldNodeArns, oldNode.ContainerInstanceArn)
	}
	fmt.Fprintf(client.stdout, "Drain old container instances and wait until no task running...\n%v\n", awsutil.Prettify(oldNodeArns))
	err = client.ECSNodeDrain(ECSNodeDrainOptions{
		Cluster:            options.Cluster,
		ContainerInstances: oldNodeArns,
		Wait:               true,
	})
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// All old container instances are drained doesn't mean all services are stable.
	// It depends on the deployment strategy of each service.
	// We should make sure all services are stable
	fmt.Fprintln(client.stdout, "Wait until all ECS services stable...")
	err = client.WaitUntilECSAllServicesStable(options.Cluster)
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// A stable state for all services does not mean that all targets are healthy.
	// We need to explicitly confirm it.
	fmt.Fprintln(client.stdout, "Wait until all targets healthy...")
	err = client.WaitUntilECSAllTargetsInService(options.Cluster)
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// Select instances to protect from scale in.
	// By setting "scale-in protection" to instances created at scale-out,
	// the intended instances (instances created before scale-in) are only terminated at scale-in process.
	protectInstanceIds, err := client.selectInstanceToProtectFromScaleIn(oldNodes, options.Cluster)
	if err != nil {
		return err
	}

	fmt.Fprintln(client.stdout, "Setting scale in protection: ", awsutil.Prettify(protectInstanceIds))
	// set "scale in protection" to instances created at scale-out.
	err = client.AutoScalingSetInstanceProtection(AutoScalingSetInstanceProtectionOptions{
		options.AsgName,
		protectInstanceIds,
		true})
	if err != nil {
		return err
	}

	// restore the desired capacity and wait until old instances are discarded
	fmt.Fprintf(client.stdout, "Update autoscaling group %s (DesiredCapacity: %d => %d)\n", options.AsgName, targetCapacity, desiredCapacity)

	err = client.AutoscalingUpdate(AutoscalingUpdateOptions{
		AsgName:         options.AsgName,
		DesiredCapacity: desiredCapacity,
		Wait:            true,
	})
	if err != nil {
		return err
	}

	// remove "scale in protection" to instances created at scale-out.
	fmt.Fprintln(client.stdout, "Removing scale in protection: ", awsutil.Prettify(protectInstanceIds))
	err = client.AutoScalingSetInstanceProtection(AutoScalingSetInstanceProtectionOptions{
		options.AsgName,
		protectInstanceIds,
		false})
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	fmt.Fprintln(client.stdout, "end: ecs node renew")
	return nil
}

// selectInstanceToProtectFromScaleIn selects instance to protect from Scale in.
// instance select rule:
//   instances after scale out - instances before scale out - instances which already set `InstanceProtection==true`
func (client *Client) selectInstanceToProtectFromScaleIn(oldNodes []*ecs.ContainerInstance, cluster string) ([]*string, error) {
	// Get a list of instance IDs before auto scaling
	var oldInstanceIds []*string
	for _, oldNode := range oldNodes {
		oldInstanceIds = append(oldInstanceIds, oldNode.Ec2InstanceId)
	}

	// Get a list of instances after auto scaling
	allNodes, err := client.findECSNodes(cluster)
	if err != nil {
		return nil, err
	}

	// Get a list of instance IDs after auto scaling
	var allInstanceIds []*string
	for _, allNode := range allNodes {
		allInstanceIds = append(allInstanceIds, allNode.Ec2InstanceId)
	}

	// get newly created nodes (allInstanceIds - oldInstanceIds)
	newInstanceIds := difference(allInstanceIds, oldInstanceIds)

	// exclude ProtectedFromScaleIn == true nodes
	params := &autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: newInstanceIds,
	}
	response, err := client.AutoScaling.DescribeAutoScalingInstances(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeAutoScalingGroups failed:")
	}

	var targetInstanceIds []*string
	for _, instance := range response.AutoScalingInstances {
		if *instance.ProtectedFromScaleIn == false {
			targetInstanceIds = append(targetInstanceIds, instance.InstanceId)
		}
	}
	return targetInstanceIds, nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []*string) []*string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[*x] = struct{}{}
	}
	var diff []*string
	for _, x := range a {
		if _, ok := mb[*x]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}
