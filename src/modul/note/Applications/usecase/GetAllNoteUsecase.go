package usecase

import (
	"errors"
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
)

type GetAllNote struct {
	repo domains.NotesRepository
}

func (h *GetAllNote) Execute(userId string) ([]entities.Note, error) {
	notes := h.repo.GetAllNote(userId)

	if len(notes) == 0 {
		return []entities.Note{}, errors.New("Tidak ada data yg bisa ditampilkan")
	}

	return notes, nil
}
