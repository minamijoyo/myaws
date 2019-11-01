package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
	funk "github.com/thoas/go-funk"
)

// findECSNodes finds ECS container instances
func (client *Client) findECSNodes(cluster string) ([]*ecs.ContainerInstance, error) {
	arns, err := client.ECS.ListContainerInstances(
		&ecs.ListContainerInstancesInput{
			Cluster: &cluster,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "ListContainerInstances failed")
	}

	if len(arns.ContainerInstanceArns) == 0 {
		return nil, errors.New("container instances not found")
	}

	instances, err := client.ECS.DescribeContainerInstances(
		&ecs.DescribeContainerInstancesInput{
			Cluster:            &cluster,
			ContainerInstances: arns.ContainerInstanceArns,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeContainerInstances failed")
	}

	if len(instances.ContainerInstances) == 0 {
		return nil, errors.New("ListContainerInstances succeed, but DescribeContainerInstances returns no instances")
	}

	return instances.ContainerInstances, nil
}

// findECSService find ECS services.
func (client *Client) findECSServices(cluster string) ([]*ecs.Service, error) {
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

	// We can specify up to 10 services to describe in a single operation.
	// So we need to divide the list by 10.
	chunks := (funk.Chunk(serviceArns, 10)).([][]*string)
	services := []*ecs.Service{}
	for _, c := range chunks {
		ss, err := client.ECS.DescribeServices(
			&ecs.DescribeServicesInput{
				Cluster:  &cluster,
				Services: c,
			},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "DescribeServices failed")
		}
		services = append(services, (ss.Services)...)
	}

	if len(services) == 0 {
		return nil, errors.New("ListServices succeed, but DescribeServices returns no services")
	}

	return services, nil
}

func (client *Client) printECSStatus(cluster string) error {
	fmt.Fprintln(client.stdout, "[Service]")
	err := client.ECSServiceLs(ECSServiceLsOptions{
		Cluster: cluster,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(client.stdout, "[Node]")
	err = client.ECSNodeLs(ECSNodeLsOptions{
		Cluster: cluster,
	})
	if err != nil {
		return err
	}

	return nil
}
