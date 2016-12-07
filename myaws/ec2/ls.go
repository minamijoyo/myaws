package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// LsOptions customize the behavior of the Ls command.
type LsOptions struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
}

// Ls describes EC2 instances.
func Ls(client *myaws.Client, options LsOptions) error {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			buildStateFilter(options.All),
			buildTagFilter(options.FilterTag),
		},
	}

	response, err := client.EC2.DescribeInstances(params)
	if err != nil {
		return errors.Wrap(err, "DescribeInstances failed")
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(formatInstance(client, options, instance))
		}
	}

	return nil
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

func formatInstance(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	formatFuncs := map[string]func(client *myaws.Client, options LsOptions, instance *ec2.Instance) string{
		"InstanceId":       formatInstanceID,
		"InstanceType":     formatInstanceType,
		"PublicIpAddress":  formatPublicIPAddress,
		"PrivateIpAddress": formatPrivateIPAddress,
		"StateName":        formatStateName,
		"LaunchTime":       formatLaunchTime,
	}

	var outputFields []string
	if options.Quiet {
		outputFields = []string{"InstanceId"}
	} else {
		outputFields = options.Fields
	}

	output := []string{}

	for _, field := range outputFields {
		value := ""
		if strings.Index(field, "Tag:") != -1 {
			key := strings.Split(field, ":")[1]
			value = formatTag(instance, key)
		} else {
			value = formatFuncs[field](client, options, instance)
		}
		output = append(output, value)
	}
	return strings.Join(output[:], "\t")
}

func formatInstanceID(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	return *instance.InstanceId
}

func formatInstanceType(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	return fmt.Sprintf("%-11s", *instance.InstanceType)
}

func formatPublicIPAddress(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	if instance.PublicIpAddress == nil {
		return "___.___.___.___"
	}
	return *instance.PublicIpAddress
}

func formatPrivateIPAddress(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	if instance.PrivateIpAddress == nil {
		return "___.___.___.___"
	}
	return *instance.PrivateIpAddress
}

func formatStateName(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	return *instance.State.Name
}

func formatLaunchTime(client *myaws.Client, options LsOptions, instance *ec2.Instance) string {
	return client.FormatTime(instance.LaunchTime)
}

func formatTag(instance *ec2.Instance, key string) string {
	var value string
	for _, t := range instance.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
