package entities

func VerifyUserLogin(user UserLogin) error {
	if user.Username == "" {
		return ErrUsernameRequired
	}

	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}
