package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// SSMParameterGetOptions customize the behavior of the ParameterGet command.
type SSMParameterGetOptions struct {
	Names          []*string
	WithDecryption bool
}

// SSMParameterGet get values from SSM parameter store with KMS decryption.
func (client *Client) SSMParameterGet(options SSMParameterGetOptions) error {
	params := &ssm.GetParametersInput{
		Names:          options.Names,
		WithDecryption: &options.WithDecryption,
	}

	response, err := client.SSM.GetParameters(params)
	if err != nil {
		return errors.Wrap(err, "GetParameters failed:")
	}

	if len(response.InvalidParameters) > 0 {
		return errors.Errorf("InvalidParameters: %v", awsutil.Prettify(response.InvalidParameters))
	}

	for _, parameter := range response.Parameters {
		fmt.Fprintln(client.stdout, formatSSMParameter(parameter))
	}

	return nil
}

func formatSSMParameter(parameter *ssm.Parameter) string {
	return *parameter.Value
}
