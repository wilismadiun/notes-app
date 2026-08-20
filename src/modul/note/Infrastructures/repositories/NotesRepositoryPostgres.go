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

func (h *NoteRepository) GetAllNote(userId string) []entities.Note {
	var notes []entities.Note

	h.DB.Where("owner = ?", userId).Find(&notes)

	return notes
}

func (h *NoteRepository) FindNoteById(note *entities.Note) error {
	err := h.DB.Where("id = ? AND owner = ?", note.ID, note.Owner).First(note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("id tidak ditemukan")
		}
		return err
	}

	return nil
}

func (h *NoteRepository) EditNoteById(note entities.Note, update map[string]any) error {
	result := h.DB.Model(&note).Where("id = ? AND owner = ?", note.ID, note.Owner).Updates(update)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("id tidak ditemukan")
	}

	return nil
}

func (h *NoteRepository) DeleteNoteById(note entities.Note) error {
	result := h.DB.Where("id = ? AND owner = ?", note.ID, note.Owner).Delete(&note)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("id tidak ditemukan")
	}

	return nil
}
