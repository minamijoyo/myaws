package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
)

// ELBPsOptions customize the behavior of the Ps command.
type ELBPsOptions struct {
	LoadBalancerName string
}

// ELBPs describes ELB's instance health status.
func (client *Client) ELBPs(options ELBPsOptions) error {
	params := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: &options.LoadBalancerName,
	}

	response, err := client.ELB.DescribeInstanceHealth(params)
	if err != nil {
		return errors.Wrap(err, "DescribeInstanceHealth failed:")
	}

	for _, state := range response.InstanceStates {
		fmt.Println(formatELBInstanceState(state))
	}

	return nil
}

func formatELBInstanceState(state *elb.InstanceState) string {
	output := []string{
		*state.InstanceId,
		*state.State,
	}

	return strings.Join(output[:], "\t")
}
