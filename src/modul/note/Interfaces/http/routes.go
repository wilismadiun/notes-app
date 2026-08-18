package http

import "github.com/gin-gonic/gin"

func NoteRouter(router *gin.Engine, h *NoteHandler) {
	router.POST("/note", h.AddNote)
	router.GET("/note/:id", h.GetNoteById)
	router.PATCH("/note/:id", h.EditNoteById)
	router.DELETE("/note/:id", h.DeleteNote)
}
