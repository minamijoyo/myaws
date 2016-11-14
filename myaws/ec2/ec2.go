package ec2

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/minamijoyo/myaws/myaws"
)

func newEC2Client() *ec2.EC2 {
	return ec2.New(
		session.New(),
		myaws.NewConfig(),
	)
}
