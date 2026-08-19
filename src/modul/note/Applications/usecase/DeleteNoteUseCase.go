package usecase

import (
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
)

type DeleteNote struct {
	Repo domains.NotesRepository
}

func (h *DeleteNote) Execute(noteId, userId string) (string, error) {
	note := entities.Note{
		ID:    noteId,
		Owner: userId,
	}

	err := h.Repo.DeleteNoteById(note)
	if err != nil {
		return "", err
	}

	return noteId, nil
}
