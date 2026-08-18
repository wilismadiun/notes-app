package usecase

import (
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
)

type GetNoteById struct {
	Repo domains.NotesRepository
}

func (h *GetNoteById) Execute(id string) (entities.Note, error) {
	note, err := h.Repo.FindNoteById(id)
	if err != nil {
		return entities.Note{}, err
	}

	return note, nil
}
