package ec2

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// SSH resolves IP address of EC2 instance and connects to it by SSH.
func SSH(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "INSTANCE_ID is required.")
	}
	hostname := resolveIPAddress(args[0])

	loginName := viper.GetString("ec2.ssh.login-name")
	identityFile := strings.Replace(viper.GetString("ec2.ssh.identity-file"), "~", os.Getenv("HOME"), 1)
	fmt.Println(identityFile)
	key, err := ioutil.ReadFile(identityFile)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: loginName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
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
		log.Fatal("Unable to put terminal in Raw Mode", err)
	}
	defer terminal.Restore(fd, oldState)

	width, height, _ := terminal.GetSize(fd)

	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		log.Fatalf("Request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal("Unable to setup stdin for session", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal("Unable to setup stdout for session", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatal("Unable to setup stderr for session", err)
	}
	go io.Copy(os.Stderr, stderr)

	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}
	session.Wait()
}

func resolveIPAddress(instanceID string) string {
	client := newEC2Client()

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	response, err := client.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	return *response.Reservations[0].Instances[0].PublicIpAddress
}
