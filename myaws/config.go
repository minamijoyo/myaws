package myaws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
)

// newConfig creates *aws.config from profile and region options.
// AWS credentials are checked in the order of
// profile, environment variables, IAM Task Role (ECS), IAM Role.
// Unlike the aws default, load profile before environment variables
// because we want to prioritize explicit arguments over the environment.
func newConfig(profile string, region string) *aws.Config {
	defaultConfig := defaults.Get().Config
	cred := newCredentials(getenv("AWS_DEFAULT_PROFILE", profile), getenv("AWS_DEFAULT_REGION", region))
	return defaultConfig.WithCredentials(cred).WithRegion(getenv("AWS_DEFAULT_REGION", region))
}

func newCredentials(profile string, region string) *credentials.Credentials {
	// temporary config to resolve RemoteCredProvider
	tmpConfig := defaults.Get().Config.WithRegion(region)
	tmpHandlers := defaults.Handlers()

	return credentials.NewChainCredentials(
		[]credentials.Provider{
			// Read profile before environment variables
			&credentials.SharedCredentialsProvider{
				Profile: profile,
			},
			&credentials.EnvProvider{},
			// for IAM Task Role (ECS) and IAM Role
			defaults.RemoteCredProvider(*tmpConfig, tmpHandlers),
		})
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
