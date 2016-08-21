package cmd

import (
	"fmt"

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
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {

			if *inst.State.Name != "running" {
				continue
			}

			var tag_name string
			for _, t := range inst.Tags {
				if *t.Key == "Name" {
					tag_name = *t.Value
					break
				}
			}

			fmt.Println(
				*inst.PublicIpAddress,
				*inst.InstanceId,
				*inst.State.Name,
				*inst.LaunchTime,
				tag_name,
			)
		}
	}
}
