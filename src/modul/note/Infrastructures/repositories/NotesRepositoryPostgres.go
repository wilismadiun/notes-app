package repositories

import (
	"log"
	"notes-app/src/modul/note/Domains/entities"

	"gorm.io/gorm"
)

type NoteRepository struct {
	db *gorm.DB
}

func (h *NoteRepository) CreateNote(note entities.Note) error {
	err := h.db.Create(&note).Error
	if err != nil {
		log.Println("=================================================================")
		log.Println("ii adalah error")
		log.Println(err)
		return err
	}

	log.Println("=================================================================")
	log.Println("ii adalah error")
	log.Println(err)

	return nil
}
