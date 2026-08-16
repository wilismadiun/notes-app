package app

import (
	notecontainer "notes-app/src/modul/note/Infrastructures"
	notehttp "notes-app/src/modul/note/Interfaces/http"
	usercontainer "notes-app/src/modul/user/Infrastructures"
	userhttp "notes-app/src/modul/user/Interfaces/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(router *gin.Engine, db *gorm.DB) {
	// handlers
	userHandler := usercontainer.UserContainer(db)
	notehandler := notecontainer.NoteContainer(db)

	// routes
	userhttp.UserRouter(router, userHandler)
	notehttp.NoteRouter(router, notehandler)
}
