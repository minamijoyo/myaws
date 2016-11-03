package myaws

import (
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

// NewConfig creates *aws.config from profile and region options.
// AWS credentials are checked in the order of
// the profile, environment variables, IAM Role.
func NewConfig() *aws.Config {
	return &aws.Config{
		Credentials: newCredentials(viper.GetString("profile")),
		Region:      getRegion(viper.GetString("region")),
	}
}

func newCredentials(profile string) *credentials.Credentials {
	return credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{
				Profile: profile,
			},
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(session.New(&aws.Config{
					HTTPClient: &http.Client{Timeout: 3000 * time.Millisecond},
				},
				)),
			},
		})
}

func getRegion(region string) *string {
	if region != "" {
		// get region from the arg
		return aws.String(region)
	}

	// get region from the environement variable
	return aws.String(os.Getenv("AWS_DEFAULT_REGION"))
}
