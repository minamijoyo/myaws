package ecr

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/minamijoyo/myaws/myaws"
)

func newECRClient() *ecr.ECR {
	return ecr.New(
		session.New(),
		myaws.NewConfig(),
	)
}
