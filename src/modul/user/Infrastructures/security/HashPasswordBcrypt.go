package security

import "golang.org/x/crypto/bcrypt"

type HashPasswordBcrypt struct{}

func (h *HashPasswordBcrypt) HashingPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hash)

	return hashedPassword, nil
}
