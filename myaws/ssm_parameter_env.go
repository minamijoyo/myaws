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
	var parameters []*ssm.Parameter
	// Since GetSSMParameters does not have a hierarchy, it is necessary to retrieve all keys at first, then filter the target keys.
	// Use the DescribeParameters API when retrieving the list of keys.
	// The rate limit of the DescribeParameters API is not publicly available.
	// When your SSM Parameter store have large number of keys, it is possible to exceed the rate limit.
	// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_DescribeParameters.html
	// The GetParametersByPath API has a hierarchy separated by '/'.
	// This API has less impact to the rate limit.
	// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_GetParametersByPath.html
	if strings.HasPrefix(options.Name, "/") {
		var err error
		parameters, err = client.GetParametersByPath(&options.Name, true)
		if err != nil {
			return err
		}
	} else {
		metadata, err := client.FindSSMParameterMetadata(options.Name)
		if err != nil {
			return err
		}

		names := []*string{}
		for _, m := range metadata {
			names = append(names, m.Name)
		}

		parameters, err = client.GetSSMParameters(names, true)
		if err != nil {
			return err
		}
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
	if suffix[0] == '.'|'/' {
		suffix = suffix[1:]
	}
	// Flatten period and slash to underscore for nested keys.
	flatten := strings.Replace(suffix, ".", "_", -1)
	flatten = strings.Replace(flatten, "/", "_", -1)
	// The name of environment variable should be uppercase.
	name := strings.ToUpper(flatten)
	outputOptionName := ""

	if dockerFormat {
		// Output in docker environment variables format such as -e KEY=VALUE
		outputOptionName = "-e "
	}
	return fmt.Sprintf("%s%s=", outputOptionName, name) + *parameter.Value
}
