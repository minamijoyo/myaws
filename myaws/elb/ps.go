package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// PsOptions customize the behavior of the Ps command.
type PsOptions struct {
	LoadBalancerName string
}

// Ps describes ELB's instance health status.
func Ps(client *myaws.Client, options PsOptions) error {
	params := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: &options.LoadBalancerName,
	}

	response, err := client.ELB.DescribeInstanceHealth(params)
	if err != nil {
		return errors.Wrap(err, "DescribeInstanceHealth failed:")
	}

	for _, state := range response.InstanceStates {
		fmt.Println(formatInstanceState(state))
	}

	return nil
}

func formatInstanceState(state *elb.InstanceState) string {
	output := []string{
		*state.InstanceId,
		*state.State,
	}

	return strings.Join(output[:], "\t")
}
