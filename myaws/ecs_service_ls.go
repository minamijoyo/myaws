package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
	funk "github.com/thoas/go-funk"
)

// ECSServiceLsOptions customize the behavior of the Ls command.
type ECSServiceLsOptions struct {
	Cluster string
}

// ECSServiceLs describes ECS services.
func (client *Client) ECSServiceLs(options ECSServiceLsOptions) error {
	serviceArns := []*string{}

	err := client.ECS.ListServicesPages(
		&ecs.ListServicesInput{
			Cluster: &options.Cluster,
		},
		func(p *ecs.ListServicesOutput, lastPage bool) bool {
			serviceArns = append(serviceArns, p.ServiceArns...)
			return true
		},
	)
	if err != nil {
		return errors.Wrapf(err, "ListServices failed")
	}

	if len(serviceArns) == 0 {
		return errors.New("services not found")
	}

	// We can specify up to 10 services to describe in a single operation.
	// So we need to divide the list by 10.
	chunks := (funk.Chunk(serviceArns, 10)).([][]*string)
	services := []*ecs.Service{}
	for _, c := range chunks {
		ss, err := client.ECS.DescribeServices(
			&ecs.DescribeServicesInput{
				Cluster:  &options.Cluster,
				Services: c,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "DescribeServices failed")
		}
		services = append(services, (ss.Services)...)
	}

	if len(services) == 0 {
		return errors.New("ListServices succeed, but DescribeServices returns no services")
	}
	for _, service := range services {
		fmt.Fprintln(client.stdout, formatECSService(client, options, service))
	}

	return nil
}

func formatECSService(client *Client, options ECSServiceLsOptions, service *ecs.Service) string {
	taskDefinitions := strings.Split(*service.TaskDefinition, "/")

	return fmt.Sprintf("%d\t%d\t%d\t%d\t%s\t%s",
		*service.DesiredCount,
		*service.RunningCount,
		*service.PendingCount,
		len(service.Deployments),
		*service.ServiceName,
		taskDefinitions[1],
	)
}
