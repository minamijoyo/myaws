package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type LsFlag struct {
	All bool
}

func Ls(flag *LsFlag) {
	svc := ec2.New(
		session.New(),
		&aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	)

	var stateFilter *ec2.Filter
	if flag.All {
		stateFilter = &ec2.Filter{}
	} else {
		stateFilter = &ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		}
	}

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			stateFilter,
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

func formatInstance(inst *ec2.Instance) string {
	output := []string{
		*inst.InstanceId,
		*inst.InstanceType,
		publicIpAddress(inst),
		*inst.PrivateIpAddress,
		*inst.State.Name,
		(*inst.LaunchTime).Format("2006-01-02 15:04:05"),
		lookupTag(inst, "Name"),
	}
	return strings.Join(output[:], "\t")
}

func publicIpAddress(inst *ec2.Instance) string {
	ip := "___.___.___.___"
	if *inst.State.Name == "running" {
		ip = *inst.PublicIpAddress
	}
	return ip
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
