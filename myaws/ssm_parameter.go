package myaws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// FindSSMParameterMetadata returns an array of parameter metadata matching the name.
func (client *Client) FindSSMParameterMetadata(name string) ([]*ssm.ParameterMetadata, error) {
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

	input := &ssm.DescribeParametersInput{
		Filters: filters,
	}

	// We need to fetch all pages to get results.
	// The request timeout should be set in the caller context,
	// but for the moment we will create a context here.
	metadata := []*ssm.ParameterMetadata{}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := client.SSM.DescribeParametersPagesWithContext(ctx,
		input,
		func(page *ssm.DescribeParametersOutput, lastPage bool) bool {
			metadata = append(metadata, page.Parameters...)
			return true
		})

	if err != nil {
		return nil, errors.Wrap(err, "DescribeParameters failed:")
	}

	return metadata, nil
}
