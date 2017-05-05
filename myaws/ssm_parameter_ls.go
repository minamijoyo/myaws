package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// SSMParameterLsOptions customize the behavior of the ParameterGet command.
type SSMParameterLsOptions struct {
	Name string
}

// SSMParameterLs get values from SSM parameter store with KMS decryption.
func (client *Client) SSMParameterLs(options SSMParameterLsOptions) error {
	filter := &ssm.ParametersFilter{
		Key: aws.String("Name"),
		Values: []*string{
			aws.String(options.Name),
		},
	}
	filters := []*ssm.ParametersFilter{filter}

	params := &ssm.DescribeParametersInput{
		Filters: filters,
	}

	response, err := client.SSM.DescribeParameters(params)
	if err != nil {
		return errors.Wrap(err, "DescribeParameters failed:")
	}

	for _, parameter := range response.Parameters {
		fmt.Fprintln(client.stdout, *parameter.Name)
	}

	return nil
}
