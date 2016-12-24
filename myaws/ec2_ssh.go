package myaws

import (
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
	FilterTag    string
	LoginName    string
	IdentityFile string
	Private      bool
	Command      string
}

// EC2SSH resolves IP address of EC2 instance and connects to it by SSH.
func (client *Client) EC2SSH(options EC2SSHOptions) error {
	instances, err := client.FindEC2Instances(options.FilterTag, false)
	if err != nil {
		return err
	}
	if len(instances) == 0 {
		return errors.Errorf("no such instance: %s", options.FilterTag)
	}

	if len(instances) >= 2 {
		return errors.New("multiple instances found")
	}

	instance := instances[0]
	hostname, err := client.resolveEC2IPAddress(instance, options.Private)
	if err != nil {
		return errors.Wrap(err, "unable to resolve IP address:")
	}

	config, err := buildSSHConfig(options.LoginName, options.IdentityFile)
	if err != nil {
		return err
	}

	if options.Command == "" {
		return startSSHSessionWithTerminal(hostname, config)
	}

	return executeSSHCommand(hostname, config, options.Command)
}

func (client *Client) resolveEC2IPAddress(instance *ec2.Instance, private bool) (string, error) {
	if private {
		return client.resolveEC2PrivateIPAddress(instance)
	}
	return client.resolveEC2PublicIPAddress(instance)
}

func (client *Client) resolveEC2PrivateIPAddress(instance *ec2.Instance) (string, error) {
	if instance.PrivateIpAddress == nil {
		return "", errors.Errorf("no private ip address: %s", instance.InstanceId)
	}
	return *instance.PrivateIpAddress, nil
}

func (client *Client) resolveEC2PublicIPAddress(instance *ec2.Instance) (string, error) {
	if instance.PublicIpAddress == nil {
		return "", errors.Errorf("no public ip address: %s", instance.InstanceId)
	}
	return *instance.PublicIpAddress, nil
}

func buildSSHConfig(loginName string, identityFile string) (*ssh.ClientConfig, error) {
	normalizedIdentityFile := strings.Replace(identityFile, "~", os.Getenv("HOME"), 1)
	key, err := ioutil.ReadFile(normalizedIdentityFile)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read private key:")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse private key:")
	}

	config := &ssh.ClientConfig{
		User: loginName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return config, nil
}

func buildSSHSessionPipe(session *ssh.Session) error {
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

	return nil
}

func startSSHSessionWithTerminal(hostname string, config *ssh.ClientConfig) error {
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

	if err := buildSSHSessionPipe(session); err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return errors.Wrap(err, "failed to start shell:")
	}
	session.Wait()

	return nil
}

func executeSSHCommand(hostname string, config *ssh.ClientConfig, command string) error {
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

	if err := buildSSHSessionPipe(session); err != nil {
		return err
	}

	session.Wait()
	if err := session.Run(command); err != nil {
		return errors.Wrapf(err, "failed to execute command: %s", command)
	}

	return nil
}
