package usecase

import (
	"errors"
	"testing"

	"notes-app/src/modul/user/Domains/entities"
	"notes-app/src/modul/user/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockHash := mocks.NewMockHashPassword(ctrl)

	uc := CreateUserUseCase{
		Repo:         mockRepo,
		HashPassword: mockHash,
	}

	user := entities.User{
		Username: "John1234",
		Password: "12345678",
	}

	mockRepo.EXPECT().VerifyUsername(gomock.Any()).Return(nil)

	mockHash.
		EXPECT().
		HashingPassword("12345678").
		Return("hashed-password", nil)

	mockRepo.
		EXPECT().
		Createuser(gomock.Any()).
		DoAndReturn(func(user *entities.User) error {

			assert.Equal(t, "hashed-password", user.Password)
			assert.Equal(t, "John1234", user.Username)

			return nil
		})

	_, err := uc.Execute(&user)

	assert.NoError(t, err)
}

func TestAddUser_VerifyUsernameFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockHash := mocks.NewMockHashPassword(ctrl)

	uc := CreateUserUseCase{
		Repo:         mockRepo,
		HashPassword: mockHash,
	}

	user := entities.User{
		Username: "John1234",
		Password: "12345678",
	}

	mockRepo.
		EXPECT().
		VerifyUsername(gomock.Any()).
		Return(errors.New("username sudah digunakan"))

	_, err := uc.Execute(&user)

	assert.Error(t, err)
	assert.EqualError(t, err, "username sudah digunakan")
}

func TestAddUser_HashPasswordFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockHash := mocks.NewMockHashPassword(ctrl)

	uc := CreateUserUseCase{
		Repo:         mockRepo,
		HashPassword: mockHash,
	}

	user := entities.User{
		Username: "John1234",
		Password: "12345678",
	}

	mockRepo.EXPECT().VerifyUsername(gomock.Any()).Return(nil)

	mockHash.
		EXPECT().
		HashingPassword("12345678").
		Return("", errors.New("hash failed"))

	_, err := uc.Execute(&user)

	assert.Error(t, err)
	assert.EqualError(t, err, "hash failed")
}

func TestAddUser_CreateUserFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockHash := mocks.NewMockHashPassword(ctrl)

	uc := CreateUserUseCase{
		Repo:         mockRepo,
		HashPassword: mockHash,
	}

	user := entities.User{
		Username: "John1234",
		Password: "12345678",
	}

	mockRepo.EXPECT().VerifyUsername(gomock.Any()).Return(nil)

	mockHash.
		EXPECT().
		HashingPassword("12345678").
		Return("hashed-password", nil)

	mockRepo.
		EXPECT().
		Createuser(gomock.Any()).
		Return(errors.New("database error"))

	_, err := uc.Execute(&user)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}
