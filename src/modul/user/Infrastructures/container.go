package infrastructures

import (
	usecase "notes-app/src/modul/user/Applications/Usecase"
	"notes-app/src/modul/user/Infrastructures/repositories"
	securityInfra "notes-app/src/modul/user/Infrastructures/security"
	http "notes-app/src/modul/user/Interfaces/http"

	"gorm.io/gorm"
)

func UserContainer(db *gorm.DB) *http.UserHandler {
	userRepoHandler := repositories.UserRepository{DB: db}
	hashHandler := securityInfra.HashPasswordBcrypt{}
	authTokenHandler := securityInfra.AuthenticationTokenJWT{}

	createUser := usecase.CreateUserUseCase{
		Repo:         &userRepoHandler,
		HashPassword: &hashHandler,
	}

	loginUser := usecase.UserLogin{
		UserRepo:  &userRepoHandler,
		Hash:      &hashHandler,
		AuthToken: &authTokenHandler,
	}

	return &http.UserHandler{
		CreateHandler:    &createUser,
		LoginUserHandler: &loginUser,
	}
}
