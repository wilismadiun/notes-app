package usecase

import (
	"errors"
	"notes-app/src/modul/user/Domains/entities"
	"notes-app/src/modul/user/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockHash := mocks.NewMockHashPassword(ctrl)
	mockAuthToken := mocks.NewMockAuthToken(ctrl)

	lu := UserLogin{
		UserRepo:  mockUserRepo,
		Hash:      mockHash,
		AuthToken: mockAuthToken,
	}

	t.Run("should be error when username or password not found", func(t *testing.T) {
		userLogin := entities.UserLogin{
			Username: "irfanaufa07",
		}

		_, err := lu.Execute(userLogin)

		assert.Error(t, err)
		assert.EqualError(t, err, "username atau Password tidak ada")
	})

	t.Run("should be error when username not found", func(t *testing.T) {
		userLogin := entities.UserLogin{
			Username: "irfanaufa07",
			Password: "12345678",
		}

		mockUserRepo.EXPECT().FindUserByUsername(gomock.Any()).Return(errors.New("username tidak ditemukan"))

		_, err := lu.Execute(userLogin)

		assert.Error(t, err)
		assert.EqualError(t, err, "username atau password salah")
	})

	t.Run("should be error when password and password hash do not match", func(t *testing.T) {
		userLogin := entities.UserLogin{
			Username: "irfanaufa07",
			Password: "12345678",
		}

		mockUserRepo.EXPECT().FindUserByUsername(gomock.Any()).Return(nil)
		mockHash.EXPECT().CompareHashPassword(gomock.Any(), gomock.Any()).Return(errors.New("password tidak cocok"))

		_, err := lu.Execute(userLogin)

		assert.Error(t, err)
		assert.EqualError(t, err, "username atau password salah")
	})

	t.Run("should be error when generate token fail", func(t *testing.T) {
		userLogin := entities.UserLogin{
			Username: "irfanaufa07",
			Password: "12345678",
		}

		mockUserRepo.EXPECT().FindUserByUsername(gomock.Any()).Return(nil)
		mockHash.EXPECT().CompareHashPassword(gomock.Any(), gomock.Any()).Return(nil)
		mockAuthToken.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("generate gagal"))

		_, err := lu.Execute(userLogin)

		assert.Error(t, err)
		assert.EqualError(t, err, "login tidak berhasil")
	})

	t.Run("login success", func(t *testing.T) {
		userLogin := entities.UserLogin{
			Username: "irfanaufa07",
			Password: "12345678",
		}

		mockUserRepo.EXPECT().FindUserByUsername(gomock.Any()).Return(nil)
		mockHash.EXPECT().CompareHashPassword(gomock.Any(), gomock.Any()).Return(nil)
		mockAuthToken.EXPECT().GenerateToken(gomock.Any()).Return("token-123", nil)

		token, err := lu.Execute(userLogin)

		assert.NoError(t, err)
		assert.Equal(t, token, "token-123")
	})
}
