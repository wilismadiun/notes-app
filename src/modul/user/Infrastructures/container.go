package infrastructures

import (
	usecase "notes-app/src/modul/user/Applications/Usecase"
	"notes-app/src/modul/user/Infrastructures/repositories"
	bcryptsecurity "notes-app/src/modul/user/Infrastructures/security"
	http "notes-app/src/modul/user/Interfaces/http"

	"gorm.io/gorm"
)

func UserContainer(db *gorm.DB) *http.UserHandler {
	repoHandler := repositories.UserRepository{DB: db}
	hashHandler := bcryptsecurity.HashPasswordBcrypt{}

	CreateUser := usecase.CreateUserUseCase{
		Repo:         &repoHandler,
		HashPassword: &hashHandler,
	}

	return &http.UserHandler{
		CreateHandler: &CreateUser,
	}
}
