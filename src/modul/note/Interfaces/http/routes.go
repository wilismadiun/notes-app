package http

import "github.com/gin-gonic/gin"

func NoteRouter(router *gin.Engine, h *NoteHandler) {
	router.POST("/note", h.AddNote)
	router.DELETE("/note/:id", h.DeleteNote)
}
