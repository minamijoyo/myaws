package myaws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Client represents myaws CLI
type Client struct {
	profile  string
	region   string
	timezone string
	humanize bool
	EC2      *ec2.EC2
}

// NewClient initializes Client instance
func NewClient(profile string, region string, timezone string, humanize bool) (*Client, error) {
	session := session.New()
	config := NewConfig()
	client := &Client{
		profile:  profile,
		region:   region,
		timezone: timezone,
		humanize: humanize,
		EC2:      ec2.New(session, config),
	}
	return client, nil
}
