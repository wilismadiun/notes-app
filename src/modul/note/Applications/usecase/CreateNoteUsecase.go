package usecase

import (
	"notes-app/src/modul/note/Applications/generator"
	domains "notes-app/src/modul/note/Domains"
	"notes-app/src/modul/note/Domains/entities"
	"time"
)

type CreateNote struct {
	Generator generator.IdGenerator
	Repo      domains.NotesRepository
}

func (h *CreateNote) Execute(note entities.Note) (entities.CreateNoteResponse, error) {
	now := time.Now()

	note.CreateAt = now
	note.UpdateAt = now
	note.ID = h.Generator.Generator()

	err := entities.VerifyNote(note)
	if err != nil {
		return entities.CreateNoteResponse{}, err
	}

	err = h.Repo.CreateNote(note)
	if err != nil {
		return entities.CreateNoteResponse{}, err
	}

	return entities.CreateNoteResponse{
		ID:    note.ID,
		Title: note.Title,
	}, nil
}
