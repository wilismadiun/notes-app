package entities

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmailRequired    = errors.New("EMAIL_REQUIRED")
	ErrUsernameRequired = errors.New("USERNAME_REQUIRED")
	ErrPasswordRequired = errors.New("PASSWORD_REQUIRED")
	ErrUsernameInvalid  = errors.New("USERNAME_INVALID")
	ErrPasswordTooShort = errors.New("PASSWORD_TOO_SHORT")
	ErrEmailInvalid     = errors.New("EMAIL_INVALID")
)

type User struct {
	ID       string
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func VerifyRegisterUser(user User) error {
	if strings.TrimSpace(user.Email) == "" {
		return ErrEmailRequired
	}

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

	// Password minimal 8 karakter
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}

	// Validasi format email sederhana
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}

	return nil
}
