package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Ls describes ELBs.
func Ls(*cobra.Command, []string) error {
	client := newELBClient()
	params := &elb.DescribeLoadBalancersInput{}

	response, err := client.DescribeLoadBalancers(params)
	if err != nil {
		return errors.Wrap(err, "DescribeLoadBalancers failed:")
	}

	for _, lb := range response.LoadBalancerDescriptions {
		fmt.Println(formatLoadBalancer(lb))
	}

	return nil
}

func formatLoadBalancer(lb *elb.LoadBalancerDescription) string {
	output := []string{
		*lb.LoadBalancerName,
	}

	return strings.Join(output[:], "\t")
}
