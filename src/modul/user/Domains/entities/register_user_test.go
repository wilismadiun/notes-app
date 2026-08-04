package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyRegisterUser(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name: "email required",
			user: User{
				Username: "johndoe",
				Email:    "",
				Password: "password123",
			},
			wantErr: ErrEmailRequired,
		},
		{
			name: "username required",
			user: User{
				Username: "",
				Email:    "john@gmail.com",
				Password: "password123",
			},
			wantErr: ErrUsernameRequired,
		},
		{
			name: "password required",
			user: User{
				Username: "johndoe",
				Email:    "john@gmail.com",
				Password: "",
			},
			wantErr: ErrPasswordRequired,
		},
		{
			name: "username invalid",
			user: User{
				Username: "john@doe",
				Email:    "john@gmail.com",
				Password: "password123",
			},
			wantErr: ErrUsernameInvalid,
		},
		{
			name: "password too short",
			user: User{
				Username: "johndoe",
				Email:    "john@gmail.com",
				Password: "1234567",
			},
			wantErr: ErrPasswordTooShort,
		},
		{
			name: "email invalid",
			user: User{
				Username: "johndoe",
				Email:    "john.gmail.com",
				Password: "password123",
			},
			wantErr: ErrEmailInvalid,
		},
		{
			name: "success",
			user: User{
				Username: "johndoe",
				Email:    "john@gmail.com",
				Password: "password123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyRegisterUser(tt.user)

			assert.Equal(t, err, tt.wantErr)
		})
	}
}
