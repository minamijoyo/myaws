package myaws

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

// IAMUserResetPasswordOptions customize the behavior of the IAMUserResetPassword command.
type IAMUserResetPasswordOptions struct {
	UserName string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomPassword(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))] // nolint: gosec
	}
	return string(b)
}

// IAMUserResetPassword reset password for IAM user.
func (client *Client) IAMUserResetPassword(options IAMUserResetPasswordOptions) error {
	user, err := client.IAMGetUser(options.UserName)
	if err != nil {
		return err
	}

	fmt.Fprintf(client.stdout, "%v\n", user)

	confirm, err := client.Confirmation("Are you sure want to reset password?")
	if err != nil {
		return err
	}

	if !confirm {
		// cancel reset password.
		fmt.Fprintln(client.stdout, "Cancelled.")
		return nil
	}

	password := generateRandomPassword(16)
	changeRequired := true

	// Check if IAM user has a login profile.
	_, err = client.IAM.GetLoginProfile(&iam.GetLoginProfileInput{
		UserName: aws.String(options.UserName),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchEntity" {
			// if IAM user has no login profile, create a login profile with initial password.
			err = client.IAMUserCreateLoginProfile(options.UserName, password, changeRequired)
			if err != nil {
				return err
			}

			fmt.Fprintf(client.stdout, "InitialPassword: %s\n", password)

			return nil
		}
		// unexpected error
		return err
	}

	// if IAM user has a login profile already, update password.
	err = client.IAMUserUpdatePassword(options.UserName, password, changeRequired)
	if err != nil {
		return err
	}

	fmt.Fprintf(client.stdout, "NewPassword: %s\n", password)

	return nil
}

// IAMGetUser returns IAM user.
func (client *Client) IAMGetUser(username string) (*iam.User, error) {
	params := &iam.GetUserInput{
		UserName: &username,
	}

	response, err := client.IAM.GetUser(params)
	if err != nil {
		return nil, errors.Wrap(err, "GetUser failed:")
	}

	return response.User, nil
}

// IAMUserCreateLoginProfile creates a login profile for IAM User with initial password.
func (client *Client) IAMUserCreateLoginProfile(username string, password string, changeRequired bool) error {
	params := &iam.CreateLoginProfileInput{
		UserName:              &username,
		Password:              &password,
		PasswordResetRequired: &changeRequired,
	}

	_, err := client.IAM.CreateLoginProfile(params)
	if err != nil {
		return errors.Wrap(err, "CreateLoginProfile failed:")
	}

	return nil
}

// IAMUserUpdatePassword updates the password of existing login profile for IAM user.
func (client *Client) IAMUserUpdatePassword(username string, password string, changeRequired bool) error {
	params := &iam.UpdateLoginProfileInput{
		UserName:              &username,
		Password:              &password,
		PasswordResetRequired: &changeRequired,
	}

	_, err := client.IAM.UpdateLoginProfile(params)
	if err != nil {
		return errors.Wrap(err, "UpdateLoginProfile failed:")
	}

	return nil
}
