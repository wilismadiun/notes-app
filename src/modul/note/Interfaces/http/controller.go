package http

import (
	"net/http"
	"notes-app/src/modul/note/Applications/usecase"
	"notes-app/src/modul/note/Domains/entities"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	CreateHandler *usecase.CreateNote
}

func (h *NoteHandler) AddNote(c *gin.Context) {
	var note entities.Note

	err := c.ShouldBindBodyWithJSON(&note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menambahkan note",
			"error":   err.Error(),
		})
		return
	}

	noteSuccess, err := h.CreateHandler.Execute(note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menambahkan note",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "note berhasil ditambahkan",
		"data":    noteSuccess,
	})

}
