package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
)

// ELBV2PsOptions customize the behavior of the Ps command.
type ELBV2PsOptions struct {
	TargetGroupName string
}

// ELBV2Ps describes ELBV2's instance health status.
func (client *Client) ELBV2Ps(options ELBV2PsOptions) error {
	targetGroupArn, err := client.findELBV2TargetGroup(options.TargetGroupName)
	if err != nil {
		return err
	}

	params := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: &targetGroupArn,
	}

	response, err := client.ELBV2.DescribeTargetHealth(params)
	if err != nil {
		return errors.Wrap(err, "DescribeTargetHealth failed:")
	}

	for _, d := range response.TargetHealthDescriptions {
		fmt.Fprintln(client.stdout, formatELBV2TargetHealthDescription(d))
	}

	return nil
}

func (client *Client) findELBV2TargetGroup(name string) (string, error) {
	params := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{&name},
	}

	response, err := client.ELBV2.DescribeTargetGroups(params)
	if err != nil {
		return "", errors.Wrap(err, "DescribeTargetGroups failed:")
	}

	if len(response.TargetGroups) != 1 {
		return "", errors.Errorf("ELBV2.DescribeTargetGroups expects to return 1 group, but found %d groups", len(response.TargetGroups))
	}

	t := response.TargetGroups[0]
	return *t.TargetGroupArn, nil
}

func formatELBV2TargetHealthDescription(d *elbv2.TargetHealthDescription) string {
	t := *d.Target
	h := *d.TargetHealth

	// If the state is healthy, a reason and a description are not provided.
	// All states and reasons are described in the API reference.
	// https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/APIReference/API_TargetHealth.html
	output := []string{
		*t.Id,
		fmt.Sprintf("%d", *t.Port),
		fmt.Sprintf("%12s", *h.State),
		aws.StringValue(h.Reason),
		aws.StringValue(h.Description),
	}

	return strings.Join(output[:], "\t")
}
