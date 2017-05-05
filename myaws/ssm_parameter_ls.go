package myaws

import "fmt"

// SSMParameterLsOptions customize the behavior of the ParameterGet command.
type SSMParameterLsOptions struct {
	Name string
}

// SSMParameterLs get values from SSM parameter store with KMS decryption.
func (client *Client) SSMParameterLs(options SSMParameterLsOptions) error {
	parameters, err := client.FindSSMParameters(options.Name)
	if err != nil {
		return err
	}

	for _, parameter := range parameters {
		fmt.Fprintln(client.stdout, *parameter.Name)
	}

	return nil
}
