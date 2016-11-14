package rds

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/minamijoyo/myaws/myaws"
)

func newRDSClient() *rds.RDS {
	return rds.New(
		session.New(),
		myaws.NewConfig(),
	)
}
