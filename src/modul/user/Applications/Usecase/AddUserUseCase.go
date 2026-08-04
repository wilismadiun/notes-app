package usecase

import (
	"notes-app/src/modul/user/Applications/security"
	domains "notes-app/src/modul/user/Domains"
	"notes-app/src/modul/user/Domains/entities"
)

type CreateUserUseCase struct {
	Repo         domains.UserRepository
	HashPassword security.HashPassword
}

func (uc *CreateUserUseCase) Execute(user *entities.User) (entities.RegisteredUser, error) {
	if err := entities.VerifyRegisterUser(*user); err != nil {
		return entities.RegisteredUser{}, err
	}

	if err := uc.Repo.VerifyUsername(user); err != nil {
		return entities.RegisteredUser{}, err
	}

	hashedPassword, err := uc.HashPassword.HashingPassword(user.Password)
	if err != nil {
		return entities.RegisteredUser{}, err
	}

	user.Password = hashedPassword

	if err := uc.Repo.Createuser(user); err != nil {
		return entities.RegisteredUser{}, err
	}

	return entities.RegisteredUser{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
