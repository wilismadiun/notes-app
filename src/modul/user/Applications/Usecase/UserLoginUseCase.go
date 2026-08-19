package usecase

import (
	"errors"
	"notes-app/src/modul/user/Applications/security"
	domains "notes-app/src/modul/user/Domains"
	"notes-app/src/modul/user/Domains/entities"
)

type UserLogin struct {
	// AuthRepo  domains.AuthenticationRepository
	UserRepo  domains.UserRepository
	Hash      security.HashPassword
	AuthToken security.AuthToken
}

func (h *UserLogin) Execute(userLogin entities.UserLogin) (string, error) {
	err := entities.VerifyUserLogin(userLogin)
	if err != nil {
		return "", errors.New("username atau Password tidak ada")
	}

	username := userLogin.Username

	user := entities.User{
		Username: username,
	}

	err = h.UserRepo.FindUserByUsername(&user)
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	err = h.Hash.CompareHashPassword(userLogin.Password, user.Password)
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	authToken, err := h.AuthToken.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("login tidak berhasil")
	}

	return authToken, nil
}
