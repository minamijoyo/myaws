# myaws

myaws is a simple command line tool for managing my aws resources.

The aws-cli is useful but too generic.
It has many arguments and options and generates huge JSON output.

But, in most cases, my interesting resources is the same.

By setting my favorite default values,
myaws provides a simple command line interface for managing my aws resources.

# Installation

```
$ go get github.com/minamijoyo/myaws
```

or

Download latest release, unzip and chmod

https://github.com/minamijoyo/myaws/releases

```
$ unzip myaws_v0.0.1_darwin_amd64.zip
$ chmod +x myaws
$ ./myaws
```

# Configuration

Configuration is not required.

myaws read default configuration from `$HOME/.myaws.yaml`

A sample configuration looks like the following:

```
region: ap-northeast-1
ec2:
  ls:
    all: false
```

# Example

```
$ myaws ec2 ls
i-0f48fxxxxxxxxxxxx     t2.micro        52.197.xxx.xxx  10.193.xxx.xxx    running 2016-07-20 02:38:05     proxy
i-0e267xxxxxxxxxxxx     t2.medium       52.198.xxx.xxx  10.193.xxx.xxx    running 2016-08-26 10:57:00     app
i-0fdaaxxxxxxxxxxxx     t2.large        52.197.xxx.xxx  10.193.xxx.xxx    running 2016-08-23 01:42:59     batch
```
