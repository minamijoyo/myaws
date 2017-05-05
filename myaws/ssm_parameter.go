package myaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// FindSSMParameters returns an array of parameters matching the name.
func (client *Client) FindSSMParameters(name string) ([]*ssm.ParameterMetadata, error) {
	var filter *ssm.ParametersFilter
	if len(name) > 0 {
		filter = &ssm.ParametersFilter{
			Key: aws.String("Name"),
			Values: []*string{
				aws.String(name),
			},
		}
	}
	filters := []*ssm.ParametersFilter{filter}

	params := &ssm.DescribeParametersInput{
		Filters: filters,
	}

	response, err := client.SSM.DescribeParameters(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeParameters failed:")
	}

	return response.Parameters, nil
}
