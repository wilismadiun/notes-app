package repositories

import (
	"notes-app/src/modul/note/Domains/entities"

	"gorm.io/gorm"
)

type NoteRepository struct {
	DB *gorm.DB
}

func (h *NoteRepository) CreateNote(note entities.Note) error {
	err := h.DB.Create(&note).Error
	if err != nil {
		return err
	}

	return nil
}
