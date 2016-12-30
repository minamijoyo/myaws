package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
)

// ELBLs describes ELBs.
func (client *Client) ELBLs() error {
	params := &elb.DescribeLoadBalancersInput{}

	response, err := client.ELB.DescribeLoadBalancers(params)
	if err != nil {
		return errors.Wrap(err, "DescribeLoadBalancers failed:")
	}

	for _, lb := range response.LoadBalancerDescriptions {
		fmt.Fprintln(client.stdout, formatLoadBalancer(lb))
	}

	return nil
}

func formatLoadBalancer(lb *elb.LoadBalancerDescription) string {
	output := []string{
		*lb.LoadBalancerName,
	}

	return strings.Join(output[:], "\t")
}
