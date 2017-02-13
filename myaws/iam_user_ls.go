package myaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

// IAMUserLs describes IAM users
func (client *Client) IAMUserLs() error {
	response, err := client.IAM.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return errors.Wrap(err, "ListUsers failed:")
	}

	// TODO: support pagenation
	for _, user := range response.Users {
		fmt.Fprintln(client.stdout, formatIAMUser(client, user))
	}
	return nil
}

func formatIAMUser(client *Client, user *iam.User) string {
	return fmt.Sprintf("%s\t%s\t%s", *user.UserName, client.FormatTime(user.CreateDate), client.FormatTime(user.PasswordLastUsed))
}
