package ecr

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"
)

// GetLogin gets docker login command with authorization token for ECR.
func GetLogin(*cobra.Command, []string) {
	client := newECRClient()
	params := &ecr.GetAuthorizationTokenInput{}

	response, err := client.GetAuthorizationToken(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(formatAuthorizationData(response.AuthorizationData))
}

func formatAuthorizationData(authDataList []*ecr.AuthorizationData) string {
	output := []string{}
	for _, authData := range authDataList {
		output = append(output, formatDockerLoginCommand(authData))
	}
	return strings.Join(output[:], "\n")
}

func formatDockerLoginCommand(authData *ecr.AuthorizationData) string {
	return fmt.Sprintf(
		"docker login -u AWS -p %s %s",
		decodePassword(*authData.AuthorizationToken),
		*authData.ProxyEndpoint,
	)
}

func decodePassword(authToken string) string {
	userAndPassword, _ := base64.StdEncoding.DecodeString(authToken)
	s := strings.Split(string(userAndPassword), ":")
	return s[1]
}
