package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

	resp, err := client.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			fmt.Println(formatInstance(inst, viper.GetString("ec2.ls.output-tags")))
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

func formatInstance(inst *ec2.Instance, outputTags string) string {
	output := []string{
		*inst.InstanceId,
		*inst.InstanceType,
		publicIpAddress(inst),
		*inst.PrivateIpAddress,
		*inst.State.Name,
		(*inst.LaunchTime).Format("2006-01-02 15:04:05"),
	}
	tags := lookupTags(inst, outputTags)
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
