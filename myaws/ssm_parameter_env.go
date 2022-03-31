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
	QuoteValue   bool
}

// SSMParameterEnv prints SSM parameters as a list of environment variables.
func (client *Client) SSMParameterEnv(options SSMParameterEnvOptions) error {
	var parameters []*ssm.Parameter
	// Since GetSSMParameters does not have a hierarchy, it is necessary to
	// retrieve all keys at first, then filter the target keys. To do this, we
	// need to call the DescribeParameters API multiple times, but its rate limit
	// seems to be quite low and undocumented.
	// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_DescribeParameters.html
	// This causes a rate limit exception if your SSM Parameter store have large
	// number of keys.
	// https://github.com/minamijoyo/myaws/issues/31
	//
	// The GetParametersByPath API has a hierarchy separated by '/'. This means
	// it has less impact to the rate limit.
	// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_GetParametersByPath.html
	//
	// So we use it if the path seems to have a hierarchy and fall back to
	// the original behavior if not.
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
		output = append(output, formatSSMParameterAsEnv(parameter, options.Name, options.DockerFormat, options.QuoteValue))
	}

	fmt.Fprint(client.stdout, strings.Join(output[:], " "))
	return nil
}

func formatSSMParameterAsEnv(parameter *ssm.Parameter, prefix string, dockerFormat bool, quoteValue bool) string {
	// Drop prefix and get suffix as a key name.
	suffix := strings.Replace(*parameter.Name, prefix, "", 1)
	// if first character is period, then drop it.
	if suffix[0] == '.' || suffix[0] == '/' {
		suffix = suffix[1:]
	}
	// Flatten period and slash to underscore for nested keys.
	flattenDot := strings.Replace(suffix, ".", "_", -1)
	flattenSlash := strings.Replace(flattenDot, "/", "_", -1)
	// The name of environment variable should be uppercase.
	name := strings.ToUpper(flattenSlash)
	outputOptionName := ""

	if dockerFormat {
		// Output in docker environment variables format such as -e KEY=VALUE
		outputOptionName = "-e "
	}

	value := *parameter.Value
	if quoteValue {
		value = "'" + *parameter.Value + "'"
	}
	return fmt.Sprintf("%s%s=%s", outputOptionName, name, value)
}
