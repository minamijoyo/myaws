package myaws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/elb"
)

// Client represents myaws CLI
type Client struct {
	profile     string
	region      string
	timezone    string
	humanize    bool
	AutoScaling *autoscaling.AutoScaling
	EC2         *ec2.EC2
	ECR         *ecr.ECR
	ELB         *elb.ELB
}

// NewClient initializes Client instance
func NewClient(profile string, region string, timezone string, humanize bool) (*Client, error) {
	session := session.New()
	config := NewConfig()
	client := &Client{
		profile:     profile,
		region:      region,
		timezone:    timezone,
		humanize:    humanize,
		AutoScaling: autoscaling.New(session, config),
		EC2:         ec2.New(session, config),
		ECR:         ecr.New(session, config),
		ELB:         elb.New(session, config),
	}
	return client, nil
}
