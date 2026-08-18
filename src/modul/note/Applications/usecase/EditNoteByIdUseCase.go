package usecase

import (
	"errors"
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
	"time"
)

type EditNoteById struct {
	Repo domains.NotesRepository
}

func (h *EditNoteById) execute(id string, payload entities.EditNoteRequest) (string, error) {
	note, err := h.Repo.FindNoteById(id)
	if err != nil {
		return "", err
	}

	update := make(map[string]any)
	if payload.Title != nil {
		update["title"] = *payload.Title
	}
	if payload.Content != nil {
		update["content"] = *payload.Content
	}

	if len(update) == 0 {
		return "", errors.New("Tidak ada data yang dikirim untuk diubah")
	}

	update["updateAt"] = time.Now().Truncate(time.Microsecond)

	err = h.Repo.EditNoteById(note, update)
	if err != nil {
		return "", err
	}

	return id, nil
}
