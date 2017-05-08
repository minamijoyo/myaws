package myaws

import "fmt"

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
		fmt.Fprintln(client.stdout, *m.Name)
	}

	return nil
}
