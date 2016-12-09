package myaws

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// EC2SSHOptions customize the behavior of the SSH command.
type EC2SSHOptions struct {
	InstanceID   string
	LoginName    string
	IdentityFile string
}

// EC2SSH resolves IP address of EC2 instance and connects to it by SSH.
func (client *Client) EC2SSH(options EC2SSHOptions) error {
	hostname, err := client.resolveEC2IPAddress(options.InstanceID)
	if err != nil {
		return errors.Wrap(err, "unable to resolve IP address:")
	}

	identityFile := strings.Replace(options.IdentityFile, "~", os.Getenv("HOME"), 1)
	key, err := ioutil.ReadFile(identityFile)
	if err != nil {
		return errors.Wrap(err, "unable to read private key:")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return errors.Wrap(err, "unable to parse private key:")
	}

	config := &ssh.ClientConfig{
		User: options.LoginName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	connection, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		return errors.Wrap(err, "unable to connect:")
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		return errors.Wrap(err, "unable to new session failed:")
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return errors.Wrap(err, "unable to put terminal in Raw Mode:")
	}
	defer terminal.Restore(fd, oldState)

	width, height, _ := terminal.GetSize(fd)

	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		return errors.Wrap(err, "request for pseudo terminal failed:")
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "unable to setup stdin for session:")
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "unable to setup stdout for session:")
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "unable to setup stderr for session:")
	}
	go io.Copy(os.Stderr, stderr)

	if err := session.Shell(); err != nil {
		return errors.Wrap(err, "failed to start shell:")
	}
	session.Wait()

	return nil
}

func (client *Client) resolveEC2IPAddress(instanceID string) (string, error) {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	response, err := client.EC2.DescribeInstances(params)
	if err != nil {
		return "", errors.Wrap(err, "DescribeInstances failed:")
	}

	instance := *response.Reservations[0].Instances[0]
	if instance.PublicIpAddress == nil {
		return "", fmt.Errorf("no public ip address: %s", instanceID)
	}

	return *instance.PublicIpAddress, nil
}
