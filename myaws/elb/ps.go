package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/spf13/cobra"
)

// Ps describes ELB's instance health status.
func Ps(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "ELB_NAME is required.")
	}

	client := newELBClient()
	params := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(args[0]),
	}

	response, err := client.DescribeInstanceHealth(params)
	if err != nil {
		panic(err)
	}

	for _, state := range response.InstanceStates {
		fmt.Println(formatInstanceState(state))
	}
}

func formatInstanceState(state *elb.InstanceState) string {
	output := []string{
		*state.InstanceId,
		*state.State,
	}

	return strings.Join(output[:], "\t")
}
