package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	tests := []struct {
		name    string
		user    UserLogin
		wantErr error
	}{
		{
			name: "username required",
			user: UserLogin{
				Username: "",
				Password: "12345678",
			},
			wantErr: ErrUsernameRequired,
		},
		{
			name: "password required",
			user: UserLogin{
				Username: "john123",
				Password: "",
			},
			wantErr: ErrPasswordRequired,
		},
		{
			name: "success",
			user: UserLogin{
				Username: "john123",
				Password: "1234567",
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := VerifyUserLogin(test.user)

			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}
