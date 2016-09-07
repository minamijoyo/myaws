package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

func newEC2Client() *ec2.EC2 {
	return ec2.New(
		session.New(),
		&aws.Config{
			Credentials: newCredentials(viper.GetString("profile")),
			Region:      aws.String(viper.GetString("region")),
		},
	)
}

func newCredentials(profile string) *credentials.Credentials {
	var cred *credentials.Credentials
	if profile != "" {
		cred = credentials.NewSharedCredentials("", profile)
	}
	return cred
}
