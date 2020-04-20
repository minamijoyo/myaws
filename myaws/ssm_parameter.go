package myaws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
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

// GetSSMParameters returns an array of parameters at once.
func (client *Client) GetSSMParameters(names []*string, withDecryption bool) ([]*ssm.Parameter, error) {
	results := []*ssm.Parameter{}

	// The AWS SSM GetPrameters API can only get 10 parameters at once.
	// To get 10 or more parameters, we need to call API multiple time.
	// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_GetParameters.html
	chunkSize := 10

	for i := 0; i < len(names); i += chunkSize {
		end := i + chunkSize
		if end > len(names) {
			end = len(names)
		}
		chunk := names[i:end]
		resultsPerChunk, err := client.getSSMParametersPerChunk(chunk, withDecryption)
		if err != nil {
			return nil, err
		}

		results = append(results, resultsPerChunk...)
	}

	return results, nil
}

// getSSMParametersPerChunk returns an array of parameters per chunk.
func (client *Client) getSSMParametersPerChunk(names []*string, withDecryption bool) ([]*ssm.Parameter, error) {
	input := &ssm.GetParametersInput{
		Names:          names,
		WithDecryption: aws.Bool(withDecryption),
	}

	response, err := client.SSM.GetParameters(input)
	if err != nil {
		return nil, errors.Wrap(err, "GetParameters failed:")
	}

	if len(response.InvalidParameters) > 0 {
		return nil, errors.Errorf("InvalidParameters: %v", awsutil.Prettify(response.InvalidParameters))
	}

	return response.Parameters, nil
}

func (client *Client) GetParametersByPath(path *string, withDecryption bool) ([]*ssm.Parameter, error) {
	input := &ssm.GetParametersByPathInput{
		Path:           path,
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(withDecryption),
	}

	var parameters []*ssm.Parameter
	err := client.SSM.GetParametersByPathPages(input,
		func(page *ssm.GetParametersByPathOutput, lastPage bool) bool {
			parameters = append(parameters, page.Parameters...)
			return true
		})
	if err != nil {
		return nil, errors.Wrap(err, "GetParameters failed:")
	}
	return parameters, nil
}
