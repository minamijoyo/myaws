package myaws

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// ECSServiceUpdateOptions customize the behavior of the Update command.
type ECSServiceUpdateOptions struct {
	Cluster      string
	Service      string
	DesiredCount *int64
	Wait         bool
	Timeout      time.Duration
	Force        bool
}

// ECSServiceUpdate update ECS services.
func (client *Client) ECSServiceUpdate(options ECSServiceUpdateOptions) error {
	input := &ecs.UpdateServiceInput{
		Cluster:            &options.Cluster,
		Service:            &options.Service,
		DesiredCount:       options.DesiredCount,
		ForceNewDeployment: &options.Force,
	}

	if options.Force {
		// When updating a ECS task definition, We want to deploy a new task with
		// new revision.
		// The ECS UpdateService API doesn't update the revision without the
		// TaskDefinition parameter. To update the revision, we set only a task
		// family without revision. If a revision is not specified, the latest
		// ACTIVE revision is used. If the current task uses the latest ACTIVE one,
		// it hasn't no side-effect.
		// Since the Force option expects to deploy and converge to the latest
		// state, when Force option is enabled, we implicitly fetch and set the
		// task family.
		family, err := client.getECSTaskDefinitionFamily(options.Cluster, options.Service)
		if err != nil {
			return err
		}
		input.TaskDefinition = &family
	}

	_, err := client.ECS.UpdateService(input)
	if err != nil {
		return errors.Wrapf(err, "UpdateService failed")
	}

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until the service stable...")
		ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
		defer cancel()
		err = client.WaitUntilECSServicesStableWithContext(ctx, options.Cluster, []string{options.Service})
		if err != nil {
			return err
		}
	}

	return nil
}

// getECSTaskDefinitionFamily returns a family name of ECS task definition used by a given ECS service.
func (client *Client) getECSTaskDefinitionFamily(cluster string, service string) (string, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []*string{&service},
	}

	resp, err := client.ECS.DescribeServices(input)
	if err != nil {
		return "", errors.Wrapf(err, "DescribeServices failed")
	}

	if len(resp.Services) == 0 {
		return "", fmt.Errorf("service not fould: cluster = %s, service = %s", cluster, service)
	}

	arn := *resp.Services[0].TaskDefinition

	// arn:aws:ecs:<region>:<account-id>:task-definition/<family>:<revison>
	re := regexp.MustCompile(`^arn:aws:ecs:.+:.+:task-definition/(.+):.+$`)
	matched := re.FindStringSubmatch(arn)
	if len(matched) != 2 || len(matched[1]) == 0 {
		return "", fmt.Errorf("failed to parse task definition arn: %s", arn)
	}
	family := matched[1]

	return family, nil
}
