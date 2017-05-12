package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

// STSID gets caller identity.
func (client *Client) STSID() error {
	response, err := client.STS.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return errors.Wrap(err, "GetCallerIdentity failed:")
	}

	fmt.Fprintln(client.stdout, formatSTSID(response))
	return nil
}

func formatSTSID(id *sts.GetCallerIdentityOutput) string {
	return fmt.Sprintf("Account: %s\nUserId: %s\nArn: %s",
		*id.Account, *id.UserId, *id.Arn)
}
