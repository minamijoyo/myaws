# MyAWS

A human friendly AWS CLI written in Go.

The official aws-cli is useful but too generic. It has many arguments and options and generates huge JSON outputs. But, in most cases, my interesting resources are the same. By setting my favorite default values, MyAWS provides a simple command line interface.

Note that MyAWS is under development and its interface is unstable.

# Installation

If you are Mac OSX user:

```bash
$ brew tap minamijoyo/myaws
$ brew install myaws
```

or

If you have Go development environment:

```bash
$ go get github.com/minamijoyo/myaws
```

or

Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/minamijoyo/myaws/releases

## OSX (64bit)

```bash
$ curl -fsSL https://github.com/minamijoyo/myaws/releases/download/v0.1.0/myaws_v0.1.0_darwin_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws
$ myaws --help
```

## Linux (64bit)

```bash
$ curl -fsSL https://github.com/minamijoyo/myaws/releases/download/v0.1.0/myaws_v0.1.0_linux_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws
$ myaws --help
```

# Configuration
## Required
MyAWS invokes AWS API call via aws-sdk-go.
Export environment variables for your AWS credentials:

```bash
$ export AWS_ACCESS_KEY_ID=XXXXXX
$ export AWS_SECRET_ACCESS_KEY=XXXXXX
$ export AWS_DEFAULT_REGION=XXXXXX
```

or set your credentials in `$HOME/.aws/credentials` :

```
[default]
aws_access_key_id = XXXXXX
aws_secret_access_key = XXXXXX
```

or IAM Task Role (ECS) or IAM Role are also available.

AWS credentials are checked in the order of
profile, environment variables, IAM Task Role (ECS), IAM Role.
Unlike the aws default, load profile before environment variables
because we want to prioritize explicit arguments over the environment.

AWS region can be set in Environment variable ( `AWS_DEFAULT_REGION` ), configuration file ( `$HOME/.myaws.yaml` ) , or command argument ( `--region` ).

## Optional

Configuration file is optional.

MyAWS read default configuration from `$HOME/.myaws.yml`

A sample configuration looks like the following:

```yaml
profile: default
region: ap-northeast-1
ec2:
  ls:
    all: false
    fields:
      - InstanceId
      - InstanceType
      - PublicIpAddress
      - PrivateIpAddress
      - StateName
      - LaunchTime
      - Tag:Name
      - Tag:attached_asg
```

# Usage

```bash
$ myaws --help
A human friendly AWS CLI written in Go.

Usage:
  myaws [command]

Available Commands:
  autoscaling Manage autoscaling resources
  ec2         Manage EC2 resources
  ecr         Manage ECR resources
  elb         Manage ELB resources
  iam         Manage IAM resources
  rds         Manage RDS resources
  ssm         Manage SSM resources
  sts         Manage STS resources
  version     Print version

Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws [command] --help" for more information about a command.
```

```bash
$ myaws ec2 --help
Manage EC2 resources

Usage:
  myaws ec2 [flags]
  myaws ec2 [command]

Available Commands:
  ls          List EC2 instances
  ssh         SSH to EC2 instances
  start       Start EC2 instances
  stop        Stop EC2 instances

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws ec2 [command] --help" for more information about a command.
```

```bash
$ myaws ec2 ls --help
List EC2 instances

Usage:
  myaws ec2 ls [flags]

Flags:
  -a, --all                 List all instances (by default, list running instances only)
  -F, --fields string       Output fields list separated by space (default "InstanceId InstanceType PublicIpAddress PrivateIpAddress StateName LaunchTime Tag:Name")
  -t, --filter-tag string   Filter instances by tag, such as "Name:app-production". The value of tag is assumed to be a partial match (default "Name:")
  -q, --quiet               Only display InstanceIDs

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")
```

```bash
$ myaws ec2 ssh --help
SSH to EC2 instances

Usage:
  myaws ec2 ssh [USER@]INSTANCE_NAME [COMMAND...] [flags]

Flags:
  -i, --identity-file string   SSH private key file (default "~/.ssh/id_rsa")
  -l, --login-name string      Login username
      --private                Use private IP to connect

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --debug             Enable debug mode
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")
```

```bash
$ myaws autoscaling --help
Manage autoscaling resources

Usage:
  myaws autoscaling [flags]
  myaws autoscaling [command]

Available Commands:
  attach      Attach instances/loadbalancers to autoscaling group
  detach      Detach instances/loadbalancers from autoscaling group
  ls          List autoscaling groups
  update      Update autoscaling group

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws autoscaling [command] --help" for more information about a command.
```

```bash
$ myaws elb --help
Manage ELB resources

Usage:
  myaws elb [flags]
  myaws elb [command]

Available Commands:
  ls          List ELB instances
  ps          Show ELB instances

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws elb [command] --help" for more information about a command.
```

```bash
$ myaws ecr --help
Manage ECR resources

Usage:
  myaws ecr [flags]
  myaws ecr [command]

Available Commands:
  get-login   Get docker login command for ECR

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws ecr [command] --help" for more information about a command.
```

```bash
$ myaws iam --help
Manage IAM resources

Usage:
  myaws iam [flags]
  myaws iam [command]

Available Commands:
  user        Manage IAM user resources

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --debug             Enable debug mode
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws iam [command] --help" for more information about a command.
```

```bash
$ myaws iam user --help
Manage IAM user resources

Usage:
  myaws iam user [flags]
  myaws iam user [command]

Available Commands:
  ls             List IAM users
  reset-password Reset login password for IAM user

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --debug             Enable debug mode
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws iam user [command] --help" for more information about a command.
```

```bash
$ myaws rds --help
Manage RDS resources

Usage:
  myaws rds [flags]
  myaws rds [command]

Available Commands:
  ls          List RDS instances

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws rds [command] --help" for more information about a command.
```

```bash
$ myaws ssm --help
Manage SSM resources

Usage:
  myaws ssm [flags]
  myaws ssm [command]

Available Commands:
  parameter   Manage SSM parameter resources

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --debug             Enable debug mode
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws ssm [command] --help" for more information about a command.
```

```bash
$ myaws ssm parameter --help
Manage SSM parameter resources

Usage:
  myaws ssm parameter [flags]
  myaws ssm parameter [command]

Available Commands:
  get         Get SSM parameter
  put         Put SSM parameter

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --debug             Enable debug mode
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws ssm parameter [command] --help" for more information about a command.
```

# Example

```bash
$ myaws ec2 ls
i-0f48fxxxxxxxxxxxx     t2.micro        52.197.xxx.xxx  10.193.xxx.xxx    running 1 minute ago    proxy
i-0e267xxxxxxxxxxxx     t2.medium       52.198.xxx.xxx  10.193.xxx.xxx    running 2 days ago      app
i-0fdaaxxxxxxxxxxxx     t2.large        52.197.xxx.xxx  10.193.xxx.xxx    running 1 month ago     batch
```

# LICENCE

MIT

