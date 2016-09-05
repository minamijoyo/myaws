package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Ls(*cobra.Command, []string) {
	svc := ec2.New(
		session.New(),
		&aws.Config{
			Credentials: getCredentials(viper.GetString("profile")),
			Region:      aws.String(viper.GetString("region")),
		},
	)

	var stateFilter *ec2.Filter
	if viper.GetBool("ec2.ls.all") {
		stateFilter = &ec2.Filter{}
	} else {
		stateFilter = &ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		}
	}

	var tagFilter *ec2.Filter
	filterTag := viper.GetString("ec2.ls.filter-tag")
	if filterTag == "" {
	} else {
		tagParts := strings.Split(filterTag, ":")
		tagFilter = &ec2.Filter{
			Name: aws.String("tag:" + tagParts[0]),
			Values: []*string{
				aws.String("*" + tagParts[1] + "*"),
			},
		}
	}

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			stateFilter,
			tagFilter,
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			fmt.Println(formatInstance(inst))
		}
	}
}

func getCredentials(profile string) *credentials.Credentials {
	var cred *credentials.Credentials
	if profile != "" {
		cred = credentials.NewSharedCredentials("", profile)
	}
	return cred
}

func formatInstance(inst *ec2.Instance) string {
	output := []string{
		*inst.InstanceId,
		*inst.InstanceType,
		publicIpAddress(inst),
		*inst.PrivateIpAddress,
		*inst.State.Name,
		(*inst.LaunchTime).Format("2006-01-02 15:04:05"),
	}
	tags := lookupTags(inst, viper.GetString("ec2.ls.output-tags"))
	output = append(output, tags...)
	return strings.Join(output[:], "\t")
}

func publicIpAddress(inst *ec2.Instance) string {
	ip := "___.___.___.___"
	if *inst.State.Name == "running" {
		ip = *inst.PublicIpAddress
	}
	return ip
}

func lookupTags(inst *ec2.Instance, keys string) []string {
	tags := strings.Split(keys, ",")
	var values []string

	for _, tag := range tags {
		values = append(values, lookupTag(inst, tag))
	}
	return values
}

func lookupTag(inst *ec2.Instance, key string) string {
	var value string
	for _, t := range inst.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
