package usecase

import (
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
)

type GetNoteById struct {
	Repo domains.NotesRepository
}

func (h *GetNoteById) Execute(noteId, userId string) (entities.Note, error) {
	note := entities.Note{
		ID:    noteId,
		Owner: userId,
	}

	err := h.Repo.FindNoteById(&note)
	if err != nil {
		return entities.Note{}, err
	}

	return note, nil
}
