package myaws

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

type IAMResetPasswordOptions struct {
	UserName string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomPassword(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// IAMResetPassword reset password for IAM user
func (client *Client) IAMResetPassword(options IAMResetPasswordOptions) error {
	password := generateRandomPassword(16)
	change_required := true
	params := &iam.UpdateLoginProfileInput{
		UserName:              &options.UserName,
		Password:              &password,
		PasswordResetRequired: &change_required,
	}

	_, err := client.IAM.UpdateLoginProfile(params)
	if err != nil {
		return errors.Wrap(err, "UpdateLoginProfile failed:")
	}

	fmt.Fprintf(client.stdout, "NewPassword: %s\n", password)

	return nil
}
