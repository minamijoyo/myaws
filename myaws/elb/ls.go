package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/spf13/cobra"
)

// Ls describes ELBs.
func Ls(*cobra.Command, []string) {
	client := newELBClient()
	params := &elb.DescribeLoadBalancersInput{}

	response, err := client.DescribeLoadBalancers(params)
	if err != nil {
		panic(err)
	}

	for _, lb := range response.LoadBalancerDescriptions {
		fmt.Println(formatLoadBalancer(lb))
	}
}

func formatLoadBalancer(lb *elb.LoadBalancerDescription) string {
	output := []string{
		*lb.LoadBalancerName,
	}

	return strings.Join(output[:], "\t")
}
