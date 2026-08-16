package repositories

import (
	"errors"
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

func (h *NoteRepository) DeleteNoteById(id string) error {
	var note entities.Note

	result := h.DB.
		Where("id = ?", id).
		Delete(&note)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("id tidak ditemukan")
	}

	return nil
}
