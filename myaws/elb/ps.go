package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Ps describes ELB's instance health status.
func Ps(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("ELB_NAME is required")
	}

	client := newELBClient()
	params := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(args[0]),
	}

	response, err := client.DescribeInstanceHealth(params)
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
