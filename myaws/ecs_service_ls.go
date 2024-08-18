package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSServiceLsOptions customize the behavior of the Ls command.
type ECSServiceLsOptions struct {
	Cluster     string
	PrintHeader bool
}

// ECSServiceLs describes ECS services.
func (client *Client) ECSServiceLs(options ECSServiceLsOptions) error {
	services, err := client.findECSServices(options.Cluster)
	if err != nil {
		return err
	}

	if options.PrintHeader {
		header := fmt.Sprintf("%s\t%s\t%s\t%s\t%-32s\t%s",
			"Desired",
			"Running",
			"Pending",
			"Deploy",
			"Service",
			"TaskDefinition",
		)
		fmt.Fprintln(client.stdout, header)
	}

	for _, service := range services {
		fmt.Fprintln(client.stdout, formatECSService(client, options, service))
	}

	return nil
}

func formatECSService(_ *Client, _ ECSServiceLsOptions, service *ecs.Service) string {
	taskDefinitions := strings.Split(*service.TaskDefinition, "/")

	return fmt.Sprintf("%d\t%d\t%d\t%d\t%-32s\t%s",
		*service.DesiredCount,
		*service.RunningCount,
		*service.PendingCount,
		len(service.Deployments),
		*service.ServiceName,
		taskDefinitions[1],
	)
}
