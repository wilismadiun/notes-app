package http

import "github.com/gin-gonic/gin"

func UserRouter(router *gin.Engine, h *UserHandler) {
	router.POST("/register", h.AddUser)
}
