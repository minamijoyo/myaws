package ec2

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/minamijoyo/myaws/myaws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func Ssh(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		myaws.UsageError(cmd, "INSTANCE_ID is required.")
	}
	hostname := resolveIpAddress(args[0])

	loginName := viper.GetString("ec2.ssh.login-name")
	identityFile := strings.Replace(viper.GetString("ec2.ssh.identity-file"), "~", os.Getenv("HOME"), 1)
	fmt.Println(identityFile)
	key, err := ioutil.ReadFile(identityFile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: loginName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}
	session.Wait()
}

func resolveIpAddress(instanceId string) string {
	client := newEC2Client()

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceId},
	}

	response, err := client.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	return *response.Reservations[0].Instances[0].PublicIpAddress
}
