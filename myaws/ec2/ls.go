package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Ls(*cobra.Command, []string) {
	client := newEC2Client()
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			buildStateFilter(viper.GetBool("ec2.ls.all")),
			buildTagFilter(viper.GetString("ec2.ls.filter-tag")),
		},
	}

	response, err := client.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(formatInstance(instance, viper.GetString("ec2.ls.output-tags")))
		}
	}
}

func buildStateFilter(all bool) *ec2.Filter {
	var stateFilter *ec2.Filter
	if !all {
		stateFilter = &ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		}
	}
	return stateFilter
}

func buildTagFilter(filterTag string) *ec2.Filter {
	var tagFilter *ec2.Filter
	if filterTag != "" {
		tagParts := strings.Split(filterTag, ":")
		tagFilter = &ec2.Filter{
			Name: aws.String("tag:" + tagParts[0]),
			Values: []*string{
				aws.String("*" + tagParts[1] + "*"),
			},
		}
	}
	return tagFilter
}

func formatInstance(instance *ec2.Instance, outputTags string) string {
	output := []string{
		*instance.InstanceId,
		formatInstanceType(instance),
		formatPublicIpAddress(instance),
		*instance.PrivateIpAddress,
		*instance.State.Name,
		myaws.FormatTime(instance.LaunchTime),
	}
	tags := lookupTags(instance, outputTags)
	output = append(output, tags...)
	return strings.Join(output[:], "\t")
}

func formatInstanceType(instance *ec2.Instance) string {
	return fmt.Sprintf("%-11s", *instance.InstanceType)
}

func formatPublicIpAddress(instance *ec2.Instance) string {
	ip := "___.___.___.___"
	if *instance.State.Name == "running" {
		ip = *instance.PublicIpAddress
	}
	return ip
}

func lookupTags(instance *ec2.Instance, keys string) []string {
	tags := strings.Split(keys, ",")
	var values []string

	for _, tag := range tags {
		values = append(values, lookupTag(instance, tag))
	}
	return values
}

func lookupTag(instance *ec2.Instance, key string) string {
	var value string
	for _, t := range instance.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
