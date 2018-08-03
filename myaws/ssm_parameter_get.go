package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSMParameterGetOptions customize the behavior of the ParameterGet command.
type SSMParameterGetOptions struct {
	Names          []*string
	WithDecryption bool
}

// SSMParameterGet get values from SSM parameter store with KMS decryption.
func (client *Client) SSMParameterGet(options SSMParameterGetOptions) error {
	parameters, err := client.GetSSMParameters(options.Names, options.WithDecryption)
	if err != nil {
		return err
	}

	for _, parameter := range parameters {
		fmt.Fprintln(client.stdout, formatSSMParameter(parameter))
	}

	return nil
}

func formatSSMParameter(parameter *ssm.Parameter) string {
	return *parameter.Value
}
