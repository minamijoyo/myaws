package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// ec2lsCmd represents the ec2ls command
var ec2lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List EC2 instances",
	Long:  `List EC2 instances`,
	Run:   ec2ls,
}

func init() {
	ec2Cmd.AddCommand(ec2lsCmd)
}

func ec2ls(cmd *cobra.Command, args []string) {
	svc := ec2.New(
		session.New(),
		&aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			fmt.Println(formatEC2Instance(inst))
		}
	}
}

func formatEC2Instance(inst *ec2.Instance) string {
	output := []string{
		*inst.PublicIpAddress,
		*inst.InstanceId,
		*inst.State.Name,
		(*inst.LaunchTime).Format("2006-01-02 15:04:05"),
		lookupEC2Tag(inst, "Name"),
	}
	return strings.Join(output[:], "\t")
}

func lookupEC2Tag(inst *ec2.Instance, key string) string {
	var value string
	for _, t := range inst.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
