package elb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/minamijoyo/myaws/myaws"
)

func newELBClient() *elb.ELB {
	return elb.New(
		session.New(),
		myaws.NewConfig(),
	)
}
