package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// SSMParameterEnvOptions customize the behavior of the ParameterEnv command.
type SSMParameterEnvOptions struct {
	Name string
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

	input := &ssm.GetParametersInput{
		Names:          names,
		WithDecryption: aws.Bool(true),
	}

	response, err := client.SSM.GetParameters(input)
	if err != nil {
		return errors.Wrap(err, "GetParameters failed:")
	}

	if len(response.InvalidParameters) > 0 {
		return errors.Errorf("InvalidParameters: %v", awsutil.Prettify(response.InvalidParameters))
	}

	output := []string{}
	for _, parameter := range response.Parameters {
		output = append(output, formatSSMParameterAsEnv(parameter, options.Name))
	}

	fmt.Fprintf(client.stdout, strings.Join(output[:], " "))
	return nil
}

func formatSSMParameterAsEnv(parameter *ssm.Parameter, prefix string) string {
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
	return fmt.Sprintf("%s=%s", name, *parameter.Value)
}
