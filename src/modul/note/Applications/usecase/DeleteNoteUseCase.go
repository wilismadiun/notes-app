package usecase

import domains "notes-app/src/modul/note/Domains"

type DeleteNote struct {
	Repo domains.NotesRepository
}

func (h *DeleteNote) Execute(id string) (string, error) {
	err := h.Repo.DeleteNoteById(id)
	if err != nil {
		return "", err
	}

	return id, nil
}
