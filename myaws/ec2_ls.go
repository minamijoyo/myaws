package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
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
		fmt.Fprintln(client.stdout, formatEC2Instance(client, options, instance))
	}
	return nil
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
