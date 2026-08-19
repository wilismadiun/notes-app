package app

import (
	"notes-app/src/commons/middleware"
	notecontainer "notes-app/src/modul/note/Infrastructures"
	notehttp "notes-app/src/modul/note/Interfaces/http"
	usercontainer "notes-app/src/modul/user/Infrastructures"
	"notes-app/src/modul/user/Infrastructures/security"
	userhttp "notes-app/src/modul/user/Interfaces/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(router *gin.Engine, db *gorm.DB) {
	tokenService := security.AuthenticationTokenJWT{}

	// handlers
	authMiddleware := middleware.Authentication(&tokenService)
	userHandler := usercontainer.UserContainer(db)
	notehandler := notecontainer.NoteContainer(db)

	// routes
	userhttp.UserRouter(router, userHandler)

	api := router.Group("/api")
	api.Use(authMiddleware)
	notehttp.NoteRouter(api, notehandler)
}
