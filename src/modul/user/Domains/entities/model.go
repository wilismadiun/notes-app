package entities

import "errors"

var (
	ErrEmailRequired    = errors.New("EMAIL_REQUIRED")
	ErrUsernameRequired = errors.New("USERNAME_REQUIRED")
	ErrPasswordRequired = errors.New("PASSWORD_REQUIRED")
	ErrUsernameInvalid  = errors.New("USERNAME_INVALID")
	ErrPasswordTooShort = errors.New("PASSWORD_TOO_SHORT")
)

type User struct {
	ID       string
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisteredUser struct {
	ID       string
	Username string
}
