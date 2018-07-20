package myaws

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

// Client represents myaws CLI
type Client struct {
	config      *aws.Config
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	profile     string
	region      string
	timezone    string
	humanize    bool
	AutoScaling *autoscaling.AutoScaling
	EC2         *ec2.EC2
	ECR         *ecr.ECR
	ELB         *elb.ELB
	IAM         *iam.IAM
	RDS         *rds.RDS
	SSM         *ssm.SSM
	STS         *sts.STS
}

// NewClient initializes Client instance
func NewClient(stdin io.Reader, stdout io.Writer, stderr io.Writer, profile string, region string, timezone string, humanize bool) (*Client, error) {
	session := session.New()
	config := newConfig(profile, region)
	client := &Client{
		config:      config,
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		profile:     profile,
		region:      region,
		timezone:    timezone,
		humanize:    humanize,
		AutoScaling: autoscaling.New(session, config),
		EC2:         ec2.New(session, config),
		ECR:         ecr.New(session, config),
		ELB:         elb.New(session, config),
		IAM:         iam.New(session, config),
		RDS:         rds.New(session, config),
		SSM:         ssm.New(session, config),
		STS:         sts.New(session, config),
	}
	return client, nil
}

// Confirmation asks user for confirmation.
// "y" and "Y" returns true and others are false.
func (client *Client) Confirmation(message string) (bool, error) {
	fmt.Fprintf(client.stdout, "%s [y/n]: ", message)

	reader := bufio.NewReader(client.stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, errors.Wrap(err, "ReadString failed:")
	}

	normalized := strings.ToLower(strings.TrimSpace(input))
	return normalized == "y", nil
}
