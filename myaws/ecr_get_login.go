package myaws

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
)

// ECRGetLoginOptions customize the behavior of the ECRGetLogin command.
type ECRGetLoginOptions struct {
	RegistryIds []*string
}

// ECRGetLogin gets docker login command with authorization token for ECR.
func (client *Client) ECRGetLogin(options ECRGetLoginOptions) error {
	params := &ecr.GetAuthorizationTokenInput{}

	if len(options.RegistryIds) > 0 {
		params.RegistryIds = options.RegistryIds
	}

	response, err := client.ECR.GetAuthorizationToken(params)
	if err != nil {
		return errors.Wrap(err, "GetAuthorizationToken failed:")
	}
	fmt.Fprintln(client.stdout, formatECRAuthorizationData(response.AuthorizationData))

	return nil
}

func formatECRAuthorizationData(authDataList []*ecr.AuthorizationData) string {
	output := []string{}
	for _, authData := range authDataList {
		output = append(output, formatECRDockerLoginCommand(authData))
	}
	return strings.Join(output[:], "\n")
}

func formatECRDockerLoginCommand(authData *ecr.AuthorizationData) string {
	return fmt.Sprintf(
		"docker login -u AWS -p %s %s",
		decodeECRPassword(*authData.AuthorizationToken),
		*authData.ProxyEndpoint,
	)
}

func decodeECRPassword(authToken string) string {
	userAndPassword, _ := base64.StdEncoding.DecodeString(authToken)
	s := strings.Split(string(userAndPassword), ":")
	return s[1]
}
