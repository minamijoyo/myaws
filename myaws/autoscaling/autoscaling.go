package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/minamijoyo/myaws/myaws"
)

func newAutoScalingClient() *autoscaling.AutoScaling {
	return autoscaling.New(
		session.New(),
		myaws.NewConfig(),
	)
}
