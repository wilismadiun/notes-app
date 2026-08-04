package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordBcrypt_HashingPassword_Success(t *testing.T) {
	hasher := &HashPasswordBcrypt{}

	password := "password123"

	hash, err := hasher.HashingPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	assert.NoError(t, err)
}
