package http

import (
	"net/http"
	usecase "notes-app/src/modul/user/Applications/Usecase"
	"notes-app/src/modul/user/Domains/entities"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	CreateHandler *usecase.CreateUserUseCase
}

func (h *UserHandler) AddUser(c *gin.Context) {
	var user entities.User

	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userSuccess, err := h.CreateHandler.Execute(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil menambahkan user",
		"data":    userSuccess,
	})
}
