package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSMParameterLsOptions customize the behavior of the ParameterLs command.
type SSMParameterLsOptions struct {
	Name string
}

// SSMParameterLs describes SSM parameters.
func (client *Client) SSMParameterLs(options SSMParameterLsOptions) error {
	metadata, err := client.FindSSMParameterMetadata(options.Name)
	if err != nil {
		return err
	}

	for _, m := range metadata {
		fmt.Fprintln(client.stdout, formatSSMParameterMetadata(m))
	}

	return nil
}

func formatSSMParameterMetadata(m *ssm.ParameterMetadata) string {
	keyid := formatSSMParameterKeyID(m)
	return fmt.Sprintf("%s\t%s\t%s", *m.Name, *m.Type, keyid)
}

func formatSSMParameterKeyID(m *ssm.ParameterMetadata) string {
	if m.KeyId == nil {
		return ""
	}
	return *m.KeyId
}
