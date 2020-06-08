package myaws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

// ECSServiceUpdateOptions customize the behavior of the Update command.
type ECSServiceUpdateOptions struct {
	Cluster      string
	Service      string
	DesiredCount int64
	Wait         bool
	Timeout      time.Duration
}

// ECSServiceUpdate update ECS services.
func (client *Client) ECSServiceUpdate(options ECSServiceUpdateOptions) error {
	input := &ecs.UpdateServiceInput{
		Cluster:      &options.Cluster,
		Service:      &options.Service,
		DesiredCount: &options.DesiredCount,
	}

	_, err := client.ECS.UpdateService(input)
	if err != nil {
		return errors.Wrapf(err, "UpdateService failed")
	}

	if options.Wait {
		fmt.Fprintln(client.stdout, "Wait until the service stable...")
		ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
		defer cancel()

		err = client.ECS.WaitUntilServicesStableWithContext(
			ctx,
			&ecs.DescribeServicesInput{
				Cluster:  &options.Cluster,
				Services: []*string{&options.Service},
			})
		if err != nil {
			return errors.Wrapf(err, "WaitUntilServicesStable failed")
		}
	}

	return nil
}
