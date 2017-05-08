package myaws

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// SSMParameterPutOptions customize the behavior of the ParameterPut command.
type SSMParameterPutOptions struct {
	Name  string
	Value string
	KeyID string
}

// SSMParameterPut put value to SSM parameter store with KMS encryption.
func (client *Client) SSMParameterPut(options SSMParameterPutOptions) error {
	overwrite := true

	var parameterType string
	var keyID *string
	if options.KeyID != "" {
		parameterType = "SecureString"
		keyID = &options.KeyID
	} else {
		parameterType = "String"
		// keyID must be nil when type is String.
		keyID = nil
	}

	input := &ssm.PutParameterInput{
		Name:      &options.Name,
		Value:     &options.Value,
		KeyId:     keyID,
		Type:      &parameterType,
		Overwrite: &overwrite,
	}

	_, err := client.SSM.PutParameter(input)
	if err != nil {
		return errors.Wrap(err, "PutParameter failed:")
	}

	return nil
}
