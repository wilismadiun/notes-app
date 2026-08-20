package http

import (
	"net/http"
	"notes-app/src/modul/note/Applications/usecase"
	"notes-app/src/modul/note/Domains/entities"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	CreateHandler       *usecase.CreateNote
	GetAllNoteshandler  *usecase.GetAllNotes
	GetNoteByIdhandler  *usecase.GetNoteById
	EditNoteByIdHandler *usecase.EditNoteById
	Deletehandler       *usecase.DeleteNote
}

func (h *NoteHandler) AddNote(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "authorization header tidak ditemukan",
		})
	}

	var note entities.Note

	err := c.ShouldBindBodyWithJSON(&note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menambahkan note",
			"error":   err.Error(),
		})
		return
	}

	note.Owner = userId.(string)

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

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "authorization header tidak ditemukan",
		})
	}

	notes, err := h.GetAllNoteshandler.Execute(userId.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "tidak ada data yang bisa ditampilkan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data berhasil ditampilan",
		"data":    notes,
	})
}

func (h *NoteHandler) GetNoteById(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "authorization header tidak ditemukan",
		})
	}
	noteId := c.Param("id")

	note, err := h.GetNoteByIdhandler.Execute(noteId, userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menampilkan note",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "note berhasil ditampilkan",
		"data":    note,
	})
}

func (h *NoteHandler) EditNoteById(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user_id tidak ditemukan",
		})
	}

	noteId := c.Param("id")

	var updateNote entities.EditNoteRequest

	err := c.ShouldBindBodyWithJSON(&updateNote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal memperbarui note",
			"error":   err.Error(),
		})
		return
	}

	exisistId, err := h.EditNoteByIdHandler.Execute(noteId, userId.(string), updateNote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal memperbarui note",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Note berhasil diperbarui",
		"data":    exisistId,
	})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user_id tidak ditemukan",
		})
	}

	noteid := c.Param("id")

	existId, err := h.Deletehandler.Execute(noteid, userId.(string))
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
