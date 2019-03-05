package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awsutil"
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

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	// All old container instances are drained doesn't mean all services are stable.
	// It depends on the deployment strategy of each service.
	// We should make sure all services are stable
	fmt.Fprintln(client.stdout, "Wait until all ECS services stable...")
	client.WaitUntilECSAllServicesStable(options.Cluster)
	if err != nil {
		return err
	}

	if err = client.printECSStatus(options.Cluster); err != nil {
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

	if err = client.printECSStatus(options.Cluster); err != nil {
		return err
	}

	fmt.Fprintln(client.stdout, "end: ecs node renew")
	return nil
}
