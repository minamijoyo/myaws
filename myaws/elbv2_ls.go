package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
)

// ELBV2Ls describes ELBV2s.
func (client *Client) ELBV2Ls() error {
	params := &elbv2.DescribeLoadBalancersInput{}

	response, err := client.ELBV2.DescribeLoadBalancers(params)
	if err != nil {
		return errors.Wrap(err, "DescribeLoadBalancers failed:")
	}

	for _, lb := range response.LoadBalancers {
		fmt.Fprintln(client.stdout, formatLoadBalancerV2(lb))
	}

	return nil
}

func formatLoadBalancerV2(lb *elbv2.LoadBalancer) string {
	output := []string{
		*lb.Type,
		*lb.LoadBalancerName,
	}

	return strings.Join(output[:], "\t")
}
