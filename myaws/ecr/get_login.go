package ecr

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"
)

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
	return fmt.Sprintf("docker login -u AWS -p %s -e none %s", *authData.AuthorizationToken, *authData.ProxyEndpoint)
}
