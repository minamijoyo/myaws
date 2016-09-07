package ec2

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
)

func newEC2Client() *ec2.EC2 {
	return ec2.New(
		session.New(),
		&aws.Config{
			Credentials: newCredentials(viper.GetString("profile")),
			Region:      getRegion(viper.GetString("region")),
		},
	)
}

func newCredentials(profile string) *credentials.Credentials {
	if profile != "" {
		return credentials.NewSharedCredentials("", profile)
	} else {
		return credentials.NewEnvCredentials()
	}
}

func getRegion(region string) *string {
	if region != "" {
		return aws.String(region)
	} else {
		return aws.String(os.Getenv("AWS_DEFAULT_REGION"))
	}
}

func formatTime(t *time.Time) string {
	location, err := time.LoadLocation(viper.GetString("timezone"))
	if err != nil {
		panic(err)
	}

	if viper.GetBool("humanize") {
		return humanize.Time(t.In(location))
	} else {
		return t.In(location).Format("2006-01-02 15:04:05")
	}
}
