package myaws

import "fmt"

// SSMParameterLsOptions customize the behavior of the ParameterLs command.
type SSMParameterLsOptions struct {
	Name string
}

// SSMParameterLs describes SSM parameters.
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
