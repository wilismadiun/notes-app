package app

import (
	usercontainer "notes-app/src/modul/user/Infrastructures"
	userhttp "notes-app/src/modul/user/Interfaces/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRouter(router *gin.Engine, db *gorm.DB) {
	// handlers
	userHandler := usercontainer.UserContainer(db)

	// routes
	userhttp.UserRouter(router, userHandler)
}
