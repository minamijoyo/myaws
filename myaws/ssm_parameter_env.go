package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSMParameterEnvOptions customize the behavior of the ParameterEnv command.
type SSMParameterEnvOptions struct {
	Name         string
	DockerFormat bool
}

// SSMParameterEnv prints SSM parameters as a list of environment variables.
func (client *Client) SSMParameterEnv(options SSMParameterEnvOptions) error {
	metadata, err := client.FindSSMParameterMetadata(options.Name)
	if err != nil {
		return err
	}

	names := []*string{}
	for _, m := range metadata {
		names = append(names, m.Name)
	}

	parameters, err := client.GetSSMParameters(names, true)
	if err != nil {
		return err
	}

	output := []string{}
	for _, parameter := range parameters {
		output = append(output, formatSSMParameterAsEnv(parameter, options.Name, options.DockerFormat))
	}

	fmt.Fprint(client.stdout, strings.Join(output[:], " "))
	return nil
}

func formatSSMParameterAsEnv(parameter *ssm.Parameter, prefix string, dockerFormat bool) string {
	// Drop prefix and get suffix as a key name.
	suffix := strings.Replace(*parameter.Name, prefix, "", 1)
	// if first character is period, then drop it.
	if suffix[0] == '.' {
		suffix = suffix[1:]
	}
	// Flatten period to underscore for nested keys.
	flatten := strings.Replace(suffix, ".", "_", -1)
	// The name of environment variable should be uppercase.
	name := strings.ToUpper(flatten)
	outputOptionName := ""

	if dockerFormat {
		// Output in docker environment variables format such as -e KEY=VALUE
		outputOptionName = "-e "
	}
	return fmt.Sprintf("%s%s=", outputOptionName, name) + *parameter.Value
}
