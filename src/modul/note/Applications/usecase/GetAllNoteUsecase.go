package usecase

import (
	"errors"
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
)

type GetAllNotes struct {
	Repo domains.NotesRepository
}

func (h *GetAllNotes) Execute(userId string) ([]entities.Note, error) {
	notes := h.Repo.GetAllNote(userId)

	if len(notes) == 0 {
		return []entities.Note{}, errors.New("Tidak ada data yg bisa ditampilkan")
	}

	return notes, nil
}
