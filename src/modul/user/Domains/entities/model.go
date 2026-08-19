package entities

import "errors"

var (
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

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
