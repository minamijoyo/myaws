package myaws

import (
	"fmt"
	"strings"

	"encoding/json"
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

	switch client.format {
	case "json":
		fmt.Fprintln(client.stdout, formatJSONEC2Instances(client, options, instances))
	default:
		for _, instance := range instances {
			fmt.Fprintln(client.stdout, formatTsvEC2Instance(client, options, instance))
		}
	}

	return nil
}

func ec2InstanceValues(client *Client, options EC2LsOptions, instance *ec2.Instance) map[string]string {
	formatFuncs := map[string]func(client *Client, options EC2LsOptions, instance *ec2.Instance) string{
		"InstanceId":       formatEC2InstanceID,
		"InstanceType":     formatEC2InstanceType,
		"PublicIpAddress":  formatEC2PublicIPAddress,
		"PrivateIpAddress": formatEC2PrivateIPAddress,
		"AvailabilityZone": formatEC2AvailabilityZone,
		"StateName":        formatEC2StateName,
		"LaunchTime":       formatEC2LaunchTime,
	}

	output := map[string]string{}

	for _, field := range options.Fields {
		value := ""
		if strings.Index(field, "Tag:") != -1 {
			key := strings.Split(field, ":")[1]
			value = formatEC2Tag(instance, key)
		} else {
			value = formatFuncs[field](client, options, instance)
		}
		output[field] = value
	}

	return output
}

func formatTsvEC2Instance(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	values := ec2InstanceValues(client, options, instance)

	var outputFields []string
	if options.Quiet {
		outputFields = []string{"InstanceId"}
	} else {
		outputFields = options.Fields
	}

	output := []string{}
	for _, field := range outputFields {
		output = append(output, values[field])
	}
	return strings.Join(output[:], "\t")
}

func formatJSONEC2Instances(client *Client, options EC2LsOptions, instances []*ec2.Instance) string {
	outputs := []map[string]string{}
	for _, instance := range instances {
		outputs = append(outputs, ec2InstanceValues(client, options, instance))
	}

	bytes, err := json.Marshal(outputs)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", bytes)
}

func formatEC2InstanceID(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	return *instance.InstanceId
}

func formatEC2InstanceType(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	if client.format == "json" {
		return *instance.InstanceType
	}
	return fmt.Sprintf("%-11s", *instance.InstanceType)
}

func formatEC2PublicIPAddress(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	if instance.PublicIpAddress == nil {
		if client.format == "json" {
			return ""
		}
		return "___.___.___.___"
	}
	return *instance.PublicIpAddress
}

func formatEC2PrivateIPAddress(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	if instance.PrivateIpAddress == nil {
		if client.format == "json" {
			return ""
		}
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

func formatEC2AvailabilityZone(client *Client, options EC2LsOptions, instance *ec2.Instance) string {
	return *instance.Placement.AvailabilityZone
}
