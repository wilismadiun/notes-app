package entities

import (
	"regexp"
	"strings"
)

func VerifyRegisterUser(user User) error {
	if strings.TrimSpace(user.Username) == "" {
		return ErrUsernameRequired
	}

	if strings.TrimSpace(user.Password) == "" {
		return ErrPasswordRequired
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_ ]+$`)
	if !usernameRegex.MatchString(user.Username) {
		return ErrUsernameInvalid
	}

	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
