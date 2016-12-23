package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// EC2LsOptions customize the behavior of the Ls command.
type EC2LsOptions struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
}

// EC2Ls describes EC2 instances.
func (client *Client) EC2Ls(options EC2LsOptions) error {
	instances, err := client.FindEC2Instances(options.FilterTag, options.All)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		fmt.Println(formatEC2Instance(client, options, instance))
	}
	return nil
}

// FindEC2Instances returns an array of instances matching the conditions.
func (client *Client) FindEC2Instances(filterTag string, all bool) ([]*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			buildEC2StateFilter(all),
			buildEC2TagFilter(filterTag),
		},
	}

	response, err := client.EC2.DescribeInstances(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeInstances failed")
	}

	var instances []*ec2.Instance
	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func buildEC2StateFilter(all bool) *ec2.Filter {
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

func buildEC2TagFilter(filterTag string) *ec2.Filter {
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

func formatEC2Instance(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	formatFuncs := map[string]func(client *Client, options EC2LsOptions, instance *ec2.Instance) string{
		"InstanceId":       formatEC2InstanceID,
		"InstanceType":     formatEC2InstanceType,
		"PublicIpAddress":  formatEC2PublicIPAddress,
		"PrivateIpAddress": formatEC2PrivateIPAddress,
		"StateName":        formatEC2StateName,
		"LaunchTime":       formatEC2LaunchTime,
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
			value = formatEC2Tag(instance, key)
		} else {
			value = formatFuncs[field](client, options, instance)
		}
		output = append(output, value)
	}
	return strings.Join(output[:], "\t")
}

func formatEC2InstanceID(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	return *instance.InstanceId
}

func formatEC2InstanceType(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	return fmt.Sprintf("%-11s", *instance.InstanceType)
}

func formatEC2PublicIPAddress(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	if instance.PublicIpAddress == nil {
		return "___.___.___.___"
	}
	return *instance.PublicIpAddress
}

func formatEC2PrivateIPAddress(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	if instance.PrivateIpAddress == nil {
		return "___.___.___.___"
	}
	return *instance.PrivateIpAddress
}

func formatEC2StateName(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	return *instance.State.Name
}

func formatEC2LaunchTime(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	return client.FormatTime(instance.LaunchTime)
}

func formatEC2Tag(instance *ec2.Instance, key string) string {
	var value string
	for _, t := range instance.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
