# MyAWS

A human friendly AWS CLI written in Go.

The official aws-cli is useful but too generic. It has many arguments and options and generates huge JSON outputs. But, in most cases, my interesting resources are the same. By setting my favorite default values, MyAWS provides a simple command line interface.

Note that MyAWS is under development and its interface is unstable.

# Installation

If you are Mac OSX user:

```bash
$ brew install minamijoyo/myaws/myaws
```

or

If you have Go 1.11+ development environment:

```bash
$ git clone https://github.com/minamijoyo/myaws
$ cd myaws
$ export GO111MODULE=on
$ go install
```

or

Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/minamijoyo/myaws/releases

## OSX (64bit)

```bash
$ curl -fsSL https://github.com/minamijoyo/myaws/releases/download/v0.3.7/myaws_v0.3.7_darwin_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws
```

## Linux (64bit)

```bash
$ curl -fsSL https://github.com/minamijoyo/myaws/releases/download/v0.3.7/myaws_v0.3.7_linux_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws
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

# Example

```bash
$ myaws ec2 ls
i-0f48fxxxxxxxxxxxx     t2.micro        52.197.xxx.xxx  10.193.xxx.xxx    running 1 minute ago    proxy
i-0e267xxxxxxxxxxxx     t2.medium       52.198.xxx.xxx  10.193.xxx.xxx    running 2 days ago      app
i-0fdaaxxxxxxxxxxxx     t2.large        52.197.xxx.xxx  10.193.xxx.xxx    running 1 month ago     batch
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
  ec2ri       Manage EC2 Reserved Instance resources
  ecr         Manage ECR resources
  ecs         Manage ECS resources
  elb         Manage ELB resources
  help        Help about any command
  iam         Manage IAM resources
  rds         Manage RDS resources
  ssm         Manage SSM resources
  sts         Manage STS resources
  version     Print version

Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --debug             Enable debug mode
  -h, --help              help for myaws
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")

Use "myaws [command] --help" for more information about a command.
```

```bash
$ myaws ec2 ls --help
List EC2 instances

Usage:
  myaws ec2 ls [flags]

Flags:
  -a, --all                 List all instances (by default, list running instances only)
  -F, --fields string       Output fields list separated by space (default "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name")
  -t, --filter-tag string   Filter instances by tag, such as "Name:app-production". The value of tag is assumed to be a partial match
  -q, --quiet               Only display InstanceIDs

Global Flags:
      --config string     config file (default $HOME/.myaws.yml)
      --humanize          Use Human friendly format for time (default true)
      --profile string    AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)
      --region string     AWS region (default none and used AWS_DEFAULT_REGION environment variable.
      --timezone string   Time zone, such as UTC, Asia/Tokyo (default "Local")
```

# LICENCE

MIT

