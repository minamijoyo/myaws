package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// ECSServiceLsOptions customize the behavior of the Ls command.
type ECSServiceLsOptions struct {
	Cluster string
}

// ECSServiceLs describes ECS services.
func (client *Client) ECSServiceLs(options ECSServiceLsOptions) error {
	arns, err := client.ECS.ListServices(
		&ecs.ListServicesInput{
			Cluster: &options.Cluster,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "ListServices failed")
	}

	if len(arns.ServiceArns) == 0 {
		return errors.New("services not found")
	}

	services, err := client.ECS.DescribeServices(
		&ecs.DescribeServicesInput{
			Cluster:  &options.Cluster,
			Services: arns.ServiceArns,
		},
	)

	if len(services.Services) == 0 {
		return errors.New("ListServices succeed, but DescribeServices returns no services")
	}

	for _, service := range services.Services {
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
