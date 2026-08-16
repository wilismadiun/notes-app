package http

import (
	"net/http"
	"notes-app/src/modul/note/Applications/usecase"
	"notes-app/src/modul/note/Domains/entities"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	CreateHandler *usecase.CreateNote
	Deletehandler *usecase.DeleteNote
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

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	existId, err := h.Deletehandler.Execute(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menghapus note",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil menghapus note",
		"data":    existId,
	})
}
