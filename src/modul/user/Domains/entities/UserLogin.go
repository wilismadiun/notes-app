package entities

func VerifyUserLogin(user User) error {
	if user.Username == "" {
		return ErrUsernameRequired
	}

	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}
