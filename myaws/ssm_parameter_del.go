package myaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// SSMParameterDelOptions customize the behavior of the ParameterDel command.
type SSMParameterDelOptions struct {
	Name string
}

// SSMParameterDel deletes SSM parameter.
func (client *Client) SSMParameterDel(options SSMParameterDelOptions) error {
	input := &ssm.DeleteParameterInput{
		Name: aws.String(options.Name),
	}

	_, err := client.SSM.DeleteParameter(input)
	if err != nil {
		return errors.Wrap(err, "DeleteParameters failed:")
	}

	return nil
}
